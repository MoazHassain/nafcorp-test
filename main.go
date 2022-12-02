package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/signup", signUp)
	// http.HandleFunc("/signup-handler", signUpHandler)
	http.HandleFunc("/login", logIn)
	// http.HandleFunc("/login-handler", logInHandler)

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./assets"))))

	fmt.Println("starting server on :8081")
	http.ListenAndServe(":8081", nil)
}

func home(res http.ResponseWriter, req *http.Request) {

	ptmp, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	ptmp.Execute(res, nil)

}

func signUp(res http.ResponseWriter, req *http.Request) {
	ptmp, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("signup.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, nil)

}

func logIn(res http.ResponseWriter, req *http.Request) {
	ptmp, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp, err = ptmp.ParseFiles("signin.html")
	if err != nil {
		fmt.Println(err.Error())

	}

	ptmp.Execute(res, nil)
}
