package pagination

import (
	"fmt"
	"time"
)

type KeysetPaginated[T any, TId fmt.Stringer] struct {
	Data       []T
	PrevCursor *KeysetCursor[TId]
	NextCursor *KeysetCursor[TId]
}

func NewKeysetPaginated[T any, TId fmt.Stringer](
	data []T,
	pagination *KeysetPagination[TId],
	id func(T) TId,
	createdAt func(T) time.Time,
) *KeysetPaginated[T, TId] {
	direction := pagination.Direction()

	hasMore := uint(len(data)) > pagination.Limit

	if hasMore {
		data = data[:len(data)-1]
	}

	var prevCursor *KeysetCursor[TId]
	var nextCursor *KeysetCursor[TId]

	if len(data) > 0 {
		switch direction {
		case Forward:
			if pagination.NextCursor != nil {
				prevCursor = &KeysetCursor[TId]{
					Id:        id(data[0]),
					CreatedAt: createdAt(data[0]),
				}
			}

			if hasMore {
				last := data[len(data)-1]

				nextCursor = &KeysetCursor[TId]{
					Id:        id(last),
					CreatedAt: createdAt(last),
				}
			}

		case Backward:
			if hasMore {
				last := data[len(data)-1]

				prevCursor = &KeysetCursor[TId]{
					Id:        id(last),
					CreatedAt: createdAt(last),
				}
			}

			if pagination.PrevCursor != nil {
				first := data[0]

				nextCursor = &KeysetCursor[TId]{
					Id:        id(first),
					CreatedAt: createdAt(first),
				}
			}

			for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

	return &KeysetPaginated[T, TId]{
		Data:       data,
		PrevCursor: prevCursor,
		NextCursor: nextCursor,
	}
}
