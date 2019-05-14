package article

import (
	"database/sql"
	"time"
)

//DbArticle struktur
type DbArticle struct {
	Db *sql.DB
}

//Article Model
type Article struct {
	idArticle   int
	Title       string
	Description string
	AddBy       string
}

//CreateArticle untuk membuat artikel baru
func (a *DbArticle) CreateArticle(title, description, addby string) error {
	QInsArticle, err := a.Db.Prepare("INSERT INTO article (Title,Description,Addby,AddDate) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	currentTime := time.Now()
	_, err = QInsArticle.Exec(title, description, addby, currentTime)
	return err
}

//GetArticles untuk ambil semua article dari DB
func (a *DbArticle) GetArticles() Article {
	article := Article{}

	a.Db.QueryRow("SELECT idArticle,Title,Description,AddBy FROM article").Scan(
		&article.idArticle, &article.Title, &article.Description, &article.AddBy,
	)
	return article
}

//GetArticleByID untuk ambil article dari DB berdasarkan ID
func (a *DbArticle) GetArticleByID(idarticle int) Article {
	article := Article{}

	a.Db.QueryRow("SELECT idArticle,Title,Description,AddBy FROM article WHERE idArticle =?", idarticle).Scan(
		&article.idArticle, &article.Title, &article.Description, &article.AddBy,
	)
	return article
}

//UpdateArticle untuk update article
func (a *DbArticle) UpdateArticle(structArticle Article) error {
	QUptArticle, err := a.Db.Prepare("UPDATE article SET Title=?, Description=? WHERE idArticle = ?")
	if err != nil {
		return err
	}

	_, err = QUptArticle.Exec(structArticle.Title, structArticle.Description, structArticle.idArticle)
	return err
}

//DeleteArticle untuk delete article
func (a *DbArticle) DeleteArticle(idarticle int) error {
	QUptArticle, err := a.Db.Prepare("DELETE article WHERE idArticle = ?")
	if err != nil {
		return err
	}

	_, err = QUptArticle.Exec(idarticle)
	return err
}
