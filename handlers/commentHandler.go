package handlers

import (
	"project-last/models"
	"project-last/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetComments(c *fiber.Ctx) error {
	comments, err := services.GetComments()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(), //"failed to get comments",
		})
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > len(comments) {
		limit = len(comments)
	}

	return c.Status(200).JSON(comments[:limit])
}

func GetCommentsByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	comment, err := services.GetCommentsByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(comment)
}

func GetCommentsByUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	comment, err := services.GetCommentsByUser(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(comment)
}

func CreateComments(c *fiber.Ctx) error {
	var comment models.Comment

	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}
	createdComments, err := services.CreateComments(comment)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(createdComments)
}

func DeleteComments(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid comment id",
		})
	}
	status, err := services.DeleteComments(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": status,
	})
}

func UpdateComments(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}

	user, err := services.UpdateComments(id, data)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

func PutComments(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedComment models.Comment
	if err := c.BodyParser(&updatedComment); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}
	result, err := services.PutComments(id, updatedComment)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}