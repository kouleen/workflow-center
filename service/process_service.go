package service

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/kouleen/common/pkg/ctxutil"
	"github.com/kouleen/common/pkg/redis"
	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/modle"
	"github.com/kouleen/workflow-center/repository"
	"github.com/kouleen/workflow-center/utils"
)

func QueryProcessPage(ctx context.Context, req *workflow.WorkflowProcessRequest) (r *workflow.WorkflowProcessPageResponse, err error) {
	list, total, err := repository.QueryProcessPage(ctx, req)
	if err != nil {
		return
	}
	respList := make([]*workflow.WorkflowProcessResponse, len(list))
	for index, item := range list {
		respList[index] = item.ConvertResp()
	}
	return &workflow.WorkflowProcessPageResponse{
		Total:   total,
		Records: respList,
	}, err
}

func QueryProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp *workflow.WorkflowProcessResponse, err error) {
	item, err := repository.QueryProcessById(ctx, req.GetId())
	if err != nil {
		return
	}
	return item.ConvertResp(), err
}

func CreateProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp *workflow.WorkflowProcessResponse, err error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	id := node.Generate().Int64()
	process := &redis.CodeProcess{}
	entity := &modle.WorkflowProcess{
		ID:        id,
		FlowNo:    process.GenerateCode(ctx, &modle.ProcessCodeRule{}),
		Module:    req.GetModule(),
		Business:  req.GetBusiness(),
		TradeNo:   req.GetTradeNo(),
		Remark:    req.GetRemark(),
		CreatedBy: ctxutil.GetUserId(ctx),
	}
	if err = repository.CreateProcess(ctx, entity); err != nil {
		return nil, err
	}
	return entity.ConvertResp(), nil
}

func UpdateProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp bool, err error) {
	if err = utils.Review(ctx, req); err != nil {
		return false, err
	}
	inter, err := repository.QueryProcessById(ctx, req.GetId())
	if err != nil {
		return false, err
	}
	inter.FlowNo = req.GetFlowNo()
	inter.Module = req.GetModule()
	inter.Business = req.GetBusiness()
	inter.TradeNo = req.GetTradeNo()
	inter.Status = req.Status
	inter.Remark = req.GetRemark()
	inter.UpdatedBy = ctxutil.GetUserId(ctx)
	if err = repository.UpdateProcess(ctx, inter); err != nil {
		return false, err
	}
	return true, nil
}

func DeleteProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp bool, err error) {
	if err = repository.BatchDeleteProcess(ctx, req.GetIdList()); err != nil {
		return
	}
	return true, nil
}
