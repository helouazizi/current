// internal/handlers/handlers.go
package handlers

import (
	"html/template"
	"net/http"
	"os"

	"forum/internal/database"
	"forum/internal/utils"
)

type Pages struct {
	All_Templates *template.Template
}
type Form struct {
	Title           string
	Button          string
	IsAuthenticated bool
}
type User struct {
	IsAuthenticated bool
}

var pages Pages

func init() {
	var err error
	path, err := utils.GetFolderPath("..", "templates")
	if err != nil {
		panic(err)
	}
	pages.All_Templates, err = template.ParseGlob(path + "/*.html")
	if err != nil {
		panic(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Page not found")
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	user := User{
		IsAuthenticated: false,
	}

	pages.All_Templates.ExecuteTemplate(w, "home.html", user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "method not allowed")
		return
	}
	user := User{
		IsAuthenticated: true,
	}

	pages.All_Templates.ExecuteTemplate(w, "home.html", user)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")

		password := r.FormValue("password")
		_, err := database.Database.Exec("INSERT INTO users (username, password) VALUES (?, ?)", email, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pages.All_Templates.ExecuteTemplate(w, "error.html", "Internal server error")
			return
		}
	}
	data := Form{
		Title:           "Create Account",
		Button:          "Create Account",
		IsAuthenticated: false,
	}

	pages.All_Templates.ExecuteTemplate(w, "login.html", data)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	user := User{
		IsAuthenticated: true,
	}
	pages.All_Templates.ExecuteTemplate(w, "createpost.html", user)
}

func Create_Account(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	user := User{
		IsAuthenticated: false,
	}
	pages.All_Templates.ExecuteTemplate(w, "home.html", user)
}

func Sign_In(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	data := Form{
		Title:           "Login",
		Button:          "Login",
		IsAuthenticated: false,
	}
	pages.All_Templates.ExecuteTemplate(w, "login.html", data)
}

func FilterPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	user := User{
		IsAuthenticated: true,
	}
	pages.All_Templates.ExecuteTemplate(w, "filter.html", user)
}

func MyPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	pages.All_Templates.ExecuteTemplate(w, "profile.html", "My Posts")
}

func LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	pages.All_Templates.ExecuteTemplate(w, "profile.html", "NO liked posts")
}

func CategorizePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	pages.All_Templates.ExecuteTemplate(w, "profile.html", "Category")
}

func Settings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	pages.All_Templates.ExecuteTemplate(w, "profile.html", "settings")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	pages.All_Templates.ExecuteTemplate(w, "profile.html", "Logout")
}

func Serve_Static(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Method not allowed")
		return
	}
	file, err := os.Stat(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Not Found")
		return
	}
	if file.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		pages.All_Templates.ExecuteTemplate(w, "error.html", "Not Found")
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
	/*path, _ := utils.GetFolderPath("..", "static")
	fs := http.FileServer(http.Dir(path))
	http.StripPrefix("/static/", fs).ServeHTTP(w, r)*/
}
