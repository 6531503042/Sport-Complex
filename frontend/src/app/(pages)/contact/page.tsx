import React from "react";
import NavBar from "../../components/navbar/navbar";
import {
  Person,
  Instagram,
  Facebook,
} from "@mui/icons-material";
import { Mail, Phone } from "lucide-react";

const page = () => {
  return (
    <div className="h-[645px]">
      <NavBar activePage="contact" />
      <div className="w-full h-full flex justify-center items-center">
        <div className="bg-gray-100 rounded-md shadow-lg inline-flex flex-col gap-5 p-10 h-1/2 w-1/2">
          <div className="inline-flex flex-row gap-4">
            <Mail />
            <span>6531503042@lamduan.mfu.ac.th</span>
          </div>
          <div className="inline-flex flex-row gap-4">
            <Person />
            <span>Nimitsu jung benji tanbooutor</span>
          </div>
          <div className="inline-flex flex-row gap-4">
            <Instagram />
            <span>plscallfrank</span>
          </div>
          <div className="inline-flex flex-row gap-4">
            <Facebook />
            <span>Klavivach Prajong</span>
          </div>
          <div className="inline-flex flex-row gap-4">
            <Phone />
            <span>080-236-85xx</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default page;
