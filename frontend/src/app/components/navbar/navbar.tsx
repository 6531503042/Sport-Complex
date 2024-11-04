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
    } else {
      router.replace("/login"); 
    }
  }, [router]);

  const truncateUserName = (name: string) => {
    return name.length > 10 ? name.slice(0, 10) + "..." : name;
  };

  const handleLogout = () => {
    localStorage.removeItem("user");
    localStorage.removeItem("access_token");
    router.replace("/login");
  };

  const getBackgroundColor = () => {
    switch (activePage) {
      default:
        return "bg-red-900";
    }
  };

  const getActiveClass = (page: string) => {
    return activePage === page
      ? "color-wave text-black shadow-md rounded-lg"
      : "text-white hover:text-yellow-300 hover:border-black hover:border-opacity-50";
  };

  return (
    <div className={`${getBackgroundColor()} justify-center flex flex-col`}>
      <header>
        <div className="NavBar_container flex flex-row items-center justify-between bg-white px-20 py-5">
          <Link
            href="/homepage"
            className="inline-flex flex-row items-center gap-3.5 w-1/5"
          >
            <img src={Logo.src} alt="Logo" className="w-7" />
            <span className="flex flex-col border-l-2 w-max whitespace-nowrap">
              <div className="ms-1">
                <span className="ms-1 inline-flex flex-row font-semibold text-xl">
                  <p className="text-black">SPORT.</p>
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
          <div className="name_user_and_sidebar flex-none w-1/12 flex justify-end items-center ms-5 me-2 gap-12">
            <span className="name_user inline-flex flex-row gap-5 items-center">
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
      <ul className="NavBar_res inline-flex flex-row px-10 py-4 gap-16 justify-center items-center font-semibold text-sm">
        {[
          {
            href: "/gym",
            label: "Gym Booking",
            icon: <GymIcon style={{ fontSize: "1.3rem" }} />,
            page: "gym",
          },
          {
            href: "/badminton",
            label: "Badminton Booking",
            icon: <BadmintonIcon style={{ fontSize: "1.3rem" }} />,
            page: "badminton",
          },
          {
            href: "/swimming",
            label: "Swimming Booking",
            icon: <SwimmingIcon style={{ fontSize: "1.3rem" }} />,
            page: "swimming",
          },
          {
            href: "/football",
            label: "Football Booking",
            icon: <FootballIcon style={{ fontSize: "1.3rem" }} />,
            page: "football",
          },
          {
            href: "/rule",
            label: "Rules",
            icon: <RuleIcon style={{ fontSize: "1.3rem" }} />,
            page: "rule",
          },
          {
            href: "/contact",
            label: "Contact",
            icon: <ContactIcon style={{ fontSize: "1.3rem" }} />,
            page: "contact",
          },
          {
            href: "/payment",
            label: "Payment",
            icon: <PaymentIcon style={{ fontSize: "1.3rem" }} />,
            page: "payment",
          },
        ].map(({ href, label, icon, page }) => (
          <li
            key={page}
            className={`${getActiveClass(page)} hover:animate-wiggle`}
          >
            <Link
              href={href}
              className="flex items-center gap-2.5 py-4 px-3 border border-transparent hover:border hover:shadow-md rounded-lg transition-all duration-300"
            >
              {icon}
              <p>{label}</p>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default NavBar;
``
