package app

import (
	"geek/internal/dbeye/pkg/ce"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, code int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"status": code,
		"data":   data,
		"msg":    ce.GetMsg(code),
	})
	return
}
