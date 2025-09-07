package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"naturedopesApi/endpoints"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/images", getImagesHandler).Methods("GET")
	router.HandleFunc("/images/{id}", getImageHandler).Methods("GET")
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	images, err := endpoints.GetImages(conn, &idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(images) == 0 {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(images[0])
}

func getImagesHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := connectToDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	images, err := endpoints.GetImages(conn, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(images)

}
