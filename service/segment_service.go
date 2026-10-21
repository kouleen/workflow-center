package service

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/kouleen/common/pkg/ctxutil"
	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/modle"
	"github.com/kouleen/workflow-center/repository"
)

func QuerySegmentPage(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp *workflow.WorkflowSegmentPageResponse, err error) {
	list, total, err := repository.QuerySegmentPage(ctx, req)
	if err != nil {
		return
	}
	respList := make([]*workflow.WorkflowSegmentResponse, len(list))
	for index, item := range list {
		respList[index] = item.ConvertResp()
	}
	return &workflow.WorkflowSegmentPageResponse{
		Total:   total,
		Records: respList,
	}, err
}

func QuerySegmentList(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp []*workflow.WorkflowSegmentResponse, err error) {
	list, err := repository.QuerySegmentList(ctx, req)
	if err != nil {
		return
	}
	respList := make([]*workflow.WorkflowSegmentResponse, len(list))
	for index, item := range list {
		respList[index] = item.ConvertResp()
	}
	return respList, err
}

func QuerySegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp *workflow.WorkflowSegmentResponse, err error) {
	item, err := repository.QuerySegmentById(ctx, req.GetId())
	if err != nil {
		return
	}
	return item.ConvertResp(), err
}

func SaveSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return false, err
	}
	id := node.Generate().Int64()
	entity := &modle.WorkflowSegment{
		ID:        id,
		BizTag:    req.GetBizTag(),
		MaxID:     req.GetMaxId(),
		Step:      req.GetStep(),
		StartID:   req.GetStartId(),
		EndID:     req.GetEndId(),
		Remark:    req.GetRemark(),
		CreatedBy: ctxutil.GetUserId(ctx),
	}
	if err = repository.CreateSegment(ctx, entity); err != nil {
		return false, err
	}
	return true, nil
}

func UpdateSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	item, err := repository.QuerySegmentById(ctx, req.GetId())
	if err != nil {
		return
	}
	item.BizTag = req.GetBizTag()
	item.MaxID = req.GetMaxId()
	item.Step = req.GetStep()
	item.StartID = req.GetStartId()
	item.EndID = req.GetEndId()
	item.Status = req.Status
	item.Remark = req.GetRemark()
	item.UpdatedBy = ctxutil.GetUserId(ctx)
	if err = repository.UpdateSegment(ctx, item); err != nil {
		return
	}
	return true, nil
}

func DeleteSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	if err = repository.BatchDeleteSegment(ctx, req.GetIdList()); err != nil {
		return
	}
	return true, nil
}

func UpdateSegmentStatus(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	item, err := repository.QuerySegmentById(ctx, req.GetId())
	if err != nil {
		return
	}
	item.Status = req.Status
	item.UpdatedBy = ctxutil.GetUserId(ctx)
	if err = repository.UpdateSegment(ctx, item); err != nil {
		return
	}
	return true, nil
}
