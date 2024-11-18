"use client";

import React, { useEffect, useState } from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";
import Analytics from "@/app/(Admins)/admin_dashboard/analytics";

const UserManagementPage = () => {
  return (
    <div className="w-screen h-screen relative flex">
      <Sidebar activePage="admin_dashboard" />
      <div className=" w-full bg-white text-black p-10 flex flex-col overflow-y-auto">
        <div className="inline-flex justify-between w-full items-end">
          <div className="text-lg font-medium">Dashboard</div>
          <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        </div>
        <br />
        <div className="bg-zinc-500 h-[1px] rounded-lg text-transparent">.</div>
        <br />
        <div>
          <Analytics />
        </div>
      </div>
    </div>
  );
};

export default UserManagementPage;
