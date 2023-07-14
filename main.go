package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// CourseResponse represents a course response with associated units and lessons
type CourseResponse struct {
	ID    int             `json:"id"`
	Title string          `json:"title"`
	Units []*UnitResponse `json:"units"`
}

// UnitResponse represents a unit response with associated lessons
type UnitResponse struct {
	ID      *int      `json:"id,omitempty"`
	Title   *string   `json:"title,omitempty"`
	Lessons []*Lesson `json:"lessons"`
}

// Lesson represents a lesson entity
type Lesson struct {
	ID         *int    `json:"id,omitempty"`
	UnitID     *int    `json:"unit_id,omitempty"`
	Title      *string `json:"title,omitempty"`
	Duration   *int    `json:"duration,omitempty"`
	Difficulty *string `json:"difficulty,omitempty"`
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
				lessons.name
			FROM
				courses
			LEFT JOIN
				units ON courses.id = units.course_id
			LEFT JOIN
				lessons ON units.id = lessons.unit_id order by courses.id asc
		`

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		courses := make(map[int]*CourseResponse)

		for rows.Next() {
			var courseID int
			var courseTitle string
			var unitID sql.NullInt64
			var unitTitle sql.NullString
			var lessonID sql.NullInt64
			var lessonUnitID sql.NullInt64
			var lessonTitle sql.NullString

			err := rows.Scan(
				&courseID,
				&courseTitle,
				&unitID,
				&unitTitle,
				&lessonID,
				&lessonUnitID,
				&lessonTitle,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if _, ok := courses[courseID]; !ok {
				courses[courseID] = &CourseResponse{
					ID:    courseID,
					Title: courseTitle,
					Units: []*UnitResponse{},
				}
			}

			unit := getUnitByID(courses[courseID].Units, int(unitID.Int64))
			if unit == nil {
				unit = &UnitResponse{
					ID:      getIntPointer(unitID),
					Title:   getStringPointer(unitTitle),
					Lessons: []*Lesson{},
				}
				courses[courseID].Units = append(courses[courseID].Units, unit)
			}

			lesson := &Lesson{
				ID:     getIntPointer(lessonID),
				UnitID: getIntPointer(lessonUnitID),
				Title:  getStringPointer(lessonTitle),
			}

			unit.Lessons = append(unit.Lessons, lesson)
		}

		var response []CourseResponse
		for _, course := range courses {
			response = append(response, *course)
		}

		// Return the courses with units and lessons as a JSON response
		c.JSON(http.StatusOK, response)
	})

	// Start the server
	router.Run(":8082")
}

// getUnitByID returns a unit with the specified ID from the given units slice
func getUnitByID(units []*UnitResponse, unitID int) *UnitResponse {
	for i := range units {
		if units[i].ID != nil && *units[i].ID == unitID {
			return units[i]
		}
	}
	return nil
}

// getIntPointer returns a pointer to an int value
func getIntPointer(val sql.NullInt64) *int {
	if val.Valid {
		value := int(val.Int64)
		return &value
	}
	return nil
}

// getStringPointer returns a pointer to a string value
func getStringPointer(val sql.NullString) *string {
	if val.Valid {
		value := val.String
		return &value
	}
	return nil
}
