"use client";

import React from "react";
import Link from "next/link";
import Logo from "../../assets/Logo.png";
import SideBar from "../../../../../front-end/src/app/components/sidebar/sidebar";
import SearhBar from "../search_bar/search_bar";
import "./navbar.css";


type NavBarProps = {
  activePage?: string;
};

const navbar: React.FC<NavBarProps> = ({ activePage }) => {
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
            href="/"
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
          <div className="flex-none w-3/5 flex me-3">
            <SearhBar />
          </div>
          <div className="login_and_sidebar flex-none w-1/12 flex justify-end ms-5 me-2 gap-12">
            <aside className="login_button border-b-4 border-transparent">
              <nav className="inline-flex">
                <div className=" hover:text-white hover:bg-orange-700 hover:shadow-lg transition-all duration-300 cursor-pointer py-2 px-5 bg-tranparent border rounded-full border-orange-700 text-orange-700 items-center">
                  <Link href="/admin_dashboard" className=" font-medium">
                    Login
                  </Link>
                </div>
              </nav>
            </aside>
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
            <i className="material-symbols-outlined mx-2.5 " style={{ fontSize: '1.3rem' }}>exercise</i>
            Gym Booking
          </Link>
        </li>
        <li className={getActiveClass("badminton")}>
          <Link
            href="/badminton"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <i className="material-symbols-outlined mx-2.5 " style={{ fontSize: '1.3rem' }}>sports_tennis</i>
            Badminton Booking
          </Link>
        </li>
        <li className={getActiveClass("swimming")}>
          <Link
            href="/swimming"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <i className="material-symbols-outlined mx-2.5 " style={{ fontSize: '1.3rem' }}>pool</i>
            Swimming Booking
          </Link>
        </li>
        <li className={getActiveClass("football")}>
          <Link
            href="/football"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <i className="material-symbols-outlined mx-2.5" style={{ fontSize: '1.3rem' }}>sports_soccer</i>
            Football Booking
          </Link>
        </li>
        <li className={getActiveClass("rule")}>
          <Link
            href="/rule"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <i className="material-symbols-outlined mx-2.5" style={{ fontSize: '1.3rem' }}>assignment_late</i>
            Rules
          </Link>
        </li>
        <li className={getActiveClass("contact")}>
          <Link
            href="/contact"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <i className="material-symbols-outlined mx-2.5" style={{ fontSize: '1.3rem' }}>mail</i>
            Contact
          </Link>
        </li>
        <li className={getActiveClass("payment")}>
          <Link
            href="/payment"
            className="text-white hover:text-gray-400 flex items-center pb-4 me-2"
          >
            <i className="material-symbols-outlined mx-2.5" style={{ fontSize: '1.3rem' }}>wallet</i>
            Payment
          </Link>
        </li>
      </ul>
    </div>
  );
};

export default navbar;
