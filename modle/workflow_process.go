package modle

import (
	"time"

	"github.com/kouleen/idl/kitex_gen/workflow"
)

const (
	PREPARE  int8 = 1 // 待审批
	REVIEW   int8 = 2 // 处理中
	APPROVED int8 = 3 // 已批准
	REJECTED int8 = 4 // 已驳回
)

// IsAction 审核操作
func IsAction(status int8) bool {
	return status == REJECTED || status == APPROVED
}

type WorkflowProcess struct {
	ID         int64      `json:"id,string" gorm:"column:id;primary_key;not null"`
	FlowNo     string     `json:"flowNo" gorm:"column:flow_no;not null;unique:uq_flow_no"`
	Module     string     `json:"module" gorm:"column:module;not null"`
	Business   string     `json:"business" gorm:"column:business;not null"`
	TradeNo    string     `json:"tradeNo" gorm:"column:trade_no;not null;index:idx_trade_no"`
	Status     *int8      `json:"status" gorm:"column:status;not null;default:1"`
	Remark     string     `json:"remark" gorm:"column:remark;default:''"`
	IsDelete   *int8      `json:"isDelete" gorm:"column:is_delete;not null;default:0"`
	CreatedBy  int64      `json:"createdBy,string" gorm:"column:created_by;not null;default:-1"`
	UpdatedBy  int64      `json:"updatedBy,string" gorm:"column:updated_by;not null;default:-1"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (p *WorkflowProcess) TableName() string {
	return "workflow_process"
}

func (p *WorkflowProcess) ConvertResp() *workflow.WorkflowProcessResponse {
	return &workflow.WorkflowProcessResponse{
		Id:         p.ID,
		FlowNo:     p.FlowNo,
		Module:     p.Module,
		Business:   p.Business,
		TradeNo:    p.TradeNo,
		Status:     p.Status,
		Remark:     p.Remark,
		IsDelete:   p.IsDelete,
		CreatedBy:  p.CreatedBy,
		UpdatedBy:  p.UpdatedBy,
		CreateTime: p.CreateTime.UnixMilli(),
	}
}

type ProcessCodeRule struct{}

func (ProcessCodeRule) GetPrefix() string {
	return "FLW"
}

func (ProcessCodeRule) GetPattern() string {
	return "20060102"
}

func (ProcessCodeRule) GetDigit() int {
	return 5
}
