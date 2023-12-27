package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func main() {
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Book{})
	fmt.Println("Database migration completed!")

	// Setup Fiber
	app := fiber.New()

	// CRUD routes
	// select all
	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(getBooks(db))
	})

	// select id
	app.Get("/book/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book := getBook(db, id)
		return c.JSON(book)
	})

	// create
	app.Post("/book", func(c *fiber.Ctx) error {
		book := new(Book)
		err := c.BodyParser(book)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err2 := createBook(db, book)
		if err2 != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	// update
	app.Put("/book/:id", func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book := new(Book)
		err2 := c.BodyParser(book)
		if err2 != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book.ID = uint(id)
		err3 := updateBook(db, book)
		if err3 != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	// delete
	app.Delete("/book/:id", func(c *fiber.Ctx) error {
		id, err1 := strconv.Atoi(c.Params("id"))
		if err1 != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err3 := deleteBook(db, id)
		if err3 != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	app.Listen(":8080")

	// insert data
	// newBook := Book{
	// 	Name:        "The Bubble",
	// 	Author:      "Blue",
	// 	Description: "Comprehensive guide to Go",
	// }
	// CreateBook(db, &newBook)

	// select data
	// book := GetBook(db, 1)
	// fmt.Println(book)

	// update data
	// book.Name = "The Go Programming Language, Updated Edition"
	// book.Price = 400
	// UpdateBook(db, book)

	// delete data (sorf delete)
	// DeleteBook(db, book.ID)

	// delete data (delete)
	// DeleteBook_unscope(db, 1)
}
