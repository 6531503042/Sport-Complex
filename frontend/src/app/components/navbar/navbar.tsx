"use client";

import React from "react";
import Link from "next/link";
import Logo from "../../assets/Logo.png";
import SideBar from "../sidebar/sidebar";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBasketball,
  faClipboard,
  faFutbol,
  faSwimmer,
  faWallet,
} from "@fortawesome/free-solid-svg-icons";
import { faDumbbell } from "@fortawesome/free-solid-svg-icons/faDumbbell";

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
      case "basketball":
        return "bg-yellow-500";
      case "soccer":
        return "bg-green-500";
      case "rule":
        return "bg-purple-500";
      case "payment":
        return "bg-orange-500";
      case "homepage":
        return "bg-red-400";
      default:
        return "bg-red-500";
    }
  };

  const getActiveClass = (page: string) => {
    return activePage === page
      ? "border-b-4 border-white"
      : "border-b-4 border-transparent";
  };

  return (
    <div className={`${getBackgroundColor()} px-8`}>
      <div className="NavBar_container flex flex-row items-center">
        <Link
          href="/"
          className="inline-flex flex-row items-center flex-none w-2/12 gap-3.5"
        >
          <img src={Logo.src} alt="Logo" className="w-7" />
          <span className="flex flex-col font-extrabold text-white text-lg">
            <span className="text-white">Sport</span>
            <span className="text-zinc-900">Complex</span>
          </span>
        </Link>
        <ul className="NavBar_res inline-flex flex-row flex-none w-9/12 gap-12 justify-start font-semibold pt-6">
          <li className={getActiveClass("gym")}>
            <Link
              href="/pages/gym"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              <FontAwesomeIcon icon={faDumbbell} className="mx-2.5" />
              Gym Booking
            </Link>
          </li>
          <li className={getActiveClass("swimming")}>
            <Link
              href="/pages/swimming"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              <FontAwesomeIcon icon={faSwimmer} className="mx-2.5" />
              Swimming Booking
            </Link>
          </li>
          <li className={getActiveClass("basketball")}>
            <Link
              href="/pages/basketball"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              <FontAwesomeIcon icon={faBasketball} className="mx-2.5" />
              Basketball Booking
            </Link>
          </li>
          <li className={getActiveClass("soccer")}>
            <Link
              href="/pages/soccer"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              <FontAwesomeIcon icon={faFutbol} className="mx-2.5" />
              Soccer Booking
            </Link>
          </li>
          <li className={getActiveClass("rule")}>
            <Link
              href="/pages/rule"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              <FontAwesomeIcon icon={faClipboard} className="mx-2.5" />
              Rules
            </Link>
          </li>
          <li className={getActiveClass("payment")}>
            <Link
              href="/pages/payment"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              <FontAwesomeIcon icon={faWallet} className="mx-2.5" />
              Payment
            </Link>
          </li>
        </ul>
        <div className="flex-none w-1/12 flex justify-between ms-5 me-2 gap-5">
          <aside className=" border-b-4 border-transparent">
            <nav className="inline-flex flex-row">
              <div className="text-white hover:text-gray-400 cursor-pointer">
                <Link href="/pages/registration">login</Link>
              </div>
            </nav>
          </aside>
          <SideBar />
        </div>
      </div>
    </div>
  );
};

export default navbar;
