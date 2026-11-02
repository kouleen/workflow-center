package repository

import (
	"context"

	"github.com/kouleen/common/pkg/mysql"
	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/modle"
	"github.com/kouleen/workflow-center/utils"
	"gorm.io/gorm"
)

func QueryProcessPage(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp []*modle.WorkflowProcess, total int64, err error) {
	query := mysql.GetReadMysqlDDB().WithContext(ctx).Model(&modle.WorkflowProcess{}).Where("is_delete = 0")
	if req.GetFlowNo() != "" {
		query.Where("flow_no like ?", req.GetFlowNo()+"%")
	}
	if req.GetModule() != "" {
		query.Where("module = ?", req.GetModule())
	}
	if req.GetBusiness() != "" {
		query.Where("business = ?", req.GetBusiness())
	}
	if req.GetTradeNo() != "" {
		query.Where("trade_no like ?", req.GetTradeNo()+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", req.GetStatus())
	}
	if req.Params != nil {
		if req.GetParams().GetBeginTime() != "" && req.GetParams().GetEndTime() != "" {
			startUTC, endUTC, err := utils.DateToLocalRange(req.GetParams().GetBeginTime(), req.GetParams().GetEndTime())
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("create_time between ? and ?", startUTC, endUTC)
		}
	}
	if err = query.Count(&total).Error; err != nil || total == 0 {
		return
	}
	query = query.Order("create_time desc,id desc")
	i := (req.GetCurrent() - 1) * req.GetSize()
	if err = query.Offset(int(i)).Limit(int(req.GetSize())).Find(&resp).Error; err != nil {
		return
	}
	return
}

func QueryProcessById(ctx context.Context, id int64) (resp *modle.WorkflowProcess, err error) {
	if err = mysql.GetReadMysqlDDB().WithContext(ctx).Model(&modle.WorkflowProcess{}).Take(&resp, id).Error; err != nil {
		return nil, err
	}
	return
}

func CreateProcess(ctx context.Context, entity *modle.WorkflowProcess) (err error) {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Create(entity).Error
}

func UpdateProcess(ctx context.Context, entity *modle.WorkflowProcess) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Where("id = ?", entity.ID).Save(entity).Error
}

func BatchDeleteProcess(ctx context.Context, idList []int64) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(idList) > 0 {
			if err := tx.Where("id in (?)", idList).Delete(&modle.WorkflowProcess{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
