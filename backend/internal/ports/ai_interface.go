package ports

import (
	"context"
	"seismic-monitor/backend/internal/models"
)

type AIProvider interface {
	GenerateSafetyAdvice(ctx context.Context, sismo models.Feature) (string, error)
}
