package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	ErrPathNotExist  = errors.New("path does not exist")
	ErrChirpNotFound = errors.New("chirp not found")
)

type DBStructure struct {
	Chirps         map[int]Chirp `json:"chirps"`
	LastInsertedId int
}

func (dbs *DBStructure) AddChirp(chirp Chirp) Chirp {
	chirp.ID = dbs.LastInsertedId + 1
	dbs.Chirps[chirp.ID] = chirp
	dbs.LastInsertedId++
	return chirp
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	return db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Println(err)
		return Chirp{}, err
	}

	if len(dbStructure.Chirps) == 0 {
		// Init DBStructure
		dbStructure.Chirps = make(map[int]Chirp)
		dbStructure.LastInsertedId = 0
	}

	chirp := dbStructure.AddChirp(Chirp{
		Body: body,
	})

	err = db.writeDB(dbStructure)

	if err != nil {
		log.Println(err)
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		return []Chirp{}, err
	}

	var chirps []Chirp
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]

	if !ok {
		return Chirp{}, ErrChirpNotFound
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("writeDB: failed marshalling the database %v", err)
	}

	file, err := os.OpenFile(db.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("writeDB: %w", err)
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("writeDB: %w", err)
	}

	return nil
}

func (db *DB) readDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	file, err := os.OpenFile(db.path, os.O_RDONLY, 0644)

	if err != nil {
		log.Println(err)
		return DBStructure{}, fmt.Errorf("readDB: %w", err)
	}

	defer file.Close()
	// Check if the file is empty
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		return DBStructure{}, fmt.Errorf("database: %w", err)
	}
	if info.Size() == 0 {
		return DBStructure{}, nil
	}

	var dbStructure DBStructure
	err = json.NewDecoder(file).Decode(&dbStructure)
	if err != nil {
		log.Println(err)
		return DBStructure{}, fmt.Errorf("database: %w", err)
	}

	return dbStructure, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)

	if os.IsNotExist(err) {
		return ErrPathNotExist
	}

	return nil
}
