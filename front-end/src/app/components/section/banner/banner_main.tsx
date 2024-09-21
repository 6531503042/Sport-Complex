import React, { useState, useEffect } from "react";
import Banner1Img from "../../../assets/dark_bg.jpg";
import Banner2Img from "../../../assets/banner_2.jpg";
import Banner3Img from "../../../assets/banner_1.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronLeft,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";
import Link from "next/link";

const BannerMain: React.FC = () => {
  const banners = [
    {
      image: Banner1Img,
      title: "Welcome to Sport Complex",
      description:
        "Reserve your spot and never miss out! Easily schedule your favorite sports activities with just a few clicks.",
    },
    {
      image: Banner2Img,
      title: "Stay Active, Stay Healthy",
      description:
        "Discover a variety of sports and fitness programs tailored to your needs.",
    },
    {
      image: Banner3Img,
      title: "Join a Community",
      description:
        "Meet like-minded individuals and engage in exciting activities.",
    },
  ];

  const [currentBanner, setCurrentBanner] = useState(0);
  const [isSliding, setIsSliding] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      handleRightClick();
    }, 60000);
    return () => clearInterval(interval);
  }, [currentBanner]);

  const handleLeftClick = () => {
    if (!isSliding) {
      setIsSliding(true);
      setTimeout(() => {
        setCurrentBanner((prev) =>
          prev === 0 ? banners.length - 1 : prev - 1
        );
        setIsSliding(false);
      }, 100);
    }
  };

  const handleRightClick = () => {
    if (!isSliding) {
      setIsSliding(true);
      setTimeout(() => {
        setCurrentBanner((prev) =>
          prev === banners.length - 1 ? 0 : prev + 1
        );
        setIsSliding(false);
      }, 100);
    }
  };

  const handleDotClick = (index: number) => {
    if (!isSliding) {
      setIsSliding(true);
      setTimeout(() => {
        setCurrentBanner(index);
        setIsSliding(false);
      }, 100);
    }
  };

  return (
    <div className="banner-container relative flex items-center overflow-hidden h-[500px] text-white">
      <div className="absolute inset-0 flex items-center justify-between z-20 px-5">
        <button
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400 transition-all"
          onClick={handleLeftClick}
        >
          <FontAwesomeIcon icon={faChevronLeft} style={{ fontSize: "1rem" }} />
        </button>

        <button
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400 transition-all"
          onClick={handleRightClick}
        >
          <FontAwesomeIcon icon={faChevronRight} style={{ fontSize: "1rem" }} />
        </button>
      </div>
      <div
        className={`flex transition-transform duration-700 ease-in-out transform relative z-10 ${
          isSliding ? "translate-x-[-100%]" : "translate-x-[0]"
        }`}
        style={{
          transform: `translateX(-${currentBanner * 100}%)`,
        }}
      >
        {banners.map((banner, index) => (
          <div
            key={index}
            className="flex-shrink-0 flex justify-center w-full h-[500px] bg-cover bg-center"
            style={{
              backgroundImage: `url(${banner.image.src})`,
            }}
          >
            <div className="flex flex-col justify-center items-center w-1/2 h-full text-center">
              <h1 className="text-6xl font-bold mb-5">{banner.title}</h1>
              <p className="text-lg mb-8 w-2/3 text-center">
                {banner.description}
              </p>
            </div>
            <div className="absolute whitespace-nowrap cursor-pointer p-5 bottom-32">
              <a href="/">
                <button type="button">LEARN MORE</button>
              </a>
            </div>
          </div>
        ))}
      </div>
      <div className="absolute bottom-5 flex space-x-2 justify-center w-full z-20">
        {banners.map((_, index) => (
          <span
            key={index}
            onClick={() => handleDotClick(index)}
            className={`cursor-pointer w-4 h-4 rounded-full border-2 ${
              currentBanner === index
                ? "bg-yellow-500 border-yellow-500"
                : "bg-transparent border-gray-200"
            }`}
          ></span>
        ))}
      </div>
    </div>
  );
};

export default BannerMain;
