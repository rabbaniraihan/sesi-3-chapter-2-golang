package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Book struct {
	Id     int
	Title  string
	Author string
	Desc   string
}

var mapBooks = make(map[int]Book, 0)
var counter int

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=rabbani11 dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func main() {
	g := gin.Default()

	g.GET("/book", getAllBook)
	g.POST("/book", addBook)
	g.DELETE("/book/:id", deleteBook)
	g.GET("/book/:id", getBookById)
	g.PUT("/book/:id", updateBook)

	g.Run(":8080")
}

func getAllBook(ctx *gin.Context) {
	query := "SELECT * FROM book"

	rows, err := db.Query(query)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	books := make([]Book, 0)

	for rows.Next() {
		var book Book
		err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.Desc)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		books = append(books, book)
	}

	ctx.JSON(http.StatusOK, books)
}

func addBook(ctx *gin.Context) {
	var newBook Book

	err := ctx.ShouldBindJSON(&newBook)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := "INSERT INTO book (title, author, description) values($1, $2, $3) RETURNING *"

	row := db.QueryRow(query, newBook.Title, newBook.Author, newBook.Desc)

	err = row.Scan(&newBook.Id, &newBook.Title, &newBook.Author, &newBook.Desc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, newBook)
}

func deleteBook(ctx *gin.Context) {
	//Ambil id dari param
	stringId := ctx.Param("id")

	//Convert string -> int
	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var deletedBook Book

	query := "DELETE FROM book WHERE id=$1 RETURNING *"

	row := db.QueryRow(query, id)

	err = row.Scan(&deletedBook.Id, &deletedBook.Title, &deletedBook.Author, &deletedBook.Desc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, deletedBook)
}

func getBookById(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var getBook Book

	query := "SELECT * FROM book WHERE id=$1"

	row := db.QueryRow(query, id)

	err = row.Scan(&getBook.Id, &getBook.Title, &getBook.Author, &getBook.Desc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, getBook)
}

func updateBook(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var updatedBook Book

	err = ctx.ShouldBindJSON(&updatedBook)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	query := "UPDATE book SET title=$1, author=$2, description=$3 WHERE id=$4 RETURNING id"

	row := db.QueryRow(query, &updatedBook.Title, &updatedBook.Author, &updatedBook.Desc, id)

	err = row.Scan(&updatedBook.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}
