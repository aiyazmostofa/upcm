package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"embed"
)

var d *gorm.DB
var tokenMap map[string]uint
var idMap map[uint]string

//go:embed all:dist
var f embed.FS

func main() {
	if len(os.Args) > 1 {
		path := os.Args[1]
		generate(path)
		return
	}

	_, err := os.ReadFile("contest.db")
	if err != nil {
		fmt.Println("contest.db not found")
		return
	}
	d, err = gorm.Open(sqlite.Open("contest.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		return
	}
	err = d.AutoMigrate(&User{}, &Submission{}, &Problem{}, &KeyValue{})
	if err != nil {
		fmt.Println("failed to migrate database")
		return
	}

	err = initGlobal()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tokenMap = make(map[string]uint)
	idMap = make(map[uint]string)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Get("/api/token", getToken)
	r.Mount("/api", mainRouter())

	sub, err := fs.Sub(f, "dist")
	if err != nil {
		fmt.Println("couldn't mount embed fs")
		return
	}

	r.Handle("/*", http.StripPrefix("/", http.FileServer((http.FS(sub)))))
	if _, err := os.Stat("data"); !os.IsNotExist(err) {
		fs := http.FileServer(http.Dir("data"))
		fmt.Println("detected data folder")
		r.Handle("/data/*", http.StripPrefix("/data/", fs))
	}
	http.ListenAndServe(":3000", r)
}

func mainRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(Auth)
	r.Mount("/users", userRouter())
	r.Mount("/submissions", submissionRouter())
	r.Mount("/problems", problemRouter())
	r.Mount("/misc", miscRouter())
	return r
}

func getToken(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if len(username) == 0 || len(password) == 0 {
		BadRequest(w, "username or password can't be blank")
		return
	}

	var user User
	res := d.First(&user, "username = ?", username)
	if res.RowsAffected < 1 {
		Unauthorized(w, "username does not exist")
		return
	}

	if user.Password != password {
		Unauthorized(w, "password is incorrect")
		return
	}

	token, ok := idMap[user.ID]
	if !ok {
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			InternalServerError(w, "failed to generate token")
			return
		}

		token = hex.EncodeToString(bytes)
		tokenMap[token] = user.ID
		idMap[user.ID] = token
	}

	render.JSON(w, r, struct {
		Token     string `json:"token"`
		AuthLevel string `json:"authLevel"`
		ID        uint   `json:"ID"`
	}{
		token,
		user.AuthLevel,
		user.ID,
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		userID, ok := tokenMap[token]
		if !ok {
			Unauthorized(w, "token not valid")
			return
		}

		var user User
		res := d.First(&user, "id = ?", userID)
		if res.RowsAffected < 1 {
			Unauthorized(w, "user doesn't exist anymore")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Unauthorized(w http.ResponseWriter, s string) {
	http.Error(w, s, http.StatusUnauthorized)
}

func BadRequest(w http.ResponseWriter, s string) {
	http.Error(w, s, http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, s string) {
	http.Error(w, s, http.StatusInternalServerError)
}

func NotFound(w http.ResponseWriter, s string) {
	http.Error(w, s, http.StatusNotFound)
}

func uintID(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint(v)
}
