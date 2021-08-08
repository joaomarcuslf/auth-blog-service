package serializers

import (
	"auth_blog_service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	Name        string             `json:"name"`
	Permissions []string           `json:"permissions"`
}

func SerializeOneRole(role models.Role) Role {
	return Role{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: role.Permissions,
	}
}

func SerializeManyRoles(roles []models.Role) []Role {
	var rolesArray []Role

	for _, role := range roles {
		rolesArray = append(rolesArray, SerializeOneRole(role))
	}

	return rolesArray
}
