package rapid7

import "time"

const IDR_VERSION string = "v2"

type Assignee struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Investigation struct {
	Assignee        Assignee  `json:"assignee"`
	CreatedTime     time.Time `json:"created_time"`
	Disposition     string    `json:"disposition"`
	FirstAlertTime  time.Time `json:"first_alert_time"`
	LastAccessed    time.Time `json:"last_accessed"`
	LatestAlertTime time.Time `json:"latest_alert_time"`
	OrganizationID  string    `json:"organization_id"`
	Priority        string    `json:"priority"`
	Responsibility  string    `json:"responsibility"`
	Rrn             string    `json:"rrn"`
	Source          string    `json:"source"`
	Status          string    `json:"status"`
	Tags            []string  `json:"tags"`
	Title           string    `json:"title"`
}

type APIError struct {
	Message       string `json:"message"`
	CorrelationID string `json:"correlation_id"`
}

type InvestigationsResponse struct {
	Data []*Investigation `json:"data"`
}

type InvestigationsQuery struct {
	AssigneeEmail string    `url:"assignee.email,omitempty"`
	EndTime       time.Time `url:"end_time,omitempty"`
	Index         int32     `url:"index,omitempty"`
	MultiCustomer bool      `url:"multi-customer,omitempty"`
	Priorities    []string  `url:"priorities,omitempty,comma"`
	Size          int32     `url:"size,omitempty"`
	Sort          string    `url:"sort,omitempty"`
	Sources       []string  `url:"sources,omitempty,comma"`
	StartTime     time.Time `url:"start_time,omitempty,comma"`
	Statuses      []string  `url:"statuses,omitempty,comma"`
	Tags          []string  `url:"tags,omitempty,comma"`
}
