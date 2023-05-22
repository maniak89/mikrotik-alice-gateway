package common

import (
	"fmt"
	"time"
)

type Router struct {
	ID               string
	UserID           string
	Name             string
	Model            string
	SWVersion        string
	Manufacturer     string
	Connected        bool
	Hosts            []*Host
	AdditionalFields map[string]string
}

func (r *Router) String() string {
	return fmt.Sprintf("%s (%s %s)", r.Name, r.Model, r.ID)
}

type Host struct {
	ID      string
	Name    string
	Online  bool
	Changed time.Time
	Updated time.Time
}

func (r *Router) Clone() *Router {
	if r == nil {
		return nil
	}
	result := Router{
		ID:               r.ID,
		UserID:           r.UserID,
		Name:             r.Name,
		Model:            r.Model,
		SWVersion:        r.SWVersion,
		Manufacturer:     r.Manufacturer,
		Connected:        r.Connected,
		Hosts:            make([]*Host, 0, len(r.Hosts)),
		AdditionalFields: make(map[string]string, len(r.AdditionalFields)),
	}
	for _, h := range r.Hosts {
		result.Hosts = append(result.Hosts, h.Clone())
	}
	for k, v := range r.AdditionalFields {
		result.AdditionalFields[k] = v
	}
	return &result
}

func (h *Host) Clone() *Host {
	if h == nil {
		return nil
	}
	return &Host{
		ID:     h.ID,
		Name:   h.Name,
		Online: h.Online,
	}
}
