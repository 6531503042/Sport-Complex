"use client";

import React, { useEffect, useState, Fragment } from "react";
import CheckIcon from "@mui/icons-material/Check";
import ClearIcon from "@mui/icons-material/Clear";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import GroupIcon from "@mui/icons-material/Group";
import NavBar from "@/app/components/navbar/navbar";
import ReportProblemIcon from "@mui/icons-material/ReportProblem";
import "./swimming.css";
import { useRouter } from 'next/navigation';

interface UserData {
  id: string;
  name: string;
}

interface UserDataParams {
  params: UserData;
}
async function getAccessToken(refreshToken: string | null) {
  if (!refreshToken) return null;

  try {
    const response = await fetch("http://localhost:1326/auth/refresh", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) throw new Error("Failed to refresh token");

    const data = await response.json();
    const newAccessToken = data.access_token;

    // Store the new access token in localStorage
    localStorage.setItem("access_token", newAccessToken);
    return newAccessToken;
  } catch (error) {
    console.error("Error refreshing access token:", error);
    return null;
  }
}

function Swimming_Booking({ params }: UserDataParams) {
  const { id } = params;
  const router = useRouter();
  const [storedRefreshToken, setStoredRefreshToken] = useState<string | null>(
    null
  );

  const [selectedCard, setSelectedCard] = useState<number | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    id: "",
    phone: "",
  });

  const [isBookingSuccessful, setIsBookingSuccessful] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false); // New state for mobile view
  const [isBookingFailed, setIsBookingFailed] = useState(false);

  const handleCardClick = (lot: any) => {
    // Allow booking if the slot is not fully booked
    const isSlotFull = lot.current_bookings >= lot.max_bookings;

    // Check if the slot is full
    if (!isSlotFull) {
      const index = slot.findIndex((s) => s._id === lot._id); // Get the index of the clicked lot
      if (index !== -1) {
        setSelectedCard(index === selectedCard ? null : index); // Toggle selected card

        if (window.innerWidth < 640) {
          setIsMobileView(true); // Switch to form view in mobile
        }
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      if (selectedCard === null || !slot[selectedCard]) {
        console.error("No slot selected");
        return;
      }

      let accessToken = localStorage.getItem("access_token");
      if (!accessToken) {
        accessToken = await getAccessToken(storedRefreshToken);
        if (!accessToken) {
          console.error("Failed to obtain access token");
          setIsBookingFailed(true);
          return;
        }
      }

      const bookingData = {
        user_id: formData.id,
        slot_id: slot[selectedCard]._id,
        status: 1,
        slot_type: "normal",
        badminton_slot_id: null,
      };

      console.log("Sending booking data:", bookingData); // Debug log

      const response = await fetch("http://localhost:1326/booking_v1/swimming/booking", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(bookingData),
      });

      if (!response.ok) {
        console.error("Booking failed with status:", response.status);
        setIsBookingFailed(true);
        return;
      }

      const result = await response.json();
      console.log("Complete booking response:", result); // Debug log

      if (!result.payment_id) {
        console.error("No payment_id in response");
        setIsBookingFailed(true);
        return;
      }

      // Store payment information in localStorage
      const paymentInfo = {
        payment_id: result.payment_id,
        booking_id: result.booking_id,
        qr_code_url: result.qr_code_url,
        status: result.status
      };
      localStorage.setItem('currentPaymentInfo', JSON.stringify(paymentInfo));
      console.log("Stored payment info:", paymentInfo); // Debug log

      setIsBookingSuccessful(true);
      setSelectedCard(null);

      // Immediate redirect to payment page
      setTimeout(() => {
        if (result.payment_id) {
            router.push(`/payment/${result.payment_id}`);
        } else {
            console.error("Payment ID is undefined, cannot redirect.");
            setIsBookingFailed(true);
        }
    }, 1000);

    } catch (error) {
      console.error("Error submitting booking:", error);
      setIsBookingFailed(true);
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
        "http://localhost:1335/facility_v1/swimming/slot_v1/slots",
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
      setSlot([]);
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

    // Retrieve refresh token and set it in state
    const storedRefreshToken = localStorage.getItem("access_token");
    setStoredRefreshToken(storedRefreshToken);
    if (storedRefreshToken)
      console.log("Stored Refresh Token:", storedRefreshToken); // Use it as needed

    // Fetch slot data on initial render and set up the interval for updating
    getSlot();
    const intervalId = setInterval(getSlot, 1000);

    return () => {
      clearInterval(intervalId);
    };
  }, [id]);

  return (
    <>
      <NavBar activePage="swimming" />
      <div className="flex flex-col items-center h-screen p-6">
        <div className="w-full max-w-[1189px] bg-[#FEFFFE] border-gray border rounded-3xl drop-shadow-2xl p-5">
          <h1 className="text-4xl font-bold my-10 text-black text-center">
            Swimming Booking
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
            // Normal screen view
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {slot?.map((lot) => {
                  const isSlotFull = lot.current_bookings >= lot.max_bookings; // Check if the slot is fully booked

                  return (
                    <div
                      key={lot._id}
                      className={`border border-gray-200 rounded-lg p-6 shadow-md transition-transform duration-300 ease-in-out
      ${
        isSlotFull
          ? "cursor-not-allowed bg-[#C1C7D4] text-white"
          : "cursor-pointer bg-[#4169E1] text-white border-blue-300 hover:scale-105 hover:shadow-lg"
      }
      ${!isSlotFull ? "hover:bg-[#000080]" : ""}
    `}
                      onClick={() => {
                        if (!isSlotFull) {
                          handleCardClick(lot);
                        }
                      }}
                    >
                      <div className="text-lg font-semibold grid grid-cols-2 justify-between items-center">
                        <div>
                          {lot.start_time} - {lot.end_time}
                          <div className="mt-2 text-sm text-white">
                            <GroupIcon
                              className="mr-2.5"
                              style={{ fontSize: "1.3rem" }}
                            />
                            {lot.current_bookings} / {lot.max_bookings}
                          </div>
                        </div>
                        <div className="ml-auto">
                          {isSlotFull ? (
                            <ClearIcon
                              className="mx-2.5"
                              style={{ fontSize: "2rem" }}
                            />
                          ) : (
                            <CheckIcon
                              className="mx-2.5"
                              style={{ fontSize: "2rem" }}
                            />
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>

              <div
                className={`hidden sm:block transition-all duration-300 ease-in-out mt-6 p-4 bg-white border border-gray-200 rounded-lg shadow-md transform ${
                  selectedCard !== null
                    ? "translate-y-0 opacity-100"
                    : "translate-y-5 opacity-0"
                }`}
              >
                {selectedCard !== null && ( // Always show the form if a card is selected
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
                          className="name-input-swimming mt-1 block w-full px-3 py-3"
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
                          className="name-input-swimming mt-1 block w-full px-3 py-3"
                        />
                      </label>
                      <div className="flex justify-center">
                        <button
                          type="submit"
                          className="font-semibold bg-[#4169E1] text-white px-6 py-3 my-5 rounded-md drop-shadow-2xl hover:bg-[#000080]"
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
              {/* Success Popup Content */}
              <h2 className="text-2xl font-bold text-gray-800 mb-2 text-center">
                Booking Successful!
              </h2>
              <p className="text-gray-600 mb-6 text-center">
                You have successfully booked the slot.
              </p>

              {/* Close Button */}
              <button
                className="w-full bg-gradient-to-r from-[#000080] via-[#2A52BE] to-[#4169E1] text-white py-3 rounded-lg font-semibold shadow-lg hover:shadow-xl transition duration-200 ease-in-out transform hover:scale-105"
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
        {isBookingFailed && (
          <div className="fixed inset-0 w-screen h-screen flex items-center justify-center z-50 bg-black bg-opacity-40 backdrop-blur-sm transition-opacity duration-300 ease-in-out">
            <div className="relative bg-white w-full max-w-sm mx-auto p-8 rounded-lg shadow-xl transform transition-all duration-500 ease-in-out scale-100">
              {/* Failure Popup Content */}
              <h2 className="text-2xl font-bold text-gray-800 mb-2 text-center">
                Booking Failed
              </h2>
              <p className="text-gray-600 mb-6 text-center">
                An error occurred while processing your booking. Please try
                again later.
              </p>

              {/* Close Button */}
              <button
                className="w-full bg-gradient-to-r from-red-600 via-red-500 to-red-400 text-white py-3 rounded-lg font-semibold shadow-lg hover:shadow-xl transition duration-200 ease-in-out transform hover:scale-105"
                onClick={() => setIsBookingFailed(false)}
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
