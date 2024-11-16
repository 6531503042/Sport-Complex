package user

import "time"

type (
	UserProfile struct {
		Id        string    `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		RoleCode  int       `json:"role_code"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserClaims struct {
		Id       string `json:"id"`
		RoleCode int    `json:"role_code"`
	}

	CreateUserReq struct {
		Email    string `json:"email" form:"email" validate:"required,email,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=32"`
		Name     string `json:"name" form:"name" validate:"required,max=32"`
		RoleCode int    `json:"role_code" form:"role_code"`
	}

	UserAnalytics struct {
		TotalUsers      int64                `json:"total_users"`
		ActiveUsers     int64                `json:"active_users"`
		UsersByPeriod   []UserActivityPeriod `json:"users_by_period"`
		PeakActivityDay PeakActivity         `json:"peak_activity_day"`
	}

	UserActivityPeriod struct {
		Period    string `json:"period"`
		Count     int64  `json:"count"`
		Active    int64  `json:"active"`
		Inactive  int64  `json:"inactive"`
	}

	PeakActivity struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	AnalyticsQuery struct {
		Period    string `query:"period" validate:"required,oneof=daily weekly monthly yearly"`
		StartDate string `query:"start_date" validate:"required"`
		EndDate   string `query:"end_date" validate:"required"`
	}
)
