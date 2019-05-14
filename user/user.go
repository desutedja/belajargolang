package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// type User struct {
// 	Db *sql.DB
// }

type Dbase struct {
	Db *sql.DB
}

type UserModel struct {
	Id        int
	UserName  string
	FirstName string
	LastName  string
	Password  string
}

func (d *Dbase) QueryUser(uname string) UserModel {
	db := d.Db
	usr := UserModel{}
	db.QueryRow("SELECT Id,UserName,FirstName,LastName,Password FROM user WHERE UserName =?", uname).Scan(
		&usr.Id, &usr.UserName, &usr.FirstName, &usr.LastName, &usr.Password,
	)
	return usr
}

func (u *Dbase) Register(uname, fname, lname, pwd string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	if len(hashedPassword) == 0 {
		return errors.New("hashed password empty")
	}

	regQuery, err := u.Db.Prepare("INSERT INTO user(UserName,FirstName,LastName,Password) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = regQuery.Exec(uname, fname, lname, hashedPassword)
	return err
}
