package store

import (
	"database/sql"

	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/internal/logg"
	_ "github.com/lib/pq"
)

type StoreDB struct {
	Args   config.Argser
	Logg   logg.Logger
	db     *sql.DB
	isOpen bool
}

func NewStoreDB(argser config.Argser, logger logg.Logger) *StoreDB {
	tmpStoreDB := &StoreDB{Args: argser, Logg: logger, isOpen: false}
	tmpStoreDB.Open()
	tmpStoreDB.Ping()

	err := createTables(tmpStoreDB)
	if err != nil {
		tmpStoreDB.Logg.RaiseFatal("NewStoreDB>createTables", err)
	}
	return tmpStoreDB
}

func (s *StoreDB) Open() {
	db, err := sql.Open("postgres", s.Args.GetDbUri())
	if err != nil {
		s.Logg.RaiseFatal("StoreDB>NewDB>Open", err)
	}
	s.db = db
	s.isOpen = true
}

func (s *StoreDB) Clouse() {
	s.db.Close()
}

func (s *StoreDB) Ping() {
	err := s.db.Ping()
	if err != nil {
		s.Logg.RaiseFatal("StoreDB>Ping", err)
	}
}

func (s *StoreDB) GetDB() any {
	return s.db
}
