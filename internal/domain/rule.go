package domain

import "time"

type Rule struct {
	ID              string
	Name            string
	ConditionType   string
	ConditionConfig map[string]any
	Priority        int
	Enabled         bool
	Tags            []string
	Metadata        map[string]any
	Actions         []Action
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Action struct {
	ID        string
	RuleID    string
	Type      string
	Config    map[string]any
	Order     int
	CreatedAt time.Time
}

type RuleFilter struct {
	Name    *string
	Enabled *bool
	Tags    []string
	Limit   int
	Offset  int
}
