package rdb

import (
	"fmt"
	"strings"
)

type Query struct {
	table      string
	attributes string
	condition  Condition
	groupBy    stringSlice
	orderBy    stringSlice
}

func New() Query {
	return Query{}
}

func (q Query) ToSql() (string, Attributes) {
	output := fmt.Sprintf("SELECT %s FROM %s", q.attributes, q.table)

	if q.condition.query != "" {
		output += " WHERE " + q.condition.query
	}

	if len(q.groupBy) > 0 {
		output += " GROUP BY " + strings.Join(q.groupBy, ", ")
	}

	if len(q.orderBy) > 0 {
		output += " ORDER BY " + strings.Join(q.orderBy, ", ")
	}

	return output, q.condition.attributes
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

func (q Query) GroupBy(attributes ...string) Query {
	q.groupBy = append(q.groupBy.Copy(), attributes...)

	return q
}

func (q Query) OrderByAsc(attribute string) Query {
	q.orderBy = append(q.orderBy.Copy(), attribute + " ASC")

	return q
}

func (q Query) OrderByDesc(attribute string) Query {
	q.orderBy = append(q.orderBy.Copy(), attribute + " DESC")

	return q
}

type Condition struct {
	query      string
	attributes Attributes
}

type Attributes []interface{}

func (a Attributes) join(b Attributes) Attributes {
	c := make(Attributes, len(a)+len(b))
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

type stringSlice []string
func (s stringSlice) Copy() stringSlice {
	out := make(stringSlice, len(s))
	copy(out, s)

	return out
}
