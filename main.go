// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// )

// type Movie struct {
// 	ID       string    `json:"id"`
// 	Isbn     string    `json:"isbn"`
// 	Title    string    `json:"title"`
// 	Director *Director `json:"director"`
// }

// type Director struct {
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// }

// var movies []Movie

// func getMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(movies)
// }

// func deleteMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range movies {

// 		if item.ID == params["id"] {
// 			movies = append(movies[:index], movies[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(movies)
// }

// func getMovie(w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Content-Type", "applications/json")
// 	params := mux.Vars(r)
// 	for _, item := range movies {
// 		if item.ID == params["id"] {
// 			json.NewEncoder(w).Encode(item)
// 			return
// 		}
// 	}

// }

// func createMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var movie Movie
// 	_ = json.NewDecoder(r.Body).Decode(&movie)
// 	movie.ID = strconv.Itoa(rand.Intn(10000000000))
// 	movies = append(movies, movie)
// 	json.NewEncoder(w).Encode(movie)

// }

// func updateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range movies {
// 		if item.ID == params["id"] {
// 			movies = append(movies[:index], movies[index+1:]...)
// 			var movie Movie
// 			_ = json.NewDecoder(r.Body).Decode(&movie)
// 			movie.ID = params["Id"]
// 			movies = append(movies, movie)
// 			json.NewEncoder(w).Encode(movie)
// 			return
// 		}
// 	}
// }

// func main() {
// 	r := mux.NewRouter()

// 	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie 1", Director: &Director{FirstName: "John", LastName: "Doe"}})
// 	movies = append(movies, Movie{ID: "2", Isbn: "234562", Title: "Movie 2", Director: &Director{FirstName: "Peter", LastName: "Parker"}})
// 	r.HandleFunc("/movies", getMovies).Methods("GET")
// 	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
// 	r.HandleFunc("/movies", createMovie).Methods("POST")
// 	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
// 	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

// 	fmt.Printf("Sarting server at port 8000 ")
// 	log.Fatal(http.ListenAndServe(":8000", r))

// }

package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	// Get the GitHub token from the environment variables
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	// Get the repository owner and name from the environment variables
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	repo := os.Getenv("GITHUB_REPOSITORY_NAME")
	issueNumber := os.Getenv("ISSUE_NUMBER")
	issueBody := os.Getenv("ISSUE_BODY")

	if owner == "" || repo == "" || issueNumber == "" {
		log.Fatal("GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY_NAME or ISSUE_NUMBER is not set")
	}

	issueNum, err := strconv.Atoi(issueNumber)
	if err != nil {
		log.Fatalf("Invalid issue number: %v", err)
	}

	// Assignee for bug issues
	assignee := "GreshmaShaji"

	// Create a new GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Determine the label based on the issue body
	var labels []string
	log.Printf("issueBody %v", issueBody)
	if issueBody == "" || strings.TrimSpace(issueBody) == "" {
		labels = []string{"action required"}
	} else if strings.Contains(strings.ToLower(issueBody), "version") {
		log.Printf("Documentation label selected")
		labels = []string{"documentation"}
	} else {
		labels = []string{"bug"}
	}

	// Add labels to the issue
	_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, issueNum, labels)
	if err != nil {
		log.Fatalf("Failed to add labels to issue: %v", err)
	}

	log.Printf("Labels %v added to issue #%d", labels, issueNum)

	// If the label is "bug", assign the issue to the specified user
	if contains(labels, "bug") {
		_, _, err = client.Issues.AddAssignees(ctx, owner, repo, issueNum, []string{assignee})
		if err != nil {
			log.Fatalf("Failed to assign issue to %s: %v", assignee, err)
		}
		log.Printf("Issue #%d assigned to %s", issueNum, assignee)
	}
}

// contains checks if a slice contains a specified string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
