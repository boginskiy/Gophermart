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
	tmpStoreDB.createTables()
	return tmpStoreDB
}

func (s *StoreDB) Open() {
	db, err := sql.Open("postgres", s.Args.GetDB())
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

func (s *StoreDB) createTables() error {
	if !s.isOpen {
		s.Logg.RaiseFatal("StoreDB>createTables", ErrOpeningDB)
	}

	// Create users
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				login VARCHAR(50) UNIQUE NOT NULL CHECK(login ~* '^[a-zA-Z0-9_]+$'),
				password TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				lastlogin_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				is_active BOOLEAN DEFAULT TRUE,
				role VARCHAR(20) UNIQUE NOT NULL);`)

	// Create orders
	_, err = s.db.Exec(`CREATE TABLE orders (
						id SERIAL PRIMARY KEY,
						status VARCHAR(20),
						accrual INTEGER,
						uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						user_id INTEGER REFERENCES users(id) ON DELETE CASCADE);`)

	// Crrate INDEXs
	_, err = s.db.Exec(`CREATE INDEX idx_users_login ON users(login);`)
	if err != nil {
		s.Logg.RaiseFatal("StoreDB>createTables", err)
		return err
	}
	return nil
}
