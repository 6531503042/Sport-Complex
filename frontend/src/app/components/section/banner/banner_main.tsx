import React, { useState, useEffect } from "react";
import Banner1Img from "../../../assets/football.png";
import Banner2Img from "../../../assets/banner_2.jpg";
import Banner3Img from "../../../assets/banner_1.jpg";
import Link from "next/link";
import NavigateBeforeIcon from "@mui/icons-material/NavigateBefore";
import NavigateNextIcon from "@mui/icons-material/NavigateNext";
import { motion } from 'framer-motion';
import { useInView } from 'react-intersection-observer';

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

  const [ref, inView] = useInView({
    threshold: 0.3,
    triggerOnce: true
  });

  const bannerVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: {
        duration: 0.8,
        ease: "easeOut"
      }
    }
  };

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
    <motion.div 
      ref={ref}
      initial="hidden"
      animate={inView ? "visible" : "hidden"}
      variants={bannerVariants}
      className="banner-container relative flex justify-center items-center overflow-hidden h-[600px] text-white"
    >
      <motion.button
        whileHover={{ scale: 1.1 }}
        whileTap={{ scale: 0.9 }}
        className="absolute top-1/2 left-10 z-20 p-3 flex items-center justify-center rounded-full cursor-pointer bg-white/20 backdrop-blur-sm text-white border-2 border-white/50 hover:bg-white/30 transition-all duration-300"
        onClick={handleLeftClick}
      >
        <NavigateBeforeIcon className="text-3xl" />
      </motion.button>

      <motion.button
        whileHover={{ scale: 1.1 }}
        whileTap={{ scale: 0.9 }}
        className="absolute top-1/2 right-10 z-20 p-3 flex items-center justify-center rounded-full cursor-pointer bg-white/20 backdrop-blur-sm text-white border-2 border-white/50 hover:bg-white/30 transition-all duration-300"
        onClick={handleRightClick}
      >
        <NavigateNextIcon className="text-3xl" />
      </motion.button>

      <div
        className="flex transition-all duration-700 ease-in-out transform relative z-10"
        style={{
          transform: `translateX(-${currentBanner * 100}%)`,
        }}
      >
        {banners.map((banner, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="flex-shrink-0 flex justify-center w-full h-[600px] bg-cover bg-center relative"
            style={{
              backgroundImage: `url(${banner.image.src})`,
            }}
          >
            <div className="absolute inset-0 bg-gradient-to-b from-black/50 to-black/70" />
            <motion.div 
              className="banner_detail flex flex-col justify-center items-center w-full h-full text-center p-10 relative z-10"
              initial={{ y: 30, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
            >
              <motion.h1 
                className="text-6xl md:text-7xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-300"
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: 0.5 }}
              >
                {banner.title}
              </motion.h1>
              <motion.p 
                className="text-xl md:text-2xl mb-8 max-w-2xl text-gray-200"
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: 0.7 }}
              >
                {banner.description}
              </motion.p>
              <motion.div
                initial={{ y: 20, opacity: 0 }}
                animate={{ y: 0, opacity: 1 }}
                transition={{ delay: 0.9 }}
              >
                <Link
                  href={banner.link}
                  className="group relative inline-flex items-center px-8 py-4 text-lg font-medium"
                >
                  <span className="absolute inset-0 w-full h-full transition duration-200 ease-out transform translate-x-1 translate-y-1 bg-red-700 group-hover:-translate-x-0 group-hover:-translate-y-0"></span>
                  <span className="absolute inset-0 w-full h-full bg-white border-2 border-red-700 group-hover:bg-red-700"></span>
                  <span className="relative text-red-700 group-hover:text-white">Learn More</span>
                  <NavigateNextIcon className="relative ml-2 text-red-700 group-hover:text-white" />
                </Link>
              </motion.div>
            </motion.div>
          </motion.div>
        ))}
      </div>

      <div className="absolute bottom-8 flex space-x-3 justify-center w-full z-20">
        {banners.map((_, index) => (
          <motion.button
            key={index}
            whileHover={{ scale: 1.2 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => handleDotClick(index)}
            className={`w-3 h-3 rounded-full transition-all duration-300 ${
              currentBanner === index
                ? "bg-red-600 scale-125"
                : "bg-white/50 hover:bg-white"
            }`}
          />
        ))}
      </div>
    </motion.div>
  );
};

export default BannerMain;
