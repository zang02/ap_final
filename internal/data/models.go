package data

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	// Tokens TokenModel
	Users UserModel
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the initialized MovieModel.
func NewModels(db *mongo.Database) Models {
	return Models{
		// Movies:      MovieModel{DB: db},
		// Permissions: PermissionModel{DB: db},
		// Tokens: TokenModel{DB: db},
		Users: UserModel{DB: db},
	}
}

// func NewMockModels() Models {
// 	return Models{
// 		// Movies: MockMovieModel{},
// 		Users: UserModel{},
// 	}
// }
