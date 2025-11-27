package utils

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

// Error 错误响应
func Error(msg string) Response {
	return Response{
		Code: 500,
		Msg:  msg,
	}
}

// ErrorWithCode 带状态码的错误响应
func ErrorWithCode(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
	}
}

// PageData 分页数据结构
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"pageSize"`
}

// NewPageData 创建分页数据
func NewPageData(list interface{}, total int64, page, size int) PageData {
	return PageData{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}
}
