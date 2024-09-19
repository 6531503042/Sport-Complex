import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useState } from "react";

const search_bar = () => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isOpen, setIsOpen] = useState(true);

  const toggleSearch = () => {
    setIsExpanded((prev) => !prev);
    setIsOpen((prev) => !prev);
  };

  return (
    <div className="flex justify-end w-full items-center">
      <div
        className={`flex items-center transition-all duration-500 ease-in-out ${
          isExpanded ? "w-64 border px-2 py-2" : "w-10"
        } border-gray-400 rounded-full w-full`}
      >
        {isExpanded && (
          <input
            type="text"
            placeholder="Search . . ."
            className="text-xs border-none w-full outline-none focus:outline-none focus:ring-0 transition-all duration-300"
          />
        )}
        <button
          onClick={toggleSearch}
          className={`flex justify-center items-center transition-all duration-300 mx-1 ${
            isOpen ? "border px-2 py-2 rounded-full border-gray-400 hover:border-orange-700 hover:text-orange-700 text-gray-500 " : "text-gray-500 hover:bg-orange-700 hover:text-white rounded-full px-2 py-2"
          }`}
        >
          <FontAwesomeIcon icon={faSearch} className="" />
        </button>
      </div>
    </div>
  );
};

export default search_bar;
