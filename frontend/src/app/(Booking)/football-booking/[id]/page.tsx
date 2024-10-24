"use client";

import React, { useEffect, useState, Fragment } from "react";
import "../football.css";
import Available from "@/app/assets/available.png";
import Unavailable from "@/app/assets/unavailable.png";
import Back from "@/app/assets/back.png";

interface UserData {
  id: string;
  name: string;
}

interface UserDataParams {
  params: UserData;
}

function Football_Booking({ params }: UserDataParams) {
  const { id } = params;
  

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

  const handleCardClick = (lot: any, status: any) => {
    if (!lot.current_bookings) {
      const index = slot.findIndex((s) => s._id === lot._id); // Get the index of the clicked lot
      if (index !== -1) {
        setSelectedCard(index === selectedCard ? null : index);
        setErrors({ name: "", id: "", phone: "" });

        if (window.innerWidth < 640) {
          setIsMobileView(true); // Switch to form view in mobile
        }
      }
    }
  };
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault(); // Prevent default form submission
    
    // Form validation: Check if the phone number is filled in
    if (!formData.phone) {
      setErrors((prevErrors) => ({
        ...prevErrors,
        phone: "Phone number is required",
      }));
      return;
    }
  
    try {
      const bookingData = {
        user_id: ownId, // User's ID from the state
        slot_id: slot[selectedCard]?._id, // Slot ID from the selected card
        status: 1, // Assuming 1 means successful booking
        slot_type: "normal", // Based on the Postman request, slot_type is 'normal'
        badminton_slot_id: null, // For football, badminton_slot_id can be null
      };
  
      const response = await fetch("http://localhost:1326/booking_v1/football/booking", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(bookingData),
      });
  
      if (!response.ok) {
        throw new Error(`Failed to book: ${response.statusText} (Status: ${response.status})`);
      }
  
      const result = await response.json();
      console.log("Booking successful:", result);
  
      // Show the booking success popup
      setIsBookingSuccessful(true);
  
      // Reset form and selected card
      setFormData({ name: "", id: "", phone: "" });
      setSelectedCard(null);
  
      // Update the list to reflect the successful booking
      setSlot((prevSlots) =>
        prevSlots.map((s) =>
          s._id === bookingData.slot_id ? { ...s, current_bookings: true } : s
        )
      );
    } catch (error) {
      console.error("Error submitting booking:", error);
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
      const resSlot = await fetch(
        "http://localhost:1335/facility_v1/football/slot_v1/slots",
        {
          method: "GET",
          cache: "no-store",
        }
      );

      if (!resSlot.ok) {
        throw new Error(
          `Failed to fetch: ${resSlot.statusText} (Status: ${resSlot.status})`
        );
      }
      const slotData = await resSlot.json();
      setSlot(slotData);
    } catch (error) {
      console.error("Error fetching slot data:", error);
    }
  };

  useEffect(() => {
    getSlot();
    getUserData(id); // Use the id from params when calling getUserData

    const intervalId = setInterval(() => {
      getSlot();
    }, 10000);

    return () => {
      clearInterval(intervalId);
    };
  }, [id]);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [userData, setUserData] = useState<any[]>([]);
  const [ownName, setOwnName] = useState("");
  const [ownId, setOwnId] = useState("");
  const getUserData = async (id: string) => {
    try {
      const resUserData = await fetch(
        `http://localhost:1325/user_v1/users/${id}`,
        {
          method: "GET",
          cache: "no-store",
        }
      );
      if (!resUserData.ok) {
        throw new Error("Failed to fetch the user");
      }
      const userData: UserData = await resUserData.json();
      setUserData(userData);
      setOwnName(userData.name);
      setOwnId(userData.id);
    } catch (error) {
      console.log(error);
    }
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
            <div className="block sm:hidden ">
  <div className="bg-white border border-gray-200 rounded-lg p-4 shadow-md">
    <form onSubmit={handleSubmit}>
      {/* Flex container for aligning back button and time at the top */}
      <div className="flex items-center justify-between my-3">
        <img
          src={Back.src}
          alt="back"
          onClick={handleBackToTimeSlots}
          className="border shadow-xl p-2 rounded-md cursor-pointer hover:bg-gray-200"
          width={40}
        />
        <h2 className="text-xl font-semibold text-start">
        {slot[selectedCard].start_time} -{" "}
        {slot[selectedCard].end_time}
        </h2>
      </div>

      <label className="block mb-4">
        <span className="block text-sm font-medium text-gray-700 py-2">
          Name 
        </span>
        <input
          type="text"
          name="name"
          value={ownName}
          className="name-input-football mt-1 block w-full px-3 py-3"
        />
      </label>

      <label className="block mb-4">
        <span className="block text-sm font-medium text-gray-700 py-2">
          Lecturer / Staff / Student ID
        </span>
        <input
          type="text"
          name="id"
          value={ownId}
          className="name-input-football mt-1 block w-full px-3 py-3"
        />
      </label>

      <label className="block mb-4">
        <span className="block text-sm font-medium text-gray-700 py-2">
          Phone Number
        </span>
        <input
          type="tel"
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

      {/* Center the Booking button */}
      <div className="flex justify-center">
        <button
          type="submit"
          className="font-bold bg-green-500 text-white px-5 py-2.5 rounded-md drop-shadow-2xl hover:bg-green-600"
        >
          Booking
        </button>
      </div>
    </form>
  </div>
</div>

          
          ) : (
            // Normal and larger screen view
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
  {slot.map((lot) => (
    <div
    key={lot._id}
    className={` border border-gray-200 rounded-lg p-6 shadow-md transition-transform duration-300 ease-in-out
      ${lot.current_bookings ? "cursor-not-allowed bg-light-gray text-white" : "cursor-pointer bg-[#5EB900] text-white border-green-300 hover:scale-105 hover:shadow-lg"}
      ${!lot.current_bookings ? "hover:bg-[#005400]" : ""}
    `}
    onClick={() => {
      if (!lot.current_bookings) {
        handleCardClick(lot, lot.status);
      }
    }}
  >
    <div className="text-lg font-semibold flex justify-between items-center">
      <div>
        {lot.start_time} - {lot.end_time}
      </div>
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

              <div
                className={`hidden sm:block transition-all duration-300 ease-in-out mt-6 p-4 bg-white border border-gray-200 rounded-lg shadow-md transform ${
                  selectedCard !== null && !slot[selectedCard]?.current_bookings // Check if current_bookings is false
                    ? "translate-y-0 opacity-100"
                    : "translate-y-5 opacity-0"
                }`}
              >
                {selectedCard !== null &&
                  !slot[selectedCard]?.current_bookings && ( // Check if current_bookings is false
                    <>
                      <h2 className="text-xl font-bold mb-4">
                        Booking for {slot[selectedCard].start_time} -{" "}
                        {slot[selectedCard].end_time}
                      </h2>
                      <form onSubmit={handleSubmit}>
                        <label className="block mb-4">
                          <span className="block text-sm font-medium text-gray-700 py-2">
                            Name
                          </span>
                          <input
                            type="text"
                            name="name"
                            value={ownName}
                            // placeholder={ownName}
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
                            value={ownId}
                            // placeholder={ownId}
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
                            type="tel"
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
