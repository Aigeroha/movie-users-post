package services

import (
	"database/sql"
	"errors"
	"project-last/database"
	"project-last/models"
)

func GetUsers() ([]models.User, error) {
	rows, err := database.DB.Query("SELECT id, name, email, password FROM users")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}

func GetUserByID(id int) (models.User, error) {
	var user models.User

	err := database.DB.QueryRow(
		"SELECT id, name, email, password FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func CreateUser(user models.User) (models.User, error) {
	err := database.DB.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name,
		user.Email,
		user.Password,
	).Scan(&user.ID)

	if err != nil {
		return models.User{}, err
	}

	return user, nil

}

func DeleteUser(id int) (string, error) {
	result, err := database.DB.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return "something went wrong", err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return "", errors.New("user not found")
	}
	return "User deleted", nil
}

func UpdateUser(id int, data map[string]interface{}) (models.User, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}
	if name, ok := data["name"].(string); ok {
		user.Name = name
	}
	if email, ok := data["email"].(string); ok {
		user.Email = email
	}
	if password, ok := data["password"].(string); ok {
		user.Password = password
	}
	_, err = database.DB.Exec(
		"UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		user.Name,
		user.Email,
		user.Password,
		id,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func PutUser(id int, updatedUser models.User) (models.User, error) {
	result, err := database.DB.Exec(
		"UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		updatedUser.Name,
		updatedUser.Email,
		updatedUser.Password,
		id,
	)
	if err != nil {
		return models.User{}, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.User{}, errors.New("user not found")
	}
	updatedUser.ID = id
	return updatedUser, nil
}

func GetUsersByPage(limit int, page int) ([]models.User, error) {
    
    offset := (page - 1) * limit

    rows, err := database.DB.Query(
        "SELECT id, name, email, password FROM users ORDER BY id LIMIT $1 OFFSET $2", 
        limit, offset,
    )

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

func GetUsersByName(limit int, searchName string) ([]models.User, error) {
    
    queryName := "%" + searchName + "%"

    rows, err := database.DB.Query(
        "SELECT id, name, email, password FROM users WHERE name ILIKE $1 ORDER BY id LIMIT $2", 
        queryName, limit,
    )

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}