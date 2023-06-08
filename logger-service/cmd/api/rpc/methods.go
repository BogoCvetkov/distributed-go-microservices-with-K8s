package rpc_server

import (
	"fmt"
	"logger-service/cmd/api/helpers"
	data "logger-service/cmd/api/models"
)

type RPCServer struct {
}

type RPCNewLog struct {
	EventName string
	Data      string
}

func (r *RPCServer) LogEntry(payload RPCNewLog, repy *helpers.JsonResponse) error {

	doc := data.LogEntry{
		Name: payload.EventName,
		Data: payload.Data,
	}

	err := doc.Insert(&doc)

	if err != nil {
		fmt.Println(err)
		return err
	}

	res := helpers.JsonResponse{
		Err:     false,
		Message: "Log inserted",
	}

	*repy = res

	return nil
}
