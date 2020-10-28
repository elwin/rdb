package rdb

import (
	"fmt"
	"strings"
)

type Query struct {
	table      string
	attributes string
	condition  Condition
}

func New() Query {
	return Query{}
}

func (q Query) ToSql() (string, attributes) {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", q.attributes, q.table, q.condition.query), q.condition.attributes
}

func (q Query) Table(table string) Query {
	q.table = table

	return q
}

func (q Query) Select(attributes ...string) Query {
	q.attributes = strings.Join(attributes, ", ")

	return q
}

func (q Query) Where(attribute, operator string, value interface{}) Query {
	q.condition = q.condition.Where(attribute, operator, value)

	return q
}

func (q Query) OrWhere(attribute, operator, value string) Query {
	q.condition = q.condition.OrWhere(attribute, operator, value)

	return q
}

func (q Query) OrWhereClause(f func(c Condition) Condition) Query {
	q.condition = q.condition.OrWhereClosure(f)

	return q
}

type Condition struct {
	query      string
	attributes attributes
}

type attributes []interface{}

func (a attributes) join(b attributes) attributes {
	c := make(attributes, len(a)+len(b))
	copy(c[:len(a)], a)
	copy(c[len(a):], b)

	return c
}

func (c Condition) join(glue string, b Condition) Condition {
	if c.query == "" {
		return b
	}

	return Condition{
		query:      fmt.Sprintf("%s %s %s", c.query, glue, b.query),
		attributes: c.attributes.join(b.attributes),
	}
}

func (c Condition) OrWhere(attribute, operator string, value interface{}) Condition {
	statement := Condition{
		query:      fmt.Sprintf("%s %s ?", attribute, operator),
		attributes: []interface{}{value},
	}

	return c.join("OR", statement)
}

func (c Condition) Where(attribute, operator string, value interface{}) Condition {
	statement := Condition{
		query:      fmt.Sprintf("%s %s ?", attribute, operator),
		attributes: []interface{}{value},
	}

	return c.join("AND", statement)
}

func (c Condition) OrWhereClosure(f func(c Condition) Condition) Condition {
	statement := f(Condition{})
	statement.query = fmt.Sprintf("(%s)", statement.query)

	return c.join("OR", statement)
}
