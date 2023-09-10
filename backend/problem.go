package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Problem struct {
	ID     uint   `json:"ID"`
	Title  string `json:"title"`
	Input  string `json:"input"`
	Output string `json:"output"`
	InputName  string `json:"inputName"`
}

type EProblem struct {
	ID    uint   `json:"ID"`
	Title string `json:"title"`
	InputName string `json:"inputName"`
}

func problemRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProblems)
	r.Get("/{problemID:[0-9]+}", getProblem)
	return r
}

func getProblems(w http.ResponseWriter, r *http.Request) {
	var problems []Problem
	res := d.Find(&problems)
	if res.Error != nil {
		InternalServerError(w, "failed get list of problems")
		return
	}

	e := make([]EProblem, len(problems))
	for i, v := range problems {
		e[i] = eProblem(v)
	}
	render.JSON(w, r, e)
}

func getProblem(w http.ResponseWriter, r *http.Request) {
	problemID := uintID(chi.URLParam(r, "problemID"))
	var problem Problem
	res := d.First(&problem, problemID)
	if res.RowsAffected < 1 {
		NotFound(w, "problem not found")
		return
	}

	if r.Context().Value("user").(User).AuthLevel == "Team" {
		render.JSON(w, r, eProblem(problem))
		return
	}
	render.JSON(w, r, problem)
}

func eProblem(problem Problem) EProblem {
	return EProblem{
		ID:    problem.ID,
		Title: problem.Title,
		InputName: problem.InputName,
	}
}
