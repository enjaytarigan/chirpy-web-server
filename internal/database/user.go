package database

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrUserNotFound           = errors.New("User Not Found")
	ErrEmailAlreadyRegistered = errors.New("email has been registered")
)

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Println(err)
		return User{}, err
	}

	user := dbStructure.AddUser(User{
		Email:    email,
		Password: password,
	})

	err = db.writeDB(dbStructure)

	if err != nil {
		log.Println(err)
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Printf("GetUserByEmail: %v", err)
		return User{}, err
	}

	var userIdByEmail int

	for _, user := range dbStructure.Users {
		if user.Email == email {
			userIdByEmail = user.ID
			break
		}
	}

	user, ok := dbStructure.Users[userIdByEmail]

	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (db *DB) CheckUserExist(email string) bool {
	_, err := db.GetUserByEmail(email)

	return err == nil
}

func (db *DB) UpdateUser(userID int, email string, password string) (User, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Printf("GetUserByEmail: %v", err)
		return User{}, err
	}

	if _, found := dbStructure.FindUser(func(u User) bool {
		return u.Email == email
	}); found {
		return User{}, ErrEmailAlreadyRegistered
	}

	user, ok := dbStructure.FindUserByID(userID)

	if !ok {
		return User{}, ErrUserNotFound
	}

	user.SetEmail(email)
	user.SetPassword(password)

	dbStructure.Users[user.ID] = user

	if err := db.writeDB(dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) SaveRefreshToken(userID int, refreshToken *UserRefreshToken) error {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Printf("SaveRefreshToken: %v", err)
		return err
	}

	user, found := dbStructure.FindUserByID(userID)

	if !found {
		return ErrUserNotFound
	}

	user.RefreshToken = refreshToken
	dbStructure.Users[user.ID] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return fmt.Errorf("SaveRefreshToken: %w", err)
	}

	return nil
}

func (db *DB) GetUserByID(userID int) (User, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Printf("SaveRefreshToken: %v", err)
		return User{}, err
	}

	user, found := dbStructure.FindUserByID(userID)

	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (db *DB) GetUserByRefreshToken(refreshToken string) (User, error) {
	dbStructure, err := db.readDB()

	if err != nil {
		log.Printf("GetUserByRefreshToken: %v", err)
		return User{}, err
	}

	user, found := dbStructure.FindUser(func(u User) bool {
		if u.RefreshToken == nil {
			return false
		}

		return u.RefreshToken.Token == refreshToken
	})

	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}
