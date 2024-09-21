import React, { useState, useEffect } from "react";
import Banner1Img from "../../../assets/dark_bg.jpg";
import Banner2Img from "../../../assets/banner_2.jpg";
import Banner3Img from "../../../assets/banner_1.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronLeft, faChevronRight } from "@fortawesome/free-solid-svg-icons";

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
    }, 60000);
    return () => clearInterval(interval);
  }, [currentBanner]);

  const handleLeftClick = () => {
    if (!isSliding) {
      setIsSliding(true);
      setTimeout(() => {
        setCurrentBanner((prev) => (prev === 0 ? banners.length - 1 : prev - 1));
        setIsSliding(false);
      }, 300);
    }
  };

  const handleRightClick = () => {
    if (!isSliding) {
      setIsSliding(true);
      setTimeout(() => {
        setCurrentBanner((prev) => (prev === banners.length - 1 ? 0 : prev + 1));
        setIsSliding(false);
      }, 300);
    }
  };

  const handleDotClick = (index: number) => {
    if (!isSliding) {
      setIsSliding(true);
      setTimeout(() => {
        setCurrentBanner(index);
        setIsSliding(false);
      }, 300);
    }
  };

  return (
    <div className="banner-container relative flex items-center overflow-hidden h-[500px] text-white">
      <div className="absolute inset-0 flex items-center justify-between z-10 px-5">
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
        className={`flex transition-transform duration-700 ease-in-out transform ${
          isSliding ? "translate-x-[-100%]" : "translate-x-[0]"
        }`}
        style={{
          transform: `translateX(-${currentBanner * 100}%)`,
        }}
      >
        {banners.map((banner, index) => (
          <div
            key={index}
            className="flex-shrink-0 w-full h-[500px] bg-cover bg-center"
            style={{
              backgroundImage: `url(${banner.image.src})`,
            }}
          >
            <div className="flex flex-col justify-center items-center h-full bg-black bg-opacity-40">
              <h1 className="text-6xl font-bold mb-5">{banner.title}</h1>
              <p className="text-lg mb-8 w-2/3 text-center">
                {banner.description}
              </p>
              <a
                href={banner.link}
                className="mt-3 p-3 bg-transparent border-2 border-stone-400 rounded-md text-white text-xs font-bold hover:text-white hover:border-transparent hover:bg-orange-700"
              >
                Learn More
              </a>
            </div>
          </div>
        ))}
      </div>
      <div className="absolute bottom-5 flex space-x-2 justify-center w-full z-10">
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
