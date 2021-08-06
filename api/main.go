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
	r.HandleFunc("/api/roles", controllers.CreateRole(connection, "role.create")).Methods("POST")
	r.HandleFunc("/api/roles/{id}", controllers.GetRoleById(connection)).Methods("GET")
	r.HandleFunc("/api/roles/{id}", controllers.UpdateRoleById(connection, "role.update")).Methods("PUT")
	r.HandleFunc("/api/roles/{id}", controllers.DeleteRoleById(connection, "role.delete")).Methods("DELETE")

	r.HandleFunc("/api/users", controllers.GetUsers(connection)).Methods("GET")
	r.HandleFunc("/api/users", controllers.CreateUser(connection, "user.create")).Methods("POST")
	r.HandleFunc("/api/users/{id}", controllers.GetUserById(connection)).Methods("GET")
	r.HandleFunc("/api/users/{id}/role", controllers.GetUserRoleById(connection)).Methods("GET")
	// r.HandleFunc("/api/users/{id}/posts", controllers.GetUserPostsById(connection)).Methods("GET")
	r.HandleFunc("/api/users/{id}", controllers.UpdateUserById(connection, "user.update")).Methods("PUT")
	r.HandleFunc("/api/users/{id}", controllers.DeleteUserById(connection, "user.delete")).Methods("DELETE")

	var port = os.Getenv("PORT")

	fmt.Println("Server ready at http://localhost:" + port + "/")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HealthResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
