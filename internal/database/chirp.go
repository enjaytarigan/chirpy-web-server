package database

import "log"

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Println(err)
		return Chirp{}, err
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
