package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Kevin-Bian/BianPhotography2.0/src/models"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// createConnection Connects to our PSQL DB
func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

// InsertPhoto Uploads a photo to collage
func InsertPhoto(photo models.Photo) int64 {

	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO photos (collageid, name, link, description) VALUES ($1, $2, $3, $4) RETURNING photoid`
	var id int64

	err := db.QueryRow(sqlStatement, photo.CollageID, photo.Name, photo.Link, photo.Description).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Println("Inserted a single record %v", id)
	return id
}

// GetPhoto Gets a photo given id
func GetPhoto(id int64) (models.Photo, error) {
	db := createConnection()
	defer db.Close()

	var photo models.Photo
	sqlStatement := `SELECT * FROM photos WHERE photoid=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&photo.ID, &photo.CollageID, &photo.Name, &photo.Link, &photo.Description)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return photo, nil
	case nil:
		return photo, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return photo, err
}

// GetAllPhotos Gets all photos uploaded
func GetAllPhotos() ([]models.Photo, error) {
	db := createConnection()
	defer db.Close()

	var photos []models.Photo
	sqlStatement := `SELECT * FROM photos`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var photo models.Photo
		err = rows.Scan(&photo.ID, &photo.CollageID, &photo.Name, &photo.Link, &photo.Description)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		photos = append(photos, photo)

	}
	return photos, err
}

// GetCollage Gets a collage given id
func GetCollage(collageid string) ([]models.Photo, error) {
	db := createConnection()
	defer db.Close()

	var photos []models.Photo

	sqlStatement := `SELECT * FROM photos where collageid=$1`

	rows, err := db.Query(sqlStatement, collageid)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var photo models.Photo
		err = rows.Scan(&photo.ID, &photo.CollageID, &photo.Name, &photo.Link, &photo.Description)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		photos = append(photos, photo)

	}
	return photos, err
}

// DeletePhoto Deletes a photo by id
func DeletePhoto(id int64) (string, error) {
	db := createConnection()
	defer db.Close()

	var photo models.Photo
	sqlStatement := `DELETE FROM photos WHERE photoid=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&photo.ID, &photo.CollageID, &photo.Name, &photo.Link, &photo.Description)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return "Success", nil
	case nil:
		return "Success", nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return "Failiure", err
}
