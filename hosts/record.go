package hosts

import "strings"

// Record represents a single Hosts record.
type Record struct {
	Raw      string   `json:"raw"`
	LocalPtr string   `json:"localPtr"`
	Hosts    []string `json:"hosts"`
	Error    error    `json:"error"`
}

// IsComment determines if the current Hosts record is a comment line.
func (r Record) IsComment() bool {
	trimLine := strings.TrimSpace(r.Raw)
	isComment := strings.HasPrefix(trimLine, commentChar)
	return isComment
}
