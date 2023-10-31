package rapid7

import (
	"time"
)

// CRITICAL, HIGH, MEDIUM, LOW, UNSPECIFIED
type InvestigationPriority string

func (i InvestigationPriority) String() string {
	return string(i)
}

// OPEN, INVESTIGATING, CLOSED
type InvestigationStatus string

func (i InvestigationStatus) String() string {
	return string(i)
}

// BENIGN, MALICIOUS, NOT_APPLICABLE, UNDECIDED
type InvestigationDisposition string

func (i InvestigationDisposition) String() string {
	return string(i)
}

// MANUAL, HUNT, ALERT
type InvestigationSource string

func (i InvestigationSource) String() string {
	return string(i)
}

const UNSPECIFIED InvestigationPriority = "UNSPECIFIED"
const LOW InvestigationPriority = "LOW"
const MEDIUM InvestigationPriority = "MEDIUM"
const HIGH InvestigationPriority = "HIGH"
const CRITICAL InvestigationPriority = "CRITICAL"

const OPEN InvestigationStatus = "OPEN"
const INVESTIGATING InvestigationStatus = "INVESTIGATING"
const CLOSED InvestigationStatus = "CLOSED"
const WAITING InvestigationStatus = "WAITING"

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
	Assignee        *Assignee                `json:"assignee"`
	CreatedTime     time.Time                `json:"created_time"`
	Disposition     InvestigationDisposition `json:"disposition"`
	FirstAlertTime  *time.Time               `json:"first_alert_time"`
	LastAccessed    time.Time                `json:"last_accessed"`
	LatestAlertTime *time.Time               `json:"latest_alert_time"`
	OrganizationID  string                   `json:"organization_id"`
	Priority        InvestigationPriority    `json:"priority"`
	Responsibility  string                   `json:"responsibility"`
	RRN             string                   `json:"rrn"`
	Source          InvestigationSource      `json:"source"`
	Status          InvestigationStatus      `json:"status"`
	Tags            []string                 `json:"tags"`
	Title           string                   `json:"title"`
}

type APIError struct {
	Message       string `json:"message"`
	CorrelationID string `json:"correlation_id"`
}

type Metadata struct {
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
}

type Rapid7PagedResponse[T any] struct {
	Data     []*T     `json:"data"`
	Metadata Metadata `json:"metadata"`
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

type Creator struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CommentAttachment struct {
	CreatedTime time.Time `json:"created_time"`
	Creator     Creator   `json:"creator"`
	FileName    string    `json:"file_name"`
	MimeType    string    `json:"mime_type"`
	RRN         string    `json:"rrn"`
	ScanStatus  string    `json:"scan_status"`
	Size        int64     `json:"size"`
}

type RRN struct {
	OrganizationID string   `json:"organizationId"`
	Partition      string   `json:"partition"`
	RegionCode     string   `json:"regionCode"`
	Resource       string   `json:"resource"`
	ResourceTypes  []string `json:"resourceTypes"`
	Service        string   `json:"service"`
}

type InvestigationCommentData struct {
	Body        string              `json:"body"`
	CreatedTime time.Time           `json:"created_time"`
	Creator     Creator             `json:"creator"`
	RRN         string              `json:"rrn"`
	Target      string              `json:"target"`
	Visibility  string              `json:"visibility"`
	Attachments []CommentAttachment `json:"attachments"`
}

type InvestigationComments struct {
	Data     []InvestigationCommentData `json:"data"`
	Metadata Metadata                   `json:"metadata"`
}

type InvestigationAssignee struct {
	Email string `json:"email"`
}

type InvestigationUpdateRequest struct {
	Assignee    *InvestigationAssignee   `json:"assignee,omitempty"`
	Disposition InvestigationDisposition `json:"disposition,omitempty"`
	Priority    InvestigationPriority    `json:"priority,omitempty"`
	Status      InvestigationStatus      `json:"status,omitempty"`
	Title       string                   `json:"title,omitempty"`
}
