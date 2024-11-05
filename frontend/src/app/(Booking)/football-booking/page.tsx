"use client";

import React, { useEffect, useState, Fragment } from "react";
import "./football.css";
import NavBar from "@/app/components/navbar/navbar";
import CheckIcon from "@mui/icons-material/Check";
import ClearIcon from "@mui/icons-material/Clear";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ReportProblemIcon from "@mui/icons-material/ReportProblem";

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
  const [isBookingSuccessful, setIsBookingSuccessful] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false); // New state for mobile view

  const handleCardClick = (lot: any) => {
    if (!lot.current_bookings) {
      const index = slot.findIndex((s) => s._id === lot._id); // Get the index of the clicked lot
      if (index !== -1) {
        setSelectedCard(index === selectedCard ? null : index);

        if (window.innerWidth < 640) {
          setIsMobileView(true); // Switch to form view in mobile
        }
      }
    }
  };
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault(); // Prevent default form submission

    // Form validation: Check if the phone number is filled in

    try {
      // Ensure the selectedCard is valid before proceeding
      if (selectedCard === null || !slot[selectedCard]) {
        console.error("No slot selected");
        return;
      }

      // Prepare booking data
      const bookingData = {
        user_id: formData.id, // Use ID from formData
        slot_id: slot[selectedCard]._id, // Get the selected slot's ID
        status: 1, // Assuming 1 means successful booking
        slot_type: "normal", // Slot type, based on your API
        badminton_slot_id: null, // Not applicable for football
      };

      // Log the booking data
      console.log("Booking Data:", bookingData);

      // Send booking request
      const response = await fetch(
        "http://localhost:1326/booking_v1/football/booking",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(bookingData),
        }
      );

      const result = await response.json();
      console.log("Booking successful:", result);

      // Show the booking success popup
      setIsBookingSuccessful(true);

      // Reset form and selected card
      setFormData((prevData) => ({
        ...prevData,
        phone: "", // Clear phone number field only
      }));
      setSelectedCard(null);

      // Update the slot to reflect the successful booking
      setSlot((prevSlots) =>
        prevSlots.map((s) =>
          s._id === bookingData.slot_id ? { ...s, current_bookings: true } : s
        )
      );
    } catch (error) {
      console.error("Error submitting booking:", error);
    }
  };

  const handleBackToTimeSlots = () => {
    setIsMobileView(false);
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [slot, setSlot] = useState<any[] | null>(null);
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
      setSlot(Array.isArray(slotData) && slotData.length ? slotData : []);
    } catch (error) {
      console.error("Error fetching slot data:", error);
      setSlot([]); // Set to empty array if there's an error fetching
    }
  };

  useEffect(() => {
    // Retrieve user data from localStorage
    const userDataName = localStorage.getItem("user");
    if (userDataName) {
      const user = JSON.parse(userDataName);

      setFormData((prevData) => ({
        ...prevData,
        name: user.name || "",
        id: user.id.replace(/^user:/, "") || "", // Remove "user:" prefix if it exists
      }));
    }
    const userDataId = localStorage.getItem("_id");
    if (userDataId) {
      const user = JSON.parse(userDataId);

      setFormData((prevData) => ({
        ...prevData,
        name: user.name || "",
        id: user.id.replace(/^user:/, "") || "", // Remove "user:" prefix if it exists
      }));
    }
    // Fetch slot data on initial render and set up the interval for updating
    getSlot();
    const intervalId = setInterval(getSlot, 10000);

    return () => {
      clearInterval(intervalId);
    };
  }, [id]);

  return (
    <>
      <NavBar activePage="football"/>
      <div className="flex flex-col items-center h-screen p-6">
        <div className="w-full max-w-[1189px] bg-[#FEFFFE] border-gray border rounded-3xl drop-shadow-2xl p-5">
          <h1 className="text-4xl font-bold my-10 text-black text-center">
            Football Booking
          </h1>

          {slot && slot.length === 0 ? (
            <div className="slot-unavailable-card text-center p-8 rounded-lg shadow-md transition-transform duration-200 ease-in-out transform hover:scale-105">
              <ReportProblemIcon
                className="slot-unavailable-icon text-red-500 mb-4"
                style={{ fontSize: "3rem" }}
              />
              <h2 className="text-3xl font-bold text-gray-800 mb-2">
                Slot Unavailable
              </h2>
              <p className="text-gray-600 text-lg">
                All slots are currently booked. Please check back later or
                contact support for more options.
              </p>
            </div>
          ) : isMobileView ? (
            // Mobile View: Show the form instead of the time slots
            <div className="block sm:hidden ">
              <div className="bg-white border border-gray-200 rounded-lg p-4 shadow-md">
                <form onSubmit={handleSubmit}>
                  {/* Flex container for aligning back button and time at the top */}
                  <div className="flex items-center justify-between my-3">
                    <ArrowBackIosNewIcon
                      className="border shadow-xl w-10 h-10 p-2 rounded-md cursor-pointer hover:bg-gray-200"
                      onClick={handleBackToTimeSlots}
                      style={{ fontSize: "2rem" }}
                    />
                    {selectedCard !== null && slot[selectedCard] && (
                      <h2 className="text-xl font-semibold text-start">
                        {slot[selectedCard].start_time} -{" "}
                        {slot[selectedCard].end_time}
                      </h2>
                    )}
                  </div>

                  <label className="block mb-4">
                    <span className="block text-sm font-medium text-gray-700 py-2">
                      Name
                    </span>
                    <input
                      type="text"
                      name="name"
                      value={formData.name}
                      readOnly
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
                      value={formData.id}
                      readOnly
                      className="name-input-football mt-1 block w-full px-3 py-3"
                    />
                  </label>

                  {/* Center the Booking button */}
                  <div className="flex justify-center">
                    <button
                      type="submit"
                      className="font-bold bg-[#5EB900] text-white px-5 py-2.5 rounded-md drop-shadow-2xl hover:bg-[#005400]"
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
                {slot?.map((lot) => (
                  <div
                    key={lot._id}
                    className={`border border-gray-200 rounded-lg p-6 shadow-md transition-transform duration-300 ease-in-out
                    ${
                      lot.current_bookings
                        ? "cursor-not-allowed bg-[#C1C7D4] text-white"
                        : "cursor-pointer bg-[#5EB900] text-white border-green-300 hover:scale-105 hover:shadow-lg"
                    }
                    ${!lot.current_bookings ? "hover:bg-[#005400]" : ""}`}
                    onClick={() => {
                      if (!lot.current_bookings) {
                        handleCardClick(lot);
                      }
                    }}
                  >
                    <div className="text-lg font-semibold flex justify-between items-center">
                      <div>
                        {lot.start_time} - {lot.end_time}
                      </div>
                      <div>
                        {lot.current_bookings ? (
                          <ClearIcon
                            className="mx-2.5"
                            style={{ fontSize: "1.3rem" }}
                          />
                        ) : (
                          <CheckIcon
                            className="mx-2.5"
                            style={{ fontSize: "1.3rem" }}
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
                            value={formData.name} // Display the name from localStorage
                            readOnly
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
                            value={formData.id} // Display the id from localStorage
                            readOnly
                            className="name-input-football mt-1 block w-full px-3 py-3"
                          />
                        </label>

                        {/* Center the button */}
                        <div className="flex justify-center">
                          <button
                            type="submit"
                            className="font-semibold bg-[#5EB900] text-white px-6 py-3 my-5 rounded-md drop-shadow-2xl hover:bg-[#005400]"
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
                onClick={() => {
                  setIsBookingSuccessful(false); // Close the popup
                  setIsMobileView(false); // Return to the time slots view
                }}
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
