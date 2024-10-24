"use client";

import React, { useState } from "react";
import NavBar from "../../components/navbar/navbar";
import BannerMain from "../../components/section/banner/banner_main";
import NewsReport from "../../components/section/new_report/report";
import LinkRelease from "../../components/section/link_to_another_website/link_release";
import Footer from "../../components/section/footer/footer";
import "../../css/banner.css";

const homepage: React.FC = () => {
  return (
    <div className="">
      <NavBar />
      <BannerMain />
      <NewsReport />
      <LinkRelease />
      <br />
      <br />
      <br />
      <Footer />
    </div>
  );
};

export default homepage;
