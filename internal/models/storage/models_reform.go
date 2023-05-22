// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package storage

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type routerTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *routerTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("routers").
func (v *routerTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *routerTableType) Columns() []string {
	return []string{
		"id",
		"user_id",
		"name",
		"address",
		"username",
		"password",
		"lease_period_check",
		"created_at",
		"updated_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *routerTableType) NewStruct() reform.Struct {
	return new(Router)
}

// NewRecord makes a new record for that table.
func (v *routerTableType) NewRecord() reform.Record {
	return new(Router)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *routerTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// RouterTable represents routers view or table in SQL database.
var RouterTable = &routerTableType{
	s: parse.StructInfo{
		Type:    "Router",
		SQLName: "routers",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "UserID", Type: "string", Column: "user_id"},
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "Address", Type: "string", Column: "address"},
			{Name: "Username", Type: "string", Column: "username"},
			{Name: "Password", Type: "string", Column: "password"},
			{Name: "LeasePeriodCheck", Type: "time.Duration", Column: "lease_period_check"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(Router).Values(),
}

// String returns a string representation of this struct or record.
func (s Router) String() string {
	res := make([]string, 9)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Address: " + reform.Inspect(s.Address, true)
	res[4] = "Username: " + reform.Inspect(s.Username, true)
	res[5] = "Password: " + reform.Inspect(s.Password, true)
	res[6] = "LeasePeriodCheck: " + reform.Inspect(s.LeasePeriodCheck, true)
	res[7] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[8] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Router) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.Name,
		s.Address,
		s.Username,
		s.Password,
		s.LeasePeriodCheck,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Router) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.Name,
		&s.Address,
		&s.Username,
		&s.Password,
		&s.LeasePeriodCheck,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *Router) View() reform.View {
	return RouterTable
}

// Table returns Table object for that record.
func (s *Router) Table() reform.Table {
	return RouterTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Router) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Router) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Router) HasPK() bool {
	return s.ID != RouterTable.z[RouterTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (s *Router) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = RouterTable
	_ reform.Struct = (*Router)(nil)
	_ reform.Table  = RouterTable
	_ reform.Record = (*Router)(nil)
	_ fmt.Stringer  = (*Router)(nil)
)

type logTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *logTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("logs").
func (v *logTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *logTableType) Columns() []string {
	return []string{
		"id",
		"router_id",
		"time",
		"level",
		"message",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *logTableType) NewStruct() reform.Struct {
	return new(Log)
}

// NewRecord makes a new record for that table.
func (v *logTableType) NewRecord() reform.Record {
	return new(Log)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *logTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// LogTable represents logs view or table in SQL database.
var LogTable = &logTableType{
	s: parse.StructInfo{
		Type:    "Log",
		SQLName: "logs",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "RouterID", Type: "string", Column: "router_id"},
			{Name: "Time", Type: "time.Time", Column: "time"},
			{Name: "Level", Type: "LogLevel", Column: "level"},
			{Name: "Message", Type: "string", Column: "message"},
		},
		PKFieldIndex: 0,
	},
	z: new(Log).Values(),
}

// String returns a string representation of this struct or record.
func (s Log) String() string {
	res := make([]string, 5)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "RouterID: " + reform.Inspect(s.RouterID, true)
	res[2] = "Time: " + reform.Inspect(s.Time, true)
	res[3] = "Level: " + reform.Inspect(s.Level, true)
	res[4] = "Message: " + reform.Inspect(s.Message, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Log) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.RouterID,
		s.Time,
		s.Level,
		s.Message,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Log) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.RouterID,
		&s.Time,
		&s.Level,
		&s.Message,
	}
}

// View returns View object for that struct.
func (s *Log) View() reform.View {
	return LogTable
}

// Table returns Table object for that record.
func (s *Log) Table() reform.Table {
	return LogTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Log) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Log) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Log) HasPK() bool {
	return s.ID != LogTable.z[LogTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (s *Log) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = LogTable
	_ reform.Struct = (*Log)(nil)
	_ reform.Table  = LogTable
	_ reform.Record = (*Log)(nil)
	_ fmt.Stringer  = (*Log)(nil)
)

type hostTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *hostTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("hosts").
func (v *hostTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *hostTableType) Columns() []string {
	return []string{
		"id",
		"router_id",
		"name",
		"address",
		"mac_address",
		"host_name",
		"last_online",
		"is_online",
		"online_timeout",
		"created_at",
		"updated_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *hostTableType) NewStruct() reform.Struct {
	return new(Host)
}

// NewRecord makes a new record for that table.
func (v *hostTableType) NewRecord() reform.Record {
	return new(Host)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *hostTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// HostTable represents hosts view or table in SQL database.
var HostTable = &hostTableType{
	s: parse.StructInfo{
		Type:    "Host",
		SQLName: "hosts",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "RouterID", Type: "string", Column: "router_id"},
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "Address", Type: "sql.NullString", Column: "address"},
			{Name: "MacAddress", Type: "sql.NullString", Column: "mac_address"},
			{Name: "HostName", Type: "sql.NullString", Column: "host_name"},
			{Name: "LastOnline", Type: "time.Time", Column: "last_online"},
			{Name: "IsOnline", Type: "bool", Column: "is_online"},
			{Name: "OnlineTimeout", Type: "time.Duration", Column: "online_timeout"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(Host).Values(),
}

// String returns a string representation of this struct or record.
func (s Host) String() string {
	res := make([]string, 11)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "RouterID: " + reform.Inspect(s.RouterID, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Address: " + reform.Inspect(s.Address, true)
	res[4] = "MacAddress: " + reform.Inspect(s.MacAddress, true)
	res[5] = "HostName: " + reform.Inspect(s.HostName, true)
	res[6] = "LastOnline: " + reform.Inspect(s.LastOnline, true)
	res[7] = "IsOnline: " + reform.Inspect(s.IsOnline, true)
	res[8] = "OnlineTimeout: " + reform.Inspect(s.OnlineTimeout, true)
	res[9] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[10] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Host) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.RouterID,
		s.Name,
		s.Address,
		s.MacAddress,
		s.HostName,
		s.LastOnline,
		s.IsOnline,
		s.OnlineTimeout,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Host) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.RouterID,
		&s.Name,
		&s.Address,
		&s.MacAddress,
		&s.HostName,
		&s.LastOnline,
		&s.IsOnline,
		&s.OnlineTimeout,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *Host) View() reform.View {
	return HostTable
}

// Table returns Table object for that record.
func (s *Host) Table() reform.Table {
	return HostTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Host) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Host) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Host) HasPK() bool {
	return s.ID != HostTable.z[HostTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (s *Host) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = HostTable
	_ reform.Struct = (*Host)(nil)
	_ reform.Table  = HostTable
	_ reform.Record = (*Host)(nil)
	_ fmt.Stringer  = (*Host)(nil)
)

func init() {
	parse.AssertUpToDate(&RouterTable.s, new(Router))
	parse.AssertUpToDate(&LogTable.s, new(Log))
	parse.AssertUpToDate(&HostTable.s, new(Host))
}