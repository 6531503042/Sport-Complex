"use client";

import Link from "next/link";
import React from "react";
import Logo from "../../assets/Logo.png";
import SpaceDashboardIcon from "@mui/icons-material/SpaceDashboard";
import PeopleAltIcon from "@mui/icons-material/PeopleAlt";
import LocationCityIcon from "@mui/icons-material/LocationCity";
import EventNoteIcon from "@mui/icons-material/EventNote";
import FeedIcon from "@mui/icons-material/Feed";

type SidebarProps = {
  activePage?: string;
};

const sidebar: React.FC<SidebarProps> = ({ activePage }) => {
  const getActiveClass = (page: string) => {
    return activePage === page
      ? "bg-black text-white py-3 px-5 shadow-lg hover:shadow-2xl rounded-lg"
      : "hover:scale-110 transition-transform duration-200 ease-in-out ms-1 hover:shadow-lg ";
  };

  return (
    <div className="bg-red-900 h-screen text-white md:w-[300px] w-[100px] flex flex-col px-5 py-10">
      <Link href="/" className="inline-flex flex-row justify-center md:gap-3.5">
        <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        <span className="flex flex-col md:border-l-2 border-l-0 w-max whitespace-nowrap ">
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
      </Link>
      <br />
      <br />
      <div className="border-b"></div>
      <br />
      <br />
      <ul className="inline-flex flex-col md:items-start items-center px-5 gap-10 font-medium text-base">
        <li className={getActiveClass("admin_dashboard")} style={{cursor:"pointer"}}>
          <Link
            href="/admin_dashboard"
            className="inline-flex flex-row items-center cursor-pointer"
          >
            <SpaceDashboardIcon className="md:hidden" />
            <p className="hidden md:block">Dashboard</p>
          </Link>
        </li>
        <li className={getActiveClass("admin_usermanagement")} style={{cursor:"pointer"}}>
          <Link
            href="/admin_usermanagement"
            className="inline-flex flex-row items-center cursor-pointer"
          >
            <PeopleAltIcon className="md:hidden" />
            <p className="hidden md:block">User Management</p>
          </Link>
        </li>
        <li className={getActiveClass("admin_facility")} style={{cursor:"pointer"}}>
          <Link
            href="/admin_facility"
            className="inline-flex flex-row items-center cursor-pointer"
          >
            <LocationCityIcon className="md:hidden" />
            <p className="hidden md:block">Facility</p>
          </Link>
        </li>
        <li className={getActiveClass("admin_queue")} style={{cursor:"pointer"}}>
          <Link
            href="/admin_queue"
            className="inline-flex flex-row items-center cursor-pointer"
          >
            <EventNoteIcon className="md:hidden" />
            <p className="hidden md:block">Queue</p>
          </Link>
        </li>
        <li className={getActiveClass("admin_report")} style={{cursor:"pointer"}}>
          <Link
            href="/admin_report"
            className="inline-flex flex-row items-center cursor-pointer"
          >
            <FeedIcon className="md:hidden" />
            <p className="hidden md:block">Report</p>
          </Link>
        </li>
      </ul>
      <div className="items-end justify-center h-full flex">Admin1</div>
    </div>
  );
};

export default sidebar;
