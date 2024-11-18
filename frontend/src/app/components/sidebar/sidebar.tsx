"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDumbbell,
  faEnvelope,
  faFutbol,
  faHome,
  faSwimmer,
  faUser,
  faWallet,
  faX,
} from "@fortawesome/free-solid-svg-icons";
import IconSidebar from "../../assets/icon_sidebar_black.png";
import "../sidebar/sidebar.css";
import { useRouter } from "next/navigation";
import LoadingScreen from "../loading_screen/loading";

type SidebarProps = {
  setLoading: React.Dispatch<React.SetStateAction<boolean>>;
};

const Sidebar: React.FC<SidebarProps> = ({ setLoading }) => {
  const [userName, setUserName] = useState<string | null>(null);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
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

  const truncateUserName = (name: string) => (name.length > 20 ? `${name.slice(0, 20)}...` : name);

  const handleLogout = () => {
    localStorage.removeItem("user");
    router.replace("/login");
  };

  const toggleSidebar = () => setIsSidebarOpen((prev) => !prev);

  const handleLinkClick = async (event: React.MouseEvent<HTMLAnchorElement>, href: string) => {
    event.preventDefault();
    setLoading(true);
    await router.push(href);
    setLoading(false); 
  };

  const menuItems = [
    { href: "/homepage", icon: faHome, label: "Home Page" },
    { href: "/gym-booking", icon: faDumbbell, label: "Gym Booking" },
    { href: "/badminton-booking", icon: faDumbbell, label: "Badminton Booking" },
    { href: "/swimming-booking", icon: faSwimmer, label: "Swimming Booking" },
    { href: "/football-booking", icon: faFutbol, label: "Football Booking" },
    { href: "/contact", icon: faEnvelope, label: "Contact" },
    { href: "/payment", icon: faWallet, label: "Payment" },
    { href: "/profile", icon: faUser, label: "Profile" },
    { href: "/admin_dashboard", icon: faUser, label: "( Admin )" },
  ];

  return (
    <>
      <div className="flex-none border-b-4 border-transparent flex items-center">
        <div className="cursor-pointer sidebar_icon" onClick={toggleSidebar}>
          <img
            src={IconSidebar.src}
            alt="Sidebar Icon"
            className="h-auto w-6 transform transition-transform duration-300 ease-in-out hover:scale-125"
          />
        </div>
        {isSidebarOpen && (
          <div className="fixed inset-0 bg-black bg-opacity-50 z-30" onClick={toggleSidebar} />
        )}
        <div
          className={`fixed top-0 right-0 h-full bg-white text-black transform ${
            isSidebarOpen ? "translate-x-0" : "translate-x-full"
          } transition-transform duration-300 ease-in-out w-80 z-50 overflow-y-auto`}
        >
          <div className="flex justify-between items-center p-8">
            <span className="text-lg font-semibold">
              {userName ? truncateUserName(userName) : "Loading..."}
            </span>
            <button onClick={toggleSidebar} className="text-black hover:text-gray-300">
              <FontAwesomeIcon icon={faX} style={{ fontSize: "1.5rem" }} />
            </button>
          </div>
          <ul className="flex flex-col ps-5 gap-6 font-medium uppercase">
            {menuItems.map((item, index) => (
              <li
                key={index}
                className="hover:text-gray-400 transition-transform duration-200 ease-in-out hover:scale-110"
              >
                <Link
                  href={item.href}
                  className="inline-flex flex-row items-center"
                  onClick={(e) => handleLinkClick(e, item.href)}
                >
                  <FontAwesomeIcon icon={item.icon} className="text-orange-600 w-14" />
                  <p>{item.label}</p>
                </Link>
              </li>
            ))}
          </ul>
          <div className="h-[1.5px] rounded-lg mx-10 bg-zinc-700 my-4"></div>
          <button onClick={handleLogout} className="w-full flex justify-center">
            <p className="hover:text-white hover:bg-red-900 border-[1.5px] border-red-900 py-1.5 px-6 rounded-lg transition-all duration-300">
              Logout
            </p>
          </button>
        </div>
      </div>
    </>
  );
};

export default Sidebar;
