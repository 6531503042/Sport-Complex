package usecase

import (
	"context"
	"errors"
	"main/modules/user"
	"main/modules/user/repository"

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

func (u *userUsecase) createUser (pctx context.Context, req *user.CreateUserReq) (*user.PlayerProfile, error) {
	if !u.userRepository.IsUniqueUser(pctx, req.Email, req.Name) {
		return nil, errors.New("email or name is already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error: failed to hash password")
	}

	//Insert one user
	userId, err := u.userRepository.InsertOneUser(pctx, &user.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
		// CreatedAt: utils.Lo,
		UserRole: []user.UserRole{
			{
				RoleTitle: "user",
				RoleCode:  0,
			}
		}
	})
	}) 

}