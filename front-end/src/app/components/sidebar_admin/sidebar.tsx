import {
  faHome,
  faDumbbell,
  faSwimmer,
  faFutbol,
  faEnvelope,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Link from "next/link";
import React from "react";
import Logo from "../../assets/Logo.png";

const sidebar = () => {
  return (
    <div className="bg-red-900 h-screen text-white w-80 flex flex-col px-5 py-10 ">
      <Link href="/" className="inline-flex flex-row justify-start gap-3.5">
        <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        <span className="flex flex-col border-l-2 w-max whitespace-nowrap">
          <div className="ms-1">
            <span className="ms-1 inline-flex flex-row font-semibold text-xl">
              <p className="text-black ">SPORT.</p>
              <p className="text-slate-200">MFU</p>
            </span>
            <hr />
            <span className="text-gray-300 ms-1 font-medium text-sm">
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
      <ul className="flex flex-col ps-5 gap-10 font-medium uppercase text-sm">
        <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
          <Link
            href="/admin_dashboard"
            className="inline-flex flex-row items-center"
          >
            <p>Dashboard</p>
          </Link>
        </li>
        <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
          <Link
            href="/admin_usermanagement"
            className="inline-flex flex-row items-center"
          >
            <p>User</p>
          </Link>
        </li>
        <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
          <Link
            href="/admin_sport"
            className="inline-flex flex-row items-center"
          >
            <p>Sport</p>
          </Link>
        </li>
        <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
          <Link
            href="/admin_queue"
            className="inline-flex flex-row items-center"
          >
            <p>Queue</p>
          </Link>
        </li>
        <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
          <Link
            href="/admin_report"
            className="inline-flex flex-row items-center"
          >
            <p>Report</p>
          </Link>
        </li>
      </ul>
      <div className="items-end justify-center h-full flex">Admin1</div>
    </div>
  );
};

export default sidebar;
