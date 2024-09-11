import React from "react";
import Banner2Img from "../../assets/banner_2.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronLeft,
  faChevronRight,
} from "@fortawesome/free-solid-svg-icons";

type BannerProps = {
  onLeftClick: () => void;
  onRightClick: () => void;
};

const Banner2: React.FC<BannerProps> = ({ onLeftClick, onRightClick }) => {
  return (
    <div>
      <div
        className="flex items-center justify-between h-[500px] text-white bg-cover bg-center px-10"
        style={{
          backgroundImage: `url(${Banner2Img.src})`,
        }}
      >
        <div
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-8xl text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400"
          onClick={onRightClick}
        >
          <FontAwesomeIcon icon={faChevronLeft} style={{ fontSize: "1rem" }} />
        </div>
        <div className="flex flex-col h-auto w-1/2 text-center items-center">
          <div className="flex flex-col p-4 items-center">
            <p className="text-6xl font-bold">
              Get Ready for an Active Lifestyle
            </p>
            <span className="mt-5 w-2/3 text-lg">
              Enhance your fitness with our diverse sports offerings. Check out
              our new facilities and book your favorite activities with ease!
            </span>
          </div>
          <div className="mt-3 p-3 bg-red-400 w-fit rounded-md text-white text-xs font-bold">
            <a href="/pages/registration">
              <button type="button">Learn More</button>
            </a>
          </div>
        </div>
        <div
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-8xl text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400"
          onClick={onRightClick}
        >
          <FontAwesomeIcon icon={faChevronRight} style={{ fontSize: "1rem" }} />
        </div>
      </div>
    </div>
  );
};

export default Banner2;
