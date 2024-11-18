package usecase

import (
	"context"
	"errors"
	"log"
	"main/modules/auth"
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
        UpdateOneUser(ctx context.Context, userId string, updateFields map[string]interface{}) error
        FindManyUser(pctx context.Context) ([]user.User, error)
        DeleteUser(ctx context.Context, userId string) error
        GetUserAnalytics(pctx context.Context, query *user.AnalyticsQuery) (*user.UserAnalytics, error)
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
        return nil, errors.New("error: failed to hash password")
    }

    roleCode := req.RoleCode
    if roleCode == 0 {
        roleCode = auth.RoleUser // Default to user role
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
                RoleTitle: getRoleTitle(roleCode),
                RoleCode:  roleCode,
            },
        },
    })
    if err != nil {
        return nil, errors.New("error: failed to create user")
    }

    return u.FindOneUserProfile(pctx, userId.Hex())
}

// Helper function to get role title
func getRoleTitle(roleCode int) string {
    switch roleCode {
    case auth.RoleAdmin:
        return "admin"
    default:
        return "user"
    }
}

func (u *userUsecase) UpdateOneUser(ctx context.Context, userId string, updateFields map[string]interface{}) error {
    // Ensure that the user exists before attempting to update
    if _, err := u.userRepository.FindOneUserProfile(ctx, userId); err != nil {
        return err
    }

    // If role_title is being updated, set the corresponding role_code
    if roleTitle, ok := updateFields["role_title"].(string); ok {
        roleCode := 0
        if roleTitle == "admin" {
            roleCode = 1
        }
        
        // Update user_roles array
        updateFields["user_roles"] = []user.UserRole{
            {
                RoleTitle: roleTitle,
                RoleCode:  roleCode,
            },
        }
        
        // Remove the individual role_title field as we're using user_roles array
        delete(updateFields, "role_title")
    }

    updateFields["updated_at"] = utils.LocalTime()
    return u.userRepository.UpdateOneUser(ctx, userId, updateFields)
}

func (u *userUsecase) DeleteUser(ctx context.Context, userId string) error {
	// Ensure that the user exists before attempting to delete
	_, err := u.userRepository.FindOneUserProfile(ctx, userId)
	if err != nil {
		return err
	}

	return u.userRepository.DeleteOneUser(ctx, userId)
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
        return nil, errors.New("error: password is invalid")
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

func (u *userUsecase) FindManyUser(pctx context.Context) ([]user.User, error) {
    users, err := u.userRepository.FindManyUser(pctx)
    if err != nil {
        return nil, err
    }

    // No transformation needed, return the raw user data
    return users, nil
}

func (u *userUsecase) GetUserAnalytics(pctx context.Context, query *user.AnalyticsQuery) (*user.UserAnalytics, error) {
    startDate, err := time.Parse("2006-01-02", query.StartDate)
    if err != nil {
        return nil, errors.New("invalid start date format")
    }

    endDate, err := time.Parse("2006-01-02", query.EndDate)
    if err != nil {
        return nil, errors.New("invalid end date format")
    }

    if endDate.Before(startDate) {
        return nil, errors.New("end date must be after start date")
    }

    return u.userRepository.GetUserAnalytics(pctx, query.Period, startDate, endDate)
}