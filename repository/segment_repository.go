package repository

import (
	"context"

	"github.com/kouleen/common/pkg/mysql"
	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/modle"
	"github.com/kouleen/workflow-center/utils"
	"gorm.io/gorm"
)

func QuerySegmentPage(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp []*modle.WorkflowSegment, total int64, err error) {
	query := getSegmentQuery(ctx, req)
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

func getSegmentQuery(ctx context.Context, req *workflow.WorkflowSegmentRequest) *gorm.DB {
	query := mysql.GetReadMysqlDDB().WithContext(ctx).Model(&modle.WorkflowSegment{}).Where("is_delete = 0")
	if req.GetBizTag() != "" {
		query.Where("biz_tag = ?", req.GetBizTag())
	}
	if req.Status != nil {
		query = query.Where("status = ?", req.GetStatus())
	}
	return query
}

func QuerySegmentList(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp []*modle.WorkflowSegment, err error) {
	query := getSegmentQuery(ctx, req)
	if req.Params != nil {
		if req.GetParams().GetBeginTime() != "" && req.GetParams().GetEndTime() != "" {
			startUTC, endUTC, err := utils.DateToLocalRange(req.GetParams().GetBeginTime(), req.GetParams().GetEndTime())
			if err != nil {
				return nil, err
			}
			query = query.Where("create_time between ? and ?", startUTC, endUTC)
		}
	}
	query = query.Order("create_time desc,id desc")
	if err = query.Find(&resp).Error; err != nil {
		return
	}
	return
}

func QuerySegmentById(ctx context.Context, id int64) (resp *modle.WorkflowSegment, err error) {
	if err = mysql.GetReadMysqlDDB().WithContext(ctx).Model(resp).First(&resp, id).Error; err != nil {
		return nil, err
	}
	return
}

func QuerySegmentByIdList(ctx context.Context, idList []int64) (resp []*modle.WorkflowSegment, err error) {
	if err = mysql.GetReadMysqlDDB().WithContext(ctx).Model(&modle.WorkflowSegment{}).Where("id in (?)", idList).Find(&resp).Error; err != nil {
		return nil, err
	}
	return
}

func QuerySegmentByBizTag(ctx context.Context, bizTag string) (resp []*modle.WorkflowSegment, err error) {
	query := mysql.GetReadMysqlDDB().WithContext(ctx).Model(&modle.WorkflowSegment{}).Where("is_delete = 0 and biz_tag = ?", bizTag)
	if err = query.Find(&resp).Error; err != nil {
		return nil, err
	}
	return
}

func CreateSegment(ctx context.Context, entity *modle.WorkflowSegment) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Create(entity).Error
}

func BatchCreateSegment(ctx context.Context, entityList []*modle.WorkflowSegment) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(entityList) > 0 {
			if err := tx.Model(&modle.WorkflowSegment{}).Create(entityList).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func UpdateSegment(ctx context.Context, entity *modle.WorkflowSegment) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Where("id = ?", entity.ID).Save(entity).Error
}

func BatchUpdateSegment(ctx context.Context, entityList []*modle.WorkflowSegment) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range entityList {
			if err := tx.Model(&modle.WorkflowSegment{}).Where("id = ?", item.ID).Save(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func BatchDeleteSegment(ctx context.Context, idList []int64) error {
	return mysql.GetWriteMysqlDDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(idList) > 0 {
			if err := tx.Where("id in (?)", idList).Delete(&modle.WorkflowSegment{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
