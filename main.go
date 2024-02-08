package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// handler function
func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

// handler function
func getBookByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, book := range books {
		if book.ID == id {
			return c.JSON(book)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

// handler function
func createBook(c *fiber.Ctx) error {
	newBook := new(Book)
	if err := c.BodyParser(newBook); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, book := range books {
		if newBook.ID == book.ID {
			return c.Status(fiber.StatusConflict).SendString("Book with the given ID already exists.")
		}
	}

	books = append(books, *newBook)
	return c.Status(fiber.StatusCreated).SendString("Book has been created successfully.")
}

func updateBook(c *fiber.Ctx) error {
	target := new(Book)
	if err := c.BodyParser(target); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID == target.ID {
			books[i] = *target
			return c.Status(fiber.StatusOK).SendString("Book has been updated successfully.")
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func deleteBook(c *fiber.Ctx) error {
	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.ID == int(targetID) {
			books = append(books[:i], books[i+1:]...)
			return c.Status(fiber.StatusOK).SendString("Book has been deleted.")
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

var books []Book

func main() {
	app := fiber.New()

	books = append(books, Book{ID: 1, Title: "GO Tutorial", Author: "Banky"})
	books = append(books, Book{ID: 2, Title: "JAVA Master", Author: "Banky"})
	books = append(books, Book{ID: 3, Title: "C++ Master", Author: "Banky 2"})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBookByID)
	app.Post("/books/create", createBook)
	app.Put("/books/update", updateBook)
	app.Delete("/books/delete/:id", deleteBook)

	app.Listen(":8080")
}
