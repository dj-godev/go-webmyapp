package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("./static/*.html")
	// MySQL connection parameters
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go-webdb")
	if err != nil {
		log.Fatal(err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	type User struct {
		Id         int
		FullName   string
		RollNumber string
	}

	type usersFullName struct {
		FullName string
	}

	// Define a route to fetch data from the database
	router.GET("/data", func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		var (
			user  User
			users []User
		)

		var (
			fullName       usersFullName
			usersFullNames []usersFullName
		)

		query := fmt.Sprintf("select id, fullName, fullName , roll_number from users LIMIT %d OFFSET %d", limit, offset)
		rows, err := db.Query(query)
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&user.Id, &user.FullName, &fullName.FullName, &user.RollNumber)
			users = append(users, user)
			usersFullNames = append(usersFullNames, fullName)
			if err != nil {
				fmt.Print(err.Error())
			}
		}

		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"page":   page,
			"limit":  limit,
			"result": users,
			"Names":  usersFullNames,
			"count":  len(users),
		})
	})

	// Define the GET endpoint to retrieve all courses with units
	router.GET("/courses", func(c *gin.Context) {
		// Fetch all courses
		courses, err := getAllCourses(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Fetch units for each course
		for i := range courses {
			units, err := getUnitsByCourseID(db, courses[i].Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			courses[i].Units = units
		}

		// Return the courses with lessons as a JSON response
		//c.JSON(http.StatusOK, courses)
		c.IndentedJSON(200, courses)
	})

	router.GET("/", func(c *gin.Context) {
		// Serve static files from the "static" directory
		//router.Static("/home", "./static")
		//c.String(http.StatusOK, "Hello from %v", "Gin")
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Start the server
	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

// Course represents a course entity
type Course struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Units []Unit `json:"units"`
}

// Unit represents a unit entity
type Unit struct {
	Id       int    `json:"id"`
	CourseId int    `json:"course_id"`
	Name     string `json:"name"`
}

// getAllCourses retrieves all courses from the database
func getAllCourses(db *sql.DB) ([]Course, error) {
	query := "SELECT id, name FROM courses"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.Id, &course.Name)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

// getUnitsByCourseID retrieves all units for a given course ID from the database
func getUnitsByCourseID(db *sql.DB, CourseId int) ([]Unit, error) {

	query := "SELECT id, name, course_id FROM units WHERE course_id = ?"
	rows, err := db.Query(query, CourseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit
	for rows.Next() {
		var unit Unit
		err := rows.Scan(&unit.Id, &unit.Name, &unit.CourseId)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}
