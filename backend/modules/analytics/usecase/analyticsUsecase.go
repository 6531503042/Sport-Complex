package usecase

import (
	"context"
	"main/modules/analytics"
	"main/modules/analytics/repository"
	"time"
)

type (
    AnalyticsUsecaseService interface {
        GetDashboardMetrics(ctx context.Context, facilityName string, timeRange string) (*analytics.DashboardMetrics, error)
        GetUserAnalytics(ctx context.Context) (*analytics.UserMetrics, error)
    }

    analyticsUsecase struct {
        analyticsRepo repository.AnalyticsRepositoryService
    }
)

func NewAnalyticsUsecase(analyticsRepo repository.AnalyticsRepositoryService) AnalyticsUsecaseService {
    return &analyticsUsecase{
        analyticsRepo: analyticsRepo,
    }
}

func (u *analyticsUsecase) GetDashboardMetrics(ctx context.Context, facilityName string, timeRange string) (*analytics.DashboardMetrics, error) {
    now := time.Now()
    var startDate time.Time
    
    switch timeRange {
    case "daily":
        startDate = now.AddDate(0, 0, -1)
    case "weekly":
        startDate = now.AddDate(0, 0, -7)
    case "monthly":
        startDate = now.AddDate(0, -1, 0)
    case "yearly":
        startDate = now.AddDate(-1, 0, 0)
    default:
        startDate = now.AddDate(0, 0, -7) // Default to weekly
    }

    query := &analytics.AnalyticsQuery{
        StartDate:  startDate.Format("2006-01-02"),
        EndDate:    now.Format("2006-01-02"),
        Period:     timeRange,
        FacilityID: facilityName,
    }

    return u.analyticsRepo.GetDashboardMetrics(ctx, query)
}

func (u *analyticsUsecase) GetUserAnalytics(ctx context.Context) (*analytics.UserMetrics, error) {
    // Use current time for the date range
    now := time.Now()
    startDate := now.AddDate(0, -1, 0) // Last month
    endDate := now

    return u.analyticsRepo.GetUserMetrics(ctx, startDate, endDate)
}
