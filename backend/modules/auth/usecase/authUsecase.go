package usecase

import (
	"context"
	"errors"
	"log"
	"main/config"
	"main/modules/auth"
	authPb "main/modules/auth/proto"
	"main/modules/auth/repository"
	"main/modules/user"
	userPb "main/modules/user/proto"
	"main/pkg/jwt"
	"main/pkg/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseService interface {
	Login(pctx context.Context, cfg *config.Config, req *auth.UserLoginReq) (*auth.ProfileIntercepter, error)
	RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error)
	AccessTokenSearch(pctx context.Context, accessToken string) (*authPb.AccessTokenSearchRes, error)
	Logout(pctx context.Context, credentialId string) (int64, error)
	RolesCount(pctx context.Context) (*authPb.RolesCountRes, error)
}

type authUsecase struct {
	authRepository repository.AuthRepositoryService
}

func NewAuthUsecase(authRepository repository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{authRepository}
}

// ComparePasswords compares a bcrypt hashed password with its possible plaintext equivalent.
func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Login handles the user login process
func (u *authUsecase) Login(pctx context.Context, cfg *config.Config, req *auth.UserLoginReq) (*auth.ProfileIntercepter, error) {
	profile, err := u.authRepository.CredentialSearch(pctx, cfg.Grpc.UserUrl, &userPb.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Error during credential search: %v", err)
		return nil, err
	}

	profile.Id = "user:" + profile.Id

	accessToken := u.authRepository.AccessToken(cfg, &jwt.Claims{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	refreshToken := u.authRepository.RefreshToken(cfg, &jwt.Claims{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	credentialId, err := u.authRepository.InsertOneUserCredential(pctx, &auth.Credential{
		UserId:      profile.Id,
		RoleCode:    int(profile.RoleCode),
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	credential, err := u.authRepository.FindOneUserCredential(pctx, credentialId.Hex())
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		UserProfile: &user.UserProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Name:      profile.Name,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			Id:           credentialId.Hex(),
			UserId:       credential.UserId,
			RoleCode:     credential.RoleCode,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}


func (u *authUsecase) RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error) {
	claims, err := jwt.ParseToken(cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: RefreshToken: %s", err.Error())
		return nil, errors.New(err.Error())
	}

	profile, err := u.authRepository.FindOneUserProfileToRefresh(pctx, cfg.Grpc.UserUrl, &userPb.FindOneUserProfileToRefreshReq{
		UserId: strings.TrimPrefix(claims.UserId, "user:"),
	})
	if err != nil {
		return nil, err
	}

	accessToken := jwt.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwt.Claims{
		UserId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwt.ReloadToken(cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwt.Claims{
		UserId: profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	if err := u.authRepository.UpdateOneUserCredential(pctx, req.CredentialId, &auth.UpdateRefreshTokenReq{
		UserId: profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	credential, err := u.authRepository.FindOneUserCredential(pctx, req.CredentialId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		UserProfile: &user.UserProfile{
			Id:        "player:" + profile.Id,
			Email:     profile.Email,
			Name:  profile.Name,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			UserId:     credential.UserId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUsecase) Logout(pctx context.Context, credentialId string) (int64, error) {
	return u.authRepository.DeleteOneUserCredential(pctx, credentialId)
}

func (u *authUsecase) AccessTokenSearch(pctx context.Context, accessToken string) (*authPb.AccessTokenSearchRes, error) {
	credential, err := u.authRepository.FindOneAccessToken(pctx, accessToken)
	if err != nil || credential == nil {
		return &authPb.AccessTokenSearchRes{
			IsValid: false,
		}, errors.New("error: access token is invalid")
	}

	return &authPb.AccessTokenSearchRes{
		IsValid: true,
	}, nil
}

func (u *authUsecase) RolesCount(pctx context.Context) (*authPb.RolesCountRes, error) {
	result, err := u.authRepository.RolesCount(pctx)
	if err != nil {
		return nil, err
	}

	return &authPb.RolesCountRes{
		Count: result,
	}, nil
}
