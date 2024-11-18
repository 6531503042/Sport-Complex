import React from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";

const page = () => {
  return (
    <div className="w-screen h-screen flex relative">
      <Sidebar activePage="admin_payment" />
      <div className="bg-white text-black w-full p-10 flex flex-col overflow-y-auto">
        <div className="inline-flex justify-between w-full items-end">
          <div className="text-lg font-medium">Payment</div>
          <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        </div>
        <br />
        <div className="bg-zinc-500 h-[1px] rounded-lg"></div>
        <br />
        <div className="w-full flex justify-center"></div>
        <br />
        <br />
      </div>
    </div>
  );
};

export default page;
