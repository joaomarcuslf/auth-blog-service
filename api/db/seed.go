package db

import (
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"auth_blog_service/models"
	repositories "auth_blog_service/repositories"
	types "auth_blog_service/types"
)

func Seed(connection *mongo.Database) {
	if os.Getenv("ENV") == "PRODUCTION" {
		return
	}

	roles, _, _ := repositories.GetRoles(connection)

	if len(roles) == 0 {
		fmt.Println("Seeding roles")

		roles := []models.Role{
			{
				Name:        "User",
				Permissions: []string{},
			},
			{
				Name: "Admin",
				Permissions: []string{
					"role.update",
					"role.create",
					"role.delete",
					"user.update",
					"user.create",
					"user.delete",
					"post.update",
					"post.create",
					"post.delete",
				},
			},
		}

		for _, role := range roles {
			repositories.InsertRole(connection, role)
		}
	}

	users, _, _ := repositories.GetUsers(connection)

	if len(users) == 0 {
		fmt.Println("Seeding users")

		userRole, _, _ := repositories.QueryRoles(connection, bson.M{"name": "User"})
		adminRole, _, _ := repositories.QueryRoles(connection, bson.M{"name": "Admin"})

		users := []models.User{
			{
				RoleID:   userRole[0].ID,
				UserName: "Test",
				Name:     "Test",
				BirthDate: types.Datetime{
					Time: time.Date(1996, 6, 26, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				RoleID:   adminRole[0].ID,
				UserName: "Admin",
				Name:     "Admin Test",
				BirthDate: types.Datetime{
					Time: time.Date(1996, 6, 26, 0, 0, 0, 0, time.UTC),
				},
			},
		}

		for _, user := range users {
			repositories.InsertUser(connection, user)
		}
	}

	posts, _, _ := repositories.GetPosts(connection)

	if len(posts) == 0 {
		fmt.Println("Seeding posts")

		user, _, _ := repositories.QueryUsers(connection, bson.M{"username": "Test"})

		posts := []models.Post{
			{
				UserID: user[0].ID,
				Title:  "Testing post 01",
				Body:   "Test body with some changes\nline",
				CreatedDate: types.Datetime{
					Time: time.Now(),
				},
			},
			{
				UserID: user[0].ID,
				Title:  "Testing post 02",
				Body:   "Test body with some changes\nline",
				CreatedDate: types.Datetime{
					Time: time.Now(),
				},
			},
		}

		for _, post := range posts {
			repositories.InsertPost(connection, post)
		}
	}
}
