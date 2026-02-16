package handlers

import (
	"project-last/models"
	"project-last/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UpdateUser(c *fiber.Ctx) error {
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

	user, err := services.UpdateUser(id, data)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	users, err := services.GetUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(), //"failed to get users",
		})
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > len(users) {
		limit = len(users)
	}

	return c.Status(200).JSON(users[:limit])
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	status, err := services.DeleteUser(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status": status,
	})
}

func GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	user, err := services.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}
	createdUser, err := services.CreateUser(user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(createdUser)
}

func PutUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid json",
		})
	}
	result, err := services.PutUser(id, updatedUser)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetUsersByPage(c *fiber.Ctx) error {

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	users, err := services.GetUsersByPage(limit, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.Status(200).JSON(users)
}

func GetUsersByName(c *fiber.Ctx) error {
    limit, err := strconv.Atoi(c.Query("limit", "100"))
    if err != nil || limit <= 0 {
        limit = 100
    }

    searchName := c.Query("name", "") 

    users, err := services.GetUsersByName(limit, searchName)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to fetch users",
        })
    }

    return c.Status(200).JSON(users)
}