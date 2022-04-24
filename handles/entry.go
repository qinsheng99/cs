package handles

import (
	"cve-sa-backend/handles/manage"
	"cve-sa-backend/handles/web"
)

type HandlesGroup struct {
	Manage manage.Handles
	Web    web.Handles
}

var Handles = new(HandlesGroup)

var (
	UploadHandle = Handles.Manage.UploadHandle
)

var (
	CveDatabaseHandle = Handles.Web.CveDatabaseHandle
	DriverHandle      = Handles.Web.DriverHandle
	HardwareHandle    = Handles.Web.HardwareHandle
	OsvHandle         = Handles.Web.OsvHandle
	SecurityHandle    = Handles.Web.SecurityHandle
	EsHandle          = Handles.Web.EsHandle
)
