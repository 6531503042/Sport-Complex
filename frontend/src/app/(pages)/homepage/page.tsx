"use client";

import React from "react";
import NavBar from "../../components/navbar/navbar";
import BannerMain from "../../components/section/banner/banner_main";
import NewsReport from "../../components/section/new_report/report";
import Footer from "../../components/section/footer/footer";
import "../../css/banner.css";
import useAuthRedirect from "../../components/hooks/useAuthRedirect";
import LoadingScreen from "../../components/loading_screen/loading"; // Import LoadingScreen component

const Homepage: React.FC = () => {
  const { loading, user } = useAuthRedirect(); // Get loading and user state from useAuthRedirect

  // While loading (checking auth), show the loading screen
  if (loading) {
    return <LoadingScreen />;
  }

  // After loading, display the homepage if user is authenticated
  return (
    <div>
      <NavBar />
      <BannerMain />
      <NewsReport />
      <br />
      <br />
      <br />
      <Footer />
    </div>
  );
};

export default Homepage;
