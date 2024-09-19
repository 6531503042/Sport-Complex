import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useState } from "react";

const search_bar = () => {
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleSearch = () => {
    setIsExpanded((prev) => !prev);
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
        <button onClick={toggleSearch} className="flex justify-center items-center">
          <FontAwesomeIcon icon={faSearch} className="text-gray-500 px-4" />
        </button>
      </div>
    </div>
  );
};

export default search_bar;
