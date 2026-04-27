package storage

import (
	"slices"
	"time"

	"github.com/lib/pq"
)

type Router struct {
	ID               string        `db:"id,pk"`
	UserID           string        `db:"user_id"`
	Name             string        `db:"name"`
	Address          string        `db:"address"`
	Username         string        `db:"username"`
	Password         string        `db:"password"`
	LeasePeriodCheck time.Duration `db:"lease_period_check"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
	Hosts            []*Host       `db:"-"`
}

func (s *Router) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Router) Equal(o *Router) bool {
	if s.ID != o.ID ||
		s.Password != o.Password ||
		s.Username != o.Username ||
		s.Address != o.Address ||
		s.Name != o.Name ||
		s.LeasePeriodCheck != o.LeasePeriodCheck ||
		len(s.Hosts) != len(o.Hosts) {
		return false
	}
	for i, sH := range s.Hosts {
		if !sH.Equal(o.Hosts[i]) {
			return false
		}
	}
	return true
}

type LogLevel string

const (
	Error LogLevel = "Error"
	Info  LogLevel = "Info"
)

type Log struct {
	ID       string    `db:"id,pk"`
	RouterID string    `db:"router_id"`
	Time     time.Time `db:"time"`
	Level    LogLevel  `db:"level"`
	Message  string    `db:"message"`
}

//db:hosts
type Host struct {
	ID            string         `db:"id,pk"`
	RouterID      string         `db:"router_id"`
	Name          string         `db:"name"`
	Address       pq.StringArray `db:"address"`
	MacAddress    pq.StringArray `db:"mac_address"`
	HostName      pq.StringArray `db:"host_name"`
	LastOnline    time.Time      `db:"last_online"`
	IsOnline      bool           `db:"is_online"`
	OnlineTimeout time.Duration  `db:"online_timeout"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

func (s *Host) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Host) Equal(o *Host) bool {
	slices.Sort(s.Address)
	slices.Sort(o.Address)
	slices.Sort(s.HostName)
	slices.Sort(o.HostName)
	slices.Sort(s.MacAddress)
	slices.Sort(o.MacAddress)
	return s.ID == o.ID &&
		slices.Equal(s.Address, o.Address) &&
		slices.Equal(s.MacAddress, o.MacAddress) &&
		slices.Equal(s.HostName, o.HostName) &&
		s.Name == o.Name
}
