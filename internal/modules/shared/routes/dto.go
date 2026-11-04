package routes

import (
	"fmt"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/pagination"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
)

type ResponseDTO[T any] struct {
	Data T `json:"data"`
}

type KeysetPaginatedResponseDto[T any] struct {
	Data       []T     `json:"data"`
	PrevCursor *string `json:"prevCursor"`
	NextCursor *string `json:"nextCursor"`
}

type KeysetPaginationRequest struct {
	Limit      *uint   `in:"query=limit" json:"limit"`
	PrevCursor *string `in:"query=prevCursor" json:"prevCursor"`
	NextCursor *string `in:"query=nextCursor" json:"nextCursor"`
}

type ErrorDTO struct {
	Error  string              `json:"error"`
	Code   string              `json:"code"`
	Errors map[string][]string `json:"errors"`
}

func NewKeysetPaginatedResponseDto[T any, TId fmt.Stringer](paginated *pagination.KeysetPaginated[T, TId]) KeysetPaginatedResponseDto[T] {
	var prevCursor *string
	var nextCursor *string

	if paginated.PrevCursor != nil {
		prevCursor = ptrutil.To(paginated.PrevCursor.String())
	}

	if paginated.NextCursor != nil {
		nextCursor = ptrutil.To(paginated.NextCursor.String())
	}

	return KeysetPaginatedResponseDto[T]{
		Data:       paginated.Data,
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
	}
}
