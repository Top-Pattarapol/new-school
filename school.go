package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Todo struct {
	Id     int
	Title  string
	Status string
}

func main() {

	r := gin.Default()

	r.GET("/api/todos", getTodos)

	r.Run(":1234")

}

func getTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer db.Close()

	stmt, err := db.Prepare(`Select id, title, status FROM todos ORDER BY id ASC`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := stmt.Query()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todos := []Todo{}

	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.Id, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}

	fmt.Println(todos)
	c.JSON(http.StatusOK, todos)
}
