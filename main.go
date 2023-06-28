package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

	// Define a route to fetch data from the database
	router.GET("/data", func(c *gin.Context) {
		var (
			user  User
			users []User
		)
		rows, err := db.Query("select id, fullName, roll_number from users;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&user.Id, &user.FullName, &user.RollNumber)
			users = append(users, user)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": users,
			"count":  len(users),
		})
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
