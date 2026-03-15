package db

import (
	"strings"

	"app/internal/shared/entity"

	"gorm.io/gorm"
)

// SQLBuilder 原生SQL构建器
type SQLBuilder struct {
	selectClause []string
	fromClause   string
	joinClauses  []string
	whereClause  []string
	params       []interface{}
	groupBy      []string
	havingClause []string
	orderBy      []string
	limit        *int
	offset       *int
}

// NewSQLBuilder 创建SQL构建器
func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

// Select 设置SELECT子句
func (sb *SQLBuilder) Select(columns ...string) *SQLBuilder {
	sb.selectClause = append(sb.selectClause, columns...)
	return sb
}

// From 设置FROM子句
func (sb *SQLBuilder) From(table string) *SQLBuilder {
	sb.fromClause = table
	return sb
}

// Join 添加JOIN子句
func (sb *SQLBuilder) Join(joinType, table, condition string) *SQLBuilder {
	join := joinType + " " + table + " ON " + condition
	sb.joinClauses = append(sb.joinClauses, join)
	return sb
}

// Where 添加WHERE条件
func (sb *SQLBuilder) Where(condition string, args ...interface{}) *SQLBuilder {
	sb.whereClause = append(sb.whereClause, condition)
	sb.params = append(sb.params, args...)
	return sb
}

// WhereIf 条件性添加WHERE条件
func (sb *SQLBuilder) WhereIf(condition bool, sql string, args ...interface{}) *SQLBuilder {
	if condition {
		return sb.Where(sql, args...)
	}
	return sb
}

// GroupBy 添加GROUP BY子句
func (sb *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	sb.groupBy = append(sb.groupBy, columns...)
	return sb
}

// Having 添加HAVING子句
func (sb *SQLBuilder) Having(condition string, args ...interface{}) *SQLBuilder {
	sb.havingClause = append(sb.havingClause, condition)
	sb.params = append(sb.params, args...)
	return sb
}

// OrderBy 添加ORDER BY子句
func (sb *SQLBuilder) OrderBy(column string, desc bool) *SQLBuilder {
	order := column
	if desc {
		order += " DESC"
	}
	sb.orderBy = append(sb.orderBy, order)
	return sb
}

// Limit 设置LIMIT
func (sb *SQLBuilder) Limit(limit int) *SQLBuilder {
	sb.limit = &limit
	return sb
}

// Offset 设置OFFSET
func (sb *SQLBuilder) Offset(offset int) *SQLBuilder {
	sb.offset = &offset
	return sb
}

// Build 构建完整SQL和参数
func (sb *SQLBuilder) Build() (string, []interface{}) {
	var parts []string

	// SELECT
	if len(sb.selectClause) > 0 {
		parts = append(parts, "SELECT "+strings.Join(sb.selectClause, ", "))
	} else {
		parts = append(parts, "SELECT *")
	}

	// FROM
	if sb.fromClause != "" {
		parts = append(parts, "FROM "+sb.fromClause)
	}

	// JOIN
	parts = append(parts, sb.joinClauses...)

	// WHERE
	if len(sb.whereClause) > 0 {
		wherePart := "WHERE " + strings.Join(sb.whereClause, " AND ")
		parts = append(parts, wherePart)
	}

	// GROUP BY
	if len(sb.groupBy) > 0 {
		parts = append(parts, "GROUP BY "+strings.Join(sb.groupBy, ", "))
	}

	// HAVING
	if len(sb.havingClause) > 0 {
		havingPart := "HAVING " + strings.Join(sb.havingClause, " AND ")
		parts = append(parts, havingPart)
	}

	// ORDER BY
	if len(sb.orderBy) > 0 {
		parts = append(parts, "ORDER BY "+strings.Join(sb.orderBy, ", "))
	}

	// LIMIT/OFFSET
	if sb.limit != nil {
		parts = append(parts, "LIMIT ?")
		sb.params = append(sb.params, *sb.limit)
	}
	if sb.offset != nil {
		parts = append(parts, "OFFSET ?")
		sb.params = append(sb.params, *sb.offset)
	}

	sql := strings.Join(parts, " ")
	return sql, sb.params
}

// Execute 执行查询
func (sb *SQLBuilder) Execute(db *gorm.DB, dest interface{}) error {
	sql, params := sb.Build()
	return db.Raw(sql, params...).Scan(dest).Error
}

// 定义通用的类型安全检查函数
func safeGetString(m map[string]interface{}, key string) (string, bool) {
	if val, exists := m[key]; exists && val != nil {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}

func safeGetBool(m map[string]interface{}, key string) (bool, bool) {
	if val, exists := m[key]; exists && val != nil {
		if b, ok := val.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// 示例：复杂多表查询
func ExampleComplexQuery(db *gorm.DB, filters map[string]interface{}) error {
	//mapstr := func(field string, mm map[string]interface{}) (string, bool) {
	//	if v, ok := filters[field]; ok {
	//		if value, ok2 := v.(string); ok2 {
	//			return value, true
	//		}
	//	}
	//	return "", false
	//}

	username, usernameExist := safeGetString(filters, "username")

	builder := NewSQLBuilder()

	// 构建查询
	builder.
		Select(
			"u.id",
			"u.username",
			"u.email",
			"r.name as role_name",
			"COUNT(p.id) as permission_count",
		).
		From("users u").
		Join("LEFT JOIN", "user_roles ur", "u.id = ur.user_id").
		Join("LEFT JOIN", "roles r", "ur.role_id = r.id").
		Join("LEFT JOIN", "role_permissions rp", "r.id = rp.role_id").
		Join("LEFT JOIN", "permissions p", "rp.permission_id = p.id").
		Where("u."+entity.FieldDeleteBy+" = ?", 0).
		WhereIf(usernameExist, "u.username LIKE ?", "%"+username+"%").
		WhereIf(filters["role"] != nil, "r.code = ?", filters["role"]).
		WhereIf(filters["active"] != nil && filters["active"].(bool), "u.status = ?", "active").
		GroupBy("u.id", "u.username", "u.email", "r.name").
		Having("COUNT(p.id) > ?", 0).
		OrderBy("u."+entity.FieldCreateTime, true).
		Limit(50)

	// 定义结果结构体
	type UserWithRole struct {
		ID              int64  `json:"id"`
		Username        string `json:"username"`
		Email           string `json:"email"`
		RoleName        string `json:"role_name"`
		PermissionCount int    `json:"permission_count"`
	}

	var results []UserWithRole
	return builder.Execute(db, &results)
}
