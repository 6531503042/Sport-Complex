"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import Image from "next/image";
import Logo from "../../assets/Logo.png";
import {
  Home,
  FitnessCenter,
  Pool,
  SportsSoccer,
  Email,
  Person2,
  SportsTennis,
  Payment,
  Close,
  Menu,
  LogoutOutlined,
  Settings,
  Dashboard,
} from "@mui/icons-material";

interface SidebarProps {
  setLoading: React.Dispatch<React.SetStateAction<boolean>>;
}

const Sidebar: React.FC<SidebarProps> = ({ setLoading }) => {
  const [userName, setUserName] = useState<string | null>(null);
  const [userRole, setUserRole] = useState<number | null>(null);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const userData = localStorage.getItem("user");
    if (userData) {
      const user = JSON.parse(userData);
      setUserName(user.name);
      const roleCode = user.user_roles && user.user_roles.length > 0 ? user.user_roles[0]?.role_code : null;
      setUserRole(roleCode);
    } else {
      router.replace("/login");
    }
  }, [router]);

  const truncateUserName = (name: string) => (name.length > 20 ? `${name.slice(0, 20)}...` : name);

  const handleLogout = () => {
    localStorage.removeItem("user");
    localStorage.removeItem("access_token");
    router.replace("/login");
  };

  const toggleSidebar = () => setIsSidebarOpen((prev) => !prev);

  const menuItems = [
    { href: "/homepage", icon: <Home className="text-red-600" />, label: "Home" },
    { href: "/gym-booking", icon: <FitnessCenter className="text-red-600" />, label: "Gym" },
    { href: "/badminton-booking", icon: <SportsTennis className="text-red-600" />, label: "Badminton" },
    { href: "/swimming-booking", icon: <Pool className="text-red-600" />, label: "Swimming" },
    { href: "/football-booking", icon: <SportsSoccer className="text-red-600" />, label: "Football" },
    { href: "/contact", icon: <Email className="text-red-600" />, label: "Contact" },
    { href: "/payment", icon: <Payment className="text-red-600" />, label: "Payment" },
    { href: "/profile", icon: <Person2 className="text-red-600" />, label: "Profile" },
  ];

  if (userRole === 1) {
    menuItems.push({
      href: "/admin_dashboard",
      icon: <Dashboard className="text-red-600" />,
      label: "Admin Dashboard",
    });
  }

  const sidebarVariants = {
    open: { x: 0, transition: { type: "spring", stiffness: 300, damping: 30 } },
    closed: { x: "100%", transition: { type: "spring", stiffness: 300, damping: 30 } },
  };

  const menuItemVariants = {
    open: {
      x: 0,
      opacity: 1,
      transition: { type: "spring", stiffness: 300, damping: 30 },
    },
    closed: { x: 20, opacity: 0 },
  };

  return (
    <>
      <button
        onClick={toggleSidebar}
        className="p-2 rounded-full hover:bg-gray-100 transition-colors duration-200"
      >
        <Menu className="w-6 h-6 text-gray-600" />
      </button>

      <AnimatePresence>
        {isSidebarOpen && (
          <>
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={toggleSidebar}
              className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40"
            />

            <motion.div
              initial="closed"
              animate="open"
              exit="closed"
              variants={sidebarVariants}
              className="fixed top-0 right-0 h-full w-80 bg-white shadow-2xl z-50"
            >
              <div className="p-6">
                <div className="flex justify-between items-center mb-8">
                  <div className="flex items-center gap-3">
                    <Image src={Logo} alt="Logo" width={40} height={40} />
                    <div>
                      <p className="font-semibold text-gray-900">
                        {userName ? truncateUserName(userName) : "Loading..."}
                      </p>
                      <p className="text-sm text-gray-500">
                        {userRole === 1 ? "Administrator" : "User"}
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={toggleSidebar}
                    className="p-2 rounded-full hover:bg-gray-100 transition-colors"
                  >
                    <Close className="w-5 h-5 text-gray-600" />
                  </button>
                </div>

                <div className="space-y-1">
                  {menuItems.map((item, index) => (
                    <motion.div
                      key={item.href}
                      variants={menuItemVariants}
                      custom={index}
                    >
                      <Link
                        href={item.href}
                        className="flex items-center gap-3 px-4 py-3 rounded-lg
                          text-gray-700 hover:bg-gray-100 transition-colors
                          group relative overflow-hidden"
                        onClick={toggleSidebar}
                      >
                        <span className="relative z-10">{item.icon}</span>
                        <span className="relative z-10 font-medium">{item.label}</span>
                        <motion.div
                          className="absolute inset-0 bg-red-50 opacity-0 group-hover:opacity-100
                            transition-opacity duration-200"
                          layoutId="highlight"
                        />
                      </Link>
                    </motion.div>
                  ))}
                </div>

                <div className="absolute bottom-8 left-0 right-0 px-6">
                  <button
                    onClick={handleLogout}
                    className="flex items-center justify-center gap-2 w-full px-4 py-3
                      text-red-600 font-medium rounded-lg
                      hover:bg-red-50 transition-colors duration-200"
                  >
                    <LogoutOutlined />
                    <span>Logout</span>
                  </button>
                </div>
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </>
  );
};

export default Sidebar;
