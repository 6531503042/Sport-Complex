import React from "react";
import Link from "next/link";
import Logo from "../../assets/Logo.png";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faHome, faDumbbell, faSwimmer, faFutbol, faEnvelope, faClipboard, faWallet, faUser } from "@fortawesome/free-solid-svg-icons";

const page = () => {
  return (
    <div className="w-screen h-screen flex flex-row">
      <div className="bg-red-500 text-white w-80 flex flex-col p-5 ">
        <Link
          href="/"
          className="inline-flex flex-row h-min justify-start gap-3.5"
        >
          <img src={Logo.src} alt="Logo" className="w-7" />
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
        <br/>
        <hr className="w-full"/>
        <br />
        <ul className="flex flex-col ps-5 gap-4 font-medium uppercase ">
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link
              href="/homepage"
              className="inline-flex flex-row items-center"
            >
              <FontAwesomeIcon
                icon={faHome}
                className=" text-orange-600 w-14"
              />
              <p>Dashboard</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/gym" className="inline-flex flex-row items-center">
              <FontAwesomeIcon
                icon={faDumbbell}
                className=" text-orange-600 w-14"
              />
              <p>User</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link
              href="/swimming"
              className="inline-flex flex-row items-center"
            >
              <FontAwesomeIcon
                icon={faSwimmer}
                className=" text-orange-600 w-14"
              />
              <p>Sport</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link
              href="/football"
              className="inline-flex flex-row items-center"
            >
              <FontAwesomeIcon
                icon={faFutbol}
                className=" text-orange-600 w-14"
              />
              <p>Queue</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/contact" className="inline-flex flex-row items-center">
              <FontAwesomeIcon
                icon={faEnvelope}
                className=" text-orange-600 w-14"
              />
              <p>Report</p>
            </Link>
          </li>
          
        </ul>
      </div>
      <div className="bg-black text-white w-full">right</div>
    </div>
  );
};

export default page;
