"use client";

import React, { useEffect, useState } from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/app/components/card/card";
import { Users, Calendar, DollarSign, TrendingUp } from "lucide-react";
import Analytics from "@/app/lib/analytics"

const UserManagementPage = () => {
  return (
    <div className="w-screen h-screen flex flex-row">
      <Sidebar activePage="admin_dashboard" />
      <div className="bg-white text-black w-full p-10 flex flex-col">
        <div className="inline-flex justify-between w-full items-end">
          <div className="text-lg font-medium">Dashboard</div>
          <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        </div>
        <br />
        <div className="bg-zinc-500 h-[1px] rounded-lg text-transparent">.</div>
        <br />
        {/* <br />
        <h1 className="text-lg font-semibold">Dashboard Analysis</h1>
        <br /> */}
        {/* <div className=" grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Users</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">1,234</div>
              <p className="text-xs text-muted-foreground">
                +20% from last month
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Total Bookings
              </CardTitle>
              <Calendar className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">5,678</div>
              <p className="text-xs text-muted-foreground">
                +15% from last month
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Revenue</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">$12,345</div>
              <p className="text-xs text-muted-foreground">
                +10% from last month
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Active Users
              </CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">789</div>
              <p className="text-xs text-muted-foreground">
                +5% from last week
              </p>
            </CardContent>
          </Card>
        </div> */}
        <div>
          <Analytics/>
        </div>
      </div>
    </div>
  );
};

export default UserManagementPage;
function setError(arg0: string) {
  throw new Error("Function not implemented.");
}
