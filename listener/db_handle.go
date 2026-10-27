package listener

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/kouleen/common/message"
	"github.com/kouleen/common/pkg/ctxutil"
	"github.com/kouleen/common/pkg/logutil"
)

func init() {
	go func() {
		cfg := canal.NewDefaultConfig()
		mysqlAddressListener := os.Getenv("MYSQL_ADDR_LISTENER")
		if mysqlAddressListener == "" {
			log.Println("mysql_addr_listener is empty")
			return
		}
		mysqlUserListener := os.Getenv("MYSQL_USER_LISTENER")
		if mysqlUserListener == "" {
			log.Println("mysql_user_listener is empty")
			return
		}

		mysqlPwdListener := os.Getenv("MYSQL_PASSWORD_LISTENER")
		if mysqlPwdListener == "" {
			log.Println("mysql_passwd_listener is empty")
			return
		}
		cfg.Addr = mysqlAddressListener
		cfg.User = mysqlUserListener
		cfg.Password = mysqlPwdListener
		cfg.ServerID = 2001
		cfg.Dump.ExecutionPath = ""

		c, err := canal.NewCanal(cfg)
		if err != nil {
			log.Printf("new canal failed: %v", err)
			return
		}
		c.SetEventHandler(&DbHandle{})

		pos, err := getMasterPos(cfg.Addr, cfg.User, cfg.Password)
		if err != nil {
			log.Printf("get master position failed: %v", err)
			return
		}
		log.Printf("sync from the current latest position: %+v", pos)

		if err := c.RunFrom(pos); err != nil {
			log.Printf("canal run failed: %v", err)
			return
		}
	}()
}

// 查询当前最新位点（首次运行、无存档时用）
func getMasterPos(addr, user, pwd string) (mysql.Position, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4", user, pwd, addr)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return mysql.Position{}, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	queries := []string{
		"SHOW MASTER STATUS",
		"SHOW BINARY LOG STATUS",
	}

	var errs []string
	for _, query := range queries {
		masterPos, queryErr := queryMasterPos(db, query)
		if queryErr == nil {
			return masterPos, nil
		}
		errs = append(errs, fmt.Sprintf("%s: %v", query, queryErr))
	}

	return mysql.Position{}, fmt.Errorf("get master position failed: %s", strings.Join(errs, "; "))
}

func queryMasterPos(db *sql.DB, query string) (mysql.Position, error) {
	rows, err := db.Query(query)
	if err != nil {
		return mysql.Position{}, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	cols, colsErr := rows.Columns()
	if colsErr != nil {
		return mysql.Position{}, colsErr
	}
	if len(cols) < 2 {
		return mysql.Position{}, fmt.Errorf("unexpected column count: %d", len(cols))
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return mysql.Position{}, err
		}
		return mysql.Position{}, sql.ErrNoRows
	}

	raw := make([]sql.RawBytes, len(cols))
	dest := make([]any, len(cols))
	for i := range raw {
		dest[i] = &raw[i]
	}
	if scanErr := rows.Scan(dest...); scanErr != nil {
		return mysql.Position{}, scanErr
	}

	file := string(raw[0])
	if file == "" {
		return mysql.Position{}, fmt.Errorf("empty binlog file")
	}

	pos64, parseErr := strconv.ParseUint(string(raw[1]), 10, 32)
	if parseErr != nil {
		return mysql.Position{}, fmt.Errorf("invalid binlog position %q: %w", string(raw[1]), parseErr)
	}

	return mysql.Position{Name: file, Pos: uint32(pos64)}, nil
}

type DbHandle struct {
	canal.DummyEventHandler
}

// rowToMap []any 行数据 → map[字段名]值
func rowToMap(table *schema.Table, row []any) map[string]any {
	ret := make(map[string]any, len(table.Columns))
	for idx, col := range table.Columns {
		if idx < len(row) {
			ret[col.Name] = row[idx]
		}
	}
	return ret
}

type BinlogMsg struct {
	DB     string         `json:"db"`
	Table  string         `json:"table"`
	Action string         `json:"action"`
	Before map[string]any `json:"before,omitempty"`
	After  map[string]any `json:"after,omitempty"`
}

func sendBinlogMsg(ctx context.Context, msg BinlogMsg) error {
	messageByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return message.SendToQueue(ctx, "workflow.data.listener", &message.Message{
		TraceID:     ctxutil.GetTraceId(ctx),
		MessageID:   strconv.Itoa(int(uuid.New().ID())),
		ContentType: "application/json",
		Body:        messageByte,
		Headers:     make(map[string]any),
	})
}

func (p *DbHandle) OnRow(e *canal.RowsEvent) error {
	ctx := metainfo.WithPersistentValue(context.Background(), "x-trace-id", uuid.New().String())
	dbName := e.Table.Schema
	tableName := e.Table.Name
	action := e.Action
	logutil.CtxInfo(ctx, "db handle: row event dbName: %s,tableName: %s,action:%s", dbName, tableName, action)

	switch action {
	case canal.InsertAction:
		for _, row := range e.Rows {
			afterMap := rowToMap(e.Table, row)
			if err := sendBinlogMsg(ctx, BinlogMsg{
				DB:     dbName,
				Table:  tableName,
				Action: action,
				After:  afterMap,
			}); err != nil {
				return err
			}
		}
	case canal.DeleteAction:
		for _, row := range e.Rows {
			beforeMap := rowToMap(e.Table, row)
			if err := sendBinlogMsg(ctx, BinlogMsg{
				DB:     dbName,
				Table:  tableName,
				Action: action,
				Before: beforeMap,
			}); err != nil {
				return err
			}
		}
	case canal.UpdateAction:
		for i := 0; i+1 < len(e.Rows); i += 2 {
			beforeMap := rowToMap(e.Table, e.Rows[i])
			afterMap := rowToMap(e.Table, e.Rows[i+1])
			if err := sendBinlogMsg(ctx, BinlogMsg{
				DB:     dbName,
				Table:  tableName,
				Action: action,
				Before: beforeMap,
				After:  afterMap,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
