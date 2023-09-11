package rapid7

import "strings"

// InvestigationIDFromRRN parses an RRN and returns the investigation ID.
// Given this RRN:
//
//	rrn:investigation:us3:e4fb5946-0fe5-407c-bf83-64bd1502769d:investigation:HZVPPFG2OB24
//
// The following Investigation ID would be returned:
//
//	e4fb5946-0fe5-407c-bf83-64bd1502769d
func InvestigationIDFromRRN(rrn string) string {
	if rrn == "" {
		return ""
	}
	if !strings.Contains(rrn, "investigation") {
		return ""
	}
	parts := strings.Split(rrn, ":")
	if len(parts) != 6 {
		return ""
	}
	return parts[3]
}
