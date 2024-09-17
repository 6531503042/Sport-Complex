"use client";

import React, { useState } from "react";
import NavBar from "../../components/navbar/navbar";
import BannerMain from "../../components/banner/banner_main"
import NewsReport from "../../components/new_report/report"
import "../../css/banner.css";


const homepage: React.FC = () => {

  return (
    <div className="">
      <NavBar />
      <BannerMain/>
      <NewsReport/>
      <footer>ass</footer>
    </div>
  );
};

export default homepage;
