package services

import (
	"database/sql"
	"errors"
	"project-last/database"
	"project-last/models"
)

//	func GetUsers() []models.User {
//		return users
//	}
func GetComments() ([]models.Comment, error) {
	rows, err := database.DB.Query("SELECT id, comment, user_id, movie_id FROM comments")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.UserID, &comment.MovieID); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil

}

func GetCommentsByUser(userId int) (models.Comment, error) {
	var comments models.Comment

	err := database.DB.QueryRow(
		"SELECT id, comment, user_id, movie_id FROM comments WHERE user_id = $1",
		userId,
	).Scan(&comments.ID, &comments.Comment, &comments.UserID, &comments.MovieID)

	if err == sql.ErrNoRows {
		return models.Comment{}, errors.New("comment not found")
	}

	if err != nil {
		return models.Comment{}, err
	}

	return comments, nil
}

func GetCommentsByID(id int) (models.Comment, error) {
	var comment models.Comment

	err := database.DB.QueryRow(
		"SELECT id, comment, user_id, movie_id FROM comments WHERE id = $1",
		id,
	).Scan(&comment.ID, &comment.Comment, &comment.UserID, &comment.MovieID)

	if err == sql.ErrNoRows {
		return models.Comment{}, errors.New("comment not found")
	}

	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func CreateComments(comment models.Comment) (models.Comment, error) {
	err := database.DB.QueryRow(
		"INSERT INTO comments (comment, user_id, movie_id) VALUES ($1, $2, $3) RETURNING id",
		comment.Comment,
		comment.UserID,
		comment.MovieID,
	).Scan(&comment.ID)

	if err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func DeleteComments(id int) (string, error) {
	result, err := database.DB.Exec(
		"DELETE FROM comments WHERE id = $1",
		id,
	)
	if err != nil {
		return "something went wrong", err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return "", errors.New("comment not found")
	}
	return "Comment deleted", nil
}

func UpdateComments(id int, data map[string]interface{}) (models.Comment, error) {
	comment, err := GetCommentsByID(id)
	if err != nil {
		return models.Comment{}, err
	}
	if commentText, ok := data["comment"].(string); ok {
		comment.Comment = commentText
	}
	_, err = database.DB.Exec(
		"UPDATE comments SET comment = $1 WHERE id = $2",
		comment.Comment,
		id,
	)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func PutComments(id int, updatedComment models.Comment) (models.Comment, error) {
	result, err := database.DB.Exec(
		"UPDATE comments SET comment = $1 WHERE id = $2",
		updatedComment.Comment,

		id,
	)
	if err != nil {
		return models.Comment{}, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.Comment{}, errors.New("comment not found")
	}
	updatedComment.ID = id
	return updatedComment, nil
}
