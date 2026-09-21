package service

// IDReq 单个ID请求
type IDReq struct {
	ID int64 `json:"id"`
}

// IDsReq 多个ID请求
type IDsReq struct {
	IDs []int64 `json:"ids"`
}

// maxPageSize 最大分页大小
const maxPageSize = 1000

// PageReq 分页请求
type PageReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// normalize 规范化分页请求
func (p *PageReq) normalize(defaultSize int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}
