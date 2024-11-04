"use client";

import React from "react";
import NavBar from "../../components/navbar/navbar";
import BannerMain from "../../components/section/banner/banner_main";
import NewsReport from "../../components/section/new_report/report";
import Footer from "../../components/section/footer/footer";
import "../../css/banner.css";
import useAuthRedirect from "../hooks/useAuthRedirect";

const Homepage: React.FC = () => {
  useAuthRedirect();

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
