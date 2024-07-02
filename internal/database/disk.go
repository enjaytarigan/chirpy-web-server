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
	Chirps              map[int]Chirp `json:"chirps"`
	LastInsertedChirpId int           `json:"lastInsertedChirpId"`
	Users               map[int]User  `json:"users"`
	LastInsertedUserId  int           `json:"lastInsertedUserId"`
}

func (dbs *DBStructure) AddChirp(chirp Chirp) Chirp {
	if len(dbs.Chirps) == 0 {
		dbs.LastInsertedChirpId = 0
		dbs.Chirps = make(map[int]Chirp)
	}

	chirp.ID = dbs.LastInsertedChirpId + 1
	dbs.Chirps[chirp.ID] = chirp
	dbs.LastInsertedChirpId++
	return chirp
}

func (dbs *DBStructure) AddUser(user User) User {
	if len(dbs.Users) == 0 {
		dbs.LastInsertedUserId = 0
		dbs.Users = make(map[int]User)
	}

	user.ID = dbs.LastInsertedUserId + 1
	dbs.Users[user.ID] = user
	dbs.LastInsertedUserId++
	return user
}

func (db *DBStructure) FindUser(f func(u User) bool) (u User, found bool) {
	if len(db.Users) == 0 {
		return User{}, false
	}

	for _, u := range db.Users {
		if f(u) {
			return u, true
		}
	}

	return User{}, false
}

func (db *DBStructure) FindUserByID(userID int) (u User, found bool) {
	u, ok := db.Users[userID]

	return u, ok
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
