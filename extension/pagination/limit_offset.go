package pagination

const (
	DEFAULT_PAGE_SIZE int = 10
	MAX_PAGE_SIZE     int = 100
)

type LimitOffsetPagination struct {
	PrimitivePage     int `json:"page" form:"page" binding:"gte=0"`         // 页码
	PrimitivePageSize int `json:"pageSize" form:"pageSize" binding:"gte=0"` // 每页数量
}

func NewPage(page, pageSize int) LimitOffsetPagination {
	return LimitOffsetPagination{
		PrimitivePage:     page,
		PrimitivePageSize: pageSize,
	}
}

func (p LimitOffsetPagination) Page() int {
	if p.PrimitivePage == 0 || p.PrimitivePage == 1 {
		return 1
	}
	return p.PrimitivePage
}

func (p LimitOffsetPagination) PageSize() int {
	if p.PrimitivePageSize <= 0 {
		return DEFAULT_PAGE_SIZE
	}
	return p.PrimitivePageSize
}

func (p LimitOffsetPagination) Offset() int {
	if p.PrimitivePage <= 1 {
		return 0
	}
	return (p.PrimitivePage - 1) * p.PrimitivePageSize
}

func (p LimitOffsetPagination) NextPage(total int) int {
	if p.Page() >= p.TotalPage(total) {
		return p.Page()
	}
	return p.Page() + 1
}

func (p LimitOffsetPagination) TotalPage(total int) int {
	if total == 0 {
		return 1
	}

	if total%p.PageSize() == 0 {
		return total / p.PageSize()
	}

	return total/p.PageSize() + 1
}
