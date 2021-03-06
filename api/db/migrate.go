package db

import (
	migrations "auth_blog_service/migrations"

	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(connection *mongo.Database) {
	for _, migration := range migrations.GetList() {
		_, err := migrations.GetMigrations(connection, migration.Name)

		if err != nil {
			migration.Implementation(connection)

			migrations.SaveMigration(connection, migration.Name)
		}
	}
}
