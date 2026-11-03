package message

import (
	"context"
	"encoding/json"

	"github.com/kouleen/common/message"
	"github.com/kouleen/common/pkg/ctxutil"
	"github.com/kouleen/common/pkg/logutil"
	"github.com/kouleen/idl/kitex_gen/common"
	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/modle"
	"github.com/kouleen/workflow-center/service"
)

type WorkflowReview struct{}

func (WorkflowReview) Queue() string {
	return common.WORKFLOW_FLOW_REVIEW_QUEUE
}

func (WorkflowReview) HandleMessage(ctx context.Context, msg *message.Message) error {
	defer func() {
		if err := recover(); err != nil {
			logutil.CtxError(ctx, "recover panic err: %#v", err)
		}
	}()
	var req workflow.WorkflowProcessRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return err
	}
	return nil
}

func init() {
	message.MustRegisterConsumer(new(WorkflowReview))
	message.MustRegisterConsumer(new(WorkflowFail))
}

type WorkflowFail struct{}

func (WorkflowFail) Queue() string {
	return "workflow.flow.fail.queue"
}

func (WorkflowFail) HandleMessage(ctx context.Context, msg *message.Message) error {
	ctx = ctxutil.SetMeta(ctx, msg.TraceID, -1)
	defer func() {
		if err := recover(); err != nil {
			logutil.CtxError(ctx, "recover panic err: %#v", err)
		}
	}()
	var req workflow.WorkflowProcessRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return err
	}
	status := modle.PREPARE
	req.Status = &status
	_, err := service.UpdateProcess(ctx, &req)
	if err != nil {
		logutil.CtxError(ctx, "update process err: %#v", err)
		return err
	}
	return nil
}
