package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"unicode"
	// "path/filepath"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var err error

type User struct {
	UserID    string
	UserName  string
	UserEmail string
	UserPhone string
}

var users User

var store = sessions.NewCookieStore([]byte("SESSION_KEY"))

func main() {

	// db, err = sql.Open("{sql-type}", "{username}:{password}@tcp({server:port})/{database-name}")
	db, err = sql.Open("mysql", "root:moaz@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("connected to database")

	http.HandleFunc("/", home)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/signup-handler", signUpHandler)
	http.HandleFunc("/login", logIn)
	http.HandleFunc("/login-handler", logInHandler)
	http.HandleFunc("/account", account)

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./assets"))))

	fmt.Println("starting server on :8081")
	http.ListenAndServe(":8081", nil)
	// http.ListenAndServe(":8881", context.ClearHandler(http.DefaultServeMux))
}

func home(res http.ResponseWriter, req *http.Request) {

	ptmp, err := template.ParseFiles("wpage/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	ptmp.Execute(res, nil)

}

func signUp(res http.ResponseWriter, req *http.Request) {
	ptmp, err := template.ParseFiles("wpage/index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("wpage/signup.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, nil)

}

func signUpHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("*****registerAuthHandler running*****")
	req.ParseForm()
	username := req.FormValue("name-input")
	email := req.FormValue("email-input")
	password := req.FormValue("pass-input")
	phone := req.FormValue("mobile-input")

	// check username for only alphaNumeric characters
	var nameAlphaNumeric = true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		// if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
		if !unicode.IsLetter(char) {
			nameAlphaNumeric = false
		}
	}

	// check username pswdLength
	var nameLength bool
	if 4 <= len(username) && len(username) <= 50 {
		nameLength = true
	}

	// check password criteria

	fmt.Println("password:", password, "\npswdLength:", len(password))
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true

	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 5 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
	if /* !pswdLowercase || !pswdUppercase || */ !pswdNumber || /* !pswdSpecial || */ !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/signup.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, "please check username and password criteria")

		// ptmp.ExecuteTemplate(res, "wpage/signup.gohtml", "please check username and password criteria")
		return
	}
	// check if username already exists for availability
	stmt := "SELECT u_name FROM users WHERE u_name=?"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/signup.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, "username already taken")

		// tpl.ExecuteTemplate(res, "register.html", "username already taken")
		return
	}

	// check if email already used
	emailStmt := "SELECT u_email FROM users WHERE u_email=?"
	emailRow := db.QueryRow(emailStmt, email)
	var eID string
	emailErr := emailRow.Scan(&eID)
	if emailErr != sql.ErrNoRows {
		fmt.Println("email already exists, err:", err)
		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/signup.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, "Email Already used")

		// tpl.ExecuteTemplate(res, "register.html", "username already taken")
		return
	}

	// create hash from password
	// var hash []byte
	// // func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	// hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	fmt.Println("bcrypt err:", err)
	// 	ptmp, err := template.ParseFiles("template/home.gohtml")
	// 	if err != nil {
	// 		fmt.Println(err.Error())

	// 	}

	// 	ptmp, err = ptmp.ParseFiles("wpage/signup.gohtml")
	// 	if err != nil {
	// 		fmt.Println(err.Error())

	// 	}

	// 	ptmp.Execute(res, "there was a problem registering account")

	// 	// tpl.ExecuteTemplate(res, "register.html", "there was a problem registering account")
	// 	return
	// }
	// fmt.Println("hash:", hash)
	// fmt.Println("string(hash):", string(hash))
	// func (db *DB) Prepare(query string) (*Stmt, error)

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users(u_name, u_email, u_password, u_phone) VALUES(?, ?, ?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/signup.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, "There was a problem registering account")

		// tpl.ExecuteTemplate(res, "register.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()
	var result sql.Result
	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err = insertStmt.Exec(username, email, password, phone)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/signup.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, "There was a problem registering account")

		// tpl.ExecuteTemplate(res, "register.html", "there was a problem registering account")
		return
	}

	var userID string
	accountStmt := "SELECT u_id FROM users WHERE u_email=?"
	accountRow := db.QueryRow(accountStmt, email)
	err = accountRow.Scan(&userID)
	// Get always returns a session, even if empty
	// returns error if exists and could not be decoded
	// Get(r *http.Request, name string) (*Session, error)
	session, _ := store.Get(req, "session")
	// session struct has field Values map[interface{}]interface{}
	session.Values["userID"] = userID

	// save before writing to response/return from handler
	session.Save(req, res)

	// users

	sessionUserID := fmt.Sprintln(session.Values["userID"])
	fmt.Println("sessionUserID:", sessionUserID)
	var sessionUserName, sessionUserEmail, sessionUserPhone string
	db.QueryRow("SELECT u_name, u_email, u_phone FROM users WHERE u_id=?", sessionUserID).Scan(&sessionUserName, &sessionUserEmail, &sessionUserPhone)
	fmt.Println("session user info:" + sessionUserName + " " + sessionUserEmail)

	users = User{
		UserID:    sessionUserID,
		UserName:  sessionUserName,
		UserEmail: sessionUserEmail,
		UserPhone: sessionUserPhone,
	}

	ptmp, err := template.ParseFiles("wpage/index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("wpage/account.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, users)
	// http.Redirect(res, req, "/dashboard", 301)
	// fmt.Fprint(res, "congrats, your account has been successfully created")
}

func logIn(res http.ResponseWriter, req *http.Request) {
	ptmp, err := template.ParseFiles("wpage/index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("wpage/signin.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, nil)
}

func logInHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	req.ParseForm()
	email := req.FormValue("email-input")
	password := req.FormValue("pass-input")
	fmt.Println("email:", email, "password:", password)

	// retrieve userid, password from db
	var userID, dbPass string
	stmt := "SELECT u_id, u_password FROM users WHERE u_email=?"
	row := db.QueryRow(stmt, email)
	err := row.Scan(&userID, &dbPass)
	// passRow := db.QueryRow("SELECT u_password FROM users WHERE u_email=?", username)
	// scanPass := passRow.Scan(&dbPass)

	fmt.Println("password from db:" + dbPass + "uid from db:" + userID)
	// if dbPass == password {
	// 	err = nil
	// }

	if err != nil {
		fmt.Println("error selecting password in db by email")

		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/signin.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, "Check username and password")

		// tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	// compare password
	if dbPass == password {
		// Get always returns a session, even if empty
		// returns error if exists and could not be decoded
		// Get(r *http.Request, name string) (*Session, error)
		session, _ := store.Get(req, "session")
		// session struct has field Values map[interface{}]interface{}
		session.Values["userID"] = userID

		// save before writing to response/return from handler
		session.Save(req, res)

		// users

		sessionUserID := fmt.Sprintln(session.Values["userID"])
		fmt.Println("sessionUserID:", sessionUserID)
		var sessionUserName, sessionUserEmail, sessionUserPhone string
		db.QueryRow("SELECT u_name, u_email, u_phone FROM users WHERE u_id=?", sessionUserID).Scan(&sessionUserName, &sessionUserEmail, &sessionUserPhone)
		fmt.Println("session user info:" + sessionUserName + " " + sessionUserEmail)

		users = User{
			UserID:    sessionUserID,
			UserName:  sessionUserName,
			UserEmail: sessionUserEmail,
			UserPhone: sessionUserPhone,
		}

		ptmp, err := template.ParseFiles("wpage/index.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp, err = ptmp.ParseFiles("wpage/account.html")
		if err != nil {
			fmt.Println(err.Error())

		}

		ptmp.Execute(res, users)

		// fmt.Fprint(w, "You have successfully logged in :)")
		return
	}
	// returns nill on succcess

	fmt.Println("incorrect password")
	ptmp, err := template.ParseFiles("wpage/index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("wpage/signin.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, "Check password")
	// tpl.ExecuteTemplate(w, "login.html", "check username and password")
}

func account(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	_, ok := session.Values["userID"]
	fmt.Println("ok:", ok)
	if !ok {
		http.Redirect(res, req, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
	ptmp, err := template.ParseFiles("wpage/index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("wpage/account.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, users)
}
