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

	h.Lock()
	defer h.Unlock()
	position := h.getIpPosition(ip)
	if position == -1 {
		rec, err := NewHostsRecord(buildRawLine(ip, hosts))
		if err != nil {
			return err
		}
		// Ip line is not in file, so we just append our new record.

		h.Records = append(h.Records, *rec)
	} else {
		// Otherwise, we replace the record in the correct position
		newHosts := h.Records[position].Hosts
		for _, addHost := range hosts {
			if itemInSlice(addHost, newHosts) {
				continue
			}

			newHosts = append(newHosts, addHost)
		}
		rec, err := NewHostsRecord(buildRawLine(ip, newHosts))
		if err != nil {
			return err
		}
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
		return errors.New(fmt.Sprintf("%q is an invalid IP address.", ip))
	}

	for _, r := range h.Records {

		// Bad lines or comments just get readded.
		if r.Error != nil || r.IsComment() || r.LocalPtr != ip {
			output = append(output, r)
			continue
		}

		var newHosts []string
		for _, checkHost := range r.Hosts {
			if !itemInSlice(checkHost, hosts) {
				newHosts = append(newHosts, checkHost)
			}
		}

		// If hosts is empty, skip the line completely.
		if len(newHosts) > 0 {
			nrRaw := r.LocalPtr

			for _, host := range newHosts {
				nrRaw = fmt.Sprintf("%s %s", nrRaw, host)
			}
			nr, err := NewHostsRecord(nrRaw)
			if err != nil {
				return err
			}

			output = append(output, *nr)
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

		if ip == rec.LocalPtr && itemInSlice(host, rec.Hosts) {
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

func itemInSlice(item string, list []string) bool {
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
