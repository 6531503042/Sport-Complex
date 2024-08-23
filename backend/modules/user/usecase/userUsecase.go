package usecase

import (
	"context"
	"errors"
	"main/modules/user"
	"main/modules/user/repository"
	"main/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type (

	UserUsecaseService interface {
		
	}

	userUsecase struct {
		userRepository repository.UserRepositoryService
	}
)

func NewUserUsecase(userRepository repository.UserRepositoryService) UserUsecaseService {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (u *userUsecase) createUser(pctx context.Context, req *user.CreateUserReq) (*user.UserProfile, error) {
    // Check if the user with the given email or name already exists
    if !u.userRepository.IsUniqueUser(pctx, req.Email, req.Name) {
        return nil, errors.New("error: email or name already existing")
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("Error: failed to hash password")
    }

    // Insert the new user
    userId, err := u.userRepository.InsertOneUser(pctx, &user.User{
        Email:     req.Email,
        Name:      req.Name,
        Password:  string(hashedPassword),
        CreatedAt: utils.LocalTime(),
        UpdatedAt: utils.LocalTime(),
        UserRoles: []user.UserRole{
            {
                RoleTitle: "user",
                RoleCode:  0,
            },
        },
	})

	return u.userRepository.FindOneUserProfile(pctx, userId.Hex())
}
