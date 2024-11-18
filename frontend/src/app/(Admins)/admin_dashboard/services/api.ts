export const API_BASE_URL = 'http://localhost:1325';

export const fetchAnalyticsData = async (period: string = 'monthly') => {
  try {
    const today = new Date();
    const startDate = '2024-01-01';
    const endDate = '2024-12-31';
    
    const url = `${API_BASE_URL}/user_v1/admin/analytics?period=${period}&start_date=${startDate}&end_date=${endDate}`;
    
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('Failed to fetch analytics data');
    }
    
    return response.json();
  } catch (error) {
    // Fallback to mock data if API call fails
    console.warn('Failed to fetch from API, using mock data:', error);
    return {
      totalUsers: 1250,
      totalBookings: 856,
      totalRevenue: 25600,
      growthRate: 15,
      userGrowthData: {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
        data: [65, 78, 90, 105, 125, 150],
      },
      facilityUsageData: {
        labels: ['Basketball', 'Football', 'Swimming', 'Tennis'],
        data: [30, 25, 20, 15],
      },
      revenueTrendData: {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
        data: [3000, 3500, 4000, 4200, 4800, 5100],
      },
    };
  }
};