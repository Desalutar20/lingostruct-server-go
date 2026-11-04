package pagination

import "fmt"

type KeysetDirection int

const (
	Backward KeysetDirection = iota
	Forward
)

type KeysetPagination[TId fmt.Stringer] struct {
	PrevCursor *KeysetCursor[TId]
	NextCursor *KeysetCursor[TId]
	Limit      uint
}

func (p *KeysetPagination[TId]) Direction() KeysetDirection {
	if p.PrevCursor != nil {
		return Backward
	}

	return Forward
}

func (p *KeysetPagination[TId]) Cursor() *KeysetCursor[TId] {
	if p.PrevCursor != nil {
		return p.PrevCursor
	}

	return p.NextCursor
}
