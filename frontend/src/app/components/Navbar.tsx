"use client";

import React from "react";
import Link from "next/link";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faDumbbell } from "@fortawesome/free-solid-svg-icons/faDumbbell";

type NavBarProps = {
  activePage?: string;
};

const NavBar: React.FC<NavBarProps> = ({ activePage }) => {
  const getBackgroundColor = () => {
    switch (activePage) {
      case "gym":
        return "bg-gray-500";
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
    return activePage === page ? "border-b-4 border-white" : "";
  };

  return (
    <div className={`${getBackgroundColor()} px-8`}>
      <div className="flex flex-row items-center ">
        <Link
          href="/"
          className="inline-flex flex-row items-center flex-none w-2/12 gap-3.5"
        >

          <span className="flex flex-col font-extrabold text-white text-lg">
            <span className="text-white">Sport</span>
            <span className="text-zinc-900">Complex</span>
          </span>
        </Link>
        <ul className="inline-flex flex-row flex-none w-9/12 gap-12 justify-start font-semibold pt-6">
          <li className={getActiveClass("gym")}>
            <Link
              href="/gym"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              Gym Booking
            </Link>
          </li>
          <li className={getActiveClass("swimming")}>
            <Link
              href="/swimming-booking"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >

              Swimming Booking
            </Link>
          </li>
          <li className={getActiveClass("basketball")}>
            <Link
              href="/basketball"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >

              Basketball Booking
            </Link>
          </li>
          <li className={getActiveClass("soccer")}>
            <Link
              href="/football-booking"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
              Soccer Booking
            </Link>
          </li>
          <li className={getActiveClass("rule")}>
            <Link
              href="/rule"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >
   
              Rules
            </Link>
          </li>
          <li className={getActiveClass("payment")}>
            <Link
              href="/payment"
              className="text-white hover:text-gray-400 flex items-center pb-6 me-2"
            >

              Payment
            </Link>
          </li>
        </ul>
        <aside className="flex-none w-1/12 flex justify-center border-b-4 border-transparent">
          <nav className="inline-flex flex-row">
            <div className="text-white hover:text-gray-400 cursor-pointer">
              <Link href="/registration">
              login</Link>
            </div>
          </nav>
        </aside>
      </div>
    </div>
  );
};

export default NavBar;