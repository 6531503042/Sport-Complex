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
    metrics := &analytics.DashboardMetrics{
        DailyStats: struct {
            TotalBookings    analytics.TimeSeriesData `json:"total_bookings"`
            Revenue         analytics.TimeSeriesData `json:"revenue"`
            UtilizationRate analytics.TimeSeriesData `json:"utilization_rate"`
        }{},
        WeeklyStats: struct {
            BookingTrends    analytics.TimeSeriesData `json:"booking_trends"`
            RevenueByFacility map[string]float64 `json:"revenue_by_facility"`
            PeakDays         map[string]int    `json:"peak_days"`
        }{
            RevenueByFacility: make(map[string]float64),
            PeakDays:         make(map[string]int),
        },
    }

    facilities := []string{"fitness", "swimming", "badminton", "football"}
    revenueByFacility := make(map[string]float64)
    
    // Get data for each facility
    for _, facility := range facilities {
        dailyStats, err := u.analyticsRepo.GetDailyStats(ctx, facility, now)
        if err != nil {
            continue // Skip if error, but don't fail completely
        }
        
        // Accumulate revenue for each facility
        revenueByFacility[facility] = float64(dailyStats.CompletedBookings) * 100 // Example price
        
        // Track peak days
        if dailyStats.TotalBookings > metrics.WeeklyStats.PeakDays[now.Weekday().String()] {
            metrics.WeeklyStats.PeakDays[now.Weekday().String()] = dailyStats.TotalBookings
        }
    }

    // Set daily stats
    todayStats, err := u.analyticsRepo.GetDailyStats(ctx, facilityName, now)
    if err == nil {
        metrics.DailyStats.TotalBookings = analytics.TimeSeriesData{
            Labels: []string{now.Format("2006-01-02")},
            Values: []float64{float64(todayStats.TotalBookings)},
        }
        metrics.DailyStats.Revenue = analytics.TimeSeriesData{
            Labels: []string{now.Format("2006-01-02")},
            Values: []float64{float64(todayStats.CompletedBookings) * 100}, // Example price
        }
        metrics.DailyStats.UtilizationRate = analytics.TimeSeriesData{
            Labels: []string{now.Format("2006-01-02")},
            Values: []float64{todayStats.UtilizationRate},
        }
    }

    // Set weekly stats
    weeklyData := make([]float64, 7)
    weekLabels := make([]string, 7)
    
    for i := 6; i >= 0; i-- {
        date := now.AddDate(0, 0, -i)
        stats, err := u.analyticsRepo.GetDailyStats(ctx, facilityName, date)
        if err != nil {
            continue
        }
        weeklyData[6-i] = float64(stats.TotalBookings)
        weekLabels[6-i] = date.Format("Mon")
    }

    metrics.WeeklyStats.BookingTrends = analytics.TimeSeriesData{
        Labels: weekLabels,
        Values: weeklyData,
    }
    metrics.WeeklyStats.RevenueByFacility = revenueByFacility

    return metrics, nil
}

func (u *analyticsUsecase) GetUserAnalytics(ctx context.Context) (*analytics.UserMetrics, error) {
    return u.analyticsRepo.GetUserMetrics(ctx)
}
