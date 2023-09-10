package main

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type User struct {
	ID        uint   `json:"ID"`
	Username  string `json:"username" gorm:"uniqueIndex"`
	Password  string `json:"password"`
	AuthLevel string `json:"authLevel"`
}

var usernameRegex, _ = regexp.Compile("[A-Za-z0-9-_]+")

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getUsers)
	r.Get("/{userID:[0-9]+}", getUser)
	r.Delete("/{userID:[0-9]+}", deleteUser)
	r.Put("/{userID:[0-9]+}", updateUser)
	r.Post("/", createUser)
	return r
}

func getUser(w http.ResponseWriter, r *http.Request) {
	signedIn := r.Context().Value("user").(User)
	userID := uintID(chi.URLParam(r, "userID"))

	if signedIn.AuthLevel == "Team" {
		if signedIn.ID != userID {
			Unauthorized(w, "team can't access other team")
			return
		}

		render.JSON(w, r, signedIn)
		return
	}

	var user User
	res := d.First(&user, userID)
	if res.RowsAffected < 1 {
		NotFound(w, "team not found")
		return
	}

	render.JSON(w, r, user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	if user.AuthLevel == "Team" {
		Unauthorized(w, "team can't view all teams")
		return
	}

	var users []User
	d.Find(&users)

	render.JSON(w, r, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user").(User).AuthLevel != "Admin" {
		Unauthorized(w, "only admin can create users")
		return
	}
	query := r.URL.Query()
	username := query.Get("username")
	if !usernameRegex.MatchString(username) {
		BadRequest(w, "username must only contain letters, numbers, dashes, or underscores")
		return
	}

	res := d.First(&User{}, "Username = ?", username)
	if res.RowsAffected >= 1 {
		BadRequest(w, "username already taken")
		return
	}

	password := query.Get("password")
	authLevel := query.Get("authLevel")
	if len(password) == 0 ||
		len(authLevel) == 0 {
		BadRequest(w, "password or authLevel can't be blank")
		return
	}

	if authLevel == "Admin" {
		Unauthorized(w, "only one admin can exist")
		return
	}

	user := User{
		Username:  username,
		Password:  password,
		AuthLevel: authLevel,
	}

	res = d.Omit("ID").Create(&user)
	if res.RowsAffected < 1 {
		InternalServerError(w, "failed to create user")
		return
	}

	render.JSON(w, r, user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user").(User).AuthLevel != "Admin" {
		Unauthorized(w, "only admin can delete users")
		return
	}
	userID := uintID(chi.URLParam(r, "userID"))

	var user User
	res := d.First(&user, userID)
	if res.RowsAffected < 1 {
		NotFound(w, "user not found")
		return
	}

	if user.AuthLevel == "Admin" {
		BadRequest(w, "you can't delete yourself")
		return
	}

	d.Delete(&user)
	d.Where("user_id = ?", user.ID).Delete(&Submission{})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user").(User).AuthLevel != "Admin" {
		Unauthorized(w, "only admin can update users")
		return
	}

	userID := uintID(chi.URLParam(r, "userID"))

	var user User
	res := d.First(&user, userID)
	if res.RowsAffected < 1 {
		NotFound(w, "user not found")
		return
	}
	query := r.URL.Query()

	username := query.Get("username")
	if len(username) > 0 {
		if !usernameRegex.MatchString(username) {
			BadRequest(w, "username must only contain letters, numbers, dashes, or underscores")
			return
		}
		user.Username = username
	}

	var temp User;
	res = d.First(&temp, "Username = ?", username)
	if res.RowsAffected >= 1 && temp.ID != user.ID {
		BadRequest(w, "username already taken")
		return
	}

	password := query.Get("password")
	if len(password) > 0 {
		user.Password = password
	}

	res = d.Save(&user)
	if res.Error != nil {
		InternalServerError(w, "failed to update user")
		return
	}

	render.JSON(w, r, user)
}
