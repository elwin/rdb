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
	assert.Equal(t, Attributes{18}, queryAttributes)

	query, queryAttributes = New().Table("users").
		Select("name", "email").
		Where("age", ">", 18).
		OrWhereClause(func(c Condition) Condition {
			return c.Where("age", ">", 16).
				Where("entitled", "=", 1)
		}).
		ToSql()
	assert.Equal(t, "SELECT name, email FROM users WHERE age > ? OR (age > ? AND entitled = ?)", query)
	assert.Equal(t, Attributes{18, 16, 1}, queryAttributes)

	query, queryAttributes = New().Table("users").
		Select("age", "count(*)").
		GroupBy("age").
		ToSql()
	assert.Equal(t, "SELECT age, count(*) FROM users GROUP BY age", query)
	assert.Equal(t, Attributes(nil), queryAttributes)

	query, queryAttributes = New().Table("users").
		Select("*").
		OrderByAsc("id").
		OrderByDesc("age").
		ToSql()
	assert.Equal(t, "SELECT * FROM users ORDER BY id ASC, age DESC", query)
	assert.Equal(t, Attributes(nil), queryAttributes)
}
