"use client";

import React, { useState } from "react";
import dynamic from 'next/dynamic';
import Sidebar from '../../components/sidebar_admin/sidebar';
import Logo from "../../assets/Logo.png";
import { Tabs, Tab, Box, useMediaQuery, useTheme, IconButton, Drawer } from '@mui/material';
import { Toaster } from 'react-hot-toast';
import ErrorBoundary from './components/ErrorBoundary';
import LoadingComponent from './components/LoadingComponent';
import MenuIcon from '@mui/icons-material/Menu';

// Use dynamic imports to fix hydration issues
const FacilitySlotManager = dynamic(() => import('./components/FacilitySlotManager'), {
  loading: () => <LoadingComponent />,
  ssr: false
});

const BadmintonManager = dynamic(() => import('./components/BadmintonManager'), {
  loading: () => <LoadingComponent />,
  ssr: false
});

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function CustomTabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`facility-tabpanel-${index}`}
      aria-labelledby={`facility-tab-${index}`}
      className="text-gray-900"
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const AdminFacilityPage = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [mobileOpen, setMobileOpen] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const isTablet = useMediaQuery(theme.breakpoints.down('lg'));

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Desktop Sidebar */}
      {!isMobile && (
        <div className="sticky top-0 h-screen flex-shrink-0">
          <Sidebar activePage="admin_facility" />
        </div>
      )}

      {/* Mobile Drawer */}
      {isMobile && (
        <Drawer
          variant="temporary"
          anchor="left"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            '& .MuiDrawer-paper': { 
              width: 280,
              boxSizing: 'border-box',
              backgroundColor: '#7f1d1d',
            },
          }}
        >
          <Sidebar activePage="admin_facility" />
        </Drawer>
      )}

      {/* Main Content */}
      <div className="flex-1 min-w-0">
        <ErrorBoundary>
          <div className="p-4 md:p-6 lg:p-8">
            <div className="bg-white rounded-lg shadow-sm p-4 md:p-6 lg:p-8">
              {/* Header */}
              <div className="flex items-center justify-between mb-6">
                {isMobile && (
                  <IconButton
                    color="inherit"
                    aria-label="open drawer"
                    edge="start"
                    onClick={handleDrawerToggle}
                    sx={{ mr: 2 }}
                  >
                    <MenuIcon />
                  </IconButton>
                )}
                <div className="text-xl md:text-2xl font-semibold">Facility Management</div>
                <img src={Logo.src} alt="Logo" className="w-6 md:w-7" />
              </div>

              {/* Tabs */}
              <Box sx={{ width: '100%', mt: 2 }}>
                <Box sx={{ 
                  borderBottom: 1, 
                  borderColor: 'divider',
                  '.MuiTabs-flexContainer': {
                    overflowX: 'auto',
                    flexWrap: isTablet ? 'nowrap' : 'wrap',
                    '::-webkit-scrollbar': {
                      height: '4px',
                    },
                    '::-webkit-scrollbar-track': {
                      background: '#f1f1f1',
                    },
                    '::-webkit-scrollbar-thumb': {
                      background: '#888',
                      borderRadius: '2px',
                    },
                  }
                }}>
                  <Tabs 
                    value={activeTab} 
                    onChange={(_, newValue) => setActiveTab(newValue)}
                    variant={isTablet ? "scrollable" : "standard"}
                    scrollButtons={isTablet ? "auto" : false}
                    aria-label="facility tabs"
                    TabIndicatorProps={{
                      style: {
                        backgroundColor: '#7f1d1d',
                        height: '3px',
                        borderRadius: '3px',
                      }
                    }}
                    sx={{
                      '& .MuiTab-root': {
                        minWidth: isTablet ? 'auto' : 120,
                        padding: isTablet ? '12px 24px' : '16px 32px',
                        fontSize: isTablet ? '0.875rem' : '1rem',
                        fontWeight: 600,
                        color: '#1f2937',
                        textTransform: 'none',
                        '&.Mui-selected': {
                          color: '#7f1d1d',
                          fontWeight: 700,
                        },
                        '&:hover': {
                          color: '#991b1b',
                          backgroundColor: '#fee2e2',
                          borderRadius: '8px 8px 0 0',
                        },
                        transition: 'all 0.2s ease-in-out',
                      },
                      '& .MuiTabs-flexContainer': {
                        gap: '8px',
                      },
                    }}
                  >
                    {[
                      { label: "Fitness", icon: "🏋️" },
                      { label: "Swimming", icon: "🏊" },
                      { label: "Badminton", icon: "🏸" },
                      { label: "Football", icon: "⚽" },
                    ].map((tab) => (
                      <Tab 
                        key={tab.label}
                        label={
                          <div className="flex items-center gap-2 text-inherit">
                            <span>{tab.icon}</span>
                            <span>{tab.label}</span>
                          </div>
                        }
                        sx={{
                          '&.Mui-selected': {
                            backgroundColor: '#fef2f2',
                            borderRadius: '8px 8px 0 0',
                            color: '#7f1d1d',
                          },
                        }}
                      />
                    ))}
                  </Tabs>
                </Box>

                {/* Tab Panels with animation */}
                <Box 
                  sx={{ 
                    mt: 3,
                    position: 'relative',
                    minHeight: '400px',
                    color: 'inherit'
                  }}
                >
                  {[
                    { component: <FacilitySlotManager facilityName="fitness" />, name: "fitness" },
                    { component: <FacilitySlotManager facilityName="swimming" />, name: "swimming" },
                    { component: <BadmintonManager />, name: "badminton" },
                    { component: <FacilitySlotManager facilityName="football" />, name: "football" },
                  ].map((panel, index) => (
                    <CustomTabPanel 
                      key={panel.name}
                      value={activeTab} 
                      index={index}
                      sx={{
                        opacity: activeTab === index ? 1 : 0,
                        transform: `translateX(${(activeTab - index) * 20}px)`,
                        transition: 'all 0.3s ease-in-out',
                        position: activeTab === index ? 'relative' : 'absolute',
                        width: '100%',
                        top: 0,
                        color: 'inherit'
                      }}
                    >
                      {panel.component}
                    </CustomTabPanel>
                  ))}
                </Box>
              </Box>
            </div>
          </div>
        </ErrorBoundary>
      </div>
    </div>
  );
};

export default AdminFacilityPage;
