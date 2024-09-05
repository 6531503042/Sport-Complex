package repository

import (
	"context"
	"main/config"
	"main/modules/auth"
	userPb "main/modules/user/proto"
	"main/pkg/jwt"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AuthRepositoryMock struct {
		mock.Mock
	}
)

func (m * AuthRepositoryMock) CredentialSearch(pctx context.Context, grpcUrl string, req *userPb.CredentialSearchReq) (*userPb.UserProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	return args.Get(0).(*userPb.UserProfile), args.Error(1)
}

func (m *AuthRepositoryMock) InsertOneUserCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	args := m.Called(pctx, req)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
} 

func (m *AuthRepositoryMock) AccessToken(cfg *config.Config, claims *jwt.Claims) string {
	args := m.Called(cfg, claims)
	return args.String(0)
}

func (m *AuthRepositoryMock) RefreshToken(cfg *config.Config, claims *jwt.Claims) string {
	args := m.Called(cfg, claims)
	return args.String(0)
}

func (m *AuthRepositoryMock) FindOneUserCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	args := m.Called(pctx, credentialId)
	return args.Get(0).(*auth.Credential), args.Error(1)
}

func (m *AuthRepositoryMock) FindOneUserProfileToRefresh(pctx context.Context, grpcUrl string, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
	return nil, nil
}

func (m *AuthRepositoryMock) UpdateOneUserCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error {
	return nil
}

func (m *AuthRepositoryMock) DeleteOneUserCredential(pctx context.Context, credentialId string) (int64, error) {
	return 0, nil
}

func (m *AuthRepositoryMock) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	return nil, nil
}

func (m *AuthRepositoryMock) RolesCount(pctx context.Context) (int64, error) {
	return 0, nil
}