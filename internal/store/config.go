package store

func createTables(s *StoreDB) error {
	// Create users
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				login VARCHAR(50) UNIQUE NOT NULL CHECK(login ~* '^[a-zA-Z0-9_]+$'),
				password TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				lastlogin_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				is_active BOOLEAN DEFAULT TRUE,
				role VARCHAR(20) NOT NULL);`)

	// Create orders
	_, err = s.db.Exec(`CREATE TABLE orders (
						id SERIAL PRIMARY KEY,
						status VARCHAR(20),
						accrual INTEGER,
						uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						user_id INTEGER REFERENCES users(id) ON DELETE CASCADE);`)

	// Crrate INDEXs
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);`)
	return err
}
