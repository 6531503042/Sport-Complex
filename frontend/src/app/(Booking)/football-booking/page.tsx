"use client";

import React, { useState } from "react";
import { Button } from "@nextui-org/react";
import "./football.css";

interface Data {
  time: string;
  isAvailable: boolean;
  member: string;
}

function Football_Booking() {
  const data: Data[] = [
    { time: "10.00 - 11.00", isAvailable: true, member: "0/1" },
    { time: "11.00 - 12.00", isAvailable: false, member: "1/1" },
    { time: "12.00 - 13.00", isAvailable: true, member: "0/1" },
    { time: "13.00 - 14.00", isAvailable: false, member: "1/1" },
    { time: "14.00 - 15.00", isAvailable: true, member: "0/1" },
    { time: "15.00 - 16.00", isAvailable: true, member: "0/1" },
    { time: "16.00 - 17.00", isAvailable: true, member: "0/1" },
    { time: "17.00 - 18.00", isAvailable: true, member: "0/1" },
    { time: "18.00 - 19.00", isAvailable: true, member: "0/1" },
    { time: "19.00 - 20.00", isAvailable: true, member: "0/1" },
    { time: "20.00 - 21.00", isAvailable: true, member: "0/1" },
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

      if (window.innerWidth < 640) {
        // Mobile screen size (sm breakpoint in Tailwind)
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
            Football Booking
          </h1>

          {isMobileView ? (
            // Mobile View: Show the form instead of the time slots
            <div className="block sm:hidden">
              <div className="bg-white border border-gray-200 rounded-lg p-4 shadow-md">
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
                      className="name-input mt-1 block w-full px-3 py-3"
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
                      className="name-input mt-1 block w-full px-3 py-3"
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
                      className="name-input mt-1 block w-full px-3 py-3"
                    />
                    {errors.phone && (
                      <span className="text-red-500 text-sm">
                        {errors.phone}
                      </span>
                    )}
                  </label>
                  <button
                    type="submit"
                    className="bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600"
                  >
                    Booking
                  </button>
                  <button
                    onClick={handleBackToTimeSlots}
                    type="submit"
                    className="bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600"
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
            // Normal and larger screen view: Show time slots and form sliding in
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {data.map((item, index) => (
                  <div
                    key={index}
                    className={`bg-white border border-gray-200 rounded-lg p-4 shadow-md transition-transform duration-300 ease-in-out
                      ${
                        item.isAvailable
                          ? "cursor-pointer"
                          : "cursor-not-allowed"
                      }
                      ${
                        selectedCard === index && item.isAvailable
                          ? "bg-green-100 border-green-500"
                          : ""
                      }
                      ${
                        !item.isAvailable
                          ? "bg-gray-200 text-gray-500"
                          : "hover:bg-green-100"
                      }`}
                    onClick={() => handleCardClick(index, item.isAvailable)}
                  >
                    <div className="text-lg font-semibold mb-2 flex justify-between">
                      <div>{item.time}</div>
                      <div
                        className={`text-sm ${
                          item.isAvailable ? "text-black" : "text-gray-500"
                        }`}
                      >
                        {item.member}
                      </div>
                    </div>

                    <div
                      className={`text-sm ${
                        item.isAvailable ? "text-green-500" : "text-gray-500"
                      }`}
                    >
                      {item.isAvailable ? "Available" : "Unavailable"}
                    </div>
                  </div>
                ))}
              </div>

              <div
                className={`hidden sm:block transition-all duration-300 ease-in-out mt-6 p-4 bg-white border border-gray-200 rounded-lg shadow-md transform ${
                  selectedCard !== null
                    ? "translate-y-0 opacity-100"
                    : "translate-y-5 opacity-0"
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
                          className="name-input mt-1 block w-full px-3 py-3"
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
                          className="name-input mt-1 block w-full px-3 py-3"
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
                          className="name-input mt-1 block w-full px-3 py-3"
                        />
                        {errors.phone && (
                          <span className="text-red-500 text-sm">
                            {errors.phone}
                          </span>
                        )}
                      </label>
                      <button
                        type="submit"
                        className="bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600"
                      >
                        Booking
                      </button>
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
              <p className="text-gray-600 mb-6 ">
                You have successfully booked the slot.
              </p>

              {/* Close Button */}
              <button
                className="w-full bg-gradient-to-r from-green-500 to-green-400 text-white py-3 rounded-lg font-semibold shadow-lg hover:shadow-xl transition duration-200 ease-in-out transform hover:scale-105"
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

export default Football_Booking;
