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

const Report: React.FC = () => {
  // Get the current date
  const currentDate = new Date();

  // Helper function to format the date as "DD MMM YYYY"
  const formatDate = (date: Date) => {
    return date.toLocaleDateString("en-GB", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  // Creating the newsItems with dynamic dates
  const newsItems: NewsItem[] = [
    {
      date: formatDate(new Date(currentDate)), // Current date
      title: "Newest",
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Aut, aliquam tempore. Rerum unde perspiciatis libero reiciendis ipsa a debitis neque?",
    },
    {
      date: formatDate(new Date(currentDate.setDate(currentDate.getDate() + 1))), // +1 day
      title: "Latest Update",
      description:
        "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Sint odit quasi molestias, officia at officiis vel saepe fugit soluta, facere quis repellat atque non ut tenetur eveniet nisi! Odit, odio?",
    },
    {
      date: formatDate(new Date(currentDate.setDate(currentDate.getDate() + 1))), // +1 more day
      title: "Further News",
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Ut voluptatibus placeat minus consequatur qui laudantium perspiciatis reiciendis accusantium adipisci quos!",
    },
    {
      date: formatDate(new Date(currentDate.setDate(currentDate.getDate() + 1))), // +1 more day
      title: "More Updates",
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Quisquam eius eos aspernatur provident nisi quia.",
    },
  ];

  const [selectedNews, setSelectedNews] = useState<NewsItem>(newsItems[0]);
  const [activeIndex, setActiveIndex] = useState<number>(0);

  const handleNewsClick = (item: NewsItem, index: number) => {
    setSelectedNews(item);
    setActiveIndex(index);
  };

  return (
    <div className="container mx-auto mt-10 px-4">
      <div className="flex flex-col lg:flex-row gap-10">
        {/* Left section (selected news content) */}
        <div className="lg:w-2/3 w-full">
          <h1 className="font-bold text-3xl mb-3">News</h1>
          <hr className="mt-3 border-zinc-900 rounded-full" />
          <div className="container_of_report_des mt-5">
            <div className="flex flex-col">
              <div className="ps-5 cursor-pointer border-l-2 border-l-transparent hover:border-l-red-700 hover:border-y-0 hover:border-l-2 hover:bg-gray-50 p-4 mb-4">
                <div className="flex flex-row justify-between">
                  <p className="pt-2 text-sm text-gray-500">{selectedNews.date}</p>
                  <p className="bg-red-600 p-1 text-sm text-center rounded-b-lg text-white">
                    <FontAwesomeIcon
                      className="text-white me-1"
                      icon={faCircleExclamation}
                    />
                    {selectedNews.title}
                  </p>
                </div>
                <span>{selectedNews.description}</span>
              </div>
            </div>
            <hr />
          </div>
        </div>

        {/* Right section (list of news items) */}
        <div
          className="lg:w-1/3 w-full flex lg:flex-col flex-row gap-4 lg:gap-0 overflow-x-auto"
        >
          {newsItems.map((item, index) => (
            <div
              key={index}
              className={`cursor-pointer border-l-2 border-l-transparent hover:border-red-500 focus:bg-gray-200 transition-transform hover:shadow-lg ${
                activeIndex === index
                  ? "bg-white border-l-red-500 shadow-lg"
                  : ""
              } flex-shrink-0 w-auto`}
              onClick={() => handleNewsClick(item, index)}
            >
              <p className="p-5 text-sm md:text-2xl text-gray-700 inline-flex flex-row items-center">
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

export default Report;
