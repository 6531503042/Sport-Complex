import React, { useState } from "react";
import {
  faCircleExclamation,
  faClockFour,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

interface NewsItem {
  date: string;
  title: string;
  description: string;
}

const report: React.FC = () => {
  const [selectedNews, setSelectedNews] = useState<NewsItem>({
    date: "06 sep 2024",
    title: "Newest",
    description:
      "Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut, aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa a debitis neque?",
  });

  const [activeIndex, setActiveIndex] = useState<number>(0);

  const newsItems: NewsItem[] = [
    {
      date: "06 sep 2024",
      title: "Newest",
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut, aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa a debitis neque?",
    },
    {
      date: "07 sep 2024",
      title: "Latest Update",
      description:
        "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Sint odit quasi molestias, officia at officiis vel saepe fugit soluta, facere quis repellat atque non ut tenetur eveniet nisi! Odit, odio?",
    },
    {
      date: "08 sep 2024",
      title: "Further News",
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Ut voluptatibus placeat minus consequatur qui laudantium perspiciatis reiciendis accusantium adipisci quos!",
    },
    {
      date: "09 sep 2024",
      title: "More Updates",
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Quisquam eius eos aspernatur provident nisi quia.",
    },
  ];
  const handleNewsClick = (item: NewsItem, index: number) => {
    setSelectedNews(item);
    setActiveIndex(index);
  };

  return (
    <div>
      <div className="gap-10 mt-10 mx-24 flex flex-row">
        <div className="flex flex-none w-2/3 flex-col">
          <h1 className="font-bold text-3xl">News</h1>
          <hr className="mt-3 border-zinc-900 rounded-full" />
          <div className="container_of_report_des">
            <div className="flex flex-col ">
              <div className="ps-5 cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50">
                <div className="flex flex-row justify-between">
                  <p className="pt-2 text-sm text-gray-500">
                    {selectedNews.date}
                  </p>
                  <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                    <FontAwesomeIcon
                      className="text-white me-1"
                      icon={faCircleExclamation}
                    />
                    {selectedNews.title}
                  </p>
                </div>
                <span>{selectedNews.description}</span>
                <div className="invisible">space</div>
              </div>
              <div className="ps-5 cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50">
                <div className="flex flex-row justify-between">
                  <p className="pt-2 text-sm text-gray-500">
                    {selectedNews.date}
                  </p>
                  <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                    <FontAwesomeIcon
                      className="text-white me-1"
                      icon={faCircleExclamation}
                    />
                    {selectedNews.title}
                  </p>
                </div>
                <span>{selectedNews.description}</span>
                <div className="invisible">space</div>
              </div>
            </div>
            <hr />
          </div>
        </div>
        <div className="flex flex-none w-1/3 flex-col border border-b-0">
          {newsItems.map((item, index) => (
            <div
              key={index}
              className={`cursor-pointer border-l-2 border-l-transparent hover:border-red-500 focus:bg-gray-200 transition-transform hover:shadow-lg ${
                activeIndex === index
                  ? "bg-white border-l-red-500 shadow-lg"
                  : ""
              }`}
              onClick={() => handleNewsClick(item, index)}
            >
              <p className="p-5 text-2xl text-gray-700">
                <FontAwesomeIcon
                  className="text-red-500 me-2"
                  icon={faClockFour}
                />
                {item.title}
              </p>
              <hr />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default report;
