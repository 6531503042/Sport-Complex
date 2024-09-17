import { faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";

const search_bar = () => {
  return (
    <div className="border-gray-400 border py-2">
      <input type="text" placeholder="Search . . ." className="text-xs border-none"/>
      <button className="px-4">
        <FontAwesomeIcon icon={faSearch} />
      </button>
    </div>
  );
};

export default search_bar;
