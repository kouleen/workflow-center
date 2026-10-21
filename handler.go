package main

import (
	"context"

	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/handle"
)

// WorkflowServiceImpl implements the last service interface defined in the IDL.
type WorkflowServiceImpl struct{}

// QueryWorkflowSegmentPage implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) QueryWorkflowSegmentPage(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp *workflow.WorkflowSegmentPageResponse, err error) {
	return handle.QuerySegmentPage(ctx, req)
}

// QueryWorkflowSegmentList implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) QueryWorkflowSegmentList(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp []*workflow.WorkflowSegmentResponse, err error) {
	return handle.QuerySegmentList(ctx, req)
}

// QueryWorkflowSegment implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) QueryWorkflowSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp *workflow.WorkflowSegmentResponse, err error) {
	return handle.QuerySegment(ctx, req)
}

// SaveWorkflowSegment implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) SaveWorkflowSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return handle.SaveSegment(ctx, req)
}

// UpdateWorkflowSegment implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) UpdateWorkflowSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return handle.UpdateSegment(ctx, req)
}

// UpdateWorkflowSegmentStatus implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) UpdateWorkflowSegmentStatus(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return handle.UpdateSegmentStatus(ctx, req)
}

// DeleteWorkflowSegment implements the WorkflowServiceImpl interface.
func (s *WorkflowServiceImpl) DeleteWorkflowSegment(ctx context.Context, req *workflow.WorkflowSegmentRequest) (resp bool, err error) {
	return handle.DeleteSegment(ctx, req)
}
