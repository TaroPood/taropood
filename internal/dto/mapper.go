package dto

import (
	"time"

	"github.com/TaroPood/taropood/internal/domain"
)

func RuleToResponse(r *domain.Rule) *RuleResponse {
	actions := make([]ActionDTO, len(r.Actions))
	for i, a := range r.Actions {
		actions[i] = ActionDTO{
			Type:   a.Type,
			Config: a.Config,
			Order:  a.Order,
		}
	}

	return &RuleResponse{
		ID:              r.ID,
		Name:            r.Name,
		ConditionType:   r.ConditionType,
		ConditionConfig: r.ConditionConfig,
		Priority:        r.Priority,
		Enabled:         r.Enabled,
		Tags:            r.Tags,
		Metadata:        r.Metadata,
		Actions:         actions,
		CreatedAt:       r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC3339),
	}
}

func CreateRequestToDomain(req *CreateRuleRequest) *domain.Rule {
	actions := make([]domain.Action, len(req.Actions))
	for i, a := range req.Actions {
		actions[i] = domain.Action{
			Type:   a.Type,
			Config: a.Config,
			Order:  a.Order,
		}
	}

	return &domain.Rule{
		Name:            req.Name,
		ConditionType:   req.ConditionType,
		ConditionConfig: req.ConditionConfig,
		Priority:        req.Priority,
		Enabled:         req.Enabled,
		Tags:            req.Tags,
		Metadata:        req.Metadata,
		Actions:         actions,
	}
}
