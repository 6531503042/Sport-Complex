package usecase

import (
	"context"
	"errors"
	"log"
	"main/modules/user"
	userPb "main/modules/user/proto"
	"main/modules/user/repository"
	"main/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (

	UserUsecaseService interface {
        CreateUser (pctx context.Context, req *user.CreateUserReq) (*user.UserProfile, error)
        FindOneUserProfile (pctx context.Context, userId string) (*user.UserProfile, error)
        FindOneUserCredential(pctx context.Context, password, email string) (*userPb.UserProfile, error)
		FindOneUserProfileToRefresh(pctx context.Context, userId string) (*userPb.UserProfile, error)
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

func (u *userUsecase) FindOneUserCredential(pctx context.Context, password, email string) (*userPb.UserProfile, error) {
    result, err := u.userRepository.FindOneUserCredential(pctx, email)
    if err != nil {
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
        log.Printf("Error: FindOneUserCredential: %s",err.Error())
        return nil, errors.New("error: wrong password")
    }

    roleCode := 0
    for _, v := range result.UserRoles {
        roleCode = v.RoleCode
    }

    loc, _ := time.LoadLocation("Asia/Bangkok")

    return &userPb.UserProfile{
        Id:        result.Id.Hex(),
        Email:     result.Email,
        Name:      result.Name,
        RoleCode:  int32(roleCode),
        CreatedAt: result.CreatedAt.In(loc).String(),
        UpdatedAt: result.CreatedAt.In(loc).String(),
    }, nil
}

func (u *userUsecase) FindOneUserProfileToRefresh(pctx context.Context, userId string) (*userPb.UserProfile, error) {
    result, err := u.userRepository.FindOneUserProfileRefresh(pctx, userId)
    if err != nil {
        return nil, err
    }

    roleCode := 0
    for _, v := range result.UserRoles {
        roleCode = v.RoleCode
    }

    loc, _ := time.LoadLocation("Asia/Bangkok")

    return &userPb.UserProfile{
        Id:        result.Id.Hex(),
        Email:     result.Email,
        Name:      result.Name,
        RoleCode:  int32(roleCode),
        CreatedAt: result.CreatedAt.In(loc).String(),
        UpdatedAt: result.CreatedAt.In(loc).String(),
    }, nil
}

