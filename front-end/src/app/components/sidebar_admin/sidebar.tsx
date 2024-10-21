"use client";

import Link from "next/link";
import React, { useState } from "react";
import Logo from "../../assets/Logo.png";
import SpaceDashboardIcon from "@mui/icons-material/SpaceDashboard";
import ManageAccountsIcon from "@mui/icons-material/ManageAccounts";
import LocationCityIcon from "@mui/icons-material/LocationCity";
import BookmarksIcon from "@mui/icons-material/Bookmarks";
import PaidIcon from "@mui/icons-material/Paid";
import ArrowLeftIcon from "@mui/icons-material/KeyboardDoubleArrowLeft";

type SidebarProps = {
  activePage?: string;
};

const Sidebar: React.FC<SidebarProps> = ({ activePage }) => {
  const [isCollapsed, setIsCollapsed] = useState(false);

  const toggleSidebar = () => {
    setIsCollapsed(!isCollapsed);
  };

  const getActiveClass = (page: string) => {
    return activePage === page
      ? "inline-flex flex-row items-center border border-white text-white font-semibold shadow-gray-800 hover:shadow-black py-3 px-5 shadow-lg hover:shadow-2xl rounded-lg hover:scale-105 transition-all duration-700 hover:bg-orange-800"
      : "inline-flex flex-row items-center hover:scale-110 transition-transform duration-1000 ease-in-out ms-1 hover:shadow-lg py-3 px-5 rounded-lg";
  };

  return (
    <div
      className={`bg-red-900 text-white flex flex-col px-5 py-10 transition-all duration-700 ${
        isCollapsed ? "w-28" : "md:w-72"
      }`}
    >
      <Link href="/" className="inline-flex flex-row justify-center md:gap-3.5">
        <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        {!isCollapsed && (
          <span className="flex flex-col md:border-l-2 border-l-0 w-max whitespace-nowrap font-bold">
            <div className="md:ps-1">
              <span className="flex-row font-semibold text-xl md:ps-1 md:inline-flex hidden">
                <p className="text-black ">SPORT.</p>
                <p className="text-slate-200">MFU</p>
              </span>
              <hr />
              <span className="text-gray-300 md:ps-1 font-medium text-sm md:block hidden">
                USER MANAGEMENT
              </span>
            </div>
          </span>
        )}
      </Link>
      <br />
      <div className="relative justify-center items-center md:inline-flex hidden">
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <div
            className={` w-full h-[2px] bg-white${isCollapsed ? "" : ""}`}
          ></div>
        </div>
        <div className="relative z-10 p-2 bg-red-900">
          <button
            onClick={toggleSidebar}
            className="relative z-10 flex bg-yellow-500 p-1 rounded-full hover:opacity-100 transition-all ease-in-out duration-500 hover:scale-110 cursor-pointer "
          >
            <ArrowLeftIcon
              className={`text-white transition-transform duration-500 ${
                isCollapsed ? "rotate-180" : ""
              }`}
            />
          </button>
        </div>
      </div>

      <br />
      <ul
        className={`inline-flex flex-col px-5 gap-10${
          isCollapsed ? "md:items-start items-center px-5 gap-10" : ""
        }`}
      >
        <li
          className={getActiveClass("admin_dashboard")}
          style={{ cursor: "pointer" }}
        >
          <Link
            href="/admin_dashboard"
            className="inline-flex flex-row items-center cursor-pointer gap-3"
          >
            <SpaceDashboardIcon />
            {!isCollapsed && <p className="hidden md:block">Dashboard</p>}
          </Link>
        </li>
        <li
          className={getActiveClass("admin_usermanagement")}
          style={{ cursor: "pointer" }}
        >
          <Link
            href="/admin_usermanagement"
            className="inline-flex flex-row items-center cursor-pointer gap-3"
          >
            <ManageAccountsIcon />
            {!isCollapsed && <p className="hidden md:block">User</p>}
          </Link>
        </li>
        <li
          className={getActiveClass("admin_facility")}
          style={{ cursor: "pointer" }}
        >
          <Link
            href="/admin_facility"
            className="inline-flex flex-row items-center cursor-pointer gap-3"
          >
            <LocationCityIcon />
            {!isCollapsed && <p className="hidden md:block">Facility</p>}
          </Link>
        </li>
        <li
          className={getActiveClass("admin_booking")}
          style={{ cursor: "pointer" }}
        >
          <Link
            href="/admin_booking"
            className="inline-flex flex-row items-center cursor-pointer gap-3"
          >
            <BookmarksIcon />
            {!isCollapsed && <p className="hidden md:block">Booking</p>}
          </Link>
        </li>
        <li
          className={getActiveClass("admin_payment")}
          style={{ cursor: "pointer" }}
        >
          <Link
            href="/admin_payment"
            className="inline-flex flex-row items-center cursor-pointer gap-3"
          >
            <PaidIcon />
            {!isCollapsed && <p className="hidden md:block">Payment</p>}
          </Link>
        </li>
      </ul>
      <div className="items-end justify-center h-full flex">
        {!isCollapsed && "Admin1"}
      </div>
    </div>
  );
};

export default Sidebar;
