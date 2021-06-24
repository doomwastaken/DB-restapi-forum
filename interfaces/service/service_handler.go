package service

import (
	"encoding/json"
	"fmt"
	"forum/application"
	"forum/domain/entity"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ServiceInfo struct {
	ServiceApp application.ServiceAppInterface
}

func NewServiceInfo(ServiceApp application.ServiceAppInterface) *ServiceInfo {
	return &ServiceInfo{
		ServiceApp: ServiceApp,
	}
}

func (serviceInfo *ServiceInfo) HandleClearData(ctx *fasthttp.RequestCtx) {
	err := serviceInfo.ServiceApp.ClearAllDate()
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf(`{"messege": "%s"}`, err.Error()),
		}
		body, err := json.Marshal(msg)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
		return
	}

	body, err := json.Marshal("")
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (serviceInfo *ServiceInfo) HandleGetDBStatus(ctx *fasthttp.RequestCtx) {
	status, err := serviceInfo.ServiceApp.GetDBStatus()
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf(`{"messege": "%s"}`, err.Error()),
		}
		body, err := json.Marshal(msg)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
		return
	}

	body, err := json.Marshal(status)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}
