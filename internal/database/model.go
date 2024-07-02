package database

import "time"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type UserRefreshToken struct {
	Token           string `json:"token"`
	ExpiriationTime time.Time
}

func (rf *UserRefreshToken) IsExpired() bool {
	if rf == nil {
		return true
	}

	return time.Now().After(rf.ExpiriationTime)
}

type User struct {
	ID           int               `json:"id"`
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	RefreshToken *UserRefreshToken `json:"refreshToken,omitempty"`
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetPassword(password string) {
	u.Password = password
}
