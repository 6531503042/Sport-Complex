package analytics

import (
	"time"
)

type (
    FacilityUsageStats struct {
        FacilityName    string    `json:"facility_name"`
        Date           time.Time `json:"date"`
        TotalBookings  int       `json:"total_bookings"`
        CompletedBookings int    `json:"completed_bookings"`
        FailedBookings   int       `json:"failed_bookings"`
        Revenue        float64   `json:"revenue"`
        PeakHours      []int    `json:"peak_hours"`
        UtilizationRate float64  `json:"utilization_rate"` // percentage of capacity used
    }

    FacilityComparison struct {
        FacilityName    string    `json:"facility_name"`
        BookingCount    int       `json:"booking_count"`
        Revenue         float64   `json:"revenue"`
        PopularityRank  int       `json:"popularity_rank"`
    }

    TimeSeriesData struct {
        Labels []string  `json:"labels"`   // Time periods (dates/months)
        Values []float64 `json:"values"`   // Corresponding values
    }

    DashboardMetrics struct {
        DailyStats struct {
            TotalBookings    TimeSeriesData `json:"total_bookings"`
            Revenue         TimeSeriesData `json:"revenue"`
            UtilizationRate TimeSeriesData `json:"utilization_rate"`
        } `json:"daily_stats"`
        
        WeeklyStats struct {
            BookingTrends    TimeSeriesData `json:"booking_trends"`
            RevenueByFacility map[string]float64 `json:"revenue_by_facility"`
            PeakDays         map[string]int    `json:"peak_days"`
        } `json:"weekly_stats"`
        
        MonthlyStats struct {
            BookingGrowth    TimeSeriesData `json:"booking_growth"`
            RevenueGrowth    TimeSeriesData `json:"revenue_growth"`
            FacilityComparison []FacilityComparison `json:"facility_comparison"`
        } `json:"monthly_stats"`
        
        YearlyStats struct {
            AnnualRevenue    TimeSeriesData `json:"annual_revenue"`
            FacilityTrends   map[string]TimeSeriesData `json:"facility_trends"`
            SeasonalPatterns map[string]float64 `json:"seasonal_patterns"`
        } `json:"yearly_stats"`
    }

    UserMetrics struct {
        TotalUsers      int     `json:"total_users"`
        ActiveUsers     int     `json:"active_users"`
        NewUsers        int     `json:"new_users"`
        UserGrowthRate  float64 `json:"user_growth_rate"`
    }
) 