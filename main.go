package main

import (
	"project-last/database"
	"project-last/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()

	app.Get("/movies", handlers.GetMovies)

	app.Patch("/movies/:id", handlers.UpdateMovie)

	app.Get("/movies/:id", handlers.GetMovieByID)

	app.Post("/movies", handlers.PostMovie)

	app.Delete("/movies/:id", handlers.DeleteMovie)

	app.Put("/movies/:id", handlers.PutMovie)

	app.Get("/movies/:page", handlers.GetMoviePaginated)

	app.Get("/movies", handlers.GetMoviesFilter)

	app.Get("/movies/stats", handlers.GetMovieStats)

	app.Get("/users", handlers.GetUsers)

	app.Get("/users/:id", handlers.GetUserByID)

	app.Post("/users", handlers.CreateUser)

	app.Patch("/users/:id", handlers.UpdateUser)

	app.Delete("/users/:id", handlers.DeleteUser)

	app.Get("/comments", handlers.GetComments)
	app.Get("/comments/:id", handlers.GetCommentsByID)
	app.Get("/users/:id/comments", handlers.GetCommentsByUser)
	app.Post("/comments", handlers.CreateComments)
	app.Put("/comments/:id", handlers.PutComments)
	app.Delete("/comments/:id", handlers.DeleteComments)

	app.Listen(":3000")
}
