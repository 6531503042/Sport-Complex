import React, { useState } from "react";
import Link from "next/link";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBars, faBasketball, faCancel, faClipboard, faDumbbell, faEnvelope, faFutbol, faHome, faSwimmer, faUser, faWallet, faX } from "@fortawesome/free-solid-svg-icons";
import IconSidebar from "../../assets/icon_sidebar_black.png"

const sidebar: React.FC = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="flex-none border-b-4 border-transparent flex items-center">
      <div
        className=" cursor-pointer "
        onClick={toggleSidebar}
      >
        <img src={IconSidebar.src} alt="" className="h-auto w-6 transform transition-transform duration-100 ease-in-out hover:scale-125"/>
      </div>
      <div
        className={`fixed top-0 right-0 h-full bg-gray-800 text-white transform ${
          isSidebarOpen ? "translate-x-0 " : "translate-x-full"
        } transition-transform duration-300 ease-in-out w-80`}
      >
        <div className="">
          <button
            className="text-white hover:text-gray-300 p-8"
            onClick={toggleSidebar}
          >
            <FontAwesomeIcon icon={faX} />
          </button>
        </div>
        <ul className="flex flex-col py-2 px-8 gap-8 ">
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/homepage"><FontAwesomeIcon icon={faHome} className="mx-3.5" />
            Home Page</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/gym"><FontAwesomeIcon icon={faDumbbell} className="mx-3.5" />
            Gym Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/swimming"><FontAwesomeIcon icon={faSwimmer} className="mx-3.5" />
            Swimming Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/football"><FontAwesomeIcon icon={faFutbol} className="mx-3.5" />
            Football Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/contact"><FontAwesomeIcon icon={faEnvelope} className="mx-3.5" />
            Contact</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/rule"><FontAwesomeIcon icon={faClipboard} className="mx-3.5" />
            Rules</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/pages/payment"><FontAwesomeIcon icon={faWallet} className="mx-3.5" />
            Payment</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/"><FontAwesomeIcon icon={faUser} className="mx-3.5" />
            Profile</Link>
          </li>
        </ul>
      </div>
    </div>
  );
};

export default sidebar;
