package models

import (
	"database/sql"
	"fmt"
	"regexp"
)

type Pass struct {
	ID           uint
	URL          string
	Password     string
	PasswordHash string
}

type PassDB interface {
	Create(pass *Pass) error
	// Update(pass *Pass) error
	// Delete(id uint) error
	Get() ([]Pass, error)
}

type PassService interface {
	PassDB
}

var _ PassDB = &passService{}

type passService struct {
	PassDB
}

func NewPassService(db *sql.DB) PassService {
	ps := &passSQL{db}

	pv := newPassValidator(ps)

	return &passService{PassDB: pv}
}

var _ PassDB = &passValidator{}

func newPassValidator(pdb PassDB) *passValidator {
	return &passValidator{
		PassDB: pdb,
		urlRegex: regexp.MustCompile(`(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?\/[a-zA-Z0-9]{2,}|((https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?)|(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})? 
		`),
	}
}

type passValidator struct {
	PassDB
	urlRegex *regexp.Regexp
}

var _ PassDB = &passSQL{}

type passSQL struct {
	db *sql.DB
}

func (pv *passValidator) Create(pass *Pass) error {
	fmt.Println(pass.URL)
	err := runModelValFuncs[Pass](pass, pv.validURL)
	if err != nil {
		return err
	}

	return pv.PassDB.Create(pass)
}

func (p *passSQL) Create(pass *Pass) error {
	stmt, err := p.db.Prepare(`INSERT INTO passes(url, password) VALUES(?, ?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pass.URL, pass.Password); err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}

func (pv *passValidator) Get() ([]Pass, error) {
	return nil, nil
}

func (p *passSQL) Get() ([]Pass, error) {
	return nil, nil
}

// Validators
func (pv *passValidator) idGreaterThanZero(pass *Pass) error {
	if pass.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

func (pv *passValidator) validURL(pass *Pass) error {
	if pass.URL == "" {
		return nil
	}
	if !pv.urlRegex.MatchString(pass.URL) {
		return ErrURLInvalid
	}
	return nil
}
