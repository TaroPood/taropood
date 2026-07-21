package repository

import (
	"context"

	"github.com/TaroPood/taropood/internal/domain"
)

type RuleRepository interface {
	Create(ctx context.Context, rule *domain.Rule) error
	GetByID(ctx context.Context, id string) (*domain.Rule, error)
	List(ctx context.Context, filter domain.RuleFilter) ([]*domain.Rule, error)
	Update(ctx context.Context, rule *domain.Rule) error
	Delete(ctx context.Context, id string) error
}

type Transactioner interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
