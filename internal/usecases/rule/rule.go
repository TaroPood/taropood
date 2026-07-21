package rule

import (
	"context"
	"errors"

	"github.com/TaroPood/taropood/internal/domain"
	"github.com/TaroPood/taropood/internal/dto"
	"github.com/TaroPood/taropood/internal/repository"
)

type UseCase struct {
	repo repository.RuleRepository
}

func NewUseCase(repo repository.RuleRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, req *dto.CreateRuleRequest) (*dto.RuleResponse, error) {
	rule := dto.CreateRequestToDomain(req)

	if err := uc.repo.Create(ctx, rule); err != nil {
		return nil, err
	}

	return dto.RuleToResponse(rule), nil
}

func (uc *UseCase) GetByID(ctx context.Context, id string) (*dto.RuleResponse, error) {
	if id == "" {
		return nil, domain.ErrNotFound
	}

	rule, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.RuleToResponse(rule), nil
}

func (uc *UseCase) CreateDomain(ctx context.Context, rule *domain.Rule) error {
	return uc.repo.Create(ctx, rule)
}

func (uc *UseCase) GetByIDDomain(ctx context.Context, id string) (*domain.Rule, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, filter domain.RuleFilter) ([]*domain.Rule, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *UseCase) Update(ctx context.Context, rule *domain.Rule) error {
	return uc.repo.Update(ctx, rule)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func IsNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func IsDuplicate(err error) bool {
	return errors.Is(err, domain.ErrDuplicateName)
}
