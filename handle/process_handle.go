package handle

import (
	"context"
	"errors"

	"github.com/kouleen/idl/kitex_gen/workflow"
	"github.com/kouleen/workflow-center/service"
)

func QueryProcessPage(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp *workflow.WorkflowProcessPageResponse, err error) {
	return service.QueryProcessPage(ctx, req)
}

func QueryProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp *workflow.WorkflowProcessResponse, err error) {
	return service.QueryProcess(ctx, req)
}

func CreateProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp *workflow.WorkflowProcessResponse, err error) {
	if err = checkCreateProcessParams(req); err != nil {
		return nil, err
	}
	return service.CreateProcess(ctx, req)
}

func UpdateProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp bool, err error) {
	if err = checkUpdateProcessParams(req); err != nil {
		return false, err
	}
	return service.UpdateProcess(ctx, req)
}

func DeleteProcess(ctx context.Context, req *workflow.WorkflowProcessRequest) (resp bool, err error) {
	if len(req.GetIdList()) == 0 {
		return false, errors.New("id list is required")
	}
	return service.DeleteProcess(ctx, req)
}

func checkCreateProcessParams(req *workflow.WorkflowProcessRequest) (err error) {
	if req.GetTradeNo() == "" {
		return errors.New("tradeNo is required")
	}
	if req.GetModule() == "" {
		return errors.New("module is required")
	}
	if req.GetBusiness() == "" {
		return errors.New("business is required")
	}
	return nil
}

func checkUpdateProcessParams(req *workflow.WorkflowProcessRequest) (err error) {
	if req.Id == nil {
		return errors.New("id is required")
	}
	if req.Status == nil {
		return errors.New("status is required")
	}
	return checkCreateProcessParams(req)
}
