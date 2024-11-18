package repository

import (
	"context"
	"fmt"
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

func (r *analyticsRepository) GetDailyStats(ctx context.Context, facilityName string, date time.Time) (*analytics.FacilityUsageStats, error) {
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endOfDay := startOfDay.Add(24 * time.Hour)

    bookingDb := r.db.Database("booking_db")
    facilityDb := r.db.Database(fmt.Sprintf("%s_facility", facilityName))

    // Updated booking pipeline to match your actual data structure
    bookingPipeline := []bson.M{
        {
            "$match": bson.M{
                "$and": []bson.M{
                    {
                        "$or": []bson.M{
                            {"slot_id": bson.M{"$exists": true}},
                            {"badminton_slot_id": bson.M{"$exists": true}},
                        },
                    },
                    {"created_at": bson.M{
                        "$gte": startOfDay,
                        "$lt":  endOfDay,
                    }},
                },
            },
        },
        {
            "$lookup": bson.M{
                "from": "facilities",
                "localField": "facility",
                "foreignField": "name",
                "as": "facility_info",
            },
        },
        {
            "$group": bson.M{
                "_id": nil,
                "total_bookings": bson.M{"$sum": 1},
                "paid_bookings": bson.M{
                    "$sum": bson.M{
                        "$cond": []interface{}{
                            bson.M{"$eq": []interface{}{"$status", "PAID"}},
                            1,
                            0,
                        },
                    },
                },
                "pending_bookings": bson.M{
                    "$sum": bson.M{
                        "$cond": []interface{}{
                            bson.M{"$eq": []interface{}{"$status", "pending"}},
                            1,
                            0,
                        },
                    },
                },
                "total_revenue": bson.M{
                    "$sum": bson.M{
                        "$cond": []interface{}{
                            bson.M{"$eq": []interface{}{"$status", "PAID"}},
                            "$price", // Use actual price from booking
                            0,
                        },
                    },
                },
            },
        },
    }

    var results []bson.M
    cursor, err := bookingDb.Collection("booking_transaction").Aggregate(ctx, bookingPipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate bookings: %w", err)
    }
    defer cursor.Close(ctx)
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode results: %w", err)
    }

    // Get slot data for utilization calculation
    slotPipeline := []bson.M{
        {
            "$match": bson.M{
                "facility_type": facilityName,
            },
        },
        {
            "$group": bson.M{
                "_id": nil,
                "total_capacity": bson.M{"$sum": "$max_bookings"},
                "used_capacity": bson.M{"$sum": "$current_bookings"},
                "peak_hours": bson.M{
                    "$push": bson.M{
                        "hour": bson.M{"$hour": "$start_time"},
                        "bookings": "$current_bookings",
                    },
                },
            },
        },
    }

    var slotResults []bson.M
    slotCursor, err := facilityDb.Collection("slots").Aggregate(ctx, slotPipeline)
    if err == nil {
        defer slotCursor.Close(ctx)
        slotCursor.All(ctx, &slotResults)
    }

    stats := &analytics.FacilityUsageStats{
        FacilityName:      facilityName,
        Date:             date,
        TotalBookings:    0,
        CompletedBookings: 0,
        CanceledBookings: 0,
        Revenue:         0,
        UtilizationRate:  0,
        PeakHours:       make([]int, 24),
    }

    if len(results) > 0 {
        if total, ok := results[0]["total_bookings"].(int32); ok {
            stats.TotalBookings = int(total)
        }
        if paid, ok := results[0]["paid_bookings"].(int32); ok {
            stats.CompletedBookings = int(paid)
        }
        if revenue, ok := results[0]["total_revenue"].(float64); ok {
            stats.Revenue = revenue
        }
    }

    if len(slotResults) > 0 {
        totalCapacity := slotResults[0]["total_capacity"].(int32)
        usedCapacity := slotResults[0]["used_capacity"].(int32)
        if totalCapacity > 0 {
            stats.UtilizationRate = float64(usedCapacity) / float64(totalCapacity) * 100
        }

        if peakHours, ok := slotResults[0]["peak_hours"].([]interface{}); ok {
            for _, ph := range peakHours {
                if hour, ok := ph.(bson.M); ok {
                    hourNum := int(hour["hour"].(int32))
                    bookings := int(hour["bookings"].(int32))
                    if hourNum >= 0 && hourNum < 24 {
                        stats.PeakHours[hourNum] = bookings
                    }
                }
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

    bookingDb := r.db.Database("booking_db")

    // Monthly booking trends pipeline
    monthlyPipeline := []bson.M{
        {
            "$match": bson.M{
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
                "bookings": bson.M{"$sum": 1},
                "revenue": bson.M{
                    "$sum": bson.M{
                        "$cond": []interface{}{
                            bson.M{"$eq": []interface{}{"$status", "PAID"}},
                            "$price",
                            0,
                        },
                    },
                },
            },
        },
        {
            "$sort": bson.M{"_id.day": 1},
        },
    }

    cursor, err := bookingDb.Collection("booking_transaction").Aggregate(ctx, monthlyPipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate monthly stats: %w", err)
    }
    defer cursor.Close(ctx)

    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode monthly results: %w", err)
    }

    // Prepare the metrics
    metrics := &analytics.DashboardMetrics{
        MonthlyStats: struct {
            BookingGrowth analytics.TimeSeriesData `json:"booking_growth"`
            RevenueGrowth analytics.TimeSeriesData `json:"revenue_growth"`
            FacilityComparison []analytics.FacilityComparison `json:"facility_comparison"`
        }{
            BookingGrowth: analytics.TimeSeriesData{
                Labels: make([]string, len(results)),
                Values: make([]float64, len(results)),
            },
            RevenueGrowth: analytics.TimeSeriesData{
                Labels: make([]string, len(results)),
                Values: make([]float64, len(results)),
            },
            FacilityComparison: make([]analytics.FacilityComparison, 0),
        },
    }

    // Fill in the data
    for i, result := range results {
        day := result["_id"].(bson.M)["day"].(int32)
        bookings := result["bookings"].(int32)
        revenue := result["revenue"].(float64)

        dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)
        metrics.MonthlyStats.BookingGrowth.Labels[i] = dateStr
        metrics.MonthlyStats.BookingGrowth.Values[i] = float64(bookings)
        metrics.MonthlyStats.RevenueGrowth.Labels[i] = dateStr
        metrics.MonthlyStats.RevenueGrowth.Values[i] = revenue
    }

    // Get facility comparison
    facilities := []string{"fitness", "swimming", "badminton", "football"}
    for _, f := range facilities {
        stats, err := r.GetDailyStats(ctx, f, time.Now())
        if err != nil {
            continue
        }
        metrics.MonthlyStats.FacilityComparison = append(
            metrics.MonthlyStats.FacilityComparison,
            analytics.FacilityComparison{
                FacilityName:   f,
                BookingCount:   stats.TotalBookings,
                Revenue:        stats.Revenue,
                PopularityRank: 0, // Calculate rank if needed
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