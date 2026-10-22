package handle

import (
	"context"
	"errors"

	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/service"
)

func QuerySegmentPage(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp *workflow.WorkflowSegmentPageResponse, err error) {
	return service.QuerySegmentPage(ctx, req)
}

func QuerySegmentList(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp []*workflow.WorkflowSegmentResponse, err error) {
	return service.QuerySegmentList(ctx, req)
}

func QuerySegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp *workflow.WorkflowSegmentResponse, err error) {
	return service.QuerySegment(ctx, req)
}

func SaveSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	if err = checkSaveSegmentParams(req); err != nil {
		return
	}
	return service.SaveSegment(ctx, req)
}

func UpdateSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	if err = checkUpdateSegmentParams(req); err != nil {
		return
	}
	return service.UpdateSegment(ctx, req)
}

func DeleteSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	if len(req.GetIdList()) == 0 {
		return false, errors.New("empty id list")
	}
	return service.DeleteSegment(ctx, req)
}

func UpdateSegmentStatus(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	if req.Id == nil || req.Status == nil {
		return false, errors.New("invalid param")
	}
	return service.UpdateSegmentStatus(ctx, req)
}

func checkSaveSegmentParams(req *workflow.WorkflowSegmentRequest) error {
	if req.GetBizTag() == "" {
		return errors.New("biz tag is required")
	}
	if req.MaxId == nil {
		return errors.New("max id is required")
	}
	if req.Step == nil {
		return errors.New("step is required")
	}
	if req.StartId == nil {
		return errors.New("start id is required")
	}
	if req.EndId == nil {
		return errors.New("end id is required")
	}
	return nil
}

func checkUpdateSegmentParams(req *workflow.WorkflowSegmentRequest) error {
	if req.Id == nil {
		return errors.New("id is required")
	}
	return checkSaveSegmentParams(req)
}
