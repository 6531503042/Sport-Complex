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

func (u *authUsecase) Login(pctx context.Context, cfg *config.Config, req *auth.UserLoginReq) (*auth.ProfileIntercepter, error) {
    log.Printf("Attempting login for email: %s", req.Email)

    // Call to CredentialSearch gRPC method
    profile, err := u.authRepository.CredentialSearch(pctx, cfg.Grpc.UserUrl, &userPb.CredentialSearchReq{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        log.Printf("Error during credential search: %v", err)
        return nil, err
    }

    log.Printf("Received profile: %+v", profile)
    profile.Id = "user:" + profile.Id

    // Generate tokens
    accessToken, err := u.authRepository.AccessToken(cfg, &jwt.Claim{
        UserId:   profile.Id,
        RoleCode: int(profile.RoleCode),
    })
    if err != nil {
        log.Printf("Error generating access token: %v", err)
        return nil, err
    }

    refreshToken, err := u.authRepository.RefreshToken(cfg, &jwt.Claim{
        UserId:   profile.Id,
        RoleCode: int(profile.RoleCode),
    })
    if err != nil {
        log.Printf("Error generating refresh token: %v", err)
        return nil, err
    }

    log.Printf("Generated AccessToken: [REDACTED], RefreshToken: [REDACTED]")

    // Prepare credential object
    credential := &auth.Credential{
        UserId:      profile.Id,
        RoleCode:    int(profile.RoleCode),
        AccessToken: accessToken,
        RefreshToken: refreshToken,
        CreatedAt:   utils.LocalTime(),
        UpdatedAt:   utils.LocalTime(),
    }

    // Insert credential into database
    credentialId, err := u.authRepository.InsertOneUserCredential(pctx, credential)
    if err != nil {
        log.Printf("Error inserting credential: %v", err)
        return nil, err
    }

    log.Printf("Credential inserted with ID: %s", credentialId.Hex())

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
		log.Printf("Error: RefreshToken : %s", err.Error())
		return nil, errors.New("failed to parse refresh token")
	}

	profile, err := u.authRepository.FindOneUserProfileToRefresh(pctx, cfg.Grpc.UserUrl, &userPb.FindOneUserProfileToRefreshReq{
		UserId: strings.TrimPrefix(claims.UserId, "user:"),
	})
	if err != nil {
		return nil, errors.New("failed to fetch user profile")
	}

	accessToken, err := jwt.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwt.Claim{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwt.Claim{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()
	if err != nil {
		return nil, err
	}

	if err := u.authRepository.UpdateOneUserCredential(pctx, req.CredentialId, &auth.UpdateRefreshTokenReq{
		UserId:       profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, errors.New("failed to update user credentials")
	}

	credential, err := u.authRepository.FindOneUserCredential(pctx, req.CredentialId)
	if err != nil {
		return nil, errors.New("failed to fetch updated user credential")
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
			Id:           credential.Id.Hex(),
			UserId:       profile.Id,
			RoleCode:     int(profile.RoleCode),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
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
