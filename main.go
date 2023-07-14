package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Course represents a course entity
type Course struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Units []Unit `json:"units"`
}

// Unit represents a unit entity
type Unit struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Lessons []Lesson `json:"lessons"`
}

// Lesson represents a lesson entity
type Lesson struct {
	ID       int    `json:"id"`
	UnitID   int    `json:"unit_id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Configure MySQL connection
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go-webdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define the GET endpoint to retrieve all courses with units and lessons
	router.GET("/courses", func(c *gin.Context) {
		query := `
			SELECT
				courses.id,
				courses.name,
				units.id,
				units.name,
				lessons.id,
				lessons.unit_id,
				lessons.name,
				lessons.duration
			FROM
				courses
			INNER JOIN
				units ON courses.id = units.course_id
			INNER JOIN
				lessons ON units.id = lessons.unit_id
		`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		courses := make(map[int]*Course)

		for rows.Next() {
			var courseID int
			var courseTitle string
			var unitID int
			var unitTitle string
			var lessonID int
			var lessonUnitID int
			var lessonTitle string
			var lessonDuration int

			err := rows.Scan(
				&courseID,
				&courseTitle,
				&unitID,
				&unitTitle,
				&lessonID,
				&lessonUnitID,
				&lessonTitle,
				&lessonDuration,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if _, ok := courses[courseID]; !ok {
				courses[courseID] = &Course{
					ID:    courseID,
					Title: courseTitle,
				}
			}

			unit := Unit{
				ID:    unitID,
				Title: unitTitle,
			}

			lesson := Lesson{
				ID:       lessonID,
				UnitID:   lessonUnitID,
				Title:    lessonTitle,
				Duration: lessonDuration,
			}

			courses[courseID].Units = append(courses[courseID].Units, unit)
			courses[courseID].Units[unitID].Lessons = append(courses[courseID].Units[unitID].Lessons, lesson)
		}

		var result []Course
		for _, course := range courses {
			result = append(result, *course)
		}

		// Return the courses with units and lessons as a JSON response
		c.JSON(http.StatusOK, result)
	})

	// Start the server
	router.Run(":8082")
}
