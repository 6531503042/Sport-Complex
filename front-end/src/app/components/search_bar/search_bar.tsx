import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useState } from "react";

const search_bar = () => {
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleSearch = () => {
    setIsExpanded((prev) => !prev);
  };

  return (
    <div className="searchbar_field flex w-full items-center justify-end">
      <div
        className={`flex justify-end items-center transition-all duration-300 ease-in-out ${
          isExpanded ? "w-64 border px-2 py-2 " : "w-10"
        } border-gray-400 rounded-full w-full`}
      >
        {isExpanded && (
          <input
            type="text"
            placeholder="Search . . ."
            className=" text-xs border-none w-full outline-none focus:outline-none focus:ring-0 transition-all duration-300"
          />
        )}
        <button
          onClick={toggleSearch}
          className={`flex justify-center items-center transition-all duration-300 mx-1 ${
            isExpanded
              ? "text-gray-500 hover:bg-orange-700 hover:text-white rounded-full px-2 py-2"
              : "border border-transparent px-2 py-2 rounded-full hover:bg-orange-700 hover:text-white text-gray-500"
          }`}
        >
          <FontAwesomeIcon icon={faSearch} />
        </button>
      </div>
    </div>
  );
};

export default search_bar;
