package handlers

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Code:    400,
		Message: message,
		Error:   errMsg,
	}
}

// PageResponse 分页响应
func PageResponse(message string, items interface{}, total int64) Response {
	return Response{
		Code:    200,
		Message: message,
		Data: map[string]interface{}{
			"items": items,
			"total": total,
		},
	}
}
// UpdateDIDRequest 定义了更新 DID 时需要的请求参数
type UpdateDIDRequest struct {
	DID                string `json:"did" binding:"required"`
	RecoveryKey        string `json:"recovery_key" binding:"required"`
	RecoveryPrivateKey string `json:"recovery_private_key" binding:"required"`
}

// DeleteDIDRequest 定义了删除 DID 时需要的请求参数
type DeleteDIDRequest struct {
	DID                string `json:"did" binding:"required"`
	RecoveryKey        string `json:"recovery_key" binding:"required"`
	RecoveryPrivateKey string `json:"recovery_private_key" binding:"required"`
}
