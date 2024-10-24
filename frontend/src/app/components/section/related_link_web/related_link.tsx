import React, { useState } from "react";

const ralate_link = () => {
  return (
    <div className="w-full h-full p-10 flex justify-center bg-gray-100 mt-10">
      <div className="flex flex-col gap-5 items-center w-full">
        <h1 className="text-4xl font-semibold">Related Links</h1>
        <div className="relative w-full h-auto overflow-x-visible">
          <ul className="grid grid-cols-5 gap-4 px-10 justify-center">
            <li className="bg-red-600 w-64 h-64">1</li>
            <li className="bg-red-600 w-64 h-64">2</li>
            <li className="bg-red-600 w-64 h-64">3</li>
            <li className="bg-red-600 w-64 h-64">4</li>
            <li className="bg-red-600 w-64 h-64">5</li>
            <li className="bg-red-600 w-64 h-64">6</li>
            <li className="bg-red-600 w-64 h-64">7</li>
            <li className="bg-red-600 w-64 h-64">8</li>
          </ul>
        </div>
        <div>
          <button>as</button>
        </div>
      </div>
    </div>
  );
};

export default ralate_link;
