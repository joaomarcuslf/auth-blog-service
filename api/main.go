package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	controllers "auth_blog_service/controllers"
	db "auth_blog_service/db"
)

var connection = db.ConnectDB()

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", HealthResponse).Methods("GET")
	r.HandleFunc("/api/roles", controllers.GetRoles(connection)).Methods("GET")
	r.HandleFunc("/api/roles", controllers.CreateRole(connection)).Methods("POST")

	var port = os.Getenv("PORT")

	fmt.Println("Server ready at http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HealthResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
