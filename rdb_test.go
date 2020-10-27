package rdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCondition_OrWhere(t *testing.T) {

	query := New().Table("users").
		Select("name", "email").
		Where("age", ">", "18").ToSql()
	assert.Equal(t, "SELECT name, email FROM users WHERE age > 18", query)

	query = New().Table("users").
		Select("name", "email").
		Where("age", ">", "18").
		OrWhereClause(func(c Condition) Condition {
			return c.Where("age", ">", "16").
				Where("entitled", "=", "1")
		}).
		ToSql()
	assert.Equal(t, "SELECT name, email FROM users WHERE age > 18 OR (age > 16 AND entitled = 1)", query)
}
