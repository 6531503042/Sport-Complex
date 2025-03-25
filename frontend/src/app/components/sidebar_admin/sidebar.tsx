"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import Logo from "../../assets/transparent.png";
import { motion, AnimatePresence } from "framer-motion";
import { ChevronLeft, ChevronRight } from "@mui/icons-material";
import { useRouter } from "next/navigation";
import styles from './sidebar.module.css';

type SidebarProps = {
  activePage?: string;
  isCollapsed?: boolean;
  onCollapse?: (collapsed: boolean) => void;
};

const menuItems = [
  {
    id: 'dashboard',
    label: 'Dashboard',
    icon: '📊',
    href: '/admin_dashboard',
    description: 'Overview & Analytics'
  },
  {
    id: 'users',
    label: 'User Management',
    icon: '👥',
    href: '/admin_usermanagement',
    description: 'Manage Users'
  },
  {
    id: 'facilities',
    label: 'Facility Management',
    icon: '🏟️',
    href: '/admin_facility',
    description: 'Manage Facilities'
  },
  {
    id: 'bookings',
    label: 'Booking Management',
    icon: '📅',
    href: '/admin_booking',
    description: 'Manage Bookings'
  },
  {
    id: 'payments',
    label: 'Payment Management',
    icon: '💰',
    href: '/admin_payment',
    description: 'Manage Payments'
  }
];

const Sidebar: React.FC<SidebarProps> = ({ 
  activePage, 
  isCollapsed = false, 
  onCollapse 
}) => {
  const [userName, setUserName] = useState<string | null>(null);
  const [showUserMenu, setShowUserMenu] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const userData = localStorage.getItem("user");
    if (userData) {
      const user = JSON.parse(userData);
      setUserName(user.name);
    } else {
      router.replace("/login");
    }
  }, [router]);

  const handleLogout = () => {
    localStorage.removeItem("user");
    localStorage.removeItem("access_token");
    router.replace("/login");
  };

  const handleCollapse = () => {
    onCollapse?.(!isCollapsed);
  };

  return (
    <motion.div 
      className={`${styles.sidebarContainer} ${isCollapsed ? styles.collapsed : ''}`}
      animate={{ width: isCollapsed ? '80px' : '280px' }}
      transition={{ duration: 0.3, ease: 'easeInOut' }}
    >
      <div className={styles.sidebarContent}>
        {/* Logo Section */}
        <div className={styles.logoSection}>
          <div className={styles.logoWrapper}>
            <Image 
              src={Logo} 
              alt="Logo" 
              width={40} 
              height={40}
              className={styles.logo} 
            />
            <AnimatePresence>
              {!isCollapsed && (
                <motion.div 
                  className={styles.logoText}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                >
                  <span className={styles.logoTitle}>
                    <span className={styles.logoTitleMain}>SPORT</span>
                    <span className={styles.logoTitleSub}>.MFU</span>
                  </span>
                  <span className={styles.logoSubtitle}>ADMIN PANEL</span>
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {/* Collapse Button */}
          <motion.button 
            className={styles.collapseButton}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={handleCollapse}
          >
            {isCollapsed ? <ChevronRight /> : <ChevronLeft />}
          </motion.button>
        </div>

        {/* Navigation Menu */}
        <nav className={styles.navigation}>
          {menuItems.map((item) => (
            <Link 
              key={item.id}
              href={item.href}
              className={`${styles.menuItem} ${activePage === item.id ? styles.active : ''}`}
            >
              <motion.div 
                className={styles.menuItemContent}
                whileHover={{ x: 5 }}
                whileTap={{ scale: 0.95 }}
              >
                <span className={styles.menuIcon}>{item.icon}</span>
                <AnimatePresence>
                  {!isCollapsed && (
                    <motion.div 
                      className={styles.menuDetails}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      exit={{ opacity: 0, x: -20 }}
                    >
                      <span className={styles.menuLabel}>{item.label}</span>
                      <span className={styles.menuDescription}>{item.description}</span>
                    </motion.div>
                  )}
                </AnimatePresence>
                {activePage === item.id && (
                  <motion.div 
                    className={styles.activeIndicator}
                    layoutId="activeIndicator"
                  />
                )}
              </motion.div>
            </Link>
          ))}
        </nav>

        {/* User Section */}
        <div className={styles.userSection}>
          <motion.div 
            className={styles.userProfile}
            onClick={() => setShowUserMenu(!showUserMenu)}
          >
            <div className={styles.userAvatar}>
              👨‍💼
            </div>
            <AnimatePresence>
              {!isCollapsed && (
                <motion.div 
                  className={styles.userInfo}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -20 }}
                >
                  <span className={styles.userName}>{userName}</span>
                  <span className={styles.userRole}>Administrator</span>
                </motion.div>
              )}
            </AnimatePresence>
          </motion.div>

          <AnimatePresence>
            {showUserMenu && !isCollapsed && (
              <motion.div 
                className={styles.userMenu}
                initial={{ opacity: 0, y: -20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
              >
                <motion.button 
                  className={`${styles.userMenuItem} ${styles.logoutButton}`}
                  onClick={handleLogout}
                  whileHover={{ x: 5 }}
                >
                  <span className={styles.menuIcon}>🚪</span>
                  <span>Logout</span>
                </motion.button>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </motion.div>
  );
};

export default Sidebar;
