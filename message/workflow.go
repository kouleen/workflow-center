package message

import (
	"context"

	"github.com/kouleen/common/message"
	"github.com/kouleen/common/pkg/logutil"
	"github.com/kouleen/idl/kitex_gen/common"
)

type WorkflowReview struct{}

func (WorkflowReview) Queue() string {
	return common.WORKFLOW_FLOW_REVIEW_QUEUE
}

func (WorkflowReview) HandleMessage(ctx context.Context, msg *message.Message) error {
	logutil.CtxInfo(ctx, "message received")
	return nil
}

func init() {
	message.MustRegisterConsumer(new(WorkflowReview))
}
