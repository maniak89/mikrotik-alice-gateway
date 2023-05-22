package device_provider

import (
	"context"
)

type DeviceProvider interface {
	Init(ctx context.Context) error
	Leases(ctx context.Context) ([]Lease, error)
	Resource(ctx context.Context) (Resource, error)
}

type Resource struct {
	Version   string
	BoardName string
	Platform  string
}

type LeaseStatus string

const (
	LeaseStatusBound   LeaseStatus = "bound"
	LeaseStatusWaiting LeaseStatus = "waiting"
)

type Lease struct {
	Address    string
	MacAddress string
	HostName   string
	Status     LeaseStatus
}
