package sql

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects"

	storageModels "mikrotik-alice-gateway/internal/models/storage"
	storagePkg "mikrotik-alice-gateway/internal/storage"
)

type storage struct {
	config     Config
	connection *sql.DB
	db         *reform.Querier
	reformDB   reform.DBTXContext
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

	t := reform.NewDB(sqlDB, dialects.ForDriver(s.driver), reform.NewPrintfLogger(logger.Printf))
	s.db = t.Querier
	s.reformDB = t
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
	rows, err := s.db.SelectAllFrom(storageModels.RouterTable, "")
	if err != nil {
		logger.Error().Err(err).Msg("Failed find routers")
		return nil, err
	}
	result := make([]*storageModels.Router, 0, len(rows))
	for _, r := range rows {
		router := r.(*storageModels.Router)
		logger := logger.With().Str("router_id", router.ID).Logger()
		result = append(result, router)
		rows, err := s.db.SelectAllFrom(storageModels.HostTable, "where router_id = "+s.db.Placeholder(1), router.ID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed find routers")
			return nil, err
		}
		for _, r := range rows {
			router.Hosts = append(router.Hosts, r.(*storageModels.Host))
		}
	}
	return result, nil
}

func (s *storage) Log(ctx context.Context, routerID string, level storageModels.LogLevel, msg string) {
	logger := log.Ctx(ctx).With().Str("router_id", routerID).Str("level", string(level)).Str("msg", msg).Logger()
	if err := s.db.WithContext(ctx).Insert(&storageModels.Log{
		RouterID: routerID,
		Level:    level,
		Message:  msg,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed add log")
		return
	}
	return
}

func (s *storage) UpdateHost(ctx context.Context, host *storageModels.Host) error {
	return s.db.WithContext(ctx).Save(host)
}
