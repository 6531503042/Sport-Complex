"use client";

import React, { useState } from "react";
import Available from "@/app/assets/available.png";
import Unavailable from "@/app/assets/unavailable.png";
import Back from '@/app/assets/back.png'
import Member from '@/app/assets/member.png'
import './swimming.css'

interface Data {
  time: string;
  isAvailable: boolean;
  member: string;
}

function Swimming_Booking() {
  const data: Data[] = [
    { time: "10.00 - 12.00", isAvailable: true, member: "0/30" },
    { time: "12.00 - 14.00", isAvailable: false, member: "30/30" },
    { time: "14.00 - 16.00", isAvailable: true, member: "0/30" },
    { time: "16.00 - 18.00", isAvailable: false, member: "30/30" },
    { time: "18.00 - 20.00", isAvailable: true, member: "0/30" },
  ];

  const [selectedCard, setSelectedCard] = useState<number | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    id: "",
    phone: "",
  });
  const [errors, setErrors] = useState({
    name: "",
    id: "",
    phone: "",
  });
  const [isBookingSuccessful, setIsBookingSuccessful] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false); // New state for mobile view

  const handleCardClick = (index: number, isAvailable: boolean) => {
    if (isAvailable) {
      setSelectedCard(index === selectedCard ? null : index);
      setErrors({ name: "", id: "", phone: "" });

      if (window.innerWidth < 640) { // Mobile screen size (sm breakpoint in Tailwind)
        setIsMobileView(true); // Switch to form view in mobile
      }
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const newErrors = {
      name: formData.name ? "" : "Name is required.",
      id: formData.id ? "" : "ID is required.",
      phone: formData.phone ? "" : "Phone number is required.",
    };
    setErrors(newErrors);

    const hasErrors = Object.values(newErrors).some((error) => error !== "");
    if (!hasErrors) {
      console.log("Form submitted successfully", formData);
      setIsBookingSuccessful(true);

      setFormData({
        name: "",
        id: "",
        phone: "",
      });
      setSelectedCard(null);
      setIsMobileView(false); // Reset view after booking
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleBackToTimeSlots = () => {
    setIsMobileView(false); // Go back to time slots in mobile
  };

  return (
    <>
      <div className="flex flex-col items-center h-screen p-6">
        <div className="w-full max-w-[1189px] bg-[#FEFFFE] border-gray border rounded-3xl drop-shadow-2xl p-5">
          <h1 className="text-4xl font-bold my-10 text-black text-center">
            Swimming Booking
          </h1>

          {isMobileView ? (
            // Mobile View: Show the form instead of the time slots
            <div className="block sm:hidden ">
              <div className="bg-white border border-gray-200 rounded-lg p-4 shadow-md">
                <form onSubmit={handleSubmit}>
                  <div className="my-3 ">
                    <img src={Back.src} alt="back" onClick={handleBackToTimeSlots} className="border shadow-xl p-2 rounded-md cursor-pointer hover:bg-gray-200" width={40}/>
                  </div>

                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700 py-2">
                      Name
                    </span>
                    <input
                      type="text"
                      name="name"
                      value={formData.name}
                      onChange={handleChange}
                      placeholder="Enter your name"
                      className="name-input-swimming mt-1 block w-full px-3 py-3"
                    />
                    {errors.name && (
                      <span className="text-red-500 text-sm">
                        {errors.name}
                      </span>
                    )}
                  </label>
                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700 py-2">
                      Lecturer / Staff / Student ID
                    </span>
                    <input
                      type="text"
                      name="id"
                      value={formData.id}
                      onChange={handleChange}
                      placeholder="Enter your ID"
                      className="name-input-swimming mt-1 block w-full px-3 py-3"
                    />
                    {errors.id && (
                      <span className="text-red-500 text-sm">{errors.id}</span>
                    )}
                  </label>
                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700 py-2">
                      Phone Number
                    </span>
                    <input
                      type="number"
                      name="phone"
                      value={formData.phone}
                      onChange={handleChange}
                      placeholder="Enter your phone number"
                      className="name-input-swimming mt-1 block w-full px-3 py-3"
                    />
                    {errors.phone && (
                      <span className="text-red-500 text-sm">
                        {errors.phone}
                      </span>
                    )}
                  </label>

                  <button
                    type="submit"
                    className="font-bold bg-blue-500 text-white px-5 py-2.5 my-5 rounded-md drop-shadow-2xl hover:bg-blue-600"
                  >
                    Booking
                  </button>

                  <h2 className="text-xl font-bold mb-4 text-end">
                    {selectedCard !== null && data[selectedCard].time}
                  </h2>
                </form>
              </div>
            </div>
          ) : (

            // Normal screen view
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 gap-6">
                {data.map((item, index) => (
                  <div
                  key={index}
                  className={`bg-white border border-gray-200 rounded-lg p-6 shadow-md transition-transform duration-300 ease-in-out
                    ${
                      item.isAvailable
                        ? "cursor-pointer hover:scale-105 hover:shadow-lg"
                        : "cursor-not-allowed"
                    }
                    ${
                      selectedCard === index && item.isAvailable
                        ? "bg-blue-200 border-blue-400 scale-105 shadow-lg"
                        : ""
                    }
                    ${
                      !item.isAvailable
                        ? "bg-gray-200 text-gray-700"
                        : "hover:bg-blue-200"
                    }`}
                    onClick={() => handleCardClick(index, item.isAvailable)}
                  >
                    <div className="text-lg font-semibold mb-2 flex justify-between">
                      <div>{item.time}</div>
                      <div
                        className={`text-l ${item.isAvailable ? "text-black" : "text-gray-500"}`}
                      >
                        <div className="flex"><img src={Member.src} alt="member" width={24}/> <span className="mx-1"></span>
                        {item.member}</div>
                        
                      </div>
                    </div>

                    <div>
                        {item.isAvailable ? (
                          <div className="flex"><img
                            src={Available.src}
                            alt="Available"
                            className="w-6 h-6"
                          /><span className="ml-2">Available</span></div>
                        ) : (
                          <div className="flex"><img
                            src={Unavailable.src}
                            alt="Unavailable"
                            className="w-6 h-6"
                          /><span className="ml-2">Unavailable</span></div>
                        )}
                      </div>
                  </div>
                ))}
              </div>

              <div
                className={`hidden sm:block transition-all duration-300 ease-in-out mt-6 p-4 bg-white border border-gray-200 rounded-lg shadow-md transform ${
                  selectedCard !== null ? "translate-y-0 opacity-100" : "translate-y-5 opacity-0"
                }`}
              >
                {selectedCard !== null && data[selectedCard].isAvailable && (
                  <>
                    <h2 className="text-xl font-bold mb-4">
                      Booking for {data[selectedCard].time}
                    </h2>
                    <form onSubmit={handleSubmit}>
                      <label className="block mb-4">
                        <span className="block text-sm font-medium text-gray-700 py-2">
                          Name
                        </span>
                        <input
                          type="text"
                          name="name"
                          value={formData.name}
                          onChange={handleChange}
                          placeholder="Enter your name"
                          className="name-input-swimming mt-1 block w-full px-3 py-3"
                        />
                        {errors.name && (
                          <span className="text-red-500 text-sm">
                            {errors.name}
                          </span>
                        )}
                      </label>
                      <label className="block mb-4">
                        <span className="block text-sm font-medium text-gray-700 py-2">
                          Lecturer / Staff / Student ID
                        </span>
                        <input
                          type="text"
                          name="id"
                          value={formData.id}
                          onChange={handleChange}
                          placeholder="Enter your ID"
                          className="name-input-swimming mt-1 block w-full px-3 py-3"
                        />
                        {errors.id && (
                          <span className="text-red-500 text-sm">
                            {errors.id}
                          </span>
                        )}
                      </label>
                      <label className="block mb-4">
                        <span className="block text-sm font-medium text-gray-700 py-2">
                          Phone Number
                        </span>
                        <input
                          type="number"
                          name="phone"
                          value={formData.phone}
                          onChange={handleChange}
                          placeholder="Enter your phone number"
                          className="name-input-swimming mt-1 block w-full px-3 py-3"
                        />
                        {errors.phone && (
                          <span className="text-red-500 text-sm">
                            {errors.phone}
                          </span>
                        )}
                      </label>
                      {/* Center the button */}
                      <div className="flex justify-center">
                        <button
                          type="submit"
                          className="font-semibold bg-blue-500 text-white px-6 py-3 my-5 rounded-md drop-shadow-2xl hover:bg-blue-600"
                        >
                          Booking
                        </button>
                      </div>
                    </form>
                  </>
                )}
              </div>
            </>
          )}
        </div>
        {isBookingSuccessful && (
          <div className="fixed inset-0 w-screen h-screen flex items-center justify-center z-50 bg-black bg-opacity-40 backdrop-blur-sm transition-opacity duration-300 ease-in-out">
            <div className="relative bg-white w-full max-w-sm mx-auto p-8 rounded-lg shadow-xl transform transition-all duration-500 ease-in-out scale-100">
              {/* Popup Content */}
              <h2 className="text-2xl font-bold text-gray-800 mb-2 text-center">
                Booking Successful!
              </h2>
              <p className="text-gray-600 mb-6 text-center">
                You have successfully booked the slot.
              </p>

              {/* Close Button */}
              <button
                className="w-full bg-gradient-to-r from-blue-500 to-blue-400 text-white py-3 rounded-lg font-semibold shadow-lg hover:shadow-xl transition duration-200 ease-in-out transform hover:scale-105"
                onClick={() => setIsBookingSuccessful(false)}
              >
                Close
              </button>
            </div>
          </div>
        )}
      </div>
    </>
  );
}

export default Swimming_Booking;
