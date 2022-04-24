package controllers

import (
	"cve-sa-backend/controllers/manage"
	"cve-sa-backend/controllers/web"
)

type Controllers struct {
	Web    web.Controller
	Manage manage.Controller
}

var Con = new(Controllers)
