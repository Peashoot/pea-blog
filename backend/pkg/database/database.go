package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func Connect(databaseURL string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	// Determine database type based on URL
	if strings.HasPrefix(databaseURL, "postgres://") {
		db, err = sql.Open("postgres", databaseURL)
	} else {
		// Default to SQLite for development
		if databaseURL == "" {
			databaseURL = "./pea_blog.db?_loc=auto&_busy_timeout=5000"
		}
		db, err = sql.Open("sqlite", databaseURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for SQLite for better concurrency
	if strings.HasPrefix(databaseURL, "./") || strings.HasPrefix(databaseURL, "file:") || !strings.Contains(databaseURL, "://") {
		_, err = db.Exec("PRAGMA journal_mode=WAL;")
		if err != nil {
			return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
		}
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	// Detect database type
	var dbType string = "sqlite"
	var rows *sql.Rows
	var err error
	if rows, err = db.Query("SELECT version()"); err == nil {
		rows.Close()
		dbType = "postgres"
	}

	var migrations []string

	if dbType == "postgres" {
		migrations = []string{
			`CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(50) UNIQUE NOT NULL,
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				avatar TEXT,
				role VARCHAR(20) DEFAULT 'user',
				fingerprint VARCHAR(255) UNIQUE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			)`,

			`CREATE TABLE IF NOT EXISTS articles (
				id SERIAL PRIMARY KEY,
				title VARCHAR(200) UNIQUE NOT NULL,
				content TEXT NOT NULL,
				summary VARCHAR(500) NOT NULL,
				tags TEXT[] DEFAULT '{}',
				author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				status VARCHAR(20) DEFAULT 'draft',
				view_count INTEGER DEFAULT 0,
				like_count INTEGER DEFAULT 0,
				comment_count INTEGER DEFAULT 0,
				cover_image TEXT,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				published_at TIMESTAMP WITH TIME ZONE,
				deleted_at TIMESTAMP WITH TIME ZONE
			)`,

			`CREATE TABLE IF NOT EXISTS comments (
				id SERIAL PRIMARY KEY,
				content TEXT NOT NULL,
				author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
				parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				deleted_at TIMESTAMP WITH TIME ZONE
			)`,

			`CREATE TABLE IF NOT EXISTS likes (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				UNIQUE(user_id, article_id)
			)`,
		}
	} else {
		// SQLite migrations
		migrations = []string{
			`CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				username VARCHAR(50) UNIQUE NOT NULL,
				email VARCHAR(255) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				avatar TEXT,
				role VARCHAR(20) DEFAULT 'user',
				fingerprint VARCHAR(255) UNIQUE,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`,

			`CREATE TABLE IF NOT EXISTS articles (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title VARCHAR(200) UNIQUE NOT NULL,
				content TEXT NOT NULL,
				summary VARCHAR(500) NOT NULL,
				tags TEXT DEFAULT '',
				author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				status VARCHAR(20) DEFAULT 'draft',
				view_count INTEGER DEFAULT 0,
				like_count INTEGER DEFAULT 0,
				comment_count INTEGER DEFAULT 0,
				cover_image TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				published_at DATETIME,
				deleted_at DATETIME
			)`,

			`CREATE TABLE IF NOT EXISTS comments (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT NOT NULL,
				author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
				parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`,

			`CREATE TABLE IF NOT EXISTS likes (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(user_id, article_id)
			)`,
		}
	}

	// Execute migrations
	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	// Add new columns to tables if they don't exist
	tables := []string{"articles", "comments"}
	for _, table := range tables {
		rows, err = db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
		if err == nil {
			defer rows.Close()
			var columnName string
			hasDeletedAtColumn := false
			for rows.Next() {
				var cid int
				var dflt_value interface{}
				var notnull bool
				var pk int
				var dataType string
				err := rows.Scan(&cid, &columnName, &dataType, &notnull, &dflt_value, &pk)
				if err == nil && columnName == "deleted_at" {
					hasDeletedAtColumn = true
					break
				}
			}
			if !hasDeletedAtColumn {
				_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN deleted_at DATETIME", table))
				if err != nil {
					return fmt.Errorf("failed to add deleted_at column to %s: %w", table, err)
				}
			}
		}
	}

	// Create indexes
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments(article_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id)`,
		`CREATE INDEX IF NOT EXISTS idx_likes_article_id ON likes(article_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_articles_title ON articles(title)`,
	}

	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			// Ignore index creation errors (might already exist)
			continue
		}
	}

	// Add fingerprint column to users table if it doesn't exist
	rows, err = db.Query("PRAGMA table_info(users)")
	if err == nil {
		defer rows.Close()
		var columnName string
		hasFingerprintColumn := false
		for rows.Next() {
			// The column name is the second column in the result set for PRAGMA table_info
			var cid int
			var dflt_value interface{}
			var notnull bool
			var pk int
			var dataType string
			err := rows.Scan(&cid, &columnName, &dataType, &notnull, &dflt_value, &pk)
			if err == nil && columnName == "fingerprint" {
				hasFingerprintColumn = true
				break
			}
		}
		if !hasFingerprintColumn {
			_, err = db.Exec("ALTER TABLE users ADD COLUMN fingerprint VARCHAR(255)")
			if err != nil {
				return fmt.Errorf("failed to add fingerprint column: %w", err)
			}
			_, err = db.Exec("UPDATE users SET fingerprint = 'guest_' || id WHERE fingerprint IS NULL")
			if err != nil {
				return fmt.Errorf("failed to update existing fingerprints: %w", err)
			}
			_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_fingerprint ON users(fingerprint)")
			if err != nil {
				return fmt.Errorf("failed to create unique index on fingerprint: %w", err)
			}
		}
	}

	// Create default admin user
	var hashedPassword string = "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi" // "password"

	if dbType == "postgres" {
		_, err := db.Exec(`INSERT INTO users (username, email, password_hash, role) 
			VALUES ($1, $2, $3, $4) ON CONFLICT (username) DO NOTHING`,
			"admin", "admin@example.com", hashedPassword, "admin")
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	} else {
		// SQLite
		_, err := db.Exec(`INSERT OR IGNORE INTO users (username, email, password_hash, role, fingerprint) 
			VALUES (?, ?, ?, ?, NULL)`,
			"admin", "admin@example.com", hashedPassword, "admin")
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	return nil
}
