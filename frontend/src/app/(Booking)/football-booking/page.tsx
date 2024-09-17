"use client";

import React, { useState } from "react";

interface Data {
  time: string;
  isAvailable: boolean;
  member: string;
}

function Football() {
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

  const handleCardClick = (index: number, isAvailable: boolean) => {
    if (isAvailable) {
      setSelectedCard(index === selectedCard ? null : index);
      setErrors({ name: "", id: "", phone: "" }); 
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

    const hasErrors = Object.values(newErrors).some(error => error !== "");
    if (!hasErrors) {
      console.log("Form submitted successfully", formData);
      setIsBookingSuccessful(true);


      setFormData({
        name: "",
        id: "",
        phone: "",
      });
      setSelectedCard(null);

    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <>
      <div className="flex flex-col items-center h-screen p-6">
        <div className="w-full max-w-[1189px] bg-[#FEFFFE] border-gray border rounded-3xl drop-shadow-2xl p-5">
          <h1 className="text-4xl font-bold my-10 text-black text-center">Football Booking</h1>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {data.map((item, index) => (
              <div
                key={index}
                className={`bg-white border border-gray-200 rounded-lg p-4 shadow-md transition-transform duration-300 ease-in-out
                  ${item.isAvailable ? "cursor-pointer" : "cursor-not-allowed"}
                  ${selectedCard === index && item.isAvailable ? "bg-green-100 border-green-500" : ""}
                  ${!item.isAvailable ? "bg-gray-200 text-gray-500" : "hover:bg-green-100"}`}
                onClick={() => handleCardClick(index, item.isAvailable)}
              >
                <div className="text-lg font-semibold mb-2 flex justify-between">
                  <div>{item.time}</div>
                  <div
                    className={`text-sm ${item.isAvailable ? "text-black" : "text-gray-500"}`}
                  >
                    {item.member}
                  </div>
                </div>

                <div
                  className={`text-sm ${item.isAvailable ? "text-green-500" : "text-gray-500"}`}
                >
                  {item.isAvailable ? "Available" : "Unavailable"}
                </div>
              </div>
            ))}
          </div>

          <div
            className={`transition-all duration-300 ease-in-out mt-6 p-4 bg-white border border-gray-200 rounded-lg shadow-md transform ${
              selectedCard !== null ? "translate-y-0 opacity-100" : "translate-y-5 opacity-0"
            }`}
          >
            {selectedCard !== null && data[selectedCard].isAvailable && (
              <>
                <h2 className="text-xl font-bold mb-4">Booking for {data[selectedCard].time}</h2>
                <form onSubmit={handleSubmit}>
                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700">Name</span>
                    <input
                      type="text"
                      name="name"
                      value={formData.name}
                      onChange={handleChange}
                      placeholder="Enter your name"
                      className={`mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500 sm:text-sm transition-transform transform ${
                        selectedCard !== null ? "translate-x-0 opacity-100" : "translate-x-10 opacity-0"
                      }`}
                    />
                    {errors.name && <span className="text-red-500 text-sm">{errors.name}</span>}
                  </label>
                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700">Lecturer / Staff / Student ID</span>
                    <input
                      type="text"
                      name="id"
                      value={formData.id}
                      onChange={handleChange}
                      placeholder="Enter your ID"
                      className={`mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500 sm:text-sm transition-transform transform ${
                        selectedCard !== null ? "translate-x-0 opacity-100" : "translate-x-10 opacity-0"
                      }`}
                    />
                    {errors.id && <span className="text-red-500 text-sm">{errors.id}</span>}
                  </label>
                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700">Phone Number</span>
                    <input
                      type="number"
                      name="phone"
                      value={formData.phone}
                      onChange={handleChange}
                      placeholder="Enter your phone number"
                      className={`mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500 sm:text-sm transition-transform transform ${
                        selectedCard !== null ? "translate-x-0 opacity-100" : "translate-x-10 opacity-0"
                      }`}
                    />
                    {errors.phone && <span className="text-red-500 text-sm">{errors.phone}</span>}
                  </label>
                  <button
                    type="submit"
                    className="bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600 transition-transform transform"
                  >
                    Booking
                  </button>
                </form>
              </>
            )}
          </div>
        </div>

        {isBookingSuccessful && (
          <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
            <div className="bg-white p-6 rounded-lg shadow-lg text-center">
              <h2 className="text-xl font-bold mb-2">Booking Successful!</h2>
              <p className="text-green-500">You have successfully booked the slot.</p>
              <button
                className="mt-4 bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600"
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

export default Football;
