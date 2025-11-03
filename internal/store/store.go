package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func (s *Store) Init() error {
	_, err := s.DB.Exec(`
        CREATE TABLE IF NOT EXISTS weather_logs (
            id SERIAL PRIMARY KEY,
            city TEXT,
            temp_c REAL,
            humidity REAL,
            retrieved_at TIMESTAMPTZ
        )
    `)
	return err
}

func (s *Store) SaveWeather(city string, temp, humidity float64) error {
	_, err := s.DB.Exec(
		`INSERT INTO weather_logs (city, temp_c, humidity, retrieved_at)
         VALUES ($1, $2, $3, NOW())`,
		city, temp, humidity,
	)
	return err
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.DB.Close()
}
