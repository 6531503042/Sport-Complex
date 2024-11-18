"use client";

import React, { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/app/components/card/card";
import { Bar, Line, Pie } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { Users, TrendingUp, DollarSign, Activity } from "lucide-react";

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

// Define types for dynamic data
type AnalyticsData = {
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
};

const Analytics = () => {
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Fetch data on component mount
  useEffect(() => {
    const fetchAnalyticsData = async () => {
      try {
        const response = await fetch('/api/analytics');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setAnalyticsData(data);
        setLoading(false);
      } catch (err) {
        console.error('Error fetching data:', err);
        setError("Error fetching data");
        setLoading(false);
      }
    };

    fetchAnalyticsData();
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>{error}</div>;
  if (!analyticsData) return <div>No data available</div>;

  return (
    <div className="space-y-6">
      <h2 className="text-3xl font-bold">Analytics Dashboard</h2>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Users</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analyticsData.totalUsers}</div>
            <p className="text-xs text-muted-foreground">+5% from last month</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Bookings</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analyticsData.totalBookings}</div>
            <p className="text-xs text-muted-foreground">+12% from last month</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${analyticsData.totalRevenue}</div>
            <p className="text-xs text-muted-foreground">+8% from last month</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Growth Rate</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analyticsData.growthRate}%</div>
            <p className="text-xs text-muted-foreground">+2% from last quarter</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>User Growth</CardTitle>
          </CardHeader>
          <CardContent>
            <Bar
              data={{
                labels: analyticsData.userGrowthData.labels,
                datasets: [
                  {
                    label: "New Users",
                    data: analyticsData.userGrowthData.data,
                    backgroundColor: "rgba(75, 192, 192, 0.6)",
                  },
                ],
              }}
            />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Facility Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <Pie
              data={{
                labels: analyticsData.facilityUsageData.labels,
                datasets: [
                  {
                    data: analyticsData.facilityUsageData.data,
                    backgroundColor: ["#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0"],
                  },
                ],
              }}
            />
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Revenue Trend</CardTitle>
        </CardHeader>
        <CardContent>
          <Line
            data={{
              labels: analyticsData.revenueTrendData.labels,
              datasets: [
                {
                  label: "Revenue",
                  data: analyticsData.revenueTrendData.data,
                  borderColor: "rgb(75, 192, 192)",
                  tension: 0.1,
                },
              ],
            }}
          />
        </CardContent>
      </Card>

      <div className="space-y-4">
        <h3 className="text-2xl font-bold">Insights for System Growth</h3>
        <ul className="list-disc pl-5 space-y-2">
          <li>
            User Acquisition: Focus on marketing channels that brought in the
            most active users.
          </li>
          <li>
            Facility Optimization: Increase availability for high-demand
            facilities during peak hours.
          </li>
          <li>
            Pricing Strategy: Analyze revenue per booking to optimize pricing
            for different facilities and time slots.
          </li>
          <li>
            User Retention: Implement loyalty programs or referral bonuses to
            increase user retention and acquisition.
          </li>
          <li>
            Seasonal Trends: Plan special promotions or events during typically
            slower periods to boost bookings.
          </li>
        </ul>
      </div>
    </div>
  );
};

export default Analytics;
