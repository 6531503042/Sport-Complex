import React, { useState } from "react";
import Link from "next/link";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBars, faBasketball, faCancel, faClipboard, faDumbbell, faFutbol, faSwimmer, faUser, faWallet, faX } from "@fortawesome/free-solid-svg-icons";

const sidebar: React.FC = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="flex-none w-1/12">
      <div
        className="text-white hover:text-gray-400 cursor-pointer"
        onClick={toggleSidebar}
      >
        <FontAwesomeIcon icon={faBars} />
      </div>
      <div
        className={`fixed top-0 right-0 h-full bg-gray-800 text-white transform ${
          isSidebarOpen ? "translate-x-0" : "translate-x-full"
        } transition-transform duration-300 ease-in-out w-64`}
      >
        <div className="p-4">
          <button
            className="text-white hover:text-gray-300 "
            onClick={toggleSidebar}
          >
            <FontAwesomeIcon icon={faX} />
          </button>
        </div>
        <ul className="flex flex-col p-4 gap-8">
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/pages/gym"><FontAwesomeIcon icon={faDumbbell} className="mx-3.5" />
            Gym Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/pages/swimming"><FontAwesomeIcon icon={faSwimmer} className="mx-3.5" />
            Swimming Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/pages/basketball"><FontAwesomeIcon icon={faBasketball} className="mx-3.5" />
            Basketball Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/pages/soccer"><FontAwesomeIcon icon={faFutbol} className="mx-3.5" />
            Soccer Booking</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/pages/rule"><FontAwesomeIcon icon={faClipboard} className="mx-3.5" />
            Rules</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/pages/payment"><FontAwesomeIcon icon={faWallet} className="mx-3.5" />
            Payment</Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer">
            <Link href="/"><FontAwesomeIcon icon={faUser} className="mx-3.5" />
            Profile</Link>
          </li>
        </ul>
      </div>
    </div>
  );
};

export default sidebar;
