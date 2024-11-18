export interface AnalyticsResponse {
    totalUsers: number;
    totalBookings: number;
    totalRevenue: number;
    growthRate: number;
    userGrowthData: {
      labels: string[];
      data: number[];
    };
    facilityUsageData: {
      labels: string[];
      data: number[];
    };
    revenueTrendData: {
      labels: string[];
      data: number[];
    };
  }