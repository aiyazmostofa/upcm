package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type KeyValue struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

var dry time.Time
var wet time.Time
var end time.Time
var title string

const TIME_FORMAT = "01/02/2006|15:04:05|MST"

func initGlobal() error {
	temp := KeyValue{}
	res := d.First(&temp, "key = ?", "dry")
	if res.RowsAffected < 1 {
		return fmt.Errorf("failed to fetch dry run start time\n")
	}
	dry, _ = time.Parse(TIME_FORMAT, temp.Value)

	temp = KeyValue{}
	res = d.First(&temp, "key = ?", "wet")
	if res.RowsAffected < 1 {
		return fmt.Errorf("failed to fetch contest start time\n")
	}
	wet, _ = time.Parse(TIME_FORMAT, temp.Value)

	temp = KeyValue{}
	res = d.First(&temp, "key = ?", "end")
	if res.RowsAffected < 1 {
		return fmt.Errorf("failed to fetch contest end time\n")
	}
	end, _ = time.Parse(TIME_FORMAT, temp.Value)

	temp = KeyValue{}
	res = d.First(&temp, "key = ?", "title")
	if res.RowsAffected < 1 {
		return fmt.Errorf("failed to fetch title\n")
	}
	title = temp.Value
	return nil
}

func miscRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/title", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, struct {
			Title string `json:"title"`
		}{title})
	})

	r.Get("/times", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, struct {
			Dry time.Time `json:"dryRunStartTime"`
			Wet time.Time `json:"contestStartTime"`
			End time.Time `json:"contestEndTime"`
		}{dry, wet, end})
	})

	r.Put("/times", updateTime)
	return r
}

func updateTime(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("user").(User).AuthLevel != "Admin" {
		Unauthorized(w, "only admin can update contest times")
		return
	}
	query := r.URL.Query()

	reqDry := query.Get("dryRunStartTime")
	if len(reqDry) > 0 {
		sampleTime, err := time.Parse(TIME_FORMAT, reqDry)
		if err != nil {
			BadRequest(w, "invalid format for dry run start time")
			return
		}
		dry = sampleTime
		d.Save(&KeyValue{"dry", reqDry})
	}

	reqWet := query.Get("contestStartTime")
	if len(reqWet) > 0 {
		sampleTime, err := time.Parse(TIME_FORMAT, reqWet)
		if err != nil {
			BadRequest(w, "invalid format for contest start time")
			return
		}
		wet = sampleTime
		d.Save(&KeyValue{"wet", reqWet})
	}

	reqEnd := query.Get("contestEndTime")
	if len(reqEnd) > 0 {
		sampleTime, err := time.Parse(TIME_FORMAT, reqEnd)
		if err != nil {
			BadRequest(w, "invalid format for contest end time")
			return
		}
		end = sampleTime
		d.Save(&KeyValue{"end", reqEnd})
	}

	render.JSON(w, r, struct {
		Dry time.Time `json:"dryRunStartTime"`
		Wet time.Time `json:"contestStartTime"`
		End time.Time `json:"contestEndTime"`
	}{dry, wet, end})
}
