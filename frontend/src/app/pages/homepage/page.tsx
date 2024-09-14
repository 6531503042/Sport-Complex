"use client";

import React, { useState } from "react";
import NavBar from "../../components/navbar/navbar";
import Banner1 from "../../components/banner/banner1";
import Banner2 from "../../components/banner/banner2";
import Banner3 from "../../components/banner/banner3";
import "../../css/banner.css";
import {
  faCircleExclamation,
  faClock,
  faClockFour,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const HomePage: React.FC = () => {
  const [currentBanner, setCurrentBanner] = useState<1 | 2 | 3>(1);

  const handleLeftClick = () => {
    setCurrentBanner((prev) => {
      if (prev === 1) return 3 as 1 | 2 | 3;
      return (prev - 1) as 1 | 2 | 3;
    });
  };

  const handleRightClick = () => {
    setCurrentBanner((prev) => {
      if (prev === 3) return 1 as 1 | 2 | 3;
      return (prev + 1) as 1 | 2 | 3;
    });
  };

  return (
    <div>
      <NavBar />
      <div className="banner-container">
        {currentBanner === 1 && (
          <Banner1
            onLeftClick={handleLeftClick}
            onRightClick={handleRightClick}
          />
        )}
        {currentBanner === 2 && (
          <Banner2
            onLeftClick={handleLeftClick}
            onRightClick={handleRightClick}
          />
        )}
        {currentBanner === 3 && (
          <Banner3
            onLeftClick={handleLeftClick}
            onRightClick={handleRightClick}
          />
        )}
      </div>
      <div className=" gap-10 mt-10 mx-24 flex flex-row">
        <div className="flex flex-none w-2/3 flex-col">
          <h1 className="font-bold text-3xl">News ! ! ! !</h1>
          <hr className="mt-3 border border-neutral-800" />
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50">
            <div className="flex flex-col ps-5">
              <div className="flex flew-row justify-between">
                <p className="pt-2 text-sm text-gray-500">06 sep 2024</p>
                <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                  <FontAwesomeIcon
                    className="text-white me-1"
                    icon={faCircleExclamation}
                  />
                  Newest
                </p>
              </div>
              <span>
                Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut,
                aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa
                a debitis neque?
              </span>
              <div className="invisible">space</div>
            </div>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50">
            <div className="flex flex-col ps-5">
              <div className="flex flew-row justify-between">
                <p className="pt-2 text-sm text-gray-500">06 sep 2024</p>
                <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                  <FontAwesomeIcon
                    className="text-white me-1"
                    icon={faCircleExclamation}
                  />
                  Newest
                </p>
              </div>
              <span>
                Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut,
                aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa
                a debitis neque?
              </span>
              <div className="invisible">space</div>
            </div>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50">
            <div className="flex flex-col ps-5">
              <div className="flex flew-row justify-between">
                <p className="pt-2 text-sm text-gray-500">06 sep 2024</p>
                <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                  <FontAwesomeIcon
                    className="text-white me-1"
                    icon={faCircleExclamation}
                  />
                  Newest
                </p>
              </div>
              <span>
                Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut,
                aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa
                a debitis neque?
              </span>
              <div className="invisible">space</div>
            </div>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50">
            <div className="flex flex-col ps-5">
              <div className="flex flew-row justify-between">
                <p className="pt-2 text-sm text-gray-500">06 sep 2024</p>
                <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                  <FontAwesomeIcon
                    className="text-white me-1"
                    icon={faCircleExclamation}
                  />
                  Newest
                </p>
              </div>
              <span>
                Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut,
                aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa
                a debitis neque?
              </span>
              <div className="invisible">space</div>
            </div>
            <hr />
          </div>
        </div>
        <div className="flex flex-none w-1/3 flex-col border border-b-0">
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
          <div className="cursor-pointer border-l-2 border-l-transparent hover:border-red-500 hover:bg-gray-50">
            <p className="p-5 text-2xl text-gray-700">
              <FontAwesomeIcon className="text-red-500 me-2" icon={faClockFour} />
              Lorem, ipsum dolor.
            </p>
            <hr />
          </div>
        </div>
      </div>
      <footer>ass</footer>
    </div>
  );
};

export default HomePage;
