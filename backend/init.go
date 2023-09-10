package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Contest struct {
	Title    string `json:"title"`
	Problems []struct {
		Title      string `json:"title"`
		InputName  string `json:"inputName"`
		OutputName string `json:"outputName"`
	} `json:"problems"`
}

func generate(path string) {
	syscall.Chdir(path)
	bytes, err := os.ReadFile("info.json")
	if err != nil {
		fmt.Println("folder needs to contain a info.json")
		return
	}

	var contest Contest
	err = json.Unmarshal(bytes, &contest)
	if err != nil {
		fmt.Println("info.json not proper format")
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

	res := d.Create(&User{
		Username:  "admin",
		Password:  "password",
		AuthLevel: "Admin",
	})

	if res.RowsAffected < 1 {
		fmt.Println("failed to create admin user")
		return
	}

	res = d.Create(&KeyValue{"title", contest.Title})
	if res.RowsAffected < 1 {
		fmt.Println("failed to add title")
		return
	}

	res = d.Create(&KeyValue{"dry", "01/01/1990|00:00:00|CST"})
	if res.RowsAffected < 1 {
		fmt.Println("failed to add dry run start time")
		return
	}

	res = d.Create(&KeyValue{"wet", "01/01/1990|00:00:00|CST"})
	if res.RowsAffected < 1 {
		fmt.Println("failed to add contest start time")
		return
	}

	res = res.Create(&KeyValue{"end", "01/01/1990|00:00:00|CST"})
	if res.RowsAffected < 1 {
		fmt.Println("failed to add contest end time")
		return
	}

	for i, v := range contest.Problems {
		problem := Problem{}
		if len(v.Title) == 0 {
			fmt.Println("title for problem " + fmt.Sprint(i) + " does not exist")
		}
		problem.Title = v.Title

		if len(v.InputName) > 0 {
			bytes, err = os.ReadFile(v.InputName)
			if err != nil {
				fmt.Println("input for problem " + fmt.Sprint(i) + " does not exist")
				return
			}
			problem.Input = string(bytes)
			problem.InputName = v.InputName
		}

		bytes, err = os.ReadFile(v.OutputName)
		if err != nil {

			fmt.Println("output for problem " + fmt.Sprint(i) + " does not exist")
			return
		}
		problem.Output = string(bytes)
		problem.ID = uint(i)

		res := d.Create(&problem)
		if res.RowsAffected < 1 {
			fmt.Println("failed to create problem " + fmt.Sprint(i))
			return
		}
	}

	fmt.Println("Admin Username: admin\nAdmin Password: password\nMake sure to change the password!")
}
