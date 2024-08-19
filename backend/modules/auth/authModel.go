package auth

import "time"

type (
	UserLoginReq struct {
		Email    string `json:"email" form:"email" validate: "required, email, max=255"`
		Password string `json:"password" form:"password" validate: "required, max=32"`
	}

	UserLoginRes struct {
		CredentialId string `json:"credential_id" form:"credential_id" validate: "required, max=60"`
	}

	RefreshTokenReq struct {
		CredentialId string `json:"credential_id" form:"credential_id" validate: "required, max=60"`
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate: "required, max=500"`
	}

	RefreshTokenRes struct {
		Id string `json:_id`
		UserId string `json:"user_id"`
		RoleCode []int `json:"role_code"`
		RefreshToken string `json:"refresh_token"`
		AccessToken string `json:"access_token" validate:"required"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	InsertUserRole struct {
		UserId string `json:"user_id" validate:"required"`
		Role  []int `json:"role_code" validate:"required"`
	}

	// ProfileIntercepter struct {
	// 	*user.UserProfile

	// }


)