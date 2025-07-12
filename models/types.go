package models

type FlaggedItem struct {
	Label  string `json:"label"`
	Reason string `json:"reason"`
}

type AgentDecision struct {
	Safe        bool     `json:"safe"`
	Actions     []string `json:"actions"`
	Objections  []string `json:"objections"`
	Explanation string   `json:"explanation"`
}

type AgenticResponse struct {
	Classification map[string]interface{} `json:"classification"`
	FlaggedItems   []FlaggedItem          `json:"flagged_items"`
	AgentDecision  AgentDecision          `json:"agent_decision"`
}

type FinalResponse struct {
	ApiStatus  string          `json:"api_status"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Response   AgenticResponse `json:"response"`
}
