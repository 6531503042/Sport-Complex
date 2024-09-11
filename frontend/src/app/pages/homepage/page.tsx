"use client"

import React, { useState } from "react";
import NavBar from "../../components/navbar/navbar";
import Banner1 from "../../components/banner/banner1";
import Banner2 from "../../components/banner/banner2";

const HomePage: React.FC = () => {
  const [currentBanner, setCurrentBanner] = useState<1 | 2>(1);

  const handleLeftClick = () => {
    setCurrentBanner((prev) => (prev === 1 ? 2 : 1));
  };

  const handleRightClick = () => {
    setCurrentBanner((prev) => (prev === 1 ? 2 : 1));
  };

  return (
    <div>
      <NavBar />
      {currentBanner === 1 ? (
        <Banner1 onLeftClick={handleLeftClick} onRightClick={handleRightClick} />
      ) : (
        <Banner2 onLeftClick={handleLeftClick} onRightClick={handleRightClick} />
      )}
            <div className="p-4">Breaking News ! ! ! !</div>
    </div>
  );
};

export default HomePage;
