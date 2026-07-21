package model

import (
	"encoding/json"
	"time"

	"github.com/TaroPood/taropood/internal/domain"
)

var Models = []any{
	&RuleModel{},
	&ActionModel{},
}

type RuleModel struct {
	ID              string          `gorm:"primaryKey;type:text;not null"`
	Name            string          `gorm:"type:text;not null;index:idx_rules_name"`
	ConditionType   string          `gorm:"type:text;not null"`
	ConditionConfig json.RawMessage `gorm:"type:jsonb"`
	Priority        int             `gorm:"default:0;not null"`
	Enabled         bool            `gorm:"default:true;not null;index:idx_rules_enabled"`
	Tags            []string        `gorm:"serializer:json;type:jsonb"`
	Metadata        json.RawMessage `gorm:"type:jsonb"`
	CreatedAt       time.Time       `gorm:"not null"`
	UpdatedAt       time.Time       `gorm:"not null"`
	Actions         []ActionModel   `gorm:"foreignKey:RuleID;constraint:OnDelete:CASCADE"`
}

func (m *RuleModel) ToDomain() *domain.Rule {
	r := &domain.Rule{
		ID:            m.ID,
		Name:          m.Name,
		ConditionType: m.ConditionType,
		Priority:      m.Priority,
		Enabled:       m.Enabled,
		Tags:          m.Tags,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	if len(m.ConditionConfig) > 0 {
		_ = json.Unmarshal(m.ConditionConfig, &r.ConditionConfig)
	}
	if len(m.Metadata) > 0 {
		_ = json.Unmarshal(m.Metadata, &r.Metadata)
	}
	for _, a := range m.Actions {
		r.Actions = append(r.Actions, *a.ToDomain())
	}
	return r
}

func RuleToModel(r *domain.Rule) (*RuleModel, error) {
	m := &RuleModel{
		ID:            r.ID,
		Name:          r.Name,
		ConditionType: r.ConditionType,
		Priority:      r.Priority,
		Enabled:       r.Enabled,
		Tags:          r.Tags,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
	if r.ConditionConfig != nil {
		data, err := json.Marshal(r.ConditionConfig)
		if err != nil {
			return nil, err
		}
		m.ConditionConfig = data
	}
	if r.Metadata != nil {
		data, err := json.Marshal(r.Metadata)
		if err != nil {
			return nil, err
		}
		m.Metadata = data
	}
	for _, a := range r.Actions {
		am, err := ActionToModel(&a)
		if err != nil {
			return nil, err
		}
		m.Actions = append(m.Actions, *am)
	}
	return m, nil
}

type ActionModel struct {
	ID        string          `gorm:"primaryKey;type:text;not null"`
	RuleID    string          `gorm:"type:text;not null;index:idx_rule_actions_rule_id"`
	Type      string          `gorm:"type:text;not null"`
	Config    json.RawMessage `gorm:"type:jsonb"`
	Order     int             `gorm:"default:0;not null"`
	CreatedAt time.Time       `gorm:"not null"`
}

func (m *ActionModel) ToDomain() *domain.Action {
	a := &domain.Action{
		ID:        m.ID,
		RuleID:    m.RuleID,
		Type:      m.Type,
		Order:     m.Order,
		CreatedAt: m.CreatedAt,
	}
	if len(m.Config) > 0 {
		_ = json.Unmarshal(m.Config, &a.Config)
	}
	return a
}

func ActionToModel(a *domain.Action) (*ActionModel, error) {
	m := &ActionModel{
		ID:        a.ID,
		RuleID:    a.RuleID,
		Type:      a.Type,
		Order:     a.Order,
		CreatedAt: a.CreatedAt,
	}
	if a.Config != nil {
		data, err := json.Marshal(a.Config)
		if err != nil {
			return nil, err
		}
		m.Config = data
	}
	return m, nil
}
