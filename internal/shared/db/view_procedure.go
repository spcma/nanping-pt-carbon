package db

import (
	"app/internal/shared/entity"
	"database/sql"

	"gorm.io/gorm"
)

// ViewBasedRepository 基于视图的仓储模式
type ViewBasedRepository struct {
	db *gorm.DB
}

func NewViewBasedRepository(db *gorm.DB) *ViewBasedRepository {
	return &ViewBasedRepository{db: db}
}

// CreateUserDetailView 创建用户详情视图
func (r *ViewBasedRepository) CreateUserDetailView() error {
	viewSQL := `
	CREATE OR REPLACE VIEW user_detail_view AS
	SELECT 
		u.id,
		u.username,
		u.email,
		u.status,
		u.create_time,
		u.update_time,
		COALESCE(STRING_AGG(r.name, ', '), '') as roles,
		COALESCE(STRING_AGG(r.code, ', '), '') as role_codes,
		COUNT(DISTINCT p.id) as permission_count
	FROM users u
	LEFT JOIN user_roles ur ON u.id = ur.user_id AND ur.` + entity.FieldDeleteBy + ` = 0
	LEFT JOIN roles r ON ur.role_id = r.id AND r.` + entity.FieldDeleteBy + ` = 0
	LEFT JOIN role_permissions rp ON r.id = rp.role_id AND rp.` + entity.FieldDeleteBy + ` = 0
	LEFT JOIN permissions p ON rp.permission_id = p.id AND p.` + entity.FieldDeleteBy + ` = 0
	WHERE u.` + entity.FieldDeleteBy + ` = 0
	GROUP BY u.id, u.username, u.email, u.status, u.create_time, u.update_time
	`

	return r.db.Exec(viewSQL).Error
}

// QueryUserDetails 查询用户详情（通过视图）
func (r *ViewBasedRepository) QueryUserDetails(filters map[string]interface{}) ([]map[string]interface{}, error) {
	query := r.db.Table("user_detail_view")

	// 动态添加过滤条件
	if username, ok := filters["username"]; ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if minPermissions, ok := filters["min_permissions"]; ok {
		query = query.Where("permission_count >= ?", minPermissions)
	}

	// 添加排序
	query = query.Order("create_time DESC")

	// 执行查询
	var results []map[string]interface{}
	err := query.Find(&results).Error
	return results, err
}

// StoredProcedureRepository 存储过程仓储
type StoredProcedureRepository struct {
	db *gorm.DB
}

func NewStoredProcedureRepository(db *gorm.DB) *StoredProcedureRepository {
	return &StoredProcedureRepository{db: db}
}

// CreateUserStatsProcedure 创建用户统计存储过程
func (r *StoredProcedureRepository) CreateUserStatsProcedure() error {
	procedureSQL := `
	CREATE OR REPLACE FUNCTION get_user_statistics(
		p_status VARCHAR DEFAULT NULL,
		p_min_role_count INTEGER DEFAULT 0,
		p_days_back INTEGER DEFAULT 30
	) RETURNS TABLE(
		status VARCHAR,
		user_count BIGINT,
		avg_role_count NUMERIC,
		created_last_days BIGINT
	) AS $$
	BEGIN
		RETURN QUERY
		SELECT 
			u.status,
			COUNT(*) as user_count,
			AVG(role_counts.role_count)::NUMERIC as avg_role_count,
			SUM(CASE WHEN u.create_time >= NOW() - INTERVAL '1 day' * p_days_back THEN 1 ELSE 0 END) as created_last_days
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as role_count
			FROM user_roles 
			WHERE ` + entity.FieldDeleteBy + ` = 0
			GROUP BY user_id
		) role_counts ON u.id = role_counts.user_id
		WHERE u.` + entity.FieldDeleteBy + ` = 0
		AND (p_status IS NULL OR u.status = p_status)
		AND (role_counts.role_count IS NULL OR role_counts.role_count >= p_min_role_count)
		GROUP BY u.status
		ORDER BY user_count DESC;
	END;
	$$ LANGUAGE plpgsql;
	`

	return r.db.Exec(procedureSQL).Error
}

// GetUserStatistics 调用存储过程获取统计信息
func (r *StoredProcedureRepository) GetUserStatistics(status *string, minRoleCount, daysBack int) ([]map[string]interface{}, error) {
	var statusParam interface{}
	if status != nil {
		statusParam = *status
	} else {
		statusParam = sql.NullString{}
	}

	rows, err := r.db.Raw(
		"SELECT * FROM get_user_statistics(?, ?, ?)",
		statusParam, minRoleCount, daysBack,
	).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 处理结果集
	columns, _ := rows.Columns()
	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		results = append(results, rowMap)
	}

	return results, nil
}
