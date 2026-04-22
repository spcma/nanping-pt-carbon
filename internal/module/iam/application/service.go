package application

import (
	"app/internal/module/iam/domain"
	"context"
)

// ===== Service Ports（给外部模块用） =====

type IdentityService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}

type AuthorizationService interface {
	HasPermission(ctx context.Context, userID int64, permissionCode string) (bool, error)
}
