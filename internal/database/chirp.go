package database

import (
	"log"
	"slices"
	"strconv"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Println(err)
		return Chirp{}, err
	}

	chirp := dbStructure.AddChirp(Chirp{
		Body:     body,
		AuthorID: authorID,
	})

	err = db.writeDB(dbStructure)

	if err != nil {
		log.Println(err)
		return Chirp{}, err
	}

	return chirp, nil
}

type QueryGetChirps struct {
	AuthorID string
	Sort     string
}

func (db *DB) GetChirps(in QueryGetChirps) ([]Chirp, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		return []Chirp{}, err
	}

	var chirps = []Chirp{}
	for _, chirp := range dbStructure.Chirps {
		if in.AuthorID != "" && in.AuthorID != strconv.Itoa(chirp.AuthorID) {
			continue
		}

		chirps = append(chirps, chirp)
	}

	slices.SortFunc(chirps, func(a Chirp, b Chirp) int {
		if in.Sort == string(SortDesc) {
			return b.ID - a.ID
		}

		return a.ID - b.ID
	})

	return chirps, nil
}

func (db *DB) DeleteChirpByID(chirpID int) error {
	dbStructure, err := db.readDB()

	if err != nil {
		return err
	}

	dbStructure.DeleteChirpByID(chirpID)

	err = db.writeDB(dbStructure)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
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
