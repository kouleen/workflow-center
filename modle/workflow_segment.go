package modle

import (
	"time"

	"github.com/kouleen/idl/kitex_gen/workflow"
)

type WorkflowSegment struct {
	ID         int64      `json:"id,string" gorm:"column:id;primary_key;not null"`
	BizTag     string     `json:"bizTag" gorm:"column:biz_tag;not null;index:idx_biz_tag"`
	MaxID      int64      `json:"maxId,string" gorm:"column:max_id;not null"`
	Step       int64      `json:"step" gorm:"column:step;not null"`
	StartID    int64      `json:"startId,string" gorm:"column:start_id;not null"`
	EndID      int64      `json:"endId,string" gorm:"column:end_id;not null"`
	Status     *int8      `json:"status" gorm:"column:status;not null;default:1"`
	Remark     string     `json:"remark" gorm:"column:remark;default:''"`
	IsDelete   *int8      `json:"isDelete" gorm:"column:is_delete;not null;default:0"`
	CreatedBy  int64      `json:"createdBy,string" gorm:"column:created_by;not null;default:-1"`
	UpdatedBy  int64      `json:"updatedBy,string" gorm:"column:updated_by;not null;default:-1"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (p *WorkflowSegment) TableName() string {
	return "workflow_segment"
}

func (p *WorkflowSegment) ConvertResp() *workflow.WorkflowSegmentResponse {
	return &workflow.WorkflowSegmentResponse{
		Id:         p.ID,
		BizTag:     p.BizTag,
		MaxId:      p.MaxID,
		Step:       p.Step,
		StartId:    p.StartID,
		EndId:      p.EndID,
		Status:     p.Status,
		Remark:     p.Remark,
		IsDelete:   p.IsDelete,
		CreatedBy:  p.CreatedBy,
		UpdatedBy:  p.UpdatedBy,
		CreateTime: p.CreateTime.UnixMilli(),
	}
}
