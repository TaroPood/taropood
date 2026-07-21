package dto

type CreateRuleRequest struct {
	Name            string         `json:"name"`
	ConditionType   string         `json:"condition_type"`
	ConditionConfig map[string]any `json:"condition_config"`
	Priority        int            `json:"priority"`
	Enabled         bool           `json:"enabled"`
	Tags            []string       `json:"tags"`
	Metadata        map[string]any `json:"metadata"`
	Actions         []ActionDTO    `json:"actions"`
}

type ActionDTO struct {
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
	Order  int            `json:"order"`
}

type RuleResponse struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	ConditionType   string         `json:"condition_type"`
	ConditionConfig map[string]any `json:"condition_config"`
	Priority        int            `json:"priority"`
	Enabled         bool           `json:"enabled"`
	Tags            []string       `json:"tags"`
	Metadata        map[string]any `json:"metadata"`
	Actions         []ActionDTO    `json:"actions"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
