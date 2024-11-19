package repository

import (
	"context"
	"fmt"
	"log"
	"main/modules/analytics"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsRepositoryService interface {
	GetDashboardMetrics(ctx context.Context, query *analytics.AnalyticsQuery) (*analytics.DashboardMetrics, error)
	GetUserMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.UserMetrics, error)
	GetBookingMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.BookingMetrics, error)
	GetRevenueMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.RevenueMetrics, error)
	GetFacilityMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.FacilityMetrics, error)
	GetTimeSeriesData(ctx context.Context, period string, startDate, endDate time.Time) (*analytics.TimeSeriesData, error)
	GetDailyStats(ctx context.Context, facilityName string, date time.Time) (*analytics.FacilityUsageStats, error)

    //payment
    getTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error)
}

type analyticsRepository struct {
	db *mongo.Client
}

func NewAnalyticsRepository(db *mongo.Client) AnalyticsRepositoryService {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetDashboardMetrics(ctx context.Context, query *analytics.AnalyticsQuery) (*analytics.DashboardMetrics, error) {
	log.Printf("Getting dashboard metrics for facility: %s, period: %s", query.FacilityID, query.Period)
	
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		log.Printf("Error parsing start date: %v", err)
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	if err := validateTimeRange(startDate, endDate); err != nil {
		return nil, err
	}

	// Initialize metrics
	metrics := &analytics.DashboardMetrics{}

	// Get daily stats
	dailyStats, err := r.getDailyStats(ctx, query.FacilityID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting daily stats: %w", err)
	}
	metrics.DailyStats.TotalBookings = dailyStats
	metrics.DailyStats.Revenue = r.getRevenueTimeSeries(ctx, startDate, endDate, "daily")
	metrics.DailyStats.UtilizationRate = r.getUtilizationTimeSeries(ctx, startDate, endDate, "daily")

	// Get weekly stats
	weeklyStats, err := r.getWeeklyStats(ctx, query.FacilityID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting weekly stats: %w", err)
	}
	metrics.WeeklyStats.BookingTrends = weeklyStats
	metrics.WeeklyStats.RevenueByFacility = r.getRevenueByFacility(ctx, startDate, endDate)
	metrics.WeeklyStats.PeakDays = r.getPeakDays(ctx, startDate, endDate)

	// Get monthly stats
	monthlyStats, err := r.getMonthlyStats(ctx, query.FacilityID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting monthly stats: %w", err)
	}
	metrics.MonthlyStats.BookingGrowth = monthlyStats
	metrics.MonthlyStats.RevenueGrowth = r.getRevenueGrowth(ctx, startDate, endDate)
	metrics.MonthlyStats.FacilityComparison = r.getFacilityComparison(ctx, startDate, endDate)

	// Get yearly stats
	metrics.YearlyStats.AnnualRevenue = r.getAnnualRevenue(ctx, startDate, endDate)
	metrics.YearlyStats.FacilityTrends = r.getFacilityTrends(ctx, startDate, endDate)
	metrics.YearlyStats.SeasonalPatterns = r.getSeasonalPatterns(ctx, startDate, endDate)

	return metrics, nil
}

// Helper methods for getting specific metrics
func (r *analyticsRepository) getDailyStats(ctx context.Context, facilityID string, startDate, endDate time.Time) (analytics.TimeSeriesData, error) {
	bookingDb := r.db.Database("booking_db")
	bookingCol := bookingDb.Collection("booking_transaction")
	paymentDb := r.db.Database("payment_db")
	paymentCol := paymentDb.Collection("payments")

	// Get daily bookings
	bookingPipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
				"facility": facilityID,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$created_at",
					},
				},
				"total_bookings": bson.M{"$sum": 1},
				"completed_bookings": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "completed"}},
							1,
							0,
						},
					},
				},
				"failed_bookings": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "failed"}},
							1,
							0,
						},
					},
				},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	// Get daily revenue
	revenuePipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
				"facility_name": facilityID,
				"status": "COMPLETED",
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$created_at",
					},
				},
				"revenue": bson.M{"$sum": "$amount"},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	// Execute booking pipeline
	bookingCursor, err := bookingCol.Aggregate(ctx, bookingPipeline)
	if err != nil {
		return analytics.TimeSeriesData{}, fmt.Errorf("failed to aggregate bookings: %w", err)
	}
	defer bookingCursor.Close(ctx)

	// Execute revenue pipeline
	revenueCursor, err := paymentCol.Aggregate(ctx, revenuePipeline)
	if err != nil {
		return analytics.TimeSeriesData{}, fmt.Errorf("failed to aggregate revenue: %w", err)
	}
	defer revenueCursor.Close(ctx)

	// Process results
	dailyMetrics := make(map[string]map[string]float64)

	// Process booking results
	for bookingCursor.Next(ctx) {
		var result struct {
			Date              string `bson:"_id"`
			TotalBookings     int    `bson:"total_bookings"`
			CompletedBookings int    `bson:"completed_bookings"`
		}
		if err := bookingCursor.Decode(&result); err != nil {
			return analytics.TimeSeriesData{}, fmt.Errorf("failed to decode booking result: %w", err)
		}

		if dailyMetrics[result.Date] == nil {
			dailyMetrics[result.Date] = make(map[string]float64)
		}
		dailyMetrics[result.Date]["total_bookings"] = float64(result.TotalBookings)
		dailyMetrics[result.Date]["completed_bookings"] = float64(result.CompletedBookings)
		
		// Calculate utilization rate
		if result.TotalBookings > 0 {
			utilizationRate := float64(result.CompletedBookings) / float64(result.TotalBookings) * 100
			dailyMetrics[result.Date]["utilization_rate"] = utilizationRate
		}
	}

	// Process revenue results
	for revenueCursor.Next(ctx) {
		var result struct {
			Date    string  `bson:"_id"`
			Revenue float64 `bson:"revenue"`
		}
		if err := revenueCursor.Decode(&result); err != nil {
			return analytics.TimeSeriesData{}, fmt.Errorf("failed to decode revenue result: %w", err)
		}

		if dailyMetrics[result.Date] == nil {
			dailyMetrics[result.Date] = make(map[string]float64)
		}
		dailyMetrics[result.Date]["revenue"] = result.Revenue
	}

	// Convert to TimeSeriesData
	var daily []analytics.MetricPoint
	for date, metrics := range dailyMetrics {
		daily = append(daily, analytics.MetricPoint{
			Date:    date,
			Metrics: metrics,
		})
	}

	// Sort daily metrics by date
	sort.Slice(daily, func(i, j int) bool {
		return daily[i].Date < daily[j].Date
	})

	return analytics.TimeSeriesData{
		Daily:   daily,
		Weekly:  aggregateToWeekly(daily),
		Monthly: aggregateToMonthly(daily),
		Yearly:  aggregateToYearly(daily),
	}, nil
}

func (r *analyticsRepository) getRevenueTimeSeries(ctx context.Context, startDate, endDate time.Time, period string) analytics.TimeSeriesData {
	// Implementation for revenue time series
	return analytics.TimeSeriesData{}
}

func (r *analyticsRepository) getUtilizationTimeSeries(ctx context.Context, startDate, endDate time.Time, period string) analytics.TimeSeriesData {
	// Implementation for utilization time series
	return analytics.TimeSeriesData{}
}

func (r *analyticsRepository) getWeeklyStats(ctx context.Context, facilityID string, startDate, endDate time.Time) (analytics.TimeSeriesData, error) {
	// Implementation for weekly stats
	return analytics.TimeSeriesData{}, nil
}

func (r *analyticsRepository) getRevenueByFacility(ctx context.Context, startDate, endDate time.Time) map[string]float64 {
	// Implementation for revenue by facility
	return make(map[string]float64)
}

func (r *analyticsRepository) getPeakDays(ctx context.Context, startDate, endDate time.Time) map[string]int {
	// Implementation for peak days
	return make(map[string]int)
}

func (r *analyticsRepository) getMonthlyStats(ctx context.Context, facilityID string, startDate, endDate time.Time) (analytics.TimeSeriesData, error) {
	// Implementation for monthly stats
	return analytics.TimeSeriesData{}, nil
}

func (r *analyticsRepository) getRevenueGrowth(ctx context.Context, startDate, endDate time.Time) analytics.TimeSeriesData {
	// Implementation for revenue growth
	return analytics.TimeSeriesData{}
}

func (r *analyticsRepository) getFacilityComparison(ctx context.Context, startDate, endDate time.Time) []analytics.FacilityComparison {
	// Implementation for facility comparison
	return []analytics.FacilityComparison{}
}

func (r *analyticsRepository) getAnnualRevenue(ctx context.Context, startDate, endDate time.Time) analytics.TimeSeriesData {
	// Implementation for annual revenue
	return analytics.TimeSeriesData{}
}

func (r *analyticsRepository) getFacilityTrends(ctx context.Context, startDate, endDate time.Time) map[string]analytics.TimeSeriesData {
	// Implementation for facility trends
	return make(map[string]analytics.TimeSeriesData)
}

func (r *analyticsRepository) getSeasonalPatterns(ctx context.Context, startDate, endDate time.Time) map[string]float64 {
	// Implementation for seasonal patterns
	return make(map[string]float64)
}

func (r *analyticsRepository) GetUserMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.UserMetrics, error) {
	userDb := r.db.Database("user_db")
	userCol := userDb.Collection("users")

	// Get total users
	totalUsers, err := userCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Get new users this month
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	newUsers, err := userCol.CountDocuments(ctx, bson.M{
		"created_at": bson.M{"$gte": firstDayOfMonth},
	})
	if err != nil {
		return nil, err
	}

	// Calculate growth rate
	previousMonth := firstDayOfMonth.AddDate(0, -1, 0)
	previousMonthUsers, err := userCol.CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": previousMonth,
			"$lt":  firstDayOfMonth,
		},
	})
	if err != nil {
		return nil, err
	}

	var growthRate float64
	if previousMonthUsers > 0 {
		growthRate = float64(newUsers-previousMonthUsers) / float64(previousMonthUsers) * 100
	}

	return &analytics.UserMetrics{
		TotalUsers:        int(totalUsers),
		UserGrowthRate:    growthRate,
		NewUsersThisMonth: int(newUsers),
		ActiveUsers:       int(totalUsers), // You might want to refine this based on your definition of active users
		UserRetentionRate: 0,              // Calculate based on your retention logic
	}, nil
}

func (r *analyticsRepository) GetBookingMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.BookingMetrics, error) {
	bookingDb := r.db.Database("booking_db")
	bookingCol := bookingDb.Collection("booking_transaction")

	// Get total bookings
	totalBookings, err := bookingCol.CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	})
	if err != nil {
		return nil, err
	}

	// Get bookings per facility
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$facility",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	bookingsPerFacility := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		bookingsPerFacility[result.ID] = result.Count
	}

	// Calculate booking growth rate
	previousPeriodStart := startDate.AddDate(0, -1, 0)
	previousPeriodEnd := endDate.AddDate(0, -1, 0)
	previousBookings, err := bookingCol.CountDocuments(ctx, bson.M{
		"created_at": bson.M{
			"$gte": previousPeriodStart,
			"$lte": previousPeriodEnd,
		},
	})
	if err != nil {
		return nil, err
	}

	var growthRate float64
	if previousBookings > 0 {
		growthRate = float64(totalBookings-previousBookings) / float64(previousBookings) * 100
	}

	return &analytics.BookingMetrics{
		TotalBookings:       int(totalBookings),
		BookingGrowthRate:   growthRate,
		BookingsPerFacility: bookingsPerFacility,
		// Add other metrics as needed
	}, nil
}

func (r *analyticsRepository) GetRevenueMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.RevenueMetrics, error) {
	paymentDb := r.db.Database("payment_db")
	paymentCol := paymentDb.Collection("payments")

	// Calculate total revenue
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
				"status": "COMPLETED",
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_revenue": bson.M{"$sum": "$amount"},
			},
		},
	}

	cursor, err := paymentCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalRevenue float64 `bson:"total_revenue"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	// Calculate revenue per facility
	facilityPipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
				"status": "COMPLETED",
			},
		},
		{
			"$group": bson.M{
				"_id": "$facility_name",
				"revenue": bson.M{"$sum": "$amount"},
			},
		},
	}

	facilityCursor, err := paymentCol.Aggregate(ctx, facilityPipeline)
	if err != nil {
		return nil, err
	}
	defer facilityCursor.Close(ctx)

	revenuePerFacility := make(map[string]float64)
	for facilityCursor.Next(ctx) {
		var facilityResult struct {
			ID      string  `bson:"_id"`
			Revenue float64 `bson:"revenue"`
		}
		if err := facilityCursor.Decode(&facilityResult); err != nil {
			return nil, err
		}
		revenuePerFacility[facilityResult.ID] = facilityResult.Revenue
	}

	return &analytics.RevenueMetrics{
		TotalRevenue:       result.TotalRevenue,
		RevenuePerFacility: revenuePerFacility,
		// Add other metrics as needed
	}, nil
}

func (r *analyticsRepository) GetFacilityMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.FacilityMetrics, error) {
	bookingDb := r.db.Database("booking_db")
	bookingCol := bookingDb.Collection("booking_transaction")

	// Get total facilities
	facilityDb := r.db.Database("facility_db")
	facilityCol := facilityDb.Collection("facilities")
	totalFacilities, err := facilityCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error counting facilities: %w", err)
	}

	// Get facility utilization
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$facility",
				"total_bookings": bson.M{"$sum": 1},
				"completed_bookings": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "completed"}},
							1,
							0,
						},
					},
				},
			},
		},
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating facility metrics: %w", err)
	}
	defer cursor.Close(ctx)

	facilityUtilization := make(map[string]float64)
	var peakHours []analytics.PeakHourMetric
	var popularFacilities []analytics.FacilityUsageMetric

	for cursor.Next(ctx) {
		var result struct {
			ID                string `bson:"_id"`
			TotalBookings    int    `bson:"total_bookings"`
			CompletedBookings int    `bson:"completed_bookings"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding facility metrics: %w", err)
		}

		// Calculate utilization rate
		utilization := float64(result.CompletedBookings) / float64(result.TotalBookings)
		facilityUtilization[result.ID] = utilization

		popularFacilities = append(popularFacilities, analytics.FacilityUsageMetric{
			FacilityName: result.ID,
			UsageRate:    utilization,
		})
	}

	// Get peak hours (example implementation)
	peakHours = []analytics.PeakHourMetric{
		{Hour: 9, Utilization: 0.8},
		{Hour: 17, Utilization: 0.9},
	}

	return &analytics.FacilityMetrics{
		TotalFacilities:     int(totalFacilities),
		FacilityUtilization: facilityUtilization,
		PeakHours:          peakHours,
		PopularFacilities:  popularFacilities,
	}, nil
}

func (r *analyticsRepository) GetTimeSeriesData(ctx context.Context, period string, startDate, endDate time.Time) (*analytics.TimeSeriesData, error) {
	bookingDb := r.db.Database("booking_db")
	bookingCol := bookingDb.Collection("booking_transaction")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date": "$created_at",
					},
				},
				"bookings": bson.M{"$sum": 1},
				"revenue": bson.M{"$sum": "$amount"},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating time series data: %w", err)
	}
	defer cursor.Close(ctx)

	var daily []analytics.MetricPoint
	for cursor.Next(ctx) {
		var result struct {
			Date     string  `bson:"_id"`
			Bookings int     `bson:"bookings"`
			Revenue  float64 `bson:"revenue"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding time series data: %w", err)
		}

		daily = append(daily, analytics.MetricPoint{
			Date: result.Date,
			Metrics: map[string]float64{
				"bookings": float64(result.Bookings),
				"revenue":  result.Revenue,
			},
		})
	}

	return &analytics.TimeSeriesData{
		Daily:   daily,
		Weekly:  aggregateToWeekly(daily),
		Monthly: aggregateToMonthly(daily),
		Yearly:  aggregateToYearly(daily),
	}, nil
}

func (r *analyticsRepository) GetDailyStats(ctx context.Context, facilityName string, date time.Time) (*analytics.FacilityUsageStats, error) {
	bookingDb := r.db.Database("booking_db")
	bookingCol := bookingDb.Collection("booking_transaction")

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startOfDay,
					"$lt":  endOfDay,
				},
				"facility": facilityName,
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_bookings": bson.M{"$sum": 1},
				"completed_bookings": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "completed"}},
							1,
							0,
						},
					},
				},
				"failed_bookings": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "failed"}},
							1,
							0,
						},
					},
				},
			},
		},
	}

	var result struct {
		TotalBookings     int `bson:"total_bookings"`
		CompletedBookings int     `bson:"completed_bookings"`
		FailedBookings    int     `bson:"failed_bookings"`
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	// Calculate utilization rate (example calculation)
	utilizationRate := float64(result.CompletedBookings) / 100.0 // Adjust based on your capacity definition

	return &analytics.FacilityUsageStats{
		FacilityName:      facilityName,
		Date:             date,
		TotalBookings:    result.TotalBookings,
		CompletedBookings: result.CompletedBookings,
		FailedBookings:   result.FailedBookings,
		UtilizationRate:  utilizationRate,
	}, nil
}

func aggregateToWeekly(daily []analytics.MetricPoint) []analytics.MetricPoint {
	weeklyData := make(map[string]map[string]float64)
	for _, point := range daily {
		date, _ := time.Parse("2006-01-02", point.Date)
		year, week := date.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		
		if weeklyData[weekKey] == nil {
			weeklyData[weekKey] = make(map[string]float64)
		}
		for metric, value := range point.Metrics {
			weeklyData[weekKey][metric] += value
		}
	}

	var result []analytics.MetricPoint
	for date, metrics := range weeklyData {
		result = append(result, analytics.MetricPoint{Date: date, Metrics: metrics})
	}
	return result
}

func aggregateToMonthly(daily []analytics.MetricPoint) []analytics.MetricPoint {
	monthlyData := make(map[string]map[string]float64)
	for _, point := range daily {
		date, _ := time.Parse("2006-01-02", point.Date)
		monthKey := date.Format("2006-01")
		
		if monthlyData[monthKey] == nil {
			monthlyData[monthKey] = make(map[string]float64)
		}
		for metric, value := range point.Metrics {
			monthlyData[monthKey][metric] += value
		}
	}

	var result []analytics.MetricPoint
	for date, metrics := range monthlyData {
		result = append(result, analytics.MetricPoint{Date: date, Metrics: metrics})
	}
	return result
}

func aggregateToYearly(daily []analytics.MetricPoint) []analytics.MetricPoint {
	yearlyData := make(map[string]map[string]float64)
	for _, point := range daily {
		date, _ := time.Parse("2006-01-02", point.Date)
		yearKey := date.Format("2006")
		
		if yearlyData[yearKey] == nil {
			yearlyData[yearKey] = make(map[string]float64)
		}
		for metric, value := range point.Metrics {
			yearlyData[yearKey][metric] += value
		}
	}

	var result []analytics.MetricPoint
	for date, metrics := range yearlyData {
		result = append(result, analytics.MetricPoint{Date: date, Metrics: metrics})
	}
	return result
}

// Add these helper methods
func (r *analyticsRepository) getFacilityDb(facilityName string) (*mongo.Database, error) {
	if facilityName == "" {
		return nil, fmt.Errorf("facility name cannot be empty")
	}
	
	db := r.db.Database(fmt.Sprintf("%s_facility", facilityName))
	
	// Verify database exists
	err := db.RunCommand(context.Background(), bson.D{{"ping", 1}}).Err()
	if err != nil {
		return nil, fmt.Errorf("facility database not accessible: %w", err)
	}
	
	return db, nil
}

func (r *analyticsRepository) getBookingDb() *mongo.Database {
	return r.db.Database("booking_db")
}

func (r *analyticsRepository) getPaymentDb() *mongo.Database {
	return r.db.Database("payment_db")
}

// Update getFacilityMetrics to handle facility-specific data
func (r *analyticsRepository) getFacilityMetrics(ctx context.Context, facilityName string, startDate, endDate time.Time) (*analytics.FacilityMetrics, error) {
	facilityDb, err := r.getFacilityDb(facilityName)
	if err != nil {
		return nil, err
	}
	slotsCol := facilityDb.Collection("slots")
	
	// Get slot utilization
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_slots": bson.M{"$sum": 1},
				"total_bookings": bson.M{"$sum": "$current_bookings"},
				"max_capacity": bson.M{"$sum": "$max_bookings"},
			},
		},
	}

	cursor, err := slotsCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		 TotalSlots    int `bson:"total_slots"`
		 TotalBookings int `bson:"total_bookings"`
		 MaxCapacity   int `bson:"max_capacity"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	// Calculate utilization rate
	utilizationRate := float64(result.TotalBookings) / float64(result.MaxCapacity) * 100

	// Get peak hours
	peakHoursPipeline := []bson.M{
		{
			"$match": bson.M{
				"current_bookings": bson.M{"$gt": 0},
			},
		},
		{
			"$group": bson.M{
				"_id": "$start_time",
				"utilization": bson.M{
					"$avg": bson.M{
						"$divide": []interface{}{"$current_bookings", "$max_bookings"},
					},
				},
			},
		},
		{
			"$sort": bson.M{"utilization": -1},
		},
		{
			"$limit": 5,
		},
	}

	peakHoursCursor, err := slotsCol.Aggregate(ctx, peakHoursPipeline)
	if err != nil {
		return nil, err
	}
	defer peakHoursCursor.Close(ctx)

	var peakHours []analytics.PeakHourMetric
	for peakHoursCursor.Next(ctx) {
		var peakHour struct {
			Hour        string  `bson:"_id"`
			Utilization float64 `bson:"utilization"`
		}
		if err := peakHoursCursor.Decode(&peakHour); err != nil {
			return nil, err
		}
		
		// Convert hour string to int (e.g., "09:00" -> 9)
		hourInt, _ := strconv.Atoi(strings.Split(peakHour.Hour, ":")[0])
		peakHours = append(peakHours, analytics.PeakHourMetric{
			Hour:        hourInt,
			Utilization: peakHour.Utilization * 100,
		})
	}

	return &analytics.FacilityMetrics{
		FacilityUtilization: map[string]float64{
			facilityName: utilizationRate,
		},
		PeakHours: peakHours,
	}, nil
}

// Add payment analytics
func (r *analyticsRepository) getPaymentMetrics(ctx context.Context, startDate, endDate time.Time) (*analytics.PaymentMetrics, error) {
	paymentDb := r.getPaymentDb()
	paymentsCol := paymentDb.Collection("payments")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$facilityname",
				"total_amount": bson.M{"$sum": "$amount"},
				"completed_payments": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "COMPLETED"}},
							1,
							0,
						},
					},
				},
				"pending_payments": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$status", "PENDING"}},
							1,
							0,
						},
					},
				},
				"total_payments": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := paymentsCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	revenueByFacility := make(map[string]float64)
	paymentStats := make(map[string]analytics.PaymentStats)

	for cursor.Next(ctx) {
		var result struct {
			FacilityName      string  `bson:"_id"`
			TotalAmount       float64 `bson:"total_amount"`
			CompletedPayments int     `bson:"completed_payments"`
			PendingPayments   int     `bson:"pending_payments"`
			TotalPayments     int     `bson:"total_payments"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		revenueByFacility[result.FacilityName] = result.TotalAmount
		successRate := float64(result.CompletedPayments) / float64(result.TotalPayments) * 100

		paymentStats[result.FacilityName] = analytics.PaymentStats{
			Count:       result.TotalPayments,
			TotalAmount: result.TotalAmount,
			SuccessRate: successRate,
		}
	}

	var totalPayments int
	var paymentGrowthRate float64

	// Calculate totals from paymentStats
	for _, stats := range paymentStats {
		totalPayments += stats.Count
	}

	// Calculate growth rate if needed
	// ... growth rate calculation ...

	return &analytics.PaymentMetrics{
		TotalPayments:     totalPayments,
		PaymentGrowthRate: paymentGrowthRate,
		PaymentMethods:    paymentStats,
	}, nil
}

// Add helper function to get facility-specific data
func (r *analyticsRepository) getFacilitySlotStats(ctx context.Context, facilityName string, startDate, endDate time.Time) (*analytics.FacilityMetrics, error) {
	facilityDb, err := r.getFacilityDb(facilityName)
	if err != nil {
		return nil, err
	}
	slotsCol := facilityDb.Collection("slots")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_slots": bson.M{"$sum": 1},
				"used_slots": bson.M{
					"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$gt": bson.A{"$current_bookings", 0}},
							1,
							0,
						},
					},
				},
				"total_capacity": bson.M{"$sum": "$max_bookings"},
				"total_bookings": bson.M{"$sum": "$current_bookings"},
			},
		},
	}

	cursor, err := slotsCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalSlots    int `bson:"total_slots"`
		UsedSlots     int `bson:"used_slots"`
		TotalCapacity int `bson:"total_capacity"`
		TotalBookings int `bson:"total_bookings"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	utilizationRate := float64(result.TotalBookings) / float64(result.TotalCapacity) * 100

	return &analytics.FacilityMetrics{
		TotalFacilities: 1,
		FacilityUtilization: map[string]float64{
			facilityName: utilizationRate,
		},
	}, nil
}

// Add helper function to aggregate data from all facilities
func (r *analyticsRepository) getAllFacilitiesStats(ctx context.Context, startDate, endDate time.Time) (*analytics.FacilityMetrics, error) {
	facilities := []string{"fitness", "swimming", "badminton", "football"}
	
	var totalFacilities int
	facilityUtilization := make(map[string]float64)
	var popularFacilities []analytics.FacilityUsageMetric

	for _, facilityName := range facilities {
		stats, err := r.getFacilitySlotStats(ctx, facilityName, startDate, endDate)
		if err != nil {
			continue
		}

		totalFacilities++
		if utilization, ok := stats.FacilityUtilization[facilityName]; ok {
			facilityUtilization[facilityName] = utilization
			popularFacilities = append(popularFacilities, analytics.FacilityUsageMetric{
				FacilityName: facilityName,
				UsageRate:    utilization,
			})
		}
	}

	// Sort popular facilities by usage rate
	sort.Slice(popularFacilities, func(i, j int) bool {
		return popularFacilities[i].UsageRate > popularFacilities[j].UsageRate
	})

	return &analytics.FacilityMetrics{
		TotalFacilities:     totalFacilities,
		FacilityUtilization: facilityUtilization,
		PopularFacilities:   popularFacilities,
	}, nil
}

// Add validation helper
func validateTimeRange(startDate, endDate time.Time) error {
	if endDate.Before(startDate) {
		return fmt.Errorf("end date must be after start date")
	}
	
	if startDate.After(time.Now()) {
		return fmt.Errorf("start date cannot be in the future")
	}
	
	return nil
}

func (r *analyticsRepository) getTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	bookingDb := r.db.Database("booking_db")
	bookingCol := bookingDb.Collection("booking_transaction")
	paymentDb := r.db.Database("payment_db")
	paymentCol := paymentDb.Collection("payments")

	// Aggregate query to calculate total revenue from payments related to bookings
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "payments",
				"localField":   "booking_id",
				"foreignField": "booking_id",
				"as":           "payment_info",
			},
		},
		{
			"$unwind": "$payment_info",
		},
		{
			"$match": bson.M{
				"payment_info.created_at": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_revenue": bson.M{"$sum": "$payment_info.amount"},
			},
		},
	}

	cursor, err := bookingCol.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to aggregate total revenue: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalRevenue float64 `bson:"total_revenue"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("failed to decode total revenue result: %w", err)
		}
	}

	return result.TotalRevenue, nil
}
