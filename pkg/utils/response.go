package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"member-pre/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	if appErr, ok := errors.AsAppError(err); ok {
		c.JSON(appErr.HTTPStatus(), Response{
			Code:    int(appErr.Code),
			Message: appErr.Message,
		})
		return
	}

	// 未知错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(errors.ErrCodeUnknown),
		Message: "内部服务器错误",
	})
}

// ErrorWithCode 错误响应（指定错误码）
func ErrorWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// Pagination 分页响应
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// SuccessWithPagination 成功响应（分页）
func SuccessWithPagination(c *gin.Context, list interface{}, page, pageSize int, total int64) {
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PaginationResponse{
			List: list,
			Pagination: Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
				Pages:    pages,
			},
		},
	})
}
