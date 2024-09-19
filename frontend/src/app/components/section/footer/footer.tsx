import React from "react";
import Link from "next/link";
import Logo from "../../../assets/Logo.png";

const footer = () => {
  return (
    <div className="bg-slate-600">
      <div className="p-10">
        <Link
          href="/"
          className="inline-flex flex-row items-center flex-none gap-3.5 w-1/5"
        >
          <img src={Logo.src} alt="Logo" className="w-7" />
          <span className="flex flex-col border-l-2 w-max">
            <div className="ms-1">
              <span className="ms-1 inline-flex flex-row font-semibold text-xl">
                <p className="text-black ">SPORT.</p>
                <p className="text-gray-500">MFU</p>
              </span>
              <hr />
              <span className="text-zinc-300 ms-1 font-medium text-sm">
                SPORT COMPLEX
              </span>
            </div>
          </span>
        </Link>
      </div>
    </div>
  );
};

export default footer;
