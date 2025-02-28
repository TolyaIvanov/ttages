package http

import (
	"github.com/gin-gonic/gin"
	"ttages/internal/file/usecase"
)

type HTTPHandler struct {
}

func NewHTTPHandler(router *gin.Engine, usecase *usecase.FileUsecase) {
}

func (h *HTTPHandler) GetFiles(c *gin.Context) {
}
