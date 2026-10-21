package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/kouleen/common/bootstrap"
	"github.com/kouleen/common/middleware"
	"github.com/kouleen/idl/kitex_gen/rpc"
	"github.com/kouleen/idl/kitex_gen/workflow/workflowservice"
)

func main() {
	bootstrap.Run(rpc.WORKFLOW_RPC_SERVER, func(option ...server.Option) server.Server {
		return workflowservice.NewServer(new(WorkflowServiceImpl), option...)
	}, bootstrap.WithServerMiddleware(middleware.RpcServerMiddleware))
}
