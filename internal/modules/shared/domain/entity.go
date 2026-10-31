package domain

import "time"

type Entity[T comparable] struct {
	id        T
	createdAt time.Time
	updatedAt time.Time
}

func (e *Entity[T]) Id() T {
	return e.id
}

func (e *Entity[T]) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity[T]) UpdatedAt() time.Time {
	return e.updatedAt
}

func NewEntity[T comparable](id T) Entity[T] {
	now := time.Now()

	return Entity[T]{
		id:        id,
		createdAt: now,
		updatedAt: now,
	}
}

func RestoreEntity[T comparable](id T, createdAt time.Time, updatedAt time.Time) Entity[T] {
	return Entity[T]{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
