package usecase

import (
	"context"
	"errors"
	"main/modules/user"
	"main/modules/user/repository"
	"main/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (

	UserUsecaseService interface {
        CreateUser (pctx context.Context, req *user.CreateUserReq) (*user.UserProfile, error)
        FindOneUserProfile (pctx context.Context, userId string) (*user.UserProfile, error)
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

func (u *userUsecase) CreateUser(pctx context.Context, req *user.CreateUserReq) (*user.UserProfile, error) {
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
    if err != nil {
        return nil, errors.New("Error: failed to create user")
    }

    return u.FindOneUserProfile(pctx, userId.Hex())
}

func (u * userUsecase) FindOneUserProfile (pctx context.Context, userId string) (*user.UserProfile, error) {
    result, err := u.userRepository.FindOneUserProfile(pctx, userId)
    if err != nil {
        return nil, err
    }

    loc, _ := time.LoadLocation("Asia/Bangkok")

    return &user.UserProfile{
        Id:        result.Id.Hex(),
        Email:     result.Email,
        Name:      result.Name,
        CreatedAt: result.CreatedAt.In(loc),
        UpdatedAt: result.UpdatedAt.In(loc),
    }, nil
}