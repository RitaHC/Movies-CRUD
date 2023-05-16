package main

import(
	"fmt"
	"log"
	//Parse data to JSON
	"encoding/json"
	"math/rand"
	// Helps create server
	"net/http"
	// Id craeted using Math.random will neet to be converted to string
	"strconv"
	"github.com/gorilla/mux"
)


//==================== STRUCTS ======================

type Movie struct{
	ID string `json:"id`
	Isbn string `json: "isbn`
	Title string `json: "title`
	// *Director is a pointer (drawing a relation between Movie & Director)
	Director *Director `json: director`
}

type Director struct{
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

// Define a SLICE variable called movies  -> which has many items of type MOVIE (as defined in Struct)
var movies []Movie

//========================= ROUTE FUNCTIONS ===================

	//----------------- GET All ---------------------- (Postman- Working)

	func getMovies (w http.ResponseWriter, r *http.Request){
		// w http.http.ResponseWriter --> This is the response we are getting back
		// r *http.Request --> This is the request that we are making
		w.Header().Set("Content-Type", "application/json")
		// Setting up the type of data and how we need to use it
		json.NewEncoder(w).Encode(movies)
	}

	
	


	//----------------- DELETE One ----------------------

	func deleteMovie(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		// Need Params to grab a specific Id - EXTRACT FROM URL
		//This line extracts the parameters from the URL of the HTTP request 
		//Using the mux.Vars method from the gorilla/mux package. params is a map that contains all the parameters in the URL.
		params := mux.Vars(r)
		//This line checks if the ID field of the current item in the loop is equal to the "id" parameter in the URL.
		// If they are equal, then we have found the movie that we want to delete.
		for index, item := range movies{
			if item.ID == params["id"]{
				// removes a the element by referencing it's Id
				movies = append(movies[:index], movies[index+1:]...)
				break
			}

				
		}
		json.NewEncoder(w).Encode(movies)

	}



	//----------------- GET One ---------------------- 
	func getMovie(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		params :=mux.Vars(r)
		for _, item := range movies {
			if item.ID == params["id"]{
				// Return statement when the ID is found
				json.NewEncoder(w).Encode(item)
				return
			}
		}
	}

	//----------------- CREATE ----------------------
	func createMovie(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		var movie Movie
		// Decode the data from body + store it in the movie variable
		_ = json.NewDecoder(r.Body).Decode(&movie)
		// Generate a random number to create a Unique ID and convert it to string
		movie.ID = strconv.Itoa(rand.Intn(100000000))
		movies = append(movies, movie)
		json.NewEncoder(w).Encode(movie)
	}


	//----------------- UPDATE ----------------------
	

	func updateMovie(w http.ResponseWriter, r * http.Request){
		// Set jSON content type
		w.Header().Set("Content-Type", "application/json")
		// Access to params
		params := mux.Vars(r)
		// Loop over movies (range over)
		// Delete the movie with the ID that we've send
		// Add a new movie --> which will be the movie we send
		for index, item:= range movies{
			if item.ID == params["id"]{
				// First --> Delete 
				movies = append(movies[:index], movies[index+1:]...)
				// Second --> Append a new Movie
				var movie Movie
				_ = json.NewDecoder(r.Body).Decode(&movie)
				movie.ID = params["id"]
				movies = append(movies, movie)
				json.NewEncoder(w).Encode(movie)
				return
			}
		}
	}

/////////////////////////////////// Main Function //////////////////////////
func main(){

//====================== AUTH MIDDLEWARE ===================================

	users := map[string]string{
		"admin": "password",
	}
	
	// Define a middleware that checks for a valid username and password
	authMiddleware := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
	
			if !ok || users[username] != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - Unauthorized"))
				return
			}
	
			handler(w, r)
		}
	}



	//------------------------- ROUTER ----------------------
	// Router Created
	r := mux.NewRouter()

	// ----------------- SEED DATA ------------------

	movies = append(movies, Movie{ID: "1", Isbn: "376376", Title:"Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "264265", Title: "Movie Two", Director: &Director{Firstname: "Kimbal", Lastname: "Joe"}})



	//========================== ROUTES REGISTRATION (Like in Python) =======

	// Define all routes here (PATH + Function-Name)
	r.HandleFunc("/movies", authMiddleware(getMovies)).Methods("GET")
	// r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// ----------------------------STARTING THE SERVER ------------

	fmt.Printf("Starting Server at Port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}