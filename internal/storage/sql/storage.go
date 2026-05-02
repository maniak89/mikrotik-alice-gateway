package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"

	storageModels "mikrotik-alice-gateway/internal/models/storage"
	storagePkg "mikrotik-alice-gateway/internal/storage"
)

type storage struct {
	config     Config
	connection *sql.DB
	db         *sqlx.DB
	driver     string
}

func New(config Config) *storage {
	return &storage{
		config: config,
	}
}

func (s *storage) Connect(ctx context.Context) error {
	logger := log.Ctx(ctx)
	parsedConnectionString, err := url.Parse(s.config.ConnectionString)
	if err != nil {
		logger.Error().Err(err).Str("connection_string", s.config.ConnectionString).Msg("Failed parse connection string")
		return err
	}
	s.driver = parsedConnectionString.Scheme
	var found bool
	for _, driver := range sql.Drivers() {
		if driver == s.driver {
			found = true
			break
		}
	}
	if !found {
		err := errors.New("not supported db driver: " + s.driver)
		logger.Error().Err(err).Msg("")
		return err
	}
	connectionString := s.config.ConnectionString
	// sqlite3 не работает, если в начале connectionString - схема.
	// а без нее не работает postgres(не парсит аргументы)
	if s.driver == "sqlite3" {
		connectionString = s.config.ConnectionString[len(s.driver)+3:]
	}
	sqlDB, err := sql.Open(s.driver, connectionString)
	if err != nil {
		logger.Error().Err(err).Msg("Failed open sql connection")
		return err
	}
	s.connection = sqlDB

	t := sqlx.NewDb(sqlDB, s.driver)
	s.db = t
	return nil
}

func (s *storage) Disconnect(ctx context.Context) error {
	logger := log.Ctx(ctx)
	if s.connection == nil {
		err := storagePkg.ErrInvalidState
		logger.Error().Err(err).Msg("Invalid state of connection")
		return err
	}
	if err := s.connection.Close(); err != nil {
		logger.Error().Err(err).Msg("Failed close of connection")
		return err
	}
	s.connection = nil
	return nil
}

func (s *storage) Routers(ctx context.Context) ([]*storageModels.Router, error) {
	logger := log.Ctx(ctx)

	hosts, err := fetchRows[storageModels.Host](ctx, s.db, "SELECT id, router_id, name, address, mac_address, host_name, last_online,"+
		" is_online, online_timeout, created_at, updated_at FROM hosts")
	if err != nil {
		logger.Error().Err(err).Msg("Failed find hosts")
		return nil, err
	}

	routers, err := fetchRows[storageModels.Router](ctx, s.db, "select id, user_id, name, address, username, password, lease_period_check,"+
		"created_at, updated_at from routers")
	if err != nil {
		logger.Error().Err(err).Msg("Failed find routers")
		return nil, err
	}

	hostMap := make(map[string][]*storageModels.Host, len(routers))
	for _, host := range hosts {
		hostMap[host.RouterID] = append(hostMap[host.RouterID], host)
	}

	for _, router := range routers {
		router.Hosts = hostMap[router.ID]
	}

	return routers, nil
}

func (s *storage) Log(ctx context.Context, routerID string, level storageModels.LogLevel, msg string) {
	logger := log.Ctx(ctx).With().Str("router_id", routerID).Str("level", string(level)).Str("msg", msg).Logger()
	switch level {
	case storageModels.Error:
		logger.Error().Msg("")
	case storageModels.Info:
		logger.Debug().Msg("")
	}
	if s.config.LogOnlyErrors && level != storageModels.Error {
		return
	}
	if _, err := s.db.NamedExecContext(ctx, `INSERT INTO logs(router_id, "time", level, message) VALUES (:router_id, :time, :level, :message)`, &storageModels.Log{
		RouterID: routerID,
		Level:    level,
		Time:     time.Now(),
		Message:  msg,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed add log")
	}
}

func (s *storage) UpdateHost(ctx context.Context, host *storageModels.Host) error {
	if _, err := s.db.NamedExecContext(ctx, `UPDATE public.hosts
	SET router_id=:router_id, name=:name, address=:address, mac_address=:mac_address, host_name=:host_name,
	last_online=:last_online, is_online=:is_online, online_timeout=:online_timeout, updated_at=now()
	WHERE id=:id`, host); err != nil {
		return fmt.Errorf("failed update: %w", err)
	}

	return nil
}

func fetchRows[T any](ctx context.Context, db *sqlx.DB, query string, args ...any) ([]*T, error) {
	logger := log.Ctx(ctx)

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		logger.Error().Err(err).Msg("Failed find routers")
		return nil, err
	}

	result := make([]*T, 0)
	for rows.Next() {
		var router T
		if err := rows.StructScan(&router); err != nil {
			return nil, fmt.Errorf("failed scan: %w", err)
		}

		result = append(result, &router)
	}

	return result, nil
}
