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
        CREATE TABLE IF NOT EXISTS weather_history (
            id SERIAL PRIMARY KEY,
            city VARCHAR(100) NOT NULL,
            temperature_c NUMERIC(5,2),
            humidity INTEGER,
            provider_count INTEGER,
            retrieved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	return err
}

func (s *Store) SaveWeather(city string, temp, humidity float64, providerCount int) error {
	_, err := s.DB.Exec(
		`INSERT INTO weather_history (city, temperature_c, humidity, provider_count, retrieved_at)
         VALUES ($1, $2, $3, $4, NOW())`,
		city, temp, humidity, providerCount,
	)
	return err
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.DB.Close()
}
