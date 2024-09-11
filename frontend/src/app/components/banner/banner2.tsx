import React from "react";
import Banner2Img from "../../assets/banner_2.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronCircleLeft,
  faChevronCircleRight,
} from "@fortawesome/free-solid-svg-icons";

type BannerProps = {
  onLeftClick: () => void;
  onRightClick: () => void;
};

const Banner2: React.FC<BannerProps> = ({ onLeftClick, onRightClick }) => {
  return (
    <div>
      <div
        className="flex items-center justify-between h-[400px] text-white bg-cover bg-center px-10"
        style={{
          backgroundImage: `url(${Banner2Img.src})`,
        }}
      >
        <button onClick={onLeftClick}>
          <FontAwesomeIcon
            icon={faChevronCircleLeft}
            style={{ fontSize: "2.5rem", color: "#303030" }}
          />
        </button>
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
        <button onClick={onRightClick}>
          <FontAwesomeIcon
            icon={faChevronCircleRight}
            style={{ fontSize: "2.5rem", color: "#303030" }}
          />
        </button>
      </div>
    </div>
  );
};

export default Banner2;
