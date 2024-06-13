// See https://help.rapid7.com/insightvm/en-us/api/integrations.html
package rapid7

import (
	"encoding/json"
	"fmt"
	"time"
)

type VMVulnerabilityStatus string

func (v VMVulnerabilityStatus) String() string {
	return string(v)
}

type VMType string

func (v VMType) String() string {
	return string(v)
}

type VMAssetSearchPageSize int

func (s VMAssetSearchPageSize) String() string {
	return fmt.Sprintf("%d", s)
}

var VM_ASSET_SEARCH_PAGE_SIZE VMAssetSearchPageSize = 100

type VMMetadata struct {
	// The index (zero-based) of the current page returned.
	Number int64 `json:"number"`
	// The maximum size of the page returned.
	Size int64 `json:"size"`
	// The stateless cursor associated with the series of page requests being made.
	Cursor string `json:"cursor"`
	// The total number of resources available across all pages.
	TotalResources int64 `json:"totalResources"`
	// The total number of pages available.
	TotalPages int64 `json:"totalPages"`
}

type VMLink struct {
	HREF string `json:"href"`
	Rel  string `json:"rel"`
}

type Rapid7VMPagedResponse[T any] struct {
	Data     []T         `json:"data"`
	Links    []VMLink    `json:"links"`
	Metadata *VMMetadata `json:"metadata"`
}

const (
	VMExceptionVulnExpl   VMVulnerabilityStatus = "EXCEPTION_VULN_EXPL"
	VMUnexpectedErr       VMVulnerabilityStatus = "UNEXPECTED_ERR"
	VMNotVulnDontStore    VMVulnerabilityStatus = "NOT_VULN_DONT_STORE"
	VMSuperseded          VMVulnerabilityStatus = "SUPERSEDED"
	VMExceptionVulnPotl   VMVulnerabilityStatus = "EXCEPTION_VULN_POTL"
	VMVulnerableExpl      VMVulnerabilityStatus = "VULNERABLE_EXPL"
	VMOverriddenVulnVers  VMVulnerabilityStatus = "OVERRIDDEN_VULN_VERS"
	VMSkippedDisabled     VMVulnerabilityStatus = "SKIPPED_DISABLED"
	VMVulnerableVers      VMVulnerabilityStatus = "VULNERABLE_VERS"
	VMVulnerablePotential VMVulnerabilityStatus = "VULNERABLE_POTENTIAL"
	VMSkippedVers         VMVulnerabilityStatus = "SKIPPED_VERS"
	VMExceptionVulnVers   VMVulnerabilityStatus = "EXCEPTION_VULN_VERS"
	VMNotVulnerable       VMVulnerabilityStatus = "NOT_VULNERABLE"
	VMUnknownStatus       VMVulnerabilityStatus = "UNKNOWN"
	VMSkippedDOS          VMVulnerabilityStatus = "SKIPPED_DOS"
)

const (
	VMHypervisor VMType = "hypervisor"
	VMMobile     VMType = "mobile"
	VMGuest      VMType = "guest"
	VMPhysical   VMType = "physical"
	VMUnknown    VMType = "unknown"
)

type VMAssetTag struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type VMUniqueIdentifier struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}

type VMVulnerability struct {
	CheckID         string                `json:"check_id"`
	FirstFound      time.Time             `json:"first_found"`
	Key             string                `json:"key"`
	LastFound       time.Time             `json:"last_found"`
	NIC             string                `json:"nic"`
	Port            int32                 `json:"port"`
	Proof           string                `json:"proof"`
	Protocol        string                `json:"protocol"`
	SolutionFix     string                `json:"solution_fix"`
	SolutionID      string                `json:"solution_id"`
	SolutionSummary string                `json:"solution_summary"`
	SolutionType    string                `json:"solution_type"`
	Status          VMVulnerabilityStatus `json:"status"`
	VulnerabilityID string                `json:"vulnerability_id"`
}

type VMCredentialAssessment struct {
	Port     int64  `json:"port"`
	Protocol string `json:"protocol"`
	Status   string `json:"status"`
}

type VMAsset struct {
	AssessedForPolicies            bool                     `json:"assessed_for_policies"`
	AssessedForVulnerabilities     bool                     `json:"assessed_for_vulnerabilities"`
	CredentialAssessments          []VMCredentialAssessment `json:"credential_assessments"`
	CriticalVulnerabilities        int32                    `json:"critical_vulnerabilities"`
	Exploits                       int32                    `json:"exploits"`
	HostName                       string                   `json:"host_name"`
	ID                             string                   `json:"id"`
	IP                             string                   `json:"ip"`
	LastAssessedForVulnerabilities time.Time                `json:"last_assessed_for_vulnerabilities"`
	LastScanEnd                    time.Time                `json:"last_scan_end"`
	LastScanStart                  time.Time                `json:"last_scan_start"`
	MAC                            string                   `json:"mac"`
	MalwareKits                    int32                    `json:"malware_kits"`
	ModerateVulnerabilities        int32                    `json:"moderate_vulnerabilities"`
	New                            []VMVulnerability        `json:"new"`
	OSArchitecture                 string                   `json:"os_architecture"`
	OSDescription                  string                   `json:"os_description"`
	OSFamily                       string                   `json:"os_family"`
	OSName                         string                   `json:"os_name"`
	OSSystemName                   string                   `json:"os_system_name"`
	OSType                         string                   `json:"os_type"`
	OSVendor                       string                   `json:"os_vendor"`
	OSVersion                      string                   `json:"os_version"`
	Remediated                     []VMVulnerability        `json:"remediated"`
	RiskScore                      float32                  `json:"risk_score"`
	Same                           []VMVulnerability        `json:"same"`
	SevereVulnerabilities          int32                    `json:"severe_vulnerabilities"`
	Tags                           []VMAssetTag             `json:"tags"`
	TotalVulnerabilities           int32                    `json:"total_vulnerabilities"`
	Type                           VMType                   `json:"type"`
	UniqueIdentifiers              []VMUniqueIdentifier     `json:"unique_identifiers"`
}

type VMAssetSearchRequest struct {
	Asset         string `json:"asset,omitempty"`
	Vulnerability string `json:"vulnerability,omitempty"`
}

type VMAssetSearchQuery struct {
	Cursor                   string        `json:"cursor,omitempty"`
	CurrentTime              time.Time     `json:"currentTime,omitempty"`
	ComparisonTime           time.Time     `json:"comparisonTime,omitempty"`
	IncludeSame              bool          `json:"includeSame,omitempty"`
	IncludeUniqueIdentifiers bool          `json:"includeUniqueIdentifiers,omitempty"`
	Page                     int           `json:"page,omitempty"`
	Size                     int           `json:"size,omitempty"`
	Sort                     SortDirection `json:"sort,omitempty"`
}

func (q VMAssetSearchQuery) Map() map[string]string {
	b, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	var jsonMap map[string]any
	err = json.Unmarshal(b, &jsonMap)
	if err != nil {
		panic(err)
	}
	result := make(map[string]string)
	for k, v := range jsonMap {
		if v == "0001-01-01T00:00:00Z" {
			continue
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
