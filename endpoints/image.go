package endpoints

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Id           int     `json:"id"`
	Species_name string  `json:"species_name"`
	Gps_long     float64 `json:"gps_long"`
	Gps_lat      float64 `json:"gps_lat"`
	Image_path   string  `json:"image_path"`
	User_id      int     `json:"user_id"`
}

func GetImages(conn *pgx.Conn) ([]Image, error) {

	rows, err := conn.Query(context.Background(), "SELECT * FROM images")
	if err != nil {
		return nil, fmt.Errorf("error getting images: %v", err)
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.Id, &image.Species_name, &image.Gps_long, &image.Gps_lat, &image.Image_path, &image.User_id)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil

}
