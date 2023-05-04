package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Course struct {
	ID         int     `json: "id"`
	Name       string  `json: "name"`
	Price      float64 `json: "price"`
	Instructor string  `json: "instructor`
}

var courseList []Course

func init() {
	CourseJSON := `[
		{
			"id":1,
			"name":"Python",
			"price":2590,
			"instructor":"Ize Phanthakarn"
		},
		{
			"id":2,
			"name":"JavaScript",
			"price":1990,
			"instructor":"Ize Phanthakarn"
		},
		{
			"id":3,
			"name":"SQL",
			"price":590,
			"instructor":"Ize Phanthakarn"
		}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &courseList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	for _, course := range courseList {
		if highestID < course.ID {
			highestID = course.ID
		}
	}
	return highestID + 1
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	courseJson, err := json.Marshal(courseList)

	switch r.Method {
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJson)

	case http.MethodPost:
		var newCourse Course
		bodybytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodybytes, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newCourse.ID = getNextID()
		courseList = append(courseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func main() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":8080", nil)
}
