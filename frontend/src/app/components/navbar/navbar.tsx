"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import Image from "next/image";
import Logo from "../../assets/Logo.png";
import SideBar from "@/app/components/sidebar/sidebar";
import SearchBar from "../search_bar/search_bar";
import {
  FitnessCenter,
  SportsTennis,
  Pool,
  SportsSoccer,
  Email,
  Payment,
  Home,
  LogoutOutlined,
} from "@mui/icons-material";
import LoadingScreen from "../loading_screen/loading";
import styles from './navbar.module.css';

type NavBarProps = {
  activePage?: string;
};

const NavBar: React.FC<NavBarProps> = ({ activePage }) => {
  const [userName, setUserName] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const router = useRouter();

  // Add scroll effect
  const [isScrolled, setIsScrolled] = useState(false);

  // Add state for screen size
  const [isSmallScreen, setIsSmallScreen] = useState(false);

  useEffect(() => {
    const userData = localStorage.getItem("user");
    if (userData) {
      const user = JSON.parse(userData);
      setUserName(user.name);
      const id = user.id?.replace("user:", "") || "";
      setUserId(id);
    } else {
      router.replace("/login");
    }
  }, [router]);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 0);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    const checkScreenSize = () => {
      setIsSmallScreen(window.innerWidth < 768); // 768px is typical tablet breakpoint
    };

    // Check on mount
    checkScreenSize();

    // Add resize listener
    window.addEventListener('resize', checkScreenSize);

    // Cleanup
    return () => window.removeEventListener('resize', checkScreenSize);
  }, []);

  const truncateUserName = (name: string) => {
    return name?.length > 9 ? name.slice(0, 9) + "..." : name;
  };

  const handleLogout = () => {
    localStorage.removeItem("user");
    localStorage.removeItem("access_token");
    router.replace("/login");
  };

  const menuItems = [
    { href: "/homepage", icon: <Home />, label: "Home", page: "home" },
    { href: "/gym-booking", icon: <FitnessCenter />, label: "Gym", page: "gym" },
    { href: "/badminton-booking", icon: <SportsTennis />, label: "Badminton", page: "badminton" },
    { href: "/swimming-booking", icon: <Pool />, label: "Swimming", page: "swimming" },
    { href: "/football-booking", icon: <SportsSoccer />, label: "Football", page: "football" },
    { href: "/contact", icon: <Email />, label: "Contact", page: "contact" },
    { href: `/payment/user/${userId}`, icon: <Payment />, label: "Payment", page: "payment" },
  ];

  return (
    <>
      {loading && <LoadingScreen />}
      <div className={styles.navbarContainer}>
        <header className={`${styles.navbarWrapper} ${isScrolled ? styles.scrolled : ''}`}>
          <div className={styles.headerSection}>
            <div className={styles.headerContent}>
              {/* Logo Section - Always visible */}
              <motion.div
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.5 }}
              >
                <Link href="/homepage" className={styles.logoSection}>
                  <div className={styles.logoWrapper}>
                    <Image src={Logo} alt="Logo" fill className="object-contain" />
                  </div>
                  <div className={styles.logoDivider}>
                    <div className={styles.logoText}>
                      <span className={styles.logoTitle}>
                        <span className={styles.logoTitleMain}>SPORT</span>
                        <span className={styles.logoTitleSub}>.MFU</span>
                      </span>
                      <span className={styles.logoSubtitle}>SPORT COMPLEX</span>
                    </div>
                  </div>
                </Link>
              </motion.div>

              {/* Search Bar - Hide on small screens */}
              {!isSmallScreen && (
                <div className={styles.searchSection}>
                  <SearchBar />
                </div>
              )}

              {/* User Section - Modified for small screens */}
              <div className={styles.userSection}>
                {!isSmallScreen && (
                  <motion.div
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.5 }}
                    className="relative"
                  >
                    <button
                      onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                      className={styles.userButton}
                    >
                      <span className={styles.userName}>
                        {userName ? truncateUserName(userName) : "Loading..."}
                      </span>
                      <span className={styles.userDivider}>|</span>
                    </button>
                    
                    <AnimatePresence>
                      {isDropdownOpen && (
                        <motion.div
                          initial={{ opacity: 0, y: 10, scale: 0.95 }}
                          animate={{ opacity: 1, y: 0, scale: 1 }}
                          exit={{ opacity: 0, y: 10, scale: 0.95 }}
                          transition={{ duration: 0.2 }}
                          className={styles.dropdownMenu}
                        >
                          <button
                            onClick={handleLogout}
                            className={`${styles.dropdownItem} ${styles.dropdownItemDanger}`}
                          >
                            <LogoutOutlined className="w-4 h-4" />
                            <span>Logout</span>
                          </button>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </motion.div>
                )}

                <SideBar setLoading={setLoading} />
              </div>
            </div>
          </div>
        </header>

        {/* Navigation Menu - Hide on small screens */}
        {!isSmallScreen && (
          <nav className={styles.navigationMenu}>
            <div className={styles.navigationContent}>
              <ul className={styles.navigationList}>
                {menuItems.map(({ href, icon, label, page }) => (
                  <motion.li
                    key={page}
                    className={styles.navigationItem}
                    whileHover={{ y: -2 }}
                    whileTap={{ y: 0 }}
                  >
                    <Link
                      href={href}
                      className={`${styles.navigationLink} ${
                        activePage === page 
                          ? styles.navigationLinkActive 
                          : styles.navigationLinkInactive
                      }`}
                    >
                      <span className={styles.navigationIcon}>{icon}</span>
                      <span>{label}</span>
                      <div className={styles.activeIndicator} />
                    </Link>
                  </motion.li>
                ))}
              </ul>
            </div>
          </nav>
        )}
      </div>
    </>
  );
};

export default NavBar;

