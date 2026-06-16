package model

type PageRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type PageResponse struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Items interface{} `json:"items"`
}

func (p *PageRequest) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

func (p *PageRequest) Limit() int {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return p.PageSize
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{Code: 0, Message: "success", Data: data}
}

func ErrorResponse(code int, msg string) APIResponse {
	return APIResponse{Code: code, Message: msg}
}
