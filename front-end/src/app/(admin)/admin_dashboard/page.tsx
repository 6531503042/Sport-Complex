import React from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";

const page = () => {
  return (
    <div className="w-screen h-screen flex flex-row">
      <Sidebar />
      <div className="bg-white text-black w-full p-10 inline-flex justify-between">
        <div></div>
        <img src={Logo.src} alt="Logo" className="w-7 h-min" />
      </div>
    </div>
  );
};

export default page;
