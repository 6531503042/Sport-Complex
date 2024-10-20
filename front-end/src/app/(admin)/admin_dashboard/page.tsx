"use client";

import React, { useEffect, useState } from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";

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
        <div className="border-b rounded-lg border-black"></div>
        <br />
        <div className="w-full flex justify-center"></div>
        <br />
        <br />
      </div>
    </div>
  );
};

export default UserManagementPage;
function setError(arg0: string) {
  throw new Error("Function not implemented.");
}
