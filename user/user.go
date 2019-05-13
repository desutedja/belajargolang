package user

import(
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type User struct{
	Db *sql.DB
}

func (u *User) Register(uname, fname, lname, pwd string) error{
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	if len(hashedPassword) == 0{
		return errors.New("hashed password empty")
	}

	regQuery, err := u.Db.Prepare("INSERT INTO user(UserName,FirstName,LastName,Password) VALUES(?,?,?,?)")
	if err != nil{
		return err
	}

	_, err = regQuery.Exec(uname, fname, lname, hashedPassword)
	return err
}