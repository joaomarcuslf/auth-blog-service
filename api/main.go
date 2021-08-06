package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	db "my_library_app/db"
)

var connection = db.ConnectDB()

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", HealthResponse).Methods("GET")

	var port = os.Getenv("PORT")

	fmt.Println("Server ready at http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HealthResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
