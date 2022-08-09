package router

import (
	"geek/global"
	apiV1 "geek/internal/dbeye/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,Token,X-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	//r.Use(gin.Recovery())
	//r.Use(middleware.Recovery())
	//r.Use(middleware.RequestID())
	//解决跨域问题
	r.Use(Cors())
	gin.SetMode(global.ServerSetting.RunMode)

	user := apiV1.NewUser()
	//openapi
	//配置路由
	apiGroup := r.Group("/api")
	//openGroup := r.Group("/geek/openapi")
	////只对apiv1的组进行token验证
	//apiGroup.Use(middleware.JWT())
	//api
	{
		apiGroup.GET("/v1/user", user.Get)
	}

	return r
}
