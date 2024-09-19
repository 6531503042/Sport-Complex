import React, { useState, useEffect } from "react";
import Banner1Img from "../../assets/dark_bg.jpg";
import Banner2Img from "../../assets/banner_2.jpg";
import Banner3Img from "../../assets/banner_1.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronLeft,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";

const BannerMain: React.FC = () => {
  const banners = [
    {
      image: Banner1Img,
      title: "Welcome to Sport Complex",
      description:
        "Reserve your spot and never miss out! Easily schedule your favorite sports activities with just a few clicks.",
      link: "/pages/registration",
    },
    {
      image: Banner2Img,
      title: "Stay Active, Stay Healthy",
      description:
        "Discover a variety of sports and fitness programs tailored to your needs.",
      link: "/",
    },
    {
      image: Banner3Img,
      title: "Join a Community",
      description:
        "Meet like-minded individuals and engage in exciting activities.",
      link: "/",
    },
  ];

  const [currentBanner, setCurrentBanner] = useState(0);
  const [isFading, setIsFading] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      setIsFading(true);
      setTimeout(() => {
        setCurrentBanner((prev) =>
          prev === banners.length - 1 ? 0 : prev + 1
        );
        setIsFading(false);
      }, 100);
    }, 60000);
    return () => clearInterval(interval);
  }, [banners.length]);

  const handleLeftClick = () => {
    setIsFading(true);
    setTimeout(() => {
      setCurrentBanner((prev) => (prev === 0 ? banners.length - 1 : prev - 1));
      setIsFading(false);
    }, 100);
  };

  const handleRightClick = () => {
    setIsFading(true);
    setTimeout(() => {
      setCurrentBanner((prev) => (prev === banners.length - 1 ? 0 : prev + 1));
      setIsFading(false);
    }, 100);
  };

  const handleDotClick = (index: number) => {
    setIsFading(true);
    setTimeout(() => {
      setCurrentBanner(index);
      setIsFading(false);
    }, 100);
  };

  return (
    <div className="banner_container flex items-center h-[500px] text-white bg-cover bg-center">
      <div
        className={`flex flex-col justify-between items-center h-full w-screen transition-opacity duration-300 ease-in-out py-5 ${
          isFading ? "opacity-0" : "opacity-100"
        }`}
        style={{
          backgroundImage: `url(${banners[currentBanner].image.src})`,
          backgroundSize: "cover",
          backgroundPosition: "center",
        }}
      >
        <div className="invisible">
          block
        </div>
        <div className="flex flex-row items-center w-full px-10 justify-between">
          <button
            className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-8xl text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400 transition-all duration-100"
            onClick={handleLeftClick}
          >
            <FontAwesomeIcon
              icon={faChevronLeft}
              style={{ fontSize: "1rem" }}
            />
          </button>

          <div className="flex flex-col h-auto w-1/2 text-center items-center">
            <div className="flex flex-col p-4 items-center">
              <p className="text-6xl font-bold">
                {banners[currentBanner].title}
              </p>
              <span className="mt-5 w-2/3 text-lg">
                {banners[currentBanner].description}
              </span>
            </div>
            <div className="cursor-pointer transition-all duration-200 mt-3 p-3 bg-transparent w-fit border-2 border-stone-400 rounded-md text-white text-xs font-bold hover:text-white hover:border-transparent hover:bg-orange-700">
              <a href={banners[currentBanner].link}>
                <button type="button">Learn More</button>
              </a>
            </div>
          </div>
          <button
            className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400 transition-all duration-100"
            onClick={handleRightClick}
          >
            <FontAwesomeIcon
              icon={faChevronRight}
              style={{ fontSize: "1rem" }}
            />
          </button>
        </div>
        <div className="transform flex space-x-2">
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
    </div>
  );
};

export default BannerMain;
