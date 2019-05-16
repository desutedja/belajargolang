package article

import (
	"database/sql"
	"fmt"
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
	AddDate     time.Time
	lok         string
}

//CreateArticle untuk membuat artikel baru
func (a *DbArticle) CreateArticle(title, description, addby string) int64 {
	QInsArticle, err := a.Db.Prepare("INSERT INTO article (Title,Description,Addby,AddDate) VALUES(?,?,?,?)")
	if err != nil {
		return -1
	}

	currentTime := time.Now()
	res, err := QInsArticle.Exec(title, description, addby, currentTime)
	if err != nil {
		fmt.Println("insert Error")
		return -1
	}

	id, err := res.LastInsertId()
	return id
}

//UploadFile (idarticle, filelocation)
func (a *DbArticle) UploadFile(idA int64, f string) error {
	_, err := a.Db.Exec("INSERT INTO image (idArticle,url) VALUES(?,?)", idA, f)
	fmt.Println("masuk upload")
	fmt.Println(f)
	fmt.Println(idA)
	return err
}

//GetArticles untuk ambil semua article dari DB
func (a *DbArticle) GetArticles() []Article {
	article := Article{}
	res := []Article{}

	selArticle, err := a.Db.Query(`SELECT t1.idArticle,Title,Description,AddBy,t2.Url
	FROM article t1
	INNER JOIN image t2 on t1.idArticle = t2.idArticle ORDER BY AddDate Desc LIMIT 10`)
	for selArticle.Next() {
		var idArticle int
		var title, description, addby, url string

		err = selArticle.Scan(&idArticle, &title, &description, &addby, &url)
		if err != nil {
			panic(err.Error())
		}
		article.idArticle = idArticle
		article.Title = title
		article.Description = description
		article.AddBy = addby
		article.lok = url
		res = append(res, article)
	}
	// a.Db.QueryRow("SELECT idArticle,Title,Description,AddBy FROM article").Scan(
	// 	&article.idArticle, &article.Title, &article.Description, &article.AddBy,

	// 	res = append(res, article)
	// )
	return res
}

//GetArticleByID untuk ambil article dari DB berdasarkan ID
func (a *DbArticle) GetArticleByID(idarticle int) Article {
	article := Article{}

	a.Db.QueryRow("SELECT idArticle,Title,Description,AddBy,AddDate FROM article WHERE idArticle =?", idarticle).Scan(
		&article.idArticle, &article.Title, &article.Description, &article.AddBy, &article.AddDate,
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
