package gin

import (
	"github.com/gin-gonic/gin"
)

type ControllerContext struct {
	ginContext *gin.Context
}

func (cc ControllerContext) Fail(statusCode int, err error) {
	cc.ginContext.JSON(statusCode, gin.H{"error": err.Error()})
	cc.ginContext.Abort()
}

func (cc ControllerContext) GetPathParam(key string) string {
	return cc.ginContext.Param(key)
}

func (cc ControllerContext) GetRequestHeader(key string) string {
	return cc.ginContext.Request.Header.Get(key)
}

func (cc ControllerContext) ParseBody(body interface{}) error {
	return cc.ginContext.BindJSON(body)
}

func (cc ControllerContext) ProcessNextMiddleware() {
	cc.ginContext.Next()
}

func (cc ControllerContext) Render(statusCode int, data interface{}) {
	cc.ginContext.JSON(statusCode, data)
}

func (cc ControllerContext) Set(key string, value interface{}) {
	cc.ginContext.Set(key, value)
}

var SharedControllerContextFactory = controllerContextFactory{}

func (ccf controllerContextFactory) CreateFromGinContext(ginContext *gin.Context) *ControllerContext {
	return &ControllerContext{ginContext: ginContext}
}

// PRIVATE INTERFACE

type controllerContextFactory struct{}
