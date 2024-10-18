"use client";

import React, { useEffect, useState, Fragment } from "react";
import "./football.css";
import Available from "@/app/assets/available.png";
import Unavailable from "@/app/assets/unavailable.png";
import Back from "@/app/assets/back.png";

interface Data {
  time: string;
  isAvailable: boolean;
}

function Football_Booking() {
  const data: Data[] = [
    { time: "10.00 - 11.00", isAvailable: true },
    { time: "11.00 - 12.00", isAvailable: false },
    { time: "12.00 - 13.00", isAvailable: true },
    { time: "13.00 - 14.00", isAvailable: false },
    { time: "14.00 - 15.00", isAvailable: true },
    { time: "15.00 - 16.00", isAvailable: true },
    { time: "16.00 - 17.00", isAvailable: true },
    { time: "17.00 - 18.00", isAvailable: true },
    { time: "18.00 - 19.00", isAvailable: true },
    { time: "19.00 - 20.00", isAvailable: true },
    { time: "20.00 - 21.00", isAvailable: true },
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
      setIsMobileView(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleBackToTimeSlots = () => {
    setIsMobileView(false);
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [slot, setSlot] = useState<any[]>([]);

  const getSlot = async () => {
  try {
    const resSlot = await fetch("http://localhost:1335/facility_v1/football/slot_v1/slots", {
      method: "GET",
      cache: "no-store",
    });

    if (!resSlot.ok) {
      throw new Error(`Failed to fetch: ${resSlot.statusText} (Status: ${resSlot.status})`);
    }
    const slotData = await resSlot.json();
    setSlot(slotData);
  } catch (error) {
    console.error('Error fetching slot data:', error);
  }
};

  useEffect(() => {
    getSlot();
    const intervalId = setInterval(() => {
      getSlot();
    }, 10000); 
    
    return () => {
      clearInterval(intervalId);
    };
  }, []);

  return (
    <>
      <div className="flex flex-col items-center h-screen p-6">
        <div className="w-full max-w-[1189px] bg-[#FEFFFE] border-gray border rounded-3xl drop-shadow-2xl p-5">
          <h1 className="text-4xl font-bold my-10 text-black text-center">
            Football Booking
          </h1>

          {isMobileView ? (
            // Mobile View: Show the form instead of the time slots
            <div className="block sm:hidden ">
              <div className="bg-white border border-gray-200 rounded-lg p-4 shadow-md">
                <form onSubmit={handleSubmit}>
                  <div className="my-3 ">
                    <img
                      src={Back.src}
                      alt="back"
                      onClick={handleBackToTimeSlots}
                      className="border shadow-xl p-2 rounded-md cursor-pointer hover:bg-gray-200"
                      width={40}
                    />
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
                      className="name-input-football mt-1 block w-full px-3 py-3"
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
                      className="name-input-football mt-1 block w-full px-3 py-3"
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
                      className="name-input-football mt-1 block w-full px-3 py-3"
                    />
                    {errors.phone && (
                      <span className="text-red-500 text-sm">
                        {errors.phone}
                      </span>
                    )}
                  </label>

                  <button
                    type="submit"
                    className="font-bold bg-green-500 text-white px-5 py-2.5 my-5 rounded-md drop-shadow-2xl hover:bg-green-600"
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
            // Normal and larger screen view
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {slot.length > 0 ? (
                  slot.map((lot) => (
                    <div key={lot._id} className="bg-white p-4 rounded shadow">
                      <ul>
                        <li>{lot.start_time}</li>
                        <li>{lot.end_time}</li>
                        <li>{lot.max_bookings}</li>
                        <li>{lot.current_bookings}</li>
                      </ul>
                    </div>
                  ))
                ) : (
                  <p className="bg-gray-300 p-3 my-3">No users available!</p>
                )}
              </div>


              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {slot.map((lot) => (
                  <div
                    key={lot._id}
                    className={`bg-white border border-gray-200 rounded-lg p-6 shadow-md transition-transform duration-300 ease-in-out
                      ${
                        lot.current_bookings
                          ? "cursor-not-allowed"
                          : "cursor-pointer hover:scale-105 hover:shadow-lg"
                      }
                      ${
                        selectedCard === lot && lot.current_bookings
                          ? ""
                          : "bg-green-200 scale-105 shadow-lg"
                      }
                      ${
                        !lot.current_bookings
                          ? "hover:bg-green-200"
                          : "bg-gray-300 text-gray-700"
                      }`}
                    onClick={() => handleCardClick(lot, lot.status)}
                  >
                    <div className="text-lg font-semibold flex justify-between items-center">
                      <div>{lot.start_time} - {lot.end_time}</div>
                      <div>
                        {lot.current_bookings ? (
                          
                          <img
                          src={Unavailable.src}
                          alt="Unavailable"
                          className="w-6 h-6"
                        />
                        ) : (
                          <img
                            src={Available.src}
                            alt="Available"
                            className="w-6 h-6"
                          />
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>



              {/* <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
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
                          ? "bg-green-200 border-green-500 scale-105 shadow-lg"
                          : ""
                      }
                      ${
                        !item.isAvailable
                          ? "bg-gray-200 text-gray-700"
                          : "hover:bg-green-200"
                      }`}
                    onClick={() => handleCardClick(index, item.isAvailable)}
                  >
                    <div className="text-lg font-semibold flex justify-between items-center">
                      <div>{item.time}</div>
                      <div>
                        {item.isAvailable ? (
                          <img
                            src={Available.src}
                            alt="Available"
                            className="w-6 h-6"
                          />
                        ) : (
                          <img
                            src={Unavailable.src}
                            alt="Unavailable"
                            className="w-6 h-6"
                          />
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div> */}

              <div
                className={`hidden sm:block transition-all duration-300 ease-in-out mt-6 p-4 bg-white border border-gray-200 rounded-lg shadow-md transform ${
                  selectedCard !== null
                    ? "translate-y-0 opacity-100"
                    : "translate-y-5 opacity-0"
                }`}
              >
                {selectedCard !== null && slot[selectedCard].isAvailable && (
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
                          className="name-input-football mt-1 block w-full px-3 py-3"
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
                          className="name-input-football mt-1 block w-full px-3 py-3"
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
                          className="name-input-football mt-1 block w-full px-3 py-3"
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
                          className="font-semibold bg-green-500 text-white px-6 py-3 my-5 rounded-md drop-shadow-2xl hover:bg-green-600"
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
