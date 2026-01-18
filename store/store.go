package store

import (
	"database/sql"
	"time"

	"llm-mock-server/models"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.InitSchema(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS rules (
		id TEXT PRIMARY KEY,
		name TEXT,
		type TEXT,
		pattern TEXT,
		response TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		enabled INTEGER
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) AddRule(rule *models.MockRule) error {
	now := time.Now()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	enabled := 0
	if rule.Enabled {
		enabled = 1
	}

	query := `INSERT INTO rules (id, name, type, pattern, response, created_at, updated_at, enabled) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, rule.ID, rule.Name, rule.Type, rule.Pattern, rule.Response, rule.CreatedAt, rule.UpdatedAt, enabled)
	return err
}

func (s *Store) GetEnabledRules() ([]models.MockRule, error) {
	query := `SELECT id, name, type, pattern, response, created_at, updated_at, enabled FROM rules WHERE enabled = 1`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.MockRule
	for rows.Next() {
		var r models.MockRule
		var enabled int
		err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Pattern, &r.Response, &r.CreatedAt, &r.UpdatedAt, &enabled)
		if err != nil {
			return nil, err
		}
		r.Enabled = enabled == 1
		rules = append(rules, r)
	}
	return rules, nil
}

func (s *Store) GetAllRules() ([]models.MockRule, error) {
	query := `SELECT id, name, type, pattern, response, created_at, updated_at, enabled FROM rules`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.MockRule
	for rows.Next() {
		var r models.MockRule
		var enabled int
		err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Pattern, &r.Response, &r.CreatedAt, &r.UpdatedAt, &enabled)
		if err != nil {
			return nil, err
		}
		r.Enabled = enabled == 1
		rules = append(rules, r)
	}
	return rules, nil
}

func (s *Store) DeleteRule(id string) error {
	query := `DELETE FROM rules WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}
