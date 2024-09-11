import React from "react";
import Banner1Img from "../../assets/dark_bg.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronLeft,
  faChevronRight,
  faWeight,
} from "@fortawesome/free-solid-svg-icons";

type BannerProps = {
  onLeftClick: () => void;
  onRightClick: () => void;
};

const Banner1: React.FC<BannerProps> = ({ onLeftClick, onRightClick }) => {
  return (
    <div>
      <div
        className="flex items-center justify-between h-[500px] text-white bg-cover bg-center px-10"
        style={{
          backgroundImage: `url(${Banner1Img.src})`,
        }}
      >
        <div
          className="flex items-center justify-center w-12 h-12 rounded-full cursor-pointer text-8xl text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400"
          onClick={onLeftClick}
        >
          <FontAwesomeIcon
            icon={faChevronLeft}
            style={{ fontSize: "1rem",}}
          />
        </div>
        <div className="flex flex-col h-auto w-1/2 text-center items-center">
          <div className="flex flex-col p-4 items-center">
            <p className="text-6xl font-bold">Welcome to Sport Complex</p>
            <span className="mt-5 w-2/3 text-lg">
              Reserve your spot and never miss out! Easily schedule your
              favorite sports activities with just a few clicks.
            </span>
          </div>
          <div className="mt-3 p-3 bg-red-400 w-fit rounded-md text-white text-xs font-bold">
            <a href="/pages/registration">
              <button type="button">Learn More</button>
            </a>
          </div>
        </div>
        <div
          className="flex items-center justify-center w-11 h-11 rounded-full cursor-pointer text-white border-white border-2 hover:border-yellow-400 hover:text-yellow-400"
          onClick={onRightClick}
        >
          <FontAwesomeIcon
            icon={faChevronRight}
            style={{ fontSize: "1rem"}}
          />
        </div>
      </div>
    </div>
  );
};

export default Banner1;
