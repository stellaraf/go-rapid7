package rapid7

import "time"

const IDR_VERSION string = "v2"

type InvestigationPriority string
type InvestigationStatus string
type InvestigationDisposition string
type InvestigationSource string

const UNSPECIFIED InvestigationPriority = "UNSPECIFIED"
const LOW InvestigationPriority = "LOW"
const MEDIUM InvestigationPriority = "MEDIUM"
const HIGH InvestigationPriority = "HIGH"
const CRITICAL InvestigationPriority = "CRITICAL"

const OPEN InvestigationStatus = "OPEN"
const INVESTIGATING InvestigationStatus = "INVESTIGATING"
const CLOSED InvestigationStatus = "CLOSED"

const BENIGN InvestigationDisposition = "BENIGN"
const MALICIOUS InvestigationDisposition = "MALICIOUS"
const NOT_APPLICABLE InvestigationDisposition = "NOT_APPLICABLE"
const UNDECIDED InvestigationDisposition = "UNDECIDED"

const MANUAL InvestigationSource = "MANUAL"
const HUNT InvestigationSource = "HUNT"
const ALERT InvestigationSource = "ALERT"

type Assignee struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Investigation struct {
	Assignee        Assignee                 `json:"assignee"`
	CreatedTime     time.Time                `json:"created_time"`
	Disposition     InvestigationDisposition `json:"disposition"`
	FirstAlertTime  time.Time                `json:"first_alert_time"`
	LastAccessed    time.Time                `json:"last_accessed"`
	LatestAlertTime time.Time                `json:"latest_alert_time"`
	OrganizationID  string                   `json:"organization_id"`
	Priority        InvestigationPriority    `json:"priority"`
	Responsibility  string                   `json:"responsibility"`
	Rrn             string                   `json:"rrn"`
	Source          InvestigationSource      `json:"source"`
	Status          InvestigationStatus      `json:"status"`
	Tags            []string                 `json:"tags"`
	Title           string                   `json:"title"`
}

type APIError struct {
	Message       string `json:"message"`
	CorrelationID string `json:"correlation_id"`
}

type Rapid7PagedResponse[T any] struct {
	Data     []*T `json:"data"`
	MetaData struct {
		// The current page, starting from 0. This value will always be provided.
		Index int `json:"index"`
		// The number of data items in the current page. This value will always be provided.
		Size int `json:"size"`
		// The total number of pages that make up the complete response. This will be provided if possible.
		TotalPages int `json:"total_pages,omitempty"`
		// The total number of data items that make up the complete response. This will be provided if possible.
		TotalData int `json:"total_data,omitempty"`
		// The attributes used to sort the complete response. This will be provided if the response is sorted.
		Sort string `json:"sort,omitempty"`
	} `json:"metadata"`
}

type InvestigationsResponse = Rapid7PagedResponse[Investigation]

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
