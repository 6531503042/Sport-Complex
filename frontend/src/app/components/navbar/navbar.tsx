"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import Logo from "../../assets/Logo.png";
import SideBar from "../../../../../frontend/src/app/components/sidebar/sidebar";
import SearhBar from "../search_bar/search_bar";
import GymIcon from "@mui/icons-material/FitnessCenter";
import BadmintonIcon from "@mui/icons-material/SportsTennis";
import SwimmingIcon from "@mui/icons-material/Pool";
import FootballIcon from "@mui/icons-material/SportsSoccer";
import RuleIcon from "@mui/icons-material/AssignmentLate";
import ContactIcon from "@mui/icons-material/Mail";
import PaymentIcon from "@mui/icons-material/Payment";

type NavBarProps = {
  activePage?: string;
};

const NavBar: React.FC<NavBarProps> = ({ activePage }) => {
  const [userName, setUserName] = useState<string | null>(null);
  const router = useRouter(); 

  useEffect(() => {
    const userData = localStorage.getItem("user");
    if (userData) {
      const user = JSON.parse(userData);
      setUserName(user.name);
    }
  }, []);

  const truncateUserName = (name: string) => {
    return name.length > 10 ? name.slice(0, 10) + "..." : name;
  };

  const handleLogout = () => {
    localStorage.removeItem("user"); 
    router.push("/login"); 
  };

  const getBackgroundColor = () => {
    switch (activePage) {
      default:
        return "bg-red-900";
    }
  };

  const getActiveClass = (page: string) => {
    return activePage === page
      ? "border-b-4 border-white "
      : "border-b-4 border-transparent";
  };

  return (
    <div className={`${getBackgroundColor()} justify-center flex flex-col`}>
      <header>
        <div className="NavBar_container flex flex-row items-center justify-between bg-white px-20 py-5">
          <Link
            href="/homepage"
            className="inline-flex flex-row items-center  gap-3.5 w-1/5"
          >
            <img src={Logo.src} alt="Logo" className="w-7" />
            <span className="flex flex-col border-l-2 w-max whitespace-nowrap">
              <div className="ms-1">
                <span className="ms-1 inline-flex flex-row font-semibold text-xl">
                  <p className="text-black ">SPORT.</p>
                  <p className="text-gray-500">MFU</p>
                </span>
                <hr />
                <span className="text-gray-500 ms-1 font-medium text-sm">
                  SPORT COMPLEX
                </span>
              </div>
            </span>
          </Link>
          <div className="flex-none w-3/6 flex me-3">
            <SearhBar />
          </div>
          <div className="login_and_sidebar flex-none w-1/12 flex justify-end items-center ms-5 me-2 gap-12">
            <span className="inline-flex flew-row gap-5 items-center">
              {userName ? truncateUserName(userName) : "Loading..."}
              <p>|</p>
              <button
                onClick={handleLogout}
                className="hover:text-white hover:bg-red-900 border-[1.8px] border-red-900 py-1 px-2 rounded-lg transition-all duration-300"
              >
                Logout
              </button>
            </span>
            <SideBar />
          </div>
        </div>
      </header>
      <ul className="NavBar_res inline-flex flex-row px-10 gap-16 justify-center items-center font-semibold pt-4 text-sm">
        <li className={getActiveClass("gym")}>
          <Link
            href="/gym"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <GymIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Gym Booking
          </Link>
        </li>
        <li className={getActiveClass("badminton")}>
          <Link
            href="/badminton"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <BadmintonIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Badminton Booking
          </Link>
        </li>
        <li className={getActiveClass("swimming")}>
          <Link
            href="/swimming"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <SwimmingIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Swimming Booking
          </Link>
        </li>
        <li className={getActiveClass("football")}>
          <Link
            href="/football"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FootballIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Football Booking
          </Link>
        </li>
        <li className={getActiveClass("rule")}>
          <Link
            href="/rule"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <RuleIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Rules
          </Link>
        </li>
        <li className={getActiveClass("contact")}>
          <Link
            href="/contact"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <ContactIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Contact
          </Link>
        </li>
        <li className={getActiveClass("payment")}>
          <Link
            href="/payment"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <PaymentIcon className="mx-2.5" style={{ fontSize: "1.3rem" }} />
            Payment
          </Link>
        </li>
      </ul>
    </div>
  );
};

export default NavBar;
