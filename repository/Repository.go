package repository

import "database/sql"

// InitializeDatabase creates tables if they don't exist
func InitializeDatabase(db *sql.DB) error {
	// Create Posts table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Create Users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
