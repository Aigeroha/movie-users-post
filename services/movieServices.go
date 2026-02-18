package services

import (
	"database/sql"
	"errors"
	"log"
	"project-last/database"
	"project-last/models"
)

var movies = []models.Movie{
	{ID: 1, Name: "Зверополис", Duration: 108, Genre: "анимация, комедия", Rating: 8.1}, // 0
	{ID: 2, Name: "28 лет спустя", Duration: 115, Genre: "ужасы, триллер", Rating: 7.4},
	{ID: 3, Name: "Вечность", Duration: 102, Genre: "драма, фантастика", Rating: 6.8},
	{ID: 4, Name: "Аватар", Duration: 162, Genre: "фантастика, приключения", Rating: 8.0},
	{ID: 5, Name: "Начало", Duration: 148, Genre: "фантастика, триллер", Rating: 8.8},
}

func GetMovies() ([]models.Movie, error) {
	rows, err := database.DB.Query("SELECT id, name, duration, genre, rating FROM movies")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Name, &movie.Duration, &movie.Genre, &movie.Rating); err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil

}

func UpdateMovie(id int, data map[string]interface{}) (models.Movie, error) {
	movie, err := GetMovieByID(id)
	if err != nil {
		return models.Movie{}, err
	}
	if movie.ID == id {
		log.Println(movie.Name, movie.Duration, movie.Genre, movie.Rating)
		if name, ok := data["title"].(string); ok {
			movie.Name = name
			log.Println("title", name)
		}
		if duration, ok := data["duration"].(float64); ok {
			movie.Duration = int(duration)
			log.Println("duration", duration)
		}
		if genre, ok := data["genre"].(string); ok {
			movie.Genre = genre
			log.Println("genre", genre)
		}
		if rating, ok := data["rating"].(float64); ok {
			movie.Rating = float64(rating)
		}
		_, err := database.DB.Exec(
			"UPDATE movies SET name = $1, duration = $2, genre = $3, rating = $4 WHERE id = $5",
			movie.Name,
			movie.Duration,
			movie.Genre,
			movie.Rating,
			id,
		)
		if err != nil {
			return models.Movie{}, err
		}
		return movie, nil
	}
	return models.Movie{}, errors.New("movie not found")
}

func GetMovieByID(id int) (models.Movie, error) {
	var movie models.Movie

	err := database.DB.QueryRow(
		"SELECT id, name, duration, genre, rating FROM movies WHERE id = $1",
		id,
	).Scan(&movie.ID, &movie.Name, &movie.Duration, &movie.Genre, &movie.Rating)

	if err == sql.ErrNoRows {
		return models.Movie{}, errors.New("movie not found")
	}

	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func PostMovie(movie models.Movie) (models.Movie, error) {
	err := database.DB.QueryRow(
		"INSERT INTO movies (name, duration, genre, rating) VALUES ($1, $2, $3, $4) RETURNING id",
		movie.Name,
		movie.Duration,
		movie.Genre,
		movie.Rating,
	).Scan(&movie.ID)

	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil

}

func DeleteMovie(id int) (string, error) {
	result, err := database.DB.Exec(
		"DELETE FROM movies WHERE id = $1",
		id,
	)
	if err != nil {
		return "something went wrong", err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return "", errors.New("movie not found")
	}
	return "Movie deleted", nil
}

func PutMovie(id int, updatedMovie models.Movie) (models.Movie, error) {
	result, err := database.DB.Exec(
		"UPDATE movies SET name = $1, duration = $2, genre = $3, rating = $4 WHERE id = $5",
		updatedMovie.Name,
		updatedMovie.Duration,
		updatedMovie.Genre,
		updatedMovie.Rating,
		id,
	)
	if err != nil {
		return models.Movie{}, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.Movie{}, errors.New("movie not found")
	}
	updatedMovie.ID = id
	return updatedMovie, nil
}

func GetMoviePaginated(limit int, page int) ([]models.Movie, error) {

	offset := (page - 1) * limit

	rows, err := database.DB.Query(
		"SELECT id, name, genre, duration, rating FROM movies ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie

		if err := rows.Scan(&m.ID, &m.Name, &m.Genre, &m.Duration, &m.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func GetMoviesFilter(title string, genre string) ([]models.Movie, error) {

	searchTitle := "%" + title + "%"
	searchGenre := "%" + genre + "%"

	query := `
        SELECT id, name, genre, duration, rating 
        FROM movies 
        WHERE name ILIKE $1 AND genre ILIKE $2 
        ORDER BY id`

	rows, err := database.DB.Query(query, searchTitle, searchGenre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Genre, &m.Duration, &m.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func GetMovieStats() (models.MovieStats, error) {
	var stats models.MovieStats

	query := `
        SELECT 
            COUNT(*), 
            COALESCE(AVG(duration), 0), MAX(duration), MIN(duration),
            COALESCE(AVG(rating), 0), MAX(rating), MIN(rating)
        FROM movies`

	err := database.DB.QueryRow(query).Scan(
		&stats.TotalMovies,      // Всего фильмов
		&stats.Duration.Average, // Средняя длительность
		&stats.Duration.Max,     // Макс. длительность
		&stats.Duration.Min,     // Мин. длительность
		&stats.Rating.Average,   // Средний рейтинг
		&stats.Rating.Max,       // Макс. рейтинг
		&stats.Rating.Min,       // Мин. рейтинг
	)

	return stats, err
}
