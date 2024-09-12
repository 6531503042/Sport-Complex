"use client";

import React, { useState } from "react";
import NavBar from "../../components/navbar/navbar";
import Banner1 from "../../components/banner/banner1";
import Banner2 from "../../components/banner/banner2";
import Banner3 from "../../components/banner/banner3";
import "../../css/banner.css"

const HomePage: React.FC = () => {
  const [currentBanner, setCurrentBanner] = useState<1 | 2 | 3>(1);

  const handleLeftClick = () => {
    setCurrentBanner((prev) => {
      if (prev === 1) return 3 as 1 | 2 | 3; 
      return (prev - 1) as 1 | 2 | 3;
    });
  };

  const handleRightClick = () => {
    setCurrentBanner((prev) => {
      if (prev === 3) return 1 as 1 | 2 | 3; 
      return (prev + 1) as 1 | 2 | 3; 
    });
  };

  return (
    <div>
      <NavBar />
      <div className="banner-container">
        {currentBanner === 1 && (
          <Banner1 onLeftClick={handleLeftClick} onRightClick={handleRightClick} />
        )}
        {currentBanner === 2 && (
          <Banner2 onLeftClick={handleLeftClick} onRightClick={handleRightClick} />
        )}
        {currentBanner === 3 && (
          <Banner3 onLeftClick={handleLeftClick} onRightClick={handleRightClick} />
        )}
      </div>
      <div className="p-4">Breaking News ! ! ! !</div>
    </div>
  );
};

export default HomePage;
