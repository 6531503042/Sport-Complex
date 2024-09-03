package usecase

import (
	"context"
	models "main/modules/user/model"
	"main/modules/user/repository"
	"time"
)

type (
	IUserUsecase interface {
		CreateUser(ctx context.Context, req models.CreateUserReq) (*models.UserProfile, error)
		GetUserByID(ctx context.Context, id string) (*models.UserProfile, error)
		GetAllUsers(ctx context.Context) ([]models.UserProfile, error)
		UpdateUser(ctx context.Context, id string, req models.CreateUserReq) (*models.UserProfile, error)
		DeleteUser(ctx context.Context, id string) error
	}

	userUsecase struct {
		userRepo repository.IUserRepository
	}
)

func NewUserUsecase(userRepo repository.IUserRepository) IUserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (uc *userUsecase) CreateUser(ctx context.Context, req models.CreateUserReq) (*models.UserProfile, error) {
	user := models.User{
		Email:     req.Email,
		Name:      req.Name,
		Password:  req.Password, // Note: In a real application, you should hash passwords before storing them
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	userProfile := &models.UserProfile{
		Id:        createdUser.Id.Hex(),
		Email:     createdUser.Email,
		Name:      createdUser.Name,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}

	return userProfile, nil
}

func (uc *userUsecase) GetUserByID(ctx context.Context, id string) (*models.UserProfile, error) {
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userProfile := &models.UserProfile{
		Id:        user.Id.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userProfile, nil
}

func (uc *userUsecase) GetAllUsers(ctx context.Context) ([]models.UserProfile, error) {
	users, err := uc.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userProfiles []models.UserProfile
	for _, user := range users {
		userProfile := models.UserProfile{
			Id:        user.Id.Hex(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		userProfiles = append(userProfiles, userProfile)
	}

	return userProfiles, nil
}

func (uc *userUsecase) UpdateUser(ctx context.Context, id string, req models.CreateUserReq) (*models.UserProfile, error) {
	user := models.User{
		Email:     req.Email,
		Name:      req.Name,
		Password:  req.Password, // Note: In a real application, you should hash passwords before storing them
		UpdatedAt: time.Now(),
	}

	updatedUser, err := uc.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, err
	}

	userProfile := &models.UserProfile{
		Id:        updatedUser.Id.Hex(),
		Email:     updatedUser.Email,
		Name:      updatedUser.Name,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	return userProfile, nil
}

func (uc *userUsecase) DeleteUser(ctx context.Context, id string) error {
	return uc.userRepo.DeleteUser(ctx, id)
}
