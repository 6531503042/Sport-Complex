"use client";

import React, { useState } from "react";
import '@/app/globals.css';

interface Data {
  time: string;
  isAvailable: boolean;
  member: string;
}

const Gyms_Booking: React.FC = () => {

  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedTime, setSelectedTime] = useState("");
  const [showNotification, setShowNotification] = useState(false);
  const [isShowSlot, setIsShowSlot] = useState(true);
  const [isPopupVisible, setIsPopupVisible] = useState(false);

  
  // State for form fields
  const [formData, setFormData] = useState({
    fullName: "",
    phoneNumber: "",
    studentId: "",
  });

  const data: Data[] = [
    { time: "8.00 - 9.15", isAvailable: true, member: "0/30" },
    { time: "9.30 - 10.45", isAvailable: false, member: "30/30" },
    { time: "11.00 - 12.15", isAvailable: true, member: "0/30" },
    { time: "12.30 - 13.45", isAvailable: false, member: "30/30" },
    { time: "14.00 - 15.15", isAvailable: true, member: "0/30" },
    { time: "15.30 - 16.45", isAvailable: true, member: "0/30" },
    { time: "17.00 - 18.15", isAvailable: true, member: "0/30" },
    { time: "18.30 - 19.45", isAvailable: true, member: "0/30" },
    { time: "20.00 - 21.15", isAvailable: true, member: "0/30" },
  ];

  const handleSlotClick = (time: string) => {
    setSelectedTime(time);
    setIsPopupOpen(true);
    setTimeout(() => {
      setIsPopupVisible(true);  // Trigger the transition after setting popup state
    }, 10); // Slight delay to ensure the visibility transition kicks in
    setIsShowSlot(false);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setShowNotification(true);
    setIsPopupVisible(false);
    setTimeout(() => {
      setIsPopupOpen(false);
      setIsShowSlot(true);
    }, 300);
    
    // Clear form data
    setFormData({
      fullName: "",
      phoneNumber: "",
      studentId: "",
    });

    setTimeout(() => {
      setShowNotification(false);
    }, 3000);
  };

  const closePopup = () => {
    setIsPopupVisible(false); // Start fade-out
    setTimeout(() => {
      setIsPopupOpen(false);  // Close popup after fade-out completes
      setIsShowSlot(true);
    }, 300);
  };






  return (
    <div className="bg-white w-full h-full flex justify-center items-center">
      <div className="bg-white w-4/6 h-auto rounded-lg shadow-2xl border border-gray-200 flex justify-center m-16">
        <div className="flex-none m-10 w-full h-full">
          <h1 className="text-4xl font-bold text-center mb-6">Sport Complex Gym Booking</h1>
          <p className="text-md text-gray-400 font-bold text-center mb-6">Book your gym slot easily and quickly!</p>
          <div className="flex flex-col m-10">
            {isShowSlot && 
            <div className="grid grid-cols-1 gap-6 mb-10 md:grid-cols-2 xl:grid-cols-3">
                {data.map((slot, index) => (
                  <div
                    key={index}
                    className={`w-auto h-auto p-2 rounded-lg flex flex-col text-center border border-gray-100 
                      ${slot.isAvailable ? 'bg-green-100 shadow-lg hover:bg-blue-100 hover:scale-105 transition-transform duration-300 cursor-pointer' : 'bg-red-100 shadow-lg cursor-not-allowed'}
                    `}
                    onClick={() => slot.isAvailable && handleSlotClick(slot.time)}
                  >
                    <div className="flex flex-row items-center justify-between w-full h-auto pl-5 pr-5">
                      <div className="text-lg font-semibold text-start m-2">{slot.time}</div>
                      <p className={`text-gray-600 ${slot.isAvailable ? 'text-green-500' : 'text-red-600'}`}>
                            {slot.isAvailable ? 
                            <div className="flex flex-row items-center justify-start ">
                              <svg className="w-8 h-8" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
                                <path fill-rule="evenodd" d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm13.707-1.293a1 1 0 0 0-1.414-1.414L11 12.586l-1.793-1.793a1 1 0 0 0-1.414 1.414l2.5 2.5a1 1 0 0 0 1.414 0l4-4Z" clip-rule="evenodd"/>
                              </svg>
                            </div>
                            : <div className="flex flex-row items-center justify-start">
                                <svg className="w-8 h-8" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
                                  <path fill-rule="evenodd" d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm5.757-1a1 1 0 1 0 0 2h8.486a1 1 0 1 0 0-2H7.757Z" clip-rule="evenodd"/>
                                </svg>
                              </div>
                            }         
                        </p>
                    </div>
                      
                    <div className="flex flex-row items-center justify-start pl-5 pr-5">
                        <svg className="w-5 h-5 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
                          <path fill-rule="evenodd" d="M12 4a4 4 0 1 0 0 8 4 4 0 0 0 0-8Zm-2 9a4 4 0 0 0-4 4v1a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2v-1a4 4 0 0 0-4-4h-4Z" clip-rule="evenodd"/>
                        </svg>
                        <div className="text-md font-semibold m-1">{slot.member}</div>
                      </div>
                    </div>
                ))}
              </div>
            }
            {isPopupOpen && (
                  <div className={`bg-white rounded-lg shadow-xl border border-gray-200 p-8 m-5 transition-all duration-500 ease-in-out transform 
                    ${isPopupVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-75'}`}>
                  <h1 className="text-2xl font-mono font-bold mb-6">Booking on this time {selectedTime}</h1>
                  <form onSubmit={handleSubmit}>
                    <div className="flex flex-row items-center">
                      <div className="font-bold text-md m-2">Name</div>
                    </div>
                    <input
                      type="text"
                      name="fullName"
                      value={formData.fullName}
                      onChange={handleChange}
                      placeholder="Enter Your Full Name !"
                      className="border border-gray-200 rounded-3xl shadow-lg p-4 mb-3 w-full"
                      required
                    />
                    <div className="flex flex-row items-center">
                      <div className="font-bold text-md m-2">Phone Number</div>
                    </div>
                    <input
                      type="tel"
                      name="phoneNumber"
                      value={formData.phoneNumber}
                      onChange={handleChange}
                      placeholder="Enter Your Phone Number !"
                      className="border border-gray-200 rounded-3xl shadow-lg p-4 mb-3 w-full"
                      required
                    />
                    <div className="flex flex-row items-center">
                      <div className="font-bold text-md m-2">Student ID</div>
                    </div>
                    <input
                      type="text"
                      name="studentId"
                      value={formData.studentId}
                      onChange={handleChange}
                      placeholder="Enter Your Student ID !"
                      className="border border-gray-200 rounded-3xl shadow-lg p-4 mb-10 w-full"
                      required
                    />
                    <div className="flex flex-row gap-8">
                      <button type="submit" className="bg-gray-600 text-white rounded-xl shadow-lg p-2 py-4 w-full font-bold">Confirm Booking</button>
                      <button onClick={closePopup} className="p-2 py-4 w-full bg-gray-600 rounded-xl shadow-lg font-bold text-white">Cancel</button>
                    </div>
                    </form>
                </div>
              )}
          </div>
        </div>
      </div>
      {showNotification && (
        <div className="fixed top-4 right-4 bg-green-500 text-white p-4 rounded-lg shadow-lg">
          <p>Success! Your booking for {selectedTime} has been confirmed.</p>
        </div>
      )}
    </div>
  );
};

export default Gyms_Booking;
