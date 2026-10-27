#!/bin/bash
set -e

# ===== 加载 .env 文件 =====
if [ -f ".env" ];then
  export $(grep -v '^#' .env | xargs)
fi

# 配置
ROOT_PASS="${MYSQL_ROOT_PASSWORD}"
REPL_USER="${REPL_USER}"
REPL_PASS="${REPL_PASS}"
MASTER_CONTAINER="mysql-master"
SLAVE3_CONTAINER="mysql-slave3"

# 函数：校验从库主从状态
check_slave_status() {
  local container_name=$1
  echo "👉 检查 $container_name 主从状态"

  ret=$(docker exec -i "${container_name}" mysql -uroot -p"${ROOT_PASS}" -e "SHOW SLAVE STATUS\G")
  io_running=$(echo "$ret" | grep 'Slave_IO_Running:' | awk '{print $2}')
  sql_running=$(echo "$ret" | grep 'Slave_SQL_Running:' | awk '{print $2}')

  echo "   Slave_IO_Running:  $io_running"
  echo "   Slave_SQL_Running: $sql_running"

  if [[ "$io_running" == "Yes" && "$sql_running" == "Yes" ]]; then
    echo "✅ $container_name 主从复制正常"
    return 0
  else
    echo "❌ $container_name 主从复制异常！"
    echo "$ret"
    exit 1
  fi
}

echo "====> 等待 MySQL Master 就绪，最多等待60s"
for i in $(seq 1 60); do
  if docker exec ${MASTER_CONTAINER} mysql -uroot -p"${ROOT_PASS}" -e "select 1;" >/dev/null 2>&1; then
    echo "✅ Master 数据库就绪"
    break
  fi
  echo "等待Master就绪 $i ..."
  sleep 1
done

echo "====> 2. 获取Master binlog File & Position"
eval $(docker exec -i ${MASTER_CONTAINER} mysql -uroot -p"${ROOT_PASS}" -N -s -e "SHOW MASTER STATUS;" | awk '{print "BINLOG_FILE="$1"; BINLOG_POS="$2}')
echo "Detect: BINLOG_FILE=${BINLOG_FILE}, BINLOG_POS=${BINLOG_POS}"

echo "====> 3. 配置 slave3"
docker exec -i ${SLAVE3_CONTAINER} mysql -uroot -p"${ROOT_PASS}" <<EOF
STOP SLAVE;
RESET SLAVE ALL;
CHANGE MASTER TO
MASTER_HOST='${MASTER_CONTAINER}',
MASTER_USER='${REPL_USER}',
MASTER_PASSWORD='${REPL_PASS}',
MASTER_LOG_FILE='${BINLOG_FILE}',
MASTER_LOG_POS=${BINLOG_POS};
START SLAVE;
EOF

echo "====> 5. Master解锁"
docker exec -i ${MASTER_CONTAINER} mysql -uroot -p"${ROOT_PASS}" -e "UNLOCK TABLES;"

echo "====> 检查slave3状态"
check_slave_status ${SLAVE3_CONTAINER}

echo "🎉 全部从库校验通过，主从集群初始化成功！"