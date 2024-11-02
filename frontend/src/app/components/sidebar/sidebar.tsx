import React, { useEffect, useState } from "react";
import Link from "next/link";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faClipboard,
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

const Sidebar: React.FC = () => {
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
  }, []);

  const truncateUserName = (name: string) => {
    return name.length > 20 ? name.slice(0, 20) + "..." : name;
  };

  const handleLogout = () => {
    localStorage.removeItem("user"); 
    router.replace("/login");
  };
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className="flex-none border-b-4 border-transparent flex items-center">
      <div className="cursor-pointer sidebar_icon" onClick={toggleSidebar}>
        <img
          src={IconSidebar.src}
          alt=""
          className="h-auto w-6 transform transition-transform duration-300 ease-in-out hover:scale-125"
        />
      </div>
      {isSidebarOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-30"
          onClick={toggleSidebar}
        />
      )}

      <div
        className={`fixed top-0 right-0 h-full bg-white text-black transform ${
          isSidebarOpen ? "translate-x-0" : "translate-x-full"
        } transition-transform duration-300 ease-in-out w-80 z-50 overflow-y-auto overflow-x-hidden`}
      >
        <div className="inline-flex flex-row w-full ">
          <div className="p-8 flex gap-5 w-full justify-between">
            <aside className="login_button border-b-4 border-transparent">
              <nav className="inline-flex">
                <span className="inline-flex  items-center">
                  {userName ? truncateUserName(userName) : "Loading..."}
                </span>
              </nav>
            </aside>
            <button
              className="text-black hover:text-gray-300 transition-all duration-200"
              onClick={toggleSidebar}
            >
              <FontAwesomeIcon
                icon={faX}
                style={{ fontSize: "1.5rem", fontWeight: "normal" }}
              />
            </button>
          </div>
        </div>

        <ul className="flex flex-col ps-5 gap-6 font-medium uppercase ">
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link
              href="/homepage"
              className="inline-flex flex-row items-center"
            >
              <FontAwesomeIcon
                icon={faHome}
                className=" text-orange-600 w-14"
              />
              <p>Home Page</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/gym" className="inline-flex flex-row items-center">
              <FontAwesomeIcon
                icon={faDumbbell}
                className=" text-orange-600 w-14"
              />
              <p>Gym Booking</p>
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
              <p>Swimming Booking</p>
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
              <p>Football Booking</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/contact" className="inline-flex flex-row items-center">
              <FontAwesomeIcon
                icon={faEnvelope}
                className=" text-orange-600 w-14"
              />
              <p>Contact</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/rule" className="inline-flex flex-row items-center">
              <FontAwesomeIcon
                icon={faClipboard}
                className=" text-orange-600 w-14"
              />
              <p>Rules</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link href="/payment" className="inline-flex flex-row items-center">
              <FontAwesomeIcon
                icon={faWallet}
                className=" text-orange-600 w-14"
              />
              <p>Payment</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link
              href="/homepage"
              className="inline-flex flex-row items-center"
            >
              <FontAwesomeIcon
                icon={faUser}
                className=" text-orange-600 w-14"
              />
              <p>Profile</p>
            </Link>
          </li>
          <li className="hover:text-gray-400 cursor-pointer transition-transform duration-200 ease-in-out hover:scale-110">
            <Link
              href="/admin_dashboard"
              className="inline-flex flex-row items-center"
            >
              <FontAwesomeIcon
                icon={faUser}
                className=" text-orange-600 w-14"
              />
              <p>AdminDashboard</p>
            </Link>
          </li>
        </ul>
        <br />
        <div className="h-[1.5px] rounded-lg mx-10 bg-zinc-700"></div>
        <br />
        <button
          onClick={handleLogout}
          className=" w-full flex justify-center"
        >
          <p className="hover:text-white hover:bg-red-900 border-[1.5px]  border-red-900 py-1.5 px-6 rounded-lg transition-all duration-300">Logout</p>
        </button>
      </div>
    </div>
  );
};

export default Sidebar;
