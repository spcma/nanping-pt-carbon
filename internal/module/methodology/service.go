package methodology

import (
	"context"
)

// ===== Service Ports（服务端口 - 给外部模块用） =====

type MethodologyService interface {
	GetMethodology(ctx context.Context, id int64) (*Methodology, error)
	GetByQuery(ctx context.Context, code string) (*Methodology, error)
}
