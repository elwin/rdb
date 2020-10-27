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

func (q Query) ToSql() string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", q.attributes, q.table, q.condition)
}

func (q Query) Table(table string) Query {
	q.table = table

	return q
}

func (q Query) Select(attributes ...string) Query {
	q.attributes = strings.Join(attributes, ", ")

	return q
}

func (q Query) Where(attribute, operator, value string) Query {
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

type Condition string

func (c Condition) OrWhere(attribute, operator, value string) Condition {
	statement := Condition(fmt.Sprintf("%s %s %s", attribute, operator, value)) // not injection safe

	if c == "" {
		return statement
	}

	return " OR " + statement
}

func (c Condition) Where(attribute, operator, value string) Condition {
	statement := Condition(fmt.Sprintf("%s %s %s", attribute, operator, value)) // not injection safe

	if c == "" {
		return statement
	}

	return c + " AND " + statement
}

func (c Condition) OrWhereClosure(f func(c Condition) Condition) Condition {
	statement := f("")

	if c == "" {
		return statement
	}

	return c + " OR (" + statement + ")"
}
