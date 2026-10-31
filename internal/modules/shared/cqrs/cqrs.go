package cqrs

import "context"

type CommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

type QueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, command TQuery) (TResult, error)
}
