package repository

import (
	"github.com/pkg/errors"
)

var ErrorNoRowsAffected = errors.New("no row affected")

type QueryCondition[T any] interface {
	Where(query string, args ...interface{}) T
	WhereOr(query string, args ...interface{}) T
	WhereGroup(query string, fn func(T) T) T
}

type UpdateCondition[T any] interface {
	Where(query string, args ...interface{}) T
	Set(query string, args ...interface{}) T
}

func sqlLike(s string) string {
	return "%" + s + "%"
}
