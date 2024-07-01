package database

import "log"

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Println(err)
		return User{}, err
	}

	user := dbStructure.AddUser(User{
		Email: email,
	})

	err = db.writeDB(dbStructure)

	if err != nil {
		log.Println(err)
		return User{}, err
	}

	return user, nil
}
