package rapid7

import (
	"fmt"
	"time"
)

// ASC, DESC
type SortDirection string

func (s SortDirection) String() string {
	return string(s)
}

// `created_time`, `priority`, `rrn`, `alerts_most_recent_created_time`, or
// `alerts_most_recent_detection_created_time`.
type SortField string

func (s SortField) String() string {
	return string(s)
}

// CRITICAL, HIGH, MEDIUM, LOW, UNSPECIFIED
type InvestigationPriority string

func (i InvestigationPriority) String() string {
	return string(i)
}

// OPEN, WAITING, INVESTIGATING, CLOSED
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

const SORT_ASCENDING SortDirection = "ASC"
const SORT_DESCENDING SortDirection = "DESC"

const SORT_CREATED_TIME SortField = "created_time"
const SORT_PRIORITY SortField = "priority"
const SORT_RRN SortField = "rrn"
const SORT_MOST_RECENT_CREATED_TIME SortField = "alerts_most_recent_created_time"
const SORT_MOST_RECENT_DETECTION_TIME SortField = "alerts_most_recent_detection_created_time"

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
	Index int32 `json:"index"`
	// The number of data items in the current page. This value will always be provided.
	Size int32 `json:"size"`
	// The attributes used to sort the complete response. This will be provided if the response is sorted.
	Sort string `json:"sort,omitempty"`
	// The total number of data items that make up the complete response. This will be provided if possible.
	TotalData int64 `json:"total_data,omitempty"`
	// The total number of pages that make up the complete response. This will be provided if possible.
	TotalPages int32 `json:"total_pages,omitempty"`
}

type Rapid7PagedResponse[T any] struct {
	Data     []*T     `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type InvestigationsResponse = Rapid7PagedResponse[Investigation]

type InvestigationsQuery struct {
	// A user's email address. Only investigations assigned to that user will be included.
	AssigneeEmail string `url:"assignee.email,omitempty"`
	// The time an investigation is closed. Only investigations whose created_time is before this
	// date will be returned by the API. Must be an ISO-formatted timestamp.
	EndTime time.Time `url:"end_time,omitempty"`
	// The 0-based index of the first page to retrieve. Must be an integer greater than 0.
	//
	// Default: 0
	Index int32 `url:"index,omitempty"`
	// Indicates whether the requester has multi-customer access. If set to true, a user API key
	// must be provided. Investigations will be returned from all organizations the calling user
	// has access to.
	//
	// Default: false
	MultiCustomer bool `url:"multi-customer,omitempty"`
	// A comma-separated list of investigation priorities to include in the result.
	Priorities []InvestigationPriority `url:"priorities,omitempty,comma"`
	// The maximum number of investigations to retrieve. Must be an integer greater than 0, or less
	// than or equal to 100.
	//
	// Default: 20
	Size int32 `url:"size,omitempty"`
	// Sort investigations by field and direction,  separated by a comma. Sortable fields are
	// `created_time`, `priority`, `rrn`, `alerts_most_recent_created_time`, and
	// `alerts_most_recent_detection_created_time`.
	//
	// Default: "priority,DESC"
	Sort string `url:"sort,omitempty"`
	// A comma-separated list of investigation sources to include in the result.
	Sources []string `url:"sources,omitempty,comma"`
	// The time an investigation is opened. Only investigations whose created_time is after this
	// date will be returned by the API. Must be an ISO-formatted timestamp.
	//
	// Default: 28 days prior to current time.
	StartTime time.Time `url:"start_time,omitempty,comma"`
	// A comma-separated list of investigation statuses to include in the result.
	Statuses []InvestigationStatus `url:"statuses,omitempty,comma"`
	// A comma-separated list of tags to include in the result. Only investigations who have all
	// specified tags will be included.
	Tags []string `url:"tags,omitempty,comma"`
}

func (q *InvestigationsQuery) SortBy(field SortField, direction SortDirection) {
	q.Sort = fmt.Sprintf("%s,%s", field, direction)
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
