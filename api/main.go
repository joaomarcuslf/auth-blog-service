package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gorilla/mux"

	controllers "auth_blog_service/controllers"
	db "auth_blog_service/db"
)

var connection = db.ConnectDB()

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		requestInfo := strings.Split(string(x), "\n")
		requestInfo = strings.Split(requestInfo[0], " ")

		fmt.Println(requestInfo[0], requestInfo[1])

		fn(w, r)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", logHandler(HealthResponse)).Methods("GET")

	r.HandleFunc("/api/roles", logHandler(controllers.GetRoles(connection))).Methods("GET")
	r.HandleFunc("/api/roles", logHandler(controllers.CreateRole(connection, "role.create"))).Methods("POST")
	r.HandleFunc("/api/roles/{id}", logHandler(controllers.GetRoleById(connection))).Methods("GET")
	r.HandleFunc("/api/roles/{id}", logHandler(controllers.UpdateRoleById(connection, "role.update"))).Methods("PUT")
	r.HandleFunc("/api/roles/{id}", logHandler(controllers.DeleteRoleById(connection, "role.delete"))).Methods("DELETE")

	r.HandleFunc("/api/users", logHandler(controllers.GetUsers(connection))).Methods("GET")
	r.HandleFunc("/api/users", logHandler(controllers.CreateUser(connection, "user.create"))).Methods("POST")
	r.HandleFunc("/api/users/{id}", logHandler(controllers.GetUserById(connection))).Methods("GET")
	r.HandleFunc("/api/users/{id}/role", logHandler(controllers.GetUserRoleById(connection))).Methods("GET")
	r.HandleFunc("/api/users/{id}/posts", logHandler(controllers.GetUserPostsById(connection))).Methods("GET")
	r.HandleFunc("/api/users/{id}", logHandler(controllers.UpdateUserById(connection, "user.update"))).Methods("PUT")
	r.HandleFunc("/api/users/{id}", logHandler(controllers.DeleteUserById(connection, "user.delete"))).Methods("DELETE")

	r.HandleFunc("/api/posts", logHandler(controllers.GetPosts(connection))).Methods("GET")
	r.HandleFunc("/api/posts", logHandler(controllers.CreatePost(connection, "post.create"))).Methods("POST")
	r.HandleFunc("/api/posts/{id}", logHandler(controllers.GetPostById(connection))).Methods("GET")
	r.HandleFunc("/api/posts/{id}", logHandler(controllers.UpdatePostById(connection, "post.update"))).Methods("PUT")
	r.HandleFunc("/api/posts/{id}", logHandler(controllers.DeletePostById(connection, "post.delete"))).Methods("DELETE")

	var port = os.Getenv("PORT")

	fmt.Println("Server ready at http://localhost:" + port + "/")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func HealthResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
