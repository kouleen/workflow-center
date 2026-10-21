package handle

import (
	"context"

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
	return service.SaveSegment(ctx, req)
}

func UpdateSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return service.UpdateSegment(ctx, req)
}

func DeleteSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return service.DeleteSegment(ctx, req)
}

func UpdateSegmentStatus(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return service.UpdateSegmentStatus(ctx, req)
}
