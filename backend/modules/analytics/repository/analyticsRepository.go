package repository

import (
	"context"
	"fmt"
	"log"
	"main/modules/analytics"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
    AnalyticsRepositoryService interface {
        GetDailyStats(ctx context.Context, facilityName string, date time.Time) (*analytics.FacilityUsageStats, error)
        GetWeeklyStats(ctx context.Context, facilityName string, startDate time.Time) (*analytics.DashboardMetrics, error)
        GetMonthlyStats(ctx context.Context, facilityName string, year int, month int) (*analytics.DashboardMetrics, error)
        GetYearlyStats(ctx context.Context, facilityName string, year int) (*analytics.DashboardMetrics, error)
        GetUserMetrics(ctx context.Context) (*analytics.UserMetrics, error)
    }

    analyticsRepository struct {
        db *mongo.Client
    }
)

func NewAnalyticsRepository(db *mongo.Client) AnalyticsRepositoryService {
    return &analyticsRepository{db: db}
}

func (r *analyticsRepository) getRevenueStats(ctx context.Context, facilityName string, startDate, endDate time.Time) (float64, error) {
    paymentDb := r.db.Database("payment_db")
    
    pipeline := []bson.M{
        {
            "$match": bson.M{
                "facility_name": facilityName,
                "status": "COMPLETED",
                "created_at": bson.M{
                    "$gte": startDate,
                    "$lt":  endDate,
                },
            },
        },
        {
            "$group": bson.M{
                "_id": nil,
                "total_revenue": bson.M{"$sum": "$amount"},
            },
        },
    }

    cursor, err := paymentDb.Collection("payments").Aggregate(ctx, pipeline)
    if err != nil {
        return 0, fmt.Errorf("failed to aggregate revenue: %w", err)
    }
    defer cursor.Close(ctx)

    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return 0, fmt.Errorf("failed to decode revenue results: %w", err)
    }

    if len(results) > 0 {
        if revenue, ok := results[0]["total_revenue"].(float64); ok {
            return revenue, nil
        }
    }

    return 0, nil
}

func (r *analyticsRepository) GetDailyStats(ctx context.Context, facilityName string, date time.Time) (*analytics.FacilityUsageStats, error) {
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endOfDay := startOfDay.Add(24 * time.Hour)

    bookingDb := r.db.Database("booking_db")

    // Updated booking pipeline to get real booking data
    bookingPipeline := []bson.M{
        {
            "$match": bson.M{
                "$or": []bson.M{
                    {"facility": facilityName},
                    {"facility_name": facilityName},
                },
                "created_at": bson.M{
                    "$gte": startOfDay,
                    "$lt":  endOfDay,
                },
            },
        },
        {
            "$group": bson.M{
                "_id": nil,
                "total_bookings": bson.M{"$sum": 1},
                "completed_bookings": bson.M{
                    "$sum": bson.M{
                        "$cond": []interface{}{
                            bson.M{"$eq": []interface{}{"$status", "PAID"}},
                            1,
                            0,
                        },
                    },
                },
                "failed_bookings": bson.M{
                    "$sum": bson.M{
                        "$cond": []interface{}{
                            bson.M{"$eq": []interface{}{"$status", "FAILED"}},
                            1,
                            0,
                        },
                    },
                },
            },
        },
    }

    // Execute booking pipeline
    cursor, err := bookingDb.Collection("booking_transaction").Aggregate(ctx, bookingPipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate bookings: %w", err)
    }
    defer cursor.Close(ctx)

    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode results: %w", err)
    }

    // Also check historical data
    historyCursor, err := bookingDb.Collection("histories_transaction").Aggregate(ctx, bookingPipeline)
    if err == nil {
        defer historyCursor.Close(ctx)
        var historyResults []bson.M
        if err = historyCursor.All(ctx, &historyResults); err == nil && len(historyResults) > 0 {
            // Merge history results with current results
            if len(results) == 0 {
                results = historyResults
            } else {
                results[0]["total_bookings"] = results[0]["total_bookings"].(int32) + historyResults[0]["total_bookings"].(int32)
                results[0]["completed_bookings"] = results[0]["completed_bookings"].(int32) + historyResults[0]["completed_bookings"].(int32)
                results[0]["failed_bookings"] = results[0]["failed_bookings"].(int32) + historyResults[0]["failed_bookings"].(int32)
            }
        }
    }

    // Get revenue from payments
    revenue, err := r.getRevenueStats(ctx, facilityName, startOfDay, endOfDay)
    if err != nil {
        log.Printf("Error getting revenue stats: %v", err)
    }

    // Initialize stats with default values
    stats := &analytics.FacilityUsageStats{
        FacilityName:      facilityName,
        Date:             date,
        TotalBookings:    0,
        CompletedBookings: 0,
        FailedBookings:   0, // Changed from CanceledBookings to FailedBookings
        Revenue:         revenue,
        UtilizationRate:  0,
        PeakHours:       make([]int, 24),
    }

    // Process results
    if len(results) > 0 {
        if total, ok := results[0]["total_bookings"].(int32); ok {
            stats.TotalBookings = int(total)
        }
        if completed, ok := results[0]["completed_bookings"].(int32); ok {
            stats.CompletedBookings = int(completed)
        }
        if failed, ok := results[0]["failed_bookings"].(int32); ok {
            stats.FailedBookings = int(failed) // Changed from CanceledBookings to FailedBookings
        }
    }

    // Calculate utilization from slots
    facilityDb := r.db.Database(fmt.Sprintf("%s_facility", facilityName))
    var slots []bson.M
    slotsCursor, err := facilityDb.Collection("slots").Find(ctx, bson.M{})
    if err == nil {
        defer slotsCursor.Close(ctx)
        if err = slotsCursor.All(ctx, &slots); err == nil {
            totalCapacity := 0
            usedCapacity := 0
            for _, slot := range slots {
                if maxBookings, ok := slot["max_bookings"].(int32); ok {
                    totalCapacity += int(maxBookings)
                }
                if currentBookings, ok := slot["current_bookings"].(int32); ok {
                    usedCapacity += int(currentBookings)
                }
                // Track peak hours
                if startTime, ok := slot["start_time"].(string); ok {
                    if t, err := time.Parse("15:04", startTime); err == nil {
                        if currentBookings, ok := slot["current_bookings"].(int32); ok {
                            stats.PeakHours[t.Hour()] += int(currentBookings)
                        }
                    }
                }
            }
            if totalCapacity > 0 {
                stats.UtilizationRate = float64(usedCapacity) / float64(totalCapacity) * 100
            }
        }
    }

    return stats, nil
}

func (r *analyticsRepository) GetWeeklyStats(ctx context.Context, facilityName string, startDate time.Time) (*analytics.DashboardMetrics, error) {
    metrics := &analytics.DashboardMetrics{
        WeeklyStats: struct {
            BookingTrends    analytics.TimeSeriesData `json:"booking_trends"`
            RevenueByFacility map[string]float64 `json:"revenue_by_facility"`
            PeakDays         map[string]int    `json:"peak_days"`
        }{
            BookingTrends: analytics.TimeSeriesData{
                Labels: make([]string, 0),
                Values: make([]float64, 0),
            },
            RevenueByFacility: make(map[string]float64),
            PeakDays:         make(map[string]int),
        },
    }
    
    // Get data for each day of the week
    for i := 0; i < 7; i++ {
        date := startDate.AddDate(0, 0, i)
        dailyStats, err := r.GetDailyStats(ctx, facilityName, date)
        if err != nil {
            continue
        }

        // Add to weekly trends
        metrics.WeeklyStats.BookingTrends.Labels = append(
            metrics.WeeklyStats.BookingTrends.Labels,
            date.Format("Mon"),
        )
        metrics.WeeklyStats.BookingTrends.Values = append(
            metrics.WeeklyStats.BookingTrends.Values,
            float64(dailyStats.TotalBookings),
        )

        // Track peak days
        metrics.WeeklyStats.PeakDays[date.Format("Monday")] = dailyStats.TotalBookings
        
        // Add revenue
        metrics.WeeklyStats.RevenueByFacility[facilityName] += dailyStats.Revenue
    }

    return metrics, nil
}

func (r *analyticsRepository) GetMonthlyStats(ctx context.Context, facilityName string, year int, month int) (*analytics.DashboardMetrics, error) {
    startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endOfMonth := startOfMonth.AddDate(0, 1, 0)

    // Get monthly revenue from payments
    paymentDb := r.db.Database("payment_db")
    revenuePipeline := []bson.M{
        {
            "$match": bson.M{
                "facility_name": facilityName,
                "status": "COMPLETED",
                "created_at": bson.M{
                    "$gte": startOfMonth,
                    "$lt":  endOfMonth,
                },
            },
        },
        {
            "$group": bson.M{
                "_id": bson.M{
                    "day": bson.M{"$dayOfMonth": "$created_at"},
                },
                "revenue": bson.M{"$sum": "$amount"},
                "transaction_count": bson.M{"$sum": 1},
            },
        },
        {
            "$sort": bson.M{"_id.day": 1},
        },
    }

    revenueCursor, err := paymentDb.Collection("payments").Aggregate(ctx, revenuePipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate revenue: %w", err)
    }
    defer revenueCursor.Close(ctx)

    var revenueResults []bson.M
    if err = revenueCursor.All(ctx, &revenueResults); err != nil {
        return nil, fmt.Errorf("failed to decode revenue results: %w", err)
    }

    metrics := &analytics.DashboardMetrics{
        MonthlyStats: struct {
            BookingGrowth analytics.TimeSeriesData `json:"booking_growth"`
            RevenueGrowth analytics.TimeSeriesData `json:"revenue_growth"`
            FacilityComparison []analytics.FacilityComparison `json:"facility_comparison"`
        }{
            BookingGrowth: analytics.TimeSeriesData{
                Labels: make([]string, len(revenueResults)),
                Values: make([]float64, len(revenueResults)),
            },
            RevenueGrowth: analytics.TimeSeriesData{
                Labels: make([]string, len(revenueResults)),
                Values: make([]float64, len(revenueResults)),
            },
            FacilityComparison: make([]analytics.FacilityComparison, 0),
        },
    }

    // Process revenue data
    for i, result := range revenueResults {
        day := result["_id"].(bson.M)["day"].(int32)
        revenue := result["revenue"].(float64)
        transactions := result["transaction_count"].(int32)

        dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)
        
        metrics.MonthlyStats.BookingGrowth.Labels[i] = dateStr
        metrics.MonthlyStats.BookingGrowth.Values[i] = float64(transactions)
        
        metrics.MonthlyStats.RevenueGrowth.Labels[i] = dateStr
        metrics.MonthlyStats.RevenueGrowth.Values[i] = revenue
    }

    // Get facility comparison data
    facilities := []string{"fitness", "swimming", "badminton", "football"}
    for _, facility := range facilities {
        revenue, err := r.getRevenueStats(ctx, facility, startOfMonth, endOfMonth)
        if err != nil {
            continue
        }

        metrics.MonthlyStats.FacilityComparison = append(
            metrics.MonthlyStats.FacilityComparison,
            analytics.FacilityComparison{
                FacilityName: facility,
                Revenue:     revenue,
                // Other fields can be populated as needed
            },
        )
    }

    return metrics, nil
}

func (r *analyticsRepository) GetYearlyStats(ctx context.Context, facilityName string, year int) (*analytics.DashboardMetrics, error) {
    // Implementation
    return &analytics.DashboardMetrics{}, nil
}

func (r *analyticsRepository) GetUserMetrics(ctx context.Context) (*analytics.UserMetrics, error) {
    userDb := r.db.Database("user_db")
    
    // Get total users
    totalUsers, err := userDb.Collection("users").CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, fmt.Errorf("failed to count users: %w", err)
    }

    // Get new users (registered in last 30 days)
    thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
    newUsers, err := userDb.Collection("users").CountDocuments(ctx, bson.M{
        "created_at": bson.M{"$gte": thirtyDaysAgo},
    })
    if err != nil {
        return nil, fmt.Errorf("failed to count new users: %w", err)
    }

    // Get active users (with bookings in last 30 days)
    bookingDb := r.db.Database("booking_db")
    activeUsers, err := bookingDb.Collection("booking_transaction").Distinct(ctx, "user_id", bson.M{
        "created_at": bson.M{"$gte": thirtyDaysAgo},
    })
    if err != nil {
        return nil, fmt.Errorf("failed to count active users: %w", err)
    }

    return &analytics.UserMetrics{
        TotalUsers:     int(totalUsers),
        ActiveUsers:    len(activeUsers),
        NewUsers:       int(newUsers),
        UserGrowthRate: float64(newUsers) / float64(totalUsers) * 100,
    }, nil
}

// Implement other repository methods... 