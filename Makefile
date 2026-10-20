thrift=thrift
idl_path = idl/workflow-center

gen_service:
	kitex -module github.com/kouleen/workflow-center -service workflow.service  $(idl_path)/workflow.thrift
