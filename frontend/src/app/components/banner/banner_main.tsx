import React, { useState, useEffect } from "react";
import Banner1Img from "../../assets/dark_bg.jpg";
import Banner2Img from "../../assets/banner_2.jpg";
import Banner3Img from "../../assets/banner_1.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronLeft, faChevronRight } from "@fortawesome/free-solid-svg-icons";

type BannerProps = {
  onLeftClick: () => void;
  onRightClick: () => void;
};

const BannerCarousel: React.FC = () => {
  const banners = [
    {
      image: Banner1Img,
      title: "Welcome to Sport Complex",
      description: "Reserve your spot and never miss out! Easily schedule your favorite sports activities with just a few clicks.",
      link: "/pages/registration",
    },
    {
      image: Banner2Img,
      title: "Stay Active, Stay Healthy",
      description: "Discover a variety of sports and fitness programs tailored to your needs.",
      link: "/",
    },
    {
      image: Banner3Img,
      title: "Join a Community",
      description: "Meet like-minded individuals and engage in exciting activities.",
      link: "/",
    },
  ];

  const [currentBanner, setCurrentBanner] = useState(0);
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentBanner((prev) => (prev === banners.length - 1 ? 0 : prev + 1));
    }, 60000); 

    return () => clearInterval(interval); 
  }, [banners.length]);

  const handleLeftClick = () => {
    setCurrentBanner((prev) => (prev === 0 ? banners.length - 1 : prev - 1));
  };

  const handleRightClick = () => {
    setCurrentBanner((prev) => (prev === banners.length - 1 ? 0 : prev + 1));
  };

  return (
    <div className="flex items-center h-[500px] text-white bg-cover bg-center px-10"
    style={{
        backgroundImage: `url(${banners[currentBanner].image.src})`,
        backgroundSize: "cover",  
        backgroundPosition: "center", 
      }}
    >
      <div className="flex flex-row items-center w-screen justify-between">
        <button
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-8xl text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400"
          onClick={handleLeftClick}
        >
          <FontAwesomeIcon icon={faChevronLeft} style={{ fontSize: "1rem" }} />
        </button>
        <div className="flex flex-col h-auto w-1/2 text-center items-center">
          <div className="flex flex-col p-4 items-center">
            <p className="text-6xl font-bold">{banners[currentBanner].title}</p>
            <span className="mt-5 w-2/3 text-lg">
              {banners[currentBanner].description}
            </span>
          </div>
          <div className="mt-3 p-3 bg-transparent w-fit border-2 border-stone-400 rounded-md text-white text-xs font-bold">
            <a href={banners[currentBanner].link}>
              <button type="button">Learn More</button>
            </a>
          </div>
        </div>
        <button
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400"
          onClick={handleRightClick}
        >
          <FontAwesomeIcon icon={faChevronRight} style={{ fontSize: "1rem" }} />
        </button>
      </div>
    </div>
  );
};

export default BannerCarousel;
