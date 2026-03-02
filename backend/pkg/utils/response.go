package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// ErrorWithDetails 带详细错误信息的响应
func ErrorWithDetails(c *fiber.Ctx, code int, message string, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Code:      code,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now().Unix(),
	})
}

// Pagination 分页响应
type Pagination struct {
	Total       int64       `json:"total"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
	Data        interface{} `json:"data"`
}

// SuccessWithPagination 分页成功响应
func SuccessWithPagination(c *fiber.Ctx, data interface{}, total int64, page, pageSize int) error {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return c.JSON(Response{
		Code:    0,
		Message: "success",
		Data: Pagination{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			Data:       data,
		},
		Timestamp: time.Now().Unix(),
	})
}
