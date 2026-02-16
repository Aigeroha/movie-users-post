package models

type Comment struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
	UserID  int    `json:"user_id"`
	MovieID int    `json:"movie_id"`
}
