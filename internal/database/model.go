package database

import "time"

type SortType string

const (
	SortAsc  SortType = "asc"
	SortDesc SortType = "desc"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

type UserRefreshToken struct {
	Token           string    `json:"token"`
	ExpiriationTime time.Time `json:"expiration_time"`
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
	RefreshToken *UserRefreshToken `json:"refresh_token,omitempty"`
	IsChirpyRed  bool              `json:"is_chirpy_red"`
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) SetIsChirpyRed(state bool) {
	u.IsChirpyRed = state
}
