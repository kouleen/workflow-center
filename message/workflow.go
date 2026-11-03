package message

import (
	"context"
	"encoding/json"

	"github.com/kouleen/common/message"
	"github.com/kouleen/common/pkg/logutil"
	"github.com/kouleen/idl/kitex_gen/common"
	"github.com/kouleen/idl/kitex_gen/workflow"
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
}
