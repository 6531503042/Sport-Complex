"use client";

import React from "react";
import Link from "next/link";
import Logo from "../../assets/Logo.png";
import SideBar from "../sidebar/sidebar";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBasketball,
  faCircle,
  faClipboard,
  faDumbbell,
  faEnvelope,
  faSwimmer,
  faWallet,
  faFutbol,
} from "@fortawesome/free-solid-svg-icons";
import SearhBar from "../../components/search_bar/search_bar";
import '../navbar/navbar.css'

type NavBarProps = {
  activePage?: string;
};

const navbar: React.FC<NavBarProps> = ({ activePage }) => {
  const getBackgroundColor = () => {
    switch (activePage) {
      case "gym":
        return "bg-red-500";
      case "swimming":
        return "bg-blue-500";
      case "football":
        return "bg-yellow-500";
      case "contact":
        return "bg-green-500";
      case "rule":
        return "bg-purple-500";
      case "payment":
        return "bg-orange-500";
      case "homepage":
        return "bg-red-400";
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
        <div className="NavBar_container flex flex-row items-center justify-between bg-white px-8 py-5">
          <Link
            href="/"
            className="inline-flex flex-row items-center flex-none gap-3.5 w-1/5"
          >
            <img src={Logo.src} alt="Logo" className="w-7" />
            <span className="flex flex-col border-l-2 w-max">
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
          <div className="flex-none w-3/5 flex me-3">
            <SearhBar />
          </div>
          <div className="login_and_sidebar flex-none w-1/12 flex justify-end ms-5 me-2 gap-12">
            <aside className="login_button border-b-4 border-transparent">
              <nav className="inline-flex">
                <div className=" hover:text-white hover:bg-orange-700 hover:shadow-lg transition-all duration-300 cursor-pointer py-2 px-5 bg-tranparent border rounded-full border-orange-700 text-orange-700 items-center">
                  <Link href="/pages/registration" className=" font-medium">
                    Login
                  </Link>
                </div>
              </nav>
            </aside>
            <SideBar />
          </div>
        </div>
      </header>
      <ul className="NavBar_res inline-flex flex-row pr-10 gap-16 justify-center font-semibold pt-4 text-sm">
        <li className={getActiveClass("gym")}>
          <Link
            href="/pages/gym"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FontAwesomeIcon icon={faDumbbell} className="mx-2.5" />
            Gym Booking
          </Link>
        </li>
        <li className={getActiveClass("swimming")}>
          <Link
            href="/pages/swimming"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FontAwesomeIcon icon={faSwimmer} className="mx-2.5" />
            Swimming Booking
          </Link>
        </li>
        <li className={getActiveClass("football")}>
          <Link
            href="/pages/football"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FontAwesomeIcon icon={faFutbol} className="mx-2.5" />
            Football booking
          </Link>
        </li>
        <li className={getActiveClass("rule")}>
          <Link
            href="/pages/rule"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FontAwesomeIcon icon={faClipboard} className="mx-2.5" />
            Rules
          </Link>
        </li>
        <li className={getActiveClass("contact")}>
          <Link
            href="/pages/contact"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FontAwesomeIcon icon={faEnvelope} className="mx-2.5" />
            Contact
          </Link>
        </li>
        <li className={getActiveClass("payment")}>
          <Link
            href="/pages/payment"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <FontAwesomeIcon icon={faWallet} className="mx-2.5" />
            Payment
          </Link>
        </li>
      </ul>
    </div>
  );
};

export default navbar;
