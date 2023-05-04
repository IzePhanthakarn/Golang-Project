package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
// var courseList []Course

const coursePath = "courses"
const apiBasePath = "/api"

type Course struct {
	CourseID   int     `json:"course_id"`
	CourseName string  `json:"course_name"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/coursedb")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func getCourseList() ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := Db.QueryContext(ctx, `SELECT course_id, course_name, price, image_url FROM courseonline`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	courses := make([]Course, 0)
	for results.Next() {
		var course Course
		results.Scan(&course.CourseID,
			&course.CourseName,
			&course.Price,
			&course.ImageURL)

		courses = append(courses, course)
	}
	return courses, nil
}

func insertProduct(course Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.ExecContext(ctx, `INSERT INTO courseonline (course_id, course_name, price, image_url) VALUE (?, ?, ?, ?)`,
		course.CourseID, course.CourseName, course.Price, course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func getCourse(course_id int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := Db.QueryRowContext(ctx, `SELECT course_id, course_name, price, image_url FROM courseonline WHERE course_id = ?`,
		course_id)
	course := &Course{}
	err := row.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

func removeCourse(course_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := Db.ExecContext(ctx, `DELETE FROM courseonline WHERE course_id = ?`, course_id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		CourseID, err := insertProduct(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"course_id}: %d}`, CourseID)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Lengh")
		handler.ServeHTTP(w, r)
	})
}

func setupRoutes(apiBasePath string) {
	courseHandler := http.HandlerFunc(handleCourse)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
	coursesHandler := http.HandlerFunc(handleCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))

}

func handleCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		err := removeCourse(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	SetupDB()
	setupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
