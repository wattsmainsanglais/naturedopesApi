package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"naturedopesApi/endpoints"

	"github.com/joho/godotenv"
)

func connectToDB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	conn, err := connectToDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer conn.Close(context.Background())

	fmt.Println("Connected to database!")

	images, err := endpoints.GetImages(conn)
	if err != nil {
		fmt.Println("Error retrieving images:", err)
		return
	}

	fmt.Println("Images:")
	for _, image := range images {
		fmt.Printf("ID: %d, SpeciesName: %s, GPSLong: %f, GPSLat: %f, ImagePath: %s, UserID: %d\n",
			image.Id, image.Species_name, image.Gps_long, image.Gps_lat, image.Image_path, image.User_id)
	}

}
