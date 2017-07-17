// hosts is used to manipulate the local DNS environment. No more facebooking kids muahaha.
//
// ::1			localhost
// www.nsa.gov 	www.facebook.com

package hosts

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const commentChar string = "#"

// Hosts is a representation of the local /etc/hosts file
type Hosts struct {
	Path    string   `json:"path"`
	Records []Record `json:"records"`
	sync.Mutex
}

// Load reads the hosts file and parses out each record.
func (h *Hosts) Load() error {
	var records []Record

	file, err := os.Open(h.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record, err := NewHostsRecord(scanner.Text())
		if err != nil || record == nil {
			continue
		}

		records = append(records, *record)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	h.Lock()
	defer h.Unlock()
	h.Records = records

	return nil
}

func (h *Hosts) IsWritable() bool {
	_, err := os.OpenFile(h.Path, os.O_WRONLY, 0660)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func NewHostsRecord(raw string) (*Record, error) {
	r := Record{
		Raw: raw,
	}
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return nil, errors.New("record was empty")
	}

	if r.IsComment() {
		return nil, errors.New("record was a comment")
	}

	rawIP := fields[0]
	if net.ParseIP(rawIP) == nil {
		r.Error = errors.Errorf("Bad hosts line: %q", raw)
	}

	r.LocalPtr = rawIP
	r.Hosts = fields[1:]

	return &r, nil
}

func (h Hosts) Flush() error {
	file, err := os.Create(h.Path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)

	for _, rec := range h.Records {
		fmt.Fprintf(w, "%s%s", rec.Raw, eol)
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return h.Load()
}

// Set an entry into the hosts file.
func (h *Hosts) Set(ip string, hosts ...string) error {
	if net.ParseIP(ip) == nil {
		return errors.Errorf("%q is an invalid IP address.", ip)
	}

	// local the hosts object
	h.Lock()
	defer h.Unlock()

	// get the position of the IP in our current Hosts records.
	// if we find a position, we will do an update to the existing
	// record, otherwise we'll append.
	position := h.getIpPosition(ip)
	switch position {
	case -1:

		// parse the line
		rec, err := NewHostsRecord(buildRawLine(ip, hosts))
		if err != nil {
			return err
		}

		// new record, true insert
		h.Records = append(h.Records, *rec)
	default:

		newHosts := h.Records[position].Hosts
		for _, addHost := range hosts {
			if in(addHost, newHosts) {
				continue
			}

			newHosts = append(newHosts, addHost)
		}

		rec, err := NewHostsRecord(buildRawLine(ip, newHosts))
		if err != nil {
			return err
		}

		// update
		h.Records[position] = *rec
	}

	return nil
}

func (h Hosts) Has(ip string, host string) bool {
	return h.getHostPosition(ip, host) != -1
}

// Remove an entry from the hosts file.
func (h *Hosts) Remove(ip string, hosts ...string) error {
	var output []Record

	if net.ParseIP(ip) == nil {
		return errors.Errorf("%q is an invalid IP address.", ip)
	}

	if len(h.Records) == 0 {
		return errors.New("Hosts file contained no records")
	}

	// create a map of the hosts for indexed lookup
	var hostMap map[string]string
	for _, h := range hosts {
		hostMap[h] = h
	}

	for _, r := range h.Records {

		if r.Error != nil || r.IsComment() || r.LocalPtr != ip {
			output = append(output, r)
			continue
		}

		// now that we know we have the right record for the IP address,
		// we need to update the host references for this IP.

		var newHosts []string
		for _, h := range r.Hosts {
			if !in(h, hosts) {
				newHosts = append(newHosts, h)
			}
		}

		if len(newHosts) > 0 {
			raw := r.LocalPtr

			for i, h := range newHosts {
				switch i {
				case 0:
					raw = fmt.Sprintf("%s\t%s", raw, h)
				default:
					raw = fmt.Sprintf("%s %s", raw, h)
				}
			}

			rec, err := NewHostsRecord(raw)
			if err != nil {
				continue
			}

			output = append(output, *rec)
		}
	}

	h.Lock()
	defer h.Unlock()
	h.Records = output

	return nil
}

func (h Hosts) getHostPosition(ip string, host string) int {
	for i, rec := range h.Records {
		if rec.IsComment() || rec.Raw == "" {
			continue
		}

		if ip == rec.LocalPtr && in(host, rec.Hosts) {
			return i
		}
	}

	return -1
}

func (h Hosts) getIpPosition(ip string) int {
	for i, r := range h.Records {
		if r.IsComment() || r.Raw == "" {
			continue
		}

		if r.LocalPtr == ip {
			return i
		}
	}

	return -1
}

func in(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

func buildRawLine(ip string, hosts []string) string {
	output := ip
	for _, host := range hosts {
		output = fmt.Sprintf("%s %s", output, host)
	}

	return output
}
