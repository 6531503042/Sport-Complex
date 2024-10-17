import React from "react";
import Sidebar from '../../components/sidebar_admin/sidebar'

const page = () => {
  return (
    <div className="w-screen h-screen flex flex-row">
      <Sidebar/>
      <div className="bg-white text-black w-full">Queue</div>
    </div>
  );
};

export default page;
