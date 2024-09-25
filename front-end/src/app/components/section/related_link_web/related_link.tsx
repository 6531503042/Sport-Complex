import React, { useState } from "react";

const ralate_link = () => {
  return (
    <div className="w-full h-full p-10 flex justify-center">
      <div className="flex flex-col gap-5 items-center w-full">
        <h1 className="text-4xl font-semibold">Related Links</h1>
        <div className="relative w-full h-auto overflow-x-hidden">
          <ul className="w-full h-40  grid grid-cols-5 gap-4">
            <li className="w-full">1</li>
            <li className="w-full ">2</li>
            <li className="w-full ">3</li>
            <li className="w-full">4</li>
            <li className="w-full">5</li>
            <li className="w-full">6</li>
            <li className="w-full">7</li>
            <li className="w-full">8</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default ralate_link;
