package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rs.fyi/go-gorm-setup/models"
	"github.com/rs.fyi/go-gorm-setup/storage"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := models.Books{}

	err := ctx.BodyParser(&book)

	if err != nil {
		// TODO - figure out why if we log the err value we cannot return the response
		// log.Fatal("Error processing request body ", err)
		return ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"isSuccess": false,
			"message":   err,
		})
	}

	// Save data to database
	result := r.DB.Create(&book)

	if result.Error != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"isSuccess": false,
			"message":   result.Error,
		})
	}

	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message":   "Books fetched successfully",
		"isSuccess": true,
		"data":      book,
	})
}

func (r *Repository) DeleteBookById(ctx *fiber.Ctx) error {
	bookModel := models.Books{}

	id := ctx.Params("id")

	log.Println("id .... ", id)

	if id == "" {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":   "Missing id params",
			"isSuccess": false,
		})
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":   err.Error,
			"isSuccess": false,
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":   "Delete book successfully",
		"isSuccess": true,
	})
}

func (r *Repository) GetBookById(ctx *fiber.Ctx) error {
	bookModel := &models.Books{}

	id := ctx.Params("id")

	if id == "" {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":   "Missing id params",
			"isSuccess": false,
		})
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":   err,
			"isSuccess": false,
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":   "fetched book successfully",
		"isSuccess": true,
		"data":      bookModel,
	})
}

// GetBooks returns all the books
func (r *Repository) GetBooks(ctx *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":   err,
			"isSuccess": false,
		})
	}

	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message":   "Books fetched successfully",
		"isSuccess": true,
		"data":      bookModels,
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("I am default home route")
	})

	api.Post("/create-books", r.CreateBook)
	api.Delete("/delete-book/:id", r.DeleteBookById)
	api.Get("/get-books/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

func main() {
	log.Println("[ Main ] ")

	err := godotenv.Load()

	if err != nil {
		log.Println("Error while loading env variables", err)
	}

	config := storage.Config{
		Host:     os.Getenv("Host"),
		Port:     os.Getenv("Port"),
		User:     os.Getenv("User"),
		Password: os.Getenv("Password"),
		DBName:   os.Getenv("DBName"),
		SSLMode:  os.Getenv("SSLMode"),
	}

	db, err := storage.NewConnection(&config)

	if err != nil {
		log.Fatal("Could not load db ", err)
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("Error while db migrations ", err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	r.SetupRoutes(app)

	log.Println("Server is running successfully !")
	app.Listen(":3000")

}
