import React, { useState, useEffect } from "react";
import Banner1Img from "../../../assets/dark_bg.jpg";
import Banner2Img from "../../../assets/banner_2.jpg";
import Banner3Img from "../../../assets/banner_1.jpg";
import Link from "next/link";
import NavigateBeforeIcon from "@mui/icons-material/NavigateBefore";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";

const BannerMain: React.FC = () => {
  const banners = [
    {
      image: Banner1Img,
      title: "Welcome to Sport Complex",
      description:
        "Reserve your spot and never miss out! Easily schedule your favorite sports activities with just a few clicks.",
      link: "/registration",
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
  const [isSliding, setIsSliding] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      handleRightClick();
    }, 10000);
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
    <div className="banner-container relative flex justify-center items-center overflow-hidden h-[500px] text-white">
      <button
        className="absolute top-1/2 left-10 z-20 p-2 flex items-center justify-center w-15 h-15 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400 transition-all"
        onClick={handleLeftClick}
      >
        <NavigateBeforeIcon style={{ fontSize: "2rem" }} />
      </button>

      <button
        className="absolute top-1/2 z-20 right-10 p-2 flex items-center justify-center w-15 h-15 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400 transition-all"
        onClick={handleRightClick}
      >
        <NavigateNextIcon style={{ fontSize: "2rem" }} />
      </button>

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
            <div className="banner_detail flex flex-col justify-center items-center w-full h-full bg-black bg-opacity-50 text-center p-5 sm:p-10">
              <h1 className="text-5xl sm:text-6xl md:text-7xl font-bold mb-3 sm:mb-5  w-2/3 sm:w-1/2">
                {banner.title}
              </h1>
              <p className="text-base sm:text-lg md:text-xl mb-5 w-1/3 sm:w-1/3 md:w-1/2 ">
                {banner.description}
              </p>
              <div className="cursor-pointer pt-5">
                <Link
                  href={banner.link}
                  className="py-2.5 px-3 sm:py-3.5 sm:px-5 lg:py-4.5 lg:px-7 transition-all duration-300 border-2 border-stone-200 text-stone-200 rounded-full text-xs sm:text-sm hover:border-transparent hover:shadow-lg hover:bg-yellow-500 hover:text-white"
                >
                  <button type="button" className="uppercase">
                    <span className="inline-flex flex-row items-center">
                      <p>Learn More</p>
                      <NavigateNextIcon className="ps-2 text-3xl" />
                    </span>
                  </button>
                </Link>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="absolute bottom-5 flex space-x-2 justify-center w-full z-20 transition-all duration-300">
        {banners.map((_, index) => (
          <span
            key={index}
            onClick={() => handleDotClick(index)}
            className={`cursor-pointer w-3 h-3 rounded-full border ${
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
