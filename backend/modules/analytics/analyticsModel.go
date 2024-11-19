package analytics

import "time"

type (
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
        TotalUsers        int     `json:"total_users"`
        UserGrowthRate    float64 `json:"user_growth_rate"`
        NewUsersThisMonth int     `json:"new_users_this_month"`
        ActiveUsers       int     `json:"active_users"`
        UserRetentionRate float64 `json:"user_retention_rate"`
    }

    BookingMetrics struct {
        TotalBookings          int                `json:"total_bookings"`
        BookingGrowthRate      float64            `json:"booking_growth_rate"`
        BookingsPerFacility    map[string]int     `json:"bookings_per_facility"`
        PopularTimeSlots       []TimeSlotMetric   `json:"popular_time_slots"`
        BookingStatusCount     map[string]int     `json:"booking_status_count"`
        AverageBookingsPerDay  float64            `json:"average_bookings_per_day"`
    }

    RevenueMetrics struct {
        TotalRevenue          float64                 `json:"total_revenue"`
        RevenueGrowthRate     float64                 `json:"revenue_growth_rate"`
        RevenuePerFacility    map[string]float64      `json:"revenue_per_facility"`
        MonthlyRevenue        []MonthlyRevenueMetric  `json:"monthly_revenue"`
        PaymentMethodStats    map[string]PaymentStats `json:"payment_method_stats"`
    }

    FacilityMetrics struct {
        TotalFacilities       int                     `json:"total_facilities"`
        FacilityUtilization   map[string]float64      `json:"facility_utilization"`
        PeakHours             []PeakHourMetric        `json:"peak_hours"`
        PopularFacilities     []FacilityUsageMetric   `json:"popular_facilities"`
        MaintenanceSchedule   []MaintenanceMetric     `json:"maintenance_schedule"`
    }

    TimeSeriesData struct {
        Daily   []MetricPoint `json:"daily"`
        Weekly  []MetricPoint `json:"weekly"`
        Monthly []MetricPoint `json:"monthly"`
        Yearly  []MetricPoint `json:"yearly"`
    }

    FacilityUsageStats struct {
        FacilityName      string    `json:"facility_name"`
        Date             time.Time `json:"date"`
        TotalBookings    int       `json:"total_bookings"`
        CompletedBookings int      `json:"completed_bookings"`
        FailedBookings   int       `json:"failed_bookings"`
        Revenue          float64   `json:"revenue"`
        PeakHours        []int     `json:"peak_hours"`
        UtilizationRate  float64   `json:"utilization_rate"`
    }

    FacilityComparison struct {
        FacilityName    string    `json:"facility_name"`
        BookingCount    int       `json:"booking_count"`
        Revenue         float64   `json:"revenue"`
        PopularityRank  int       `json:"popularity_rank"`
    }

    // Supporting types
    TimeSlotMetric struct {
        TimeSlot string  `json:"time_slot"`
        Count    int     `json:"count"`
        Trend    float64 `json:"trend"`
    }

    MonthlyRevenueMetric struct {
        Month    string  `json:"month"`
        Revenue  float64 `json:"revenue"`
        Growth   float64 `json:"growth"`
    }

	PaymentMetrics struct {
		TotalPayments int `json:"total_payments"`
		PaymentGrowthRate float64 `json:"payment_growth_rate"`
		PaymentMethods map[string]PaymentStats `json:"payment_methods"`
	}

    PaymentStats struct {
        Count       int
        TotalAmount float64
        SuccessRate float64
    }

    PeakHourMetric struct {
        Hour        int     `json:"hour"`
        Utilization float64 `json:"utilization"`
    }

    FacilityUsageMetric struct {
        FacilityName    string  `json:"facility_name"`
        UsageRate       float64 `json:"usage_rate"`
        Revenue         float64 `json:"revenue"`
    }

    MaintenanceMetric struct {
        FacilityName string `json:"facility_name"`
        NextDate     string `json:"next_date"`
        Status       string `json:"status"`
    }

    MetricPoint struct {
        Date     string             `json:"date"`
        Metrics  map[string]float64 `json:"metrics"`
    }

    AnalyticsQuery struct {
        StartDate   string `query:"start_date" validate:"required"`
        EndDate     string `query:"end_date" validate:"required"`
        Period      string `query:"period" validate:"required,oneof=daily weekly monthly yearly"`
        FacilityID  string `query:"facility_id"`
    }
) 