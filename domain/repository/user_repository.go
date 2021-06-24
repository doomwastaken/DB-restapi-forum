package repository

import "forum/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
	CheckIfUserExists(nickname string) (string, error)
	GetUserByNickname(nickname string) (*entity.User, error)
	UpdateUser(newUser *entity.User) (*entity.User, error)
	GetUserNicknameWithEmail(email string) (string, error)
	GetUsersWithNicknameAndEmail(nickname, email string) ([]entity.User, error)
}
//type UserRepository interface {
//	CreateUser(*entity.User) (int, error)           // Create user, returns created user's ID
//	SaveUser(*entity.User) error                    // Save changed user to database
//	DeleteUser(int) error                           // Delete user with passed userID from database
//	GetUser(int) (*entity.User, error)              // Get user by his ID
//	GetUsers() ([]entity.User, error)               // Get all users
//	GetUserByUsername(string) (*entity.User, error) // Get user by his username
//	Follow(int, int) error                          // Make first user follow second
//	Unfollow(int, int) error                        // Make first user unfollow second
//	CheckIfFollowed(int, int) (bool, error)         // Check if first user follows second
//}
