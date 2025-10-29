package store

import (
	"database/sql"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	_ "github.com/lib/pq"
)

type StoreDB struct {
	Config config.Config
	Logger logg.Logger
	db     *sql.DB
	isOpen bool
}

func NewStoreDB(config config.Config, logger logg.Logger) *StoreDB {
	tmpStoreDB := &StoreDB{Config: config, Logger: logger, isOpen: false}
	tmpStoreDB.Open()
	tmpStoreDB.Ping()

	err := createTables(tmpStoreDB)
	if err != nil {
		tmpStoreDB.Logger.RaiseFatal("NewStoreDB>createTables", err)
	}
	return tmpStoreDB
}

func (s *StoreDB) Open() {
	db, err := sql.Open("postgres", s.Config.GetDBURI())
	if err != nil {
		s.Logger.RaiseFatal("StoreDB>NewDB>Open", err)
	}
	s.db = db
	s.isOpen = true
}

func (s *StoreDB) Close() {
	s.db.Close()
}

func (s *StoreDB) Ping() {
	err := s.db.Ping()
	if err != nil {
		s.Logger.RaiseFatal("StoreDB>Ping", err)
	}
}

func (s *StoreDB) GetDB() any {
	return s.db
}
