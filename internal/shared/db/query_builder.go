package db

import (
	"app/internal/shared/entity"
	"strings"

	"gorm.io/gorm"
)

// QueryBuilder 动态查询构建器
type QueryBuilder struct {
	db      *gorm.DB
	tables  []string
	joins   []string
	where   []string
	params  []interface{}
	orderBy []string
	groupBy []string
	having  []string
	limit   *int
	offset  *int
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		db: db,
	}
}

// Table 设置主表
func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.tables = append(qb.tables, table)
	return qb
}

// Join 添加JOIN语句
func (qb *QueryBuilder) Join(join string, args ...interface{}) *QueryBuilder {
	// 简化处理，实际项目中可能需要更复杂的参数处理
	qb.joins = append(qb.joins, join)
	return qb
}

// LeftJoin 添加LEFT JOIN
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	join := "LEFT JOIN " + table + " ON " + condition
	qb.joins = append(qb.joins, join)
	return qb
}

// InnerJoin 添加INNER JOIN
func (qb *QueryBuilder) InnerJoin(table, condition string) *QueryBuilder {
	join := "INNER JOIN " + table + " ON " + condition
	qb.joins = append(qb.joins, join)
	return qb
}

// Where 添加WHERE条件
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.where = append(qb.where, condition)
	if len(args) > 0 {
		qb.params = append(qb.params, args...)
	}
	return qb
}

// WhereIf 条件性添加WHERE条件
func (qb *QueryBuilder) WhereIf(condition bool, sql string, args ...interface{}) *QueryBuilder {
	if condition {
		return qb.Where(sql, args...)
	}
	return qb
}

// OrWhere 添加OR条件
func (qb *QueryBuilder) OrWhere(condition string, args ...interface{}) *QueryBuilder {
	if len(qb.where) > 0 {
		qb.where[len(qb.where)-1] = "(" + qb.where[len(qb.where)-1] + ")"
		condition = "OR " + condition
	}
	return qb.Where(condition, args...)
}

// OrderBy 添加排序
func (qb *QueryBuilder) OrderBy(column string, desc bool) *QueryBuilder {
	order := column
	if desc {
		order += " DESC"
	}
	qb.orderBy = append(qb.orderBy, order)
	return qb
}

// GroupBy 添加分组
func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, columns...)
	return qb
}

// Having 添加HAVING条件
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, condition)
	if len(args) > 0 {
		qb.params = append(qb.params, args...)
	}
	return qb
}

// Limit 设置LIMIT
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = &limit
	return qb
}

// Offset 设置OFFSET
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = &offset
	return qb
}

// Build 构建最终查询
func (qb *QueryBuilder) Build() *gorm.DB {
	query := qb.db

	// 设置表
	if len(qb.tables) > 0 {
		query = query.Table(qb.tables[0])
	}

	// 添加JOIN
	for _, join := range qb.joins {
		query = query.Joins(join)
	}

	// 添加WHERE条件
	for i, condition := range qb.where {
		if i == 0 {
			query = query.Where(condition, qb.params...)
		} else {
			query = query.Where(condition)
		}
	}

	// 添加GROUP BY
	if len(qb.groupBy) > 0 {
		query = query.Group(strings.Join(qb.groupBy, ", "))
	}

	// 添加HAVING
	for _, having := range qb.having {
		query = query.Having(having)
	}

	// 添加ORDER BY
	if len(qb.orderBy) > 0 {
		query = query.Order(strings.Join(qb.orderBy, ", "))
	}

	// 添加LIMIT和OFFSET
	if qb.limit != nil {
		query = query.Limit(*qb.limit)
	}
	if qb.offset != nil {
		query = query.Offset(*qb.offset)
	}

	return query
}

// 示例：用户角色关联查询
func ExampleUserWithRolesQuery(db *gorm.DB, userID int64, roleName string, isActive bool) *gorm.DB {
	builder := NewQueryBuilder(db)

	return builder.
		Table("users u").
		LeftJoin("user_roles ur", "u.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		Where("u."+entity.FieldDeleteBy+" = 0").
		WhereIf(userID > 0, "u.id = ?", userID).
		WhereIf(roleName != "", "r.name = ?", roleName).
		WhereIf(isActive, "u.status = ?", "active").
		OrderBy("u."+entity.FieldCreateTime, true).
		Limit(100).
		Build()
}
