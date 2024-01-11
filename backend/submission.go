package main

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type Submission struct {
	ID        uint      `json:"ID"`
	UserID    uint      `json:"userID" gorm:"index"`
	ProblemID uint      `json:"problemID" gorm:"index"`
	Content   string    `json:"content"`
	FileName  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	Verdict   string    `json:"verdict"`
}

type ESubmission struct {
	ID        uint      `json:"ID"`
	UserID    uint      `json:"userID"`
	ProblemID uint      `json:"problemID"`
	FileName  string    `json:"filename"`
	Timestamp time.Time `json:"timestamp"`
	Verdict   string    `json:"verdict"`
}

func submissionRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getSubmissions)
	r.Get("/best", getSubmissionsBest)
	r.Get("/{submissionID:[0-9]+}", getSubmission)
	r.Post("/", createSubmission)
	return r
}

func getSubmissions(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	var submissions []Submission
	if user.AuthLevel == "Team" {
		d.Find(&submissions, "user_id = ?", user.ID)
	} else {
		d.Find(&submissions)
	}

	e := make([]ESubmission, len(submissions))
	for i, v := range submissions {
		e[i] = eSubmission(v)
	}

	render.JSON(w, r, e)
}

func getSubmissionsBest(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	var submissions []struct {
		ID        uint      `json:"ID"`
		UserID    uint      `json:"userID"`
		Username  string    `json:"username"`
		ProblemID uint      `json:"problemID"`
		FileName  string    `json:"file_name"`
		Title     string    `json:"title"`
		Timestamp time.Time `json:"timestamp"`
		Verdict   string    `json:"verdict"`
	}

	var res *gorm.DB
	if user.AuthLevel == "Team" {
		res = d.Raw("SELECT submissions.id,user_id,username,title,problem_id,timestamp,verdict,file_name FROM submissions INNER JOIN users ON users.id = submissions.user_id INNER JOIN problems ON problems.id = submissions.problem_id WHERE users.id = ?", user.ID).Scan(&submissions)
	} else {
		res = d.Raw("SELECT submissions.id,user_id,username,title,problem_id,timestamp,verdict,file_name FROM submissions INNER JOIN users ON users.id = submissions.user_id INNER JOIN problems ON problems.id = submissions.problem_id").Scan(&submissions)
	}

	if res.Error != nil {
		InternalServerError(w, res.Error.Error())
		return
	}

	if res.RowsAffected == 0 {
		submissions = make([]struct {
			ID        uint      `json:"ID"`
			UserID    uint      `json:"userID"`
			Username  string    `json:"username"`
			ProblemID uint      `json:"problemID"`
			FileName  string    `json:"file_name"`
			Title     string    `json:"title"`
			Timestamp time.Time `json:"timestamp"`
			Verdict   string    `json:"verdict"`
		}, 0)
	}

	render.JSON(w, r, submissions)
}

func getSubmission(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	submissionID := uintID(chi.URLParam(r, "submissionID"))

	var submission Submission
	res := d.First(&submission, submissionID)
	if res.RowsAffected < 1 {
		NotFound(w, "submission not found")
		return
	}

	if user.AuthLevel == "Team" {
		if user.ID != submission.UserID {
			Unauthorized(w, "team can't access other team's submissions")
			return
		}
	}

	render.JSON(w, r, submission)
}

func createSubmission(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now()
	user := r.Context().Value("user").(User)
	if user.AuthLevel == "Admin" {
		Unauthorized(w, "admins can't create submissions")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 200<<10)
	if err := r.ParseMultipartForm(200 << 10); err != nil {
		BadRequest(w, "request too big or request not proper multipart form")
		return
	}

	content, fileName, err := readMultipartFile(r, "content")
	if err != nil {
		BadRequest(w, "content file field invalid")
		return
	}

	temp, err := strconv.ParseUint(r.FormValue("problemID"), 10, 32)
	if err != nil {
		BadRequest(w, "problemID not valid positive integer")
		return
	}
	problemID := uint(temp)

	res := d.First(&Problem{}, problemID)
	if res.RowsAffected < 1 {
		NotFound(w, "problem not found")
		return
	}

	if problemID == 1 {
		if timestamp.Before(dry) || timestamp.After(end) {
			BadRequest(w, "submissions for this problem are currently not accepted")
			return
		}
	} else {
		if timestamp.Before(wet) || timestamp.After(end) {
			BadRequest(w, "submissions for this problem are currently not accepted")
			return
		}
	}

	submission := Submission{
		Content:   content,
		FileName:  fileName,
		ProblemID: problemID,
		UserID:    user.ID,
		Timestamp: timestamp,
		Verdict:   "TBD",
	}

	res = d.Omit("ID").Create(&submission)
	if res.Error != nil {
		InternalServerError(w, "failed to create submission")
		return
	}

	go start(submission.ID)
	render.JSON(w, r, eSubmission(submission))
}

func readMultipartFile(r *http.Request, name string) (string, string, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return "", "", err
	}

	defer file.Close()
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return "", "", err
	}

	return string(buffer.Bytes()), handler.Filename, nil
}

func eSubmission(submission Submission) ESubmission {
	return ESubmission{
		ID:        submission.ID,
		UserID:    submission.UserID,
		ProblemID: submission.ProblemID,
		FileName:  submission.FileName,
		Timestamp: submission.Timestamp,
		Verdict:   submission.Verdict,
	}
}
