package pagination

type PageResponse struct {
	Total int64 `json:"total" `
	List  any   `json:"list" `
}

func NewPageResponse[T any](total int64, list []T) *PageResponse {
	return &PageResponse{Total: total, List: list}
}
