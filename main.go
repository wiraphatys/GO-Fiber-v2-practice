package main

import (
	"log"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"github.com/gofiber/swagger"
  _ "github.com/wiraphatys/GO-Fiber-v2-practice/docs" // load generated docs
)

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	books = append(books, Book{ID: 1, Title: "GO Tutorial", Author: "Banky"})
	books = append(books, Book{ID: 2, Title: "JAVA Master", Author: "Banky"})
	books = append(books, Book{ID: 3, Title: "C++ Master", Author: "Banky 2"})

	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Use(checkMiddleWare)

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBookByID)
	app.Post("/books/create", createBook)
	app.Put("/books/update", updateBook)
	app.Delete("/books/delete/:id", deleteBook)

	app.Post("/upload", uploadFile)

	app.Listen(":8080")
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err = c.SaveFile(file, "./uploads/"+file.Filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User{
	Email:    "bank",
	Password: "1234",
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email == memberUser.Email && user.Password == memberUser.Password {
		// Create the Claims
		claims := jwt.MapClaims{
			"email": user.Email,
			"role":  "admin",
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"token":   t,
		})
	} else {
		return fiber.ErrUnauthorized
	}
}

func checkMiddleWare(c *fiber.Ctx) error {
	user := c.Locals("user").(*)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}
