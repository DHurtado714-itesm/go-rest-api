package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type book struct {
	ID string `json:"id"`
	Title string`json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}

func connectDatabase() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err 
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")

	err = createTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS books (id VARCHAR(255) PRIMARY KEY, title VARCHAR(255), author VARCHAR(255), quantity INT)")
	if err != nil {
		return fmt.Errorf("could not create table: %v", err)
	}

	fmt.Println("Table created successfully")

	return nil
}

func getBooks(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []book
	for rows.Next() {
		var b book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, b)
	}

	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context, db *sql.DB) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = ?)", newBook.ID).Scan(&exists)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "book with this ID already exists"})
		return
	}


	_, err = db.Exec("INSERT INTO books (id, title, author, quantity) VALUES (?, ?, ?, ?)",
		newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

		c.IndentedJSON(http.StatusCreated, newBook)
}

func editBook(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var updatedBook book

	if err := c.BindJSON(&updatedBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("UPDATE books SET title = ?, author = ?, quantity = ? WHERE id = ?", updatedBook.Title, updatedBook.Author, updatedBook.Quantity, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	

	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
}

func deleteBook(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
}

func main() {
	// Connect to the database
	db, err := connectDatabase()
	if err != nil {
		log.Fatal("Could not connect to the database", err)
	}
	defer db.Close()

	// Create a new gin router
	router := gin.Default()

	// Define the routes
	router.GET("/books", func(c *gin.Context) {
		getBooks(c, db)
	})
	router.POST("/books", func(c *gin.Context) {
		createBook(c, db)
	})
	router.PUT("/books/:id", func(c *gin.Context) {
		editBook(c, db)
	})
	router.DELETE("/books/:id", func(c *gin.Context) {
		deleteBook(c, db)
	})

	router.Run("localhost:8080")
}