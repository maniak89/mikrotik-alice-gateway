//go:generate reform
package storage

import (
	"database/sql"
	"time"
)

//reform:routers
type Router struct {
	ID               string        `reform:"id,pk"`
	UserID           string        `reform:"user_id"`
	Name             string        `reform:"name"`
	Address          string        `reform:"address"`
	Username         string        `reform:"username"`
	Password         string        `reform:"password"`
	LeasePeriodCheck time.Duration `reform:"lease_period_check"`
	CreatedAt        time.Time     `reform:"created_at"`
	UpdatedAt        time.Time     `reform:"updated_at"`
	Hosts            []*Host       `reform:"-"`
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

//reform:logs
type Log struct {
	ID       string    `reform:"id,pk"`
	RouterID string    `reform:"router_id"`
	Time     time.Time `reform:"time"`
	Level    LogLevel  `reform:"level"`
	Message  string    `reform:"message"`
}

//reform:hosts
type Host struct {
	ID            string         `reform:"id,pk"`
	RouterID      string         `reform:"router_id"`
	Name          string         `reform:"name"`
	Address       sql.NullString `reform:"address"`
	MacAddress    sql.NullString `reform:"mac_address"`
	HostName      sql.NullString `reform:"host_name"`
	LastOnline    time.Time      `reform:"last_online"`
	IsOnline      bool           `reform:"is_online"`
	OnlineTimeout time.Duration  `reform:"online_timeout"`
	CreatedAt     time.Time      `reform:"created_at"`
	UpdatedAt     time.Time      `reform:"updated_at"`
}

func (s *Host) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Host) Equal(o *Host) bool {
	return s.ID == o.ID &&
		s.Address.String == o.Address.String &&
		s.MacAddress.String == o.MacAddress.String &&
		s.HostName.String == o.HostName.String &&
		s.Name == o.Name
}
