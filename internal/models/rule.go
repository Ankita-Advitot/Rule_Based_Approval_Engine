package models

type Rule struct {
	ID          int64                  `json:"id"`
	RequestType string                 `json:"request_type"`
	Condition   map[string]interface{} `json:"condition"`
	Action      string                 `json:"action"`
	GradeID     int64                  `json:"grade_id"`
	Active      bool                   `json:"active"`
}
