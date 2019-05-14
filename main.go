
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

	//if users == (userModel{})  {
	if users.UserName != "" {
		err := user.Register(uname, fname, lname, pwd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	}

	db := DbConn()

	user := user.Dbase{
		Db : db,
	}

	uname := r.FormValue("username")
	pwd := r.FormValue("password")
	users := user.QueryUser(uname)
	//users := queryUser(uname)

	pwdCompare := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(pwd))

	if pwdCompare == nil {
		//success
		session := sessions.Start(w, r)
		session.Set("username", users.UserName)
		session.Set("name", users.FirstName)
		http.Redirect(w, r, "/", 302)
	} else {
		//fail
		http.Redirect(w, r, "/login", 302)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/login", 302)
	}

	data := map[string]string{
		"username": session.GetString("username"),
		"message":  "Welcome on Go !",
	}

	t, err := template.ParseFiles("views/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, data)
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}

/*func routes() {
	//routesnya
	r.HandleFunc("/register", register).Methods("POST")
	//r.HandleFunc("/login", login).Methods("POST")
	fmt.Println("Routes Running")
}*/

func main() {
	dbConn := DbConn()

	r := mux.NewRouter()
	r.HandleFunc("/register", register)
	r.HandleFunc("/login", login)
	r.HandleFunc("/", home)

	createArticle(dbConn)

	fmt.Println("Server running on port :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func createArticle(db *sql.DB){
	article := article.DbArticle{
		Db : db,
	}

	if err := article.CreateArticle("test","test","test"); err != nil{
		log.Fatal(err)
	}
}