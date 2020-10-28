package rdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCondition_OrWhere(t *testing.T) {
	query, queryAttributes := New().Table("users").
		Select("name", "email").
		Where("age", ">", 18).ToSql()
	assert.Equal(t, "SELECT name, email FROM users WHERE age > ?", query)
	assert.Equal(t, attributes{18}, queryAttributes)

	query, queryAttributes = New().Table("users").
		Select("name", "email").
		Where("age", ">", 18).
		OrWhereClause(func(c Condition) Condition {
			return c.Where("age", ">", 16).
				Where("entitled", "=", 1)
		}).
		ToSql()
	assert.Equal(t, "SELECT name, email FROM users WHERE age > ? OR (age > ? AND entitled = ?)", query)
	assert.Equal(t, attributes{18, 16, 1}, queryAttributes)
}
