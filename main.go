package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	ISBN     uint32    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMovie Movie
	reqBody, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &newMovie)

	if err != nil {
		log.Fatal(err)
	}

	newMovie.ID = strconv.Itoa(rand.Intn(100000))
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(newMovie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var newMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&newMovie)
	newMovie.ID = params["id"]

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], newMovie)
			movies = append(movies, movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(newMovie)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies,
		Movie{ID: "1010", Title: "The Shawshank Redemption", ISBN: 851709680,
			Director: &Director{FirstName: "X", LastName: "Y", Age: 40},
		})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))

}
