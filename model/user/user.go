package user

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"mend/db"
	"mend/db/tp"
)

// User - object with data descripted one user
type User struct {
	ID        tp.ID  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

// NewUser - create new, empty User object.
func NewUser() *User {
	return new(User)
}

func (u *User) String() string {
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Error(err)
		return "?"
	}
	return string(data)
}

func (u *User) asRow() tp.Row {
	return tp.Row{
		"id":         u.ID,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"age":        u.Age,
	}
}

func Count() int64 {
	return db.Db().Count("user")
}

// AddToDatabase adds the user to database.
func (u *User) AddToDatabase() error {
	if u.ID > 0 {
		return fmt.Errorf("can't add user with specified id (can be zero)")
	}
	id, err := db.Db().Add("user", u.asRow())
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (u *User) UpdateInDatabase() error {
	return db.Db().Update("user", u.asRow())
}

func GetAll() ([]tp.Row, error) {
	return db.Db().GetAll("user")
}

func GetWithID(id tp.ID) (tp.Row, error) {
	return db.Db().GetWithID("user", id)
}

// DeleteFromDatabase removes from database user with passed ID.
func DeleteFromDatabase(id tp.ID) error {
	return db.Db().Delete("user", id)
}
