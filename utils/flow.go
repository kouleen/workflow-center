package utils

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
	"github.com/kouleen/common/message"
	"github.com/kouleen/common/pkg/ctxutil"
	"github.com/kouleen/idl/kitex_gen/rpc"
	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/modle"
)

var moduleHandler = map[string]func() string{
	rpc.USER_RPC_SERVER: func() string {
		return "user.flow.review.queue"
	},
	rpc.SYSTEM_RPC_SERVER: func() string {
		return "system.flow.review.queue"
	},
	rpc.ROOM_RPC_SERVER: func() string {
		return "room.flow.review.queue"
	},
	rpc.WORKFLOW_RPC_SERVER: func() string {
		return "workflow.flow.review.queue"
	},
	rpc.PAYMENT_RPC_SERVER: func() string {
		return "payment.flow.review.queue"
	},
	rpc.GUILD_RPC_SERVER: func() string {
		return "guild.flow.review.queue"
	},
	rpc.GIFT_RPC_SERVER: func() string {
		return "gift.flow.review.queue"
	},
	rpc.STREAM_RPC_SERVER: func() string {
		return "stream.flow.review.queue"
	},
}

func Review(ctx context.Context, req *workflow.WorkflowProcessRequest) (err error) {
	// 不是审批操作
	if !modle.IsAction(req.GetStatus()) {
		return nil
	}
	traceId := ctxutil.GetTraceId(ctx)
	userId := ctxutil.GetUserId(ctx)
	req.UpdatedBy = &userId
	messageByte, err := json.Marshal(req)
	if err != nil {
		return err
	}
	module := req.GetModule()
	if h, ok := moduleHandler[module]; ok {
		topic := h()
		go func() {
			if msgErr := message.SendToQueue(context.Background(), topic, &message.Message{
				TraceID:     traceId,
				MessageID:   strconv.Itoa(int(uuid.New().ID())),
				ContentType: "application/json",
				Body:        messageByte,
				Headers:     make(map[string]any),
			}); msgErr != nil {
				logger.CtxErrorf(ctx, "send login log error: %v", err)
			}
		}()
		return
	}
	return errors.New("module not found")
}
