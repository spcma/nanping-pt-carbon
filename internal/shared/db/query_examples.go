package db

import (
	"app/internal/shared/entity"

	"gorm.io/gorm"
)

// MultiTableQueryExample 多表关联查询示例
type MultiTableQueryExample struct {
	db *gorm.DB
}

func NewMultiTableQueryExample(database *gorm.DB) *MultiTableQueryExample {
	return &MultiTableQueryExample{db: database}
}

// GetUserWithRoles 获取用户及其角色信息
func (m *MultiTableQueryExample) GetUserWithRoles(userID int64, filters map[string]interface{}) ([]map[string]interface{}, error) {
	// 使用QueryBuilder方式
	builder := NewQueryBuilder(m.db)

	query := builder.
		Table("users u").
		LeftJoin("user_roles ur", "u.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		Where("u."+entity.FieldDeleteBy+" = ?", 0).
		Where("r."+entity.FieldDeleteBy+" = ?", 0)

	// 动态条件
	if userID > 0 {
		query = query.Where("u.id = ?", userID)
	}

	if username, ok := filters["username"]; ok && username != "" {
		query = query.Where("u.username LIKE ?", "%"+username.(string)+"%")
	}

	if roleName, ok := filters["role_name"]; ok && roleName != "" {
		query = query.Where("r.name = ?", roleName)
	}

	query = query.OrderBy("u."+entity.FieldCreateTime, true)

	// 执行查询
	var results []map[string]interface{}
	err := query.Build().Find(&results).Error
	return results, err
}

// GetUserStatistics 使用原生SQL方式获取统计信息
func (m *MultiTableQueryExample) GetUserStatistics(statusFilter string) ([]map[string]interface{}, error) {
	// 构建复杂SQL查询
	sql := `
		SELECT 
			u.status,
			COUNT(*) as user_count,
			COALESCE(AVG(role_stats.role_count), 0) as avg_roles_per_user,
			COUNT(CASE WHEN u.create_time >= NOW() - INTERVAL '30 days' THEN 1 END) as recent_users
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as role_count
			FROM user_roles 
			WHERE delete_by = ?
			GROUP BY user_id
		) role_stats ON u.id = role_stats.user_id
		WHERE u.delete_by = ?
	`

	params := []interface{}{0, 0}

	// 动态添加条件
	if statusFilter != "" {
		sql += " AND u.status = ?"
		params = append(params, statusFilter)
	}

	sql += " GROUP BY u.status ORDER BY user_count DESC"

	// 执行查询
	var results []map[string]interface{}
	err := m.db.Raw(sql, params...).Scan(&results).Error
	return results, err
}

// GetUserDetailReport 使用视图方式获取详细报告
func (m *MultiTableQueryExample) GetUserDetailReport(filters map[string]interface{}) ([]map[string]interface{}, error) {
	query := m.db.Table("user_detail_view") // 假设已创建视图

	// 动态过滤条件
	if username, ok := filters["username"]; ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}

	if email, ok := filters["email"]; ok && email != "" {
		query = query.Where("email LIKE ?", "%"+email.(string)+"%")
	}

	if minPermissions, ok := filters["min_permissions"]; ok {
		query = query.Where("permission_count >= ?", minPermissions)
	}

	query = query.Order("create_time DESC").Limit(100)

	var results []map[string]interface{}
	err := query.Find(&results).Error
	return results, err
}

// BestPracticeRecommendation 最佳实践建议
/*
根据不同场景选择合适的方案：

1. 简单多表查询 → QueryBuilder模式
   - 优势：类型安全、易于维护
   - 适用：日常业务查询

2. 复杂统计查询 → 原生SQL + 参数化
   - 优势：性能最优、灵活性最高
   - 适用：报表统计、复杂聚合

3. 高频复杂查询 → 视图/存储过程
   - 优势：数据库层面优化、缓存友好
   - 适用：核心业务报表、大数据量场景

4. 混合使用策略：
   - CRUD操作：继续使用GORM链式调用
   - 复杂查询：使用上述工具类
   - 性能敏感：考虑数据库层面优化
*/
