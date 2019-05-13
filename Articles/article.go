package article

import (
	"fmt"
	"database/sql"
)

type Article struct{
	Db *sql.DB
}

func (a *Article)CreateArticle(title, description string) error {
	fmt.Printf("%+v", a.Db)
	return nil
}