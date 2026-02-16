package handlers

import (
	"project-last/models"
	"project-last/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetMovies(c *fiber.Ctx) error {
	movies, err := services.GetMovies()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = len(movies)
	}

	return c.Status(200).JSON(movies[:limit])
}

func UpdateMovie(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid ID",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}

	movie, err := services.UpdateMovie(id, data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(movie)
}

func GetMovieByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	movie, err := services.GetMovieByID(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Movie not found",
		})
	}
	return c.JSON(movie)
}

func PostMovie(c *fiber.Ctx) error {
	var movie models.Movie

	if err := c.BodyParser(&movie); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}
	newMovie, err := services.PostMovie(movie)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(201).JSON(newMovie)
}

func DeleteMovie(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	status, err := services.DeleteMovie(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": status,
	})
}

func PutMovie(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var updatedMovie models.Movie
	if err := c.BodyParser(&updatedMovie); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}
	result, err := services.PutMovie(id, updatedMovie)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetMoviePaginated(c *fiber.Ctx) error {

	limit, errLimit := strconv.Atoi(c.Query("limit", "10"))
	page, errPage := strconv.Atoi(c.Query("page", "1"))

	if errLimit != nil || limit <= 0 {
		limit = 10
	}
	if errPage != nil || page <= 0 {
		page = 1
	}

	movies, err := services.GetMoviePaginated(limit, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch movies",
		})
	}

	return c.JSON(movies)
}

func GetMoviesFilter(c *fiber.Ctx) error {

	title := c.Query("title", "")
	genre := c.Query("genre", "")

	movies, err := services.GetMoviesFilter(title, genre)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch movies",
		})
	}

	return c.JSON(movies)
}

func GetMovieStats(c *fiber.Ctx) error {
	stats, err := services.GetMovieStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get stats"})
	}
	return c.JSON(stats)
}
