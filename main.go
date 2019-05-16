
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
	"time"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/kataras/go-sessions"
	"github.com/blog_denny/article"
	"github.com/blog_denny/user"
)

type userModel struct {
	id        int
	UserName  string
	FirstName string
	LastName  string
	Password  string
}

type modelarticle struct {
	idArticle   int
	Title       string
	Description string
	AddBy       string
	AddDate     time.Time
}

//DbConn nantinya mungkin akan berada diluar main.go
func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "myblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		fmt.Println("error db")
		panic(err.Error())
	}
	return db
}



func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}

	return true
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "view/register.html")
		return
	}

	db := DbConn()
	uname := r.FormValue("username")
	fname := r.FormValue("firstname")
	lname := r.FormValue("lastname")
	pwd := r.FormValue("password")

	user := user.Dbase{
		Db : db,
	}

	users := user.QueryUser(uname)

	fmt.Println(userModel{}.UserName)
	fmt.Println(users.UserName)

	if users.UserName == uname {		
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	err := user.Register(uname, fname, lname, pwd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 302)
	defer db.Close()
}

func login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)

	if r.Method != "POST" {
		http.ServeFile(w, r, "view/login.html")
		return
	}

	if len(session.GetString("username")) != 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	db := DbConn()

	user := user.Dbase{
		Db : db,
	}

	uname := r.FormValue("username")
	pwd := r.FormValue("password")
	users := user.QueryUser(uname)

	pwdCompare := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(pwd))

	if pwdCompare == nil {
		//success
		session := sessions.Start(w, r)
		session.Set("username", users.UserName)
		session.Set("name", users.FirstName)
		http.Redirect(w, r, "/dashboard", 302)
	} else {
		//fail
		http.Redirect(w, r, "/login", 302)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/login", 302)
		return
	}

	homeDb := article.DbArticle{
		Db : DbConn(),
	}

	allArticle := homeDb.GetArticles()

	fmt.Printf("%+v",allArticle)

	// data := map[string] interface{}{
	// 	"idArticle":"allArticle.idArticle",
	// 	"username": session.GetString("username"),
	// 	"addby": allArticle.AddBy,
	// 	"title":allArticle.Title,
	// 	"description":allArticle.Description,
	// 	"adddate":"today",
	// }

	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, allArticle)
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}

func dashboard(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.ServeFile(w, r, "view/article.html")
		return
	}

	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/login", 302)
		return
	}

	db := DbConn()
	title := r.FormValue("title")
	description := r.FormValue("description")

	article := article.DbArticle{
		Db : db,
	}

	newid := article.CreateArticle(title, description, session.GetString("username"))
	if newid < 1 {
		http.Error(w, "Terjadi Kesalahan Pada Server", http.StatusInternalServerError)
		return
	}
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	filelok := "/test/"+handler.Filename
	errs := article.UploadFile(newid,filelok)
	if errs != nil {
		return
	}
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
	defer db.Close()	
}

var r = mux.NewRouter()

func routes() {
	//routesnya
	//r := mux.NewRouter()
	r.HandleFunc("/register", register)
	r.HandleFunc("/login", login)
	r.HandleFunc("/", home)
	r.HandleFunc("/logout", logout)

	r.HandleFunc("/dashboard", dashboard)

	fmt.Println("Server running on port :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func main() {
	//dbConn := DbConn()
	routes()

	//createArticle(dbConn)
}