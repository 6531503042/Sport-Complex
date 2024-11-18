"use client";

import React, { useEffect, useState } from "react";
import "./badminton.css";
import NavBar from "@/app/components/navbar/navbar";
import CheckIcon from "@mui/icons-material/Check";
import ClearIcon from "@mui/icons-material/Clear";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ReportProblemIcon from "@mui/icons-material/ReportProblem";
import { useRouter } from "next/navigation";
import SportsTennisIcon from '@mui/icons-material/SportsTennis';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import EventAvailableIcon from '@mui/icons-material/EventAvailable';

interface UserData {
  id: string;
  name: string;
}

interface UserDataParams {
  params: UserData;
}

interface Slot {
  _id: string;
  start_time: string;
  end_time: string;
  current_bookings: number;
  max_bookings: number;
  status: string;
}

interface Court {
  _id: string;
  court_number: number;
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

// Add new interface for BookingFormProps
interface BookingFormProps {
  formData: {
    name: string;
    id: string;
    phone: string;
  };
  handleSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  onClose: () => void;
  selectedTime?: string;
  selectedCourt?: number;
}

const BookingForm = ({ formData, handleSubmit, onClose, selectedTime, selectedCourt }: BookingFormProps) => (
  <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div className="bg-white rounded-xl p-8 max-w-md w-full transform transition-all duration-300 ease-in-out">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Complete Booking</h2>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-600 transition-colors"
        >
          ✕
        </button>
      </div>

      {selectedTime && selectedCourt && (
        <div className="mb-6 p-4 bg-blue-50 rounded-lg">
          <p className="text-blue-800 font-medium">
            Court {selectedCourt} • {selectedTime}
          </p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Name
          </label>
          <input
            type="text"
            value={formData.name}
            readOnly
            className="w-full px-4 py-3 rounded-lg border border-gray-300 bg-gray-50"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            ID Number
          </label>
          <input
            type="text"
            value={formData.id}
            readOnly
            className="w-full px-4 py-3 rounded-lg border border-gray-300 bg-gray-50"
          />
        </div>
        <div className="flex gap-4 pt-4">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 px-4 py-3 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            Confirm Booking
          </button>
        </div>
      </form>
    </div>
  </div>
);

function Badminton_Booking({ params }: UserDataParams) {
  const { id } = params;
  const [storedRefreshToken, setStoredRefreshToken] = useState<string | null>(
    null
  );
  const router = useRouter();
  const [selectedCard, setSelectedCard] = useState<number | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    id: "",
    phone: "",
  });
  const [isBookingSuccessful, setIsBookingSuccessful] = useState(false);
  const [isBookingFailed, setIsBookingFailed] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false); // New state for mobile view
  const [slot, setSlot] = useState<Slot[]>([]);
  const [court, setCourt] = useState<Court[]>([]);
  const [showBookingForm, setShowBookingForm] = useState(false);

  const handleCardClick = (slot: Slot, slotIndex: number) => {
    if (!isSlotBooked(slot)) {
      setSelectedCard(slotIndex);
      setShowBookingForm(true);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      if (selectedCard === null || !slot || !slot[selectedCard]) {
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
        slot_id: null,
        status: 1,
        slot_type: "badminton",
        badminton_slot_id: slot[selectedCard]._id,
      };

      const response = await fetch(
        "http://localhost:1326/booking_v1/badminton/booking",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(bookingData),
        }
      );

      if (!response.ok) {
        setIsBookingFailed(true);
        return;
      }

      // Parse the complete response
      const result = await response.json();
      console.log("Complete booking response:", result);

      // Store all payment-related information in localStorage
      const paymentInfo = {
        payment_id: result.payment_id,
        booking_id: result.booking_id,
        qr_code_url: result.qr_code_url,
        status: result.status,
      };
      localStorage.setItem("currentPaymentInfo", JSON.stringify(paymentInfo));

      setIsBookingSuccessful(true);
      setSelectedCard(null);

      // Redirect to payment page after a short delay
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

  const getSlot = async () => {
    try {
      const resSlot = await fetch(
        "http://localhost:1335/facility_v1/badminton_v1/slots",
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
      console.log("Fetched slots:", slotData); // Debug log
    } catch (error) {
      console.error("Error fetching slots:", error);
      setSlot([]); 
    }
  };

  // Helper function to check if a slot is fully booked
  const isSlotBooked = (slot: Slot): boolean => {
    return slot.current_bookings >= slot.max_bookings;
  };

  // Helper function to get booking status text and style
  const getBookingStatus = (slot: Slot) => {
    const isBooked = isSlotBooked(slot);
    const remainingSpots = slot.max_bookings - slot.current_bookings;

    return {
      text: isBooked 
        ? "Fully Booked" 
        : `${remainingSpots} ${remainingSpots === 1 ? 'spot' : 'spots'} left`,
      className: isBooked
        ? 'bg-red-100 text-red-800'
        : remainingSpots <= 2
        ? 'bg-yellow-100 text-yellow-800'
        : 'bg-green-100 text-green-800'
    };
  };

  // Add close handler
  const handleCloseForm = () => {
    setShowBookingForm(false);
    setSelectedCard(null);
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
    const storedRefreshToken = localStorage.getItem("refresh_token");
    setStoredRefreshToken(storedRefreshToken);
    // Fetch slot data on initial render and set up the interval for updating
    getSlot();
    const intervalId = setInterval(getSlot, 10000);

    return () => {
      clearInterval(intervalId);
    };
  }, [id]);

  return (
    <>
      <NavBar activePage="badminton" />
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-white py-8 px-4">
        <div className="max-w-7xl mx-auto">
          {/* Enhanced Header Section */}
          <div className="text-center mb-12">
            <div className="inline-block p-3 rounded-full bg-blue-100 mb-4">
              <SportsTennisIcon className="text-blue-600 text-4xl" />
            </div>
            <h1 className="text-4xl font-bold text-gray-900 mb-3 bg-gradient-to-r from-blue-600 to-indigo-600 inline-block text-transparent bg-clip-text">
              Badminton Court Booking
            </h1>
            <p className="text-gray-600 text-lg max-w-2xl mx-auto">
              Reserve your preferred court and time slot for an amazing badminton experience
            </p>
          </div>

          {/* Time Slots Legend */}
          <div className="flex justify-center gap-6 mb-8">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-green-400 rounded-full animate-pulse" />
              <span className="text-sm text-gray-600">Available</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-yellow-400 rounded-full" />
              <span className="text-sm text-gray-600">Limited Spots</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-red-400 rounded-full" />
              <span className="text-sm text-gray-600">Fully Booked</span>
            </div>
          </div>

          <div className="bg-white rounded-2xl shadow-xl p-6 md:p-8 backdrop-blur-lg bg-opacity-90">
            {/* Show unavailable message if no slots */}
            {slot && slot.length === 0 ? (
              <div className="slot-unavailable-card">
                <ReportProblemIcon className="slot-unavailable-icon" />
                <h2 className="text-2xl font-bold text-gray-800 mb-2">
                  No Available Slots
                </h2>
                <p className="text-gray-600">
                  All courts are currently booked. Please check back later.
                </p>
              </div>
            ) : isMobileView ? (
              // Mobile booking form
              <div className="block sm:hidden">
                <div className="bg-white rounded-xl shadow-md p-6">
                  <div className="flex items-center justify-between mb-6">
                    <button
                      onClick={handleBackToTimeSlots}
                      className="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
                    >
                      <ArrowBackIosNewIcon className="w-5 h-5 mr-2" />
                      <span>Back to Courts</span>
                    </button>
                    {selectedCard !== null && slot[selectedCard] && (
                      <div className="text-lg font-semibold text-blue-600">
                        Court {court[selectedCard].court_number}
                      </div>
                    )}
                  </div>
                  <BookingForm 
                    formData={formData}
                    handleSubmit={handleSubmit}
                  />
                </div>
              </div>
            ) : (
              // Desktop view with courts grid
              <div className="space-y-8">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {slot &&
                    slot.length > 0 &&
                    slot
                      .reduce((groups: Slot[][], current: Slot, index: number) => {
                        const groupIndex = Math.floor(index / 4);
                        if (!groups[groupIndex]) groups[groupIndex] = [];
                        groups[groupIndex].push(current);
                        return groups;
                      }, [])
                      .map((group: Slot[], groupIndex: number) => (
                        <div
                          key={`group-${groupIndex}`}
                          className="bg-white rounded-xl shadow-lg overflow-hidden transition-all duration-300 hover:shadow-xl border border-gray-100"
                        >
                          {/* Enhanced Time slot header */}
                          <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white p-4">
                            <div className="flex items-center justify-center gap-2 mb-1">
                              <AccessTimeIcon className="text-blue-200" />
                              <h3 className="text-xl font-semibold">
                                {group[0]?.start_time} - {group[0]?.end_time}
                              </h3>
                            </div>
                            <p className="text-blue-100 text-sm text-center">
                              {4 - group.filter(s => isSlotBooked(s)).length} courts available
                            </p>
                          </div>

                          {/* Enhanced Courts grid */}
                          <div className="grid grid-cols-2 gap-4 p-4">
                            {group.map((slot: Slot, slotIndex: number) => {
                              const courtNumber = (slotIndex % 4) + 1;
                              const bookingStatus = getBookingStatus(slot);
                              const isBooked = isSlotBooked(slot);

                              return (
                                <button
                                  key={slot._id}
                                  onClick={() => !isBooked && handleCardClick(slot, slotIndex)}
                                  className={`
                                    relative p-6 rounded-xl text-center transition-all duration-300
                                    ${isBooked 
                                      ? 'bg-gray-50 cursor-not-allowed' 
                                      : 'court-card-available hover:transform hover:scale-105'
                                    }
                                    ${selectedCard === slotIndex ? 'ring-2 ring-blue-500 shadow-lg' : 'hover:shadow-md'}
                                  `}
                                  disabled={isBooked}
                                >
                                  <div className="space-y-3">
                                    <div className="flex items-center justify-center gap-2">
                                      <EventAvailableIcon className={`
                                        ${isBooked ? 'text-gray-400' : 'text-blue-500'}
                                      `} />
                                      <span className="font-bold text-lg text-gray-900">
                                        Court {courtNumber}
                                      </span>
                                    </div>

                                    <div className={`
                                      inline-flex px-3 py-1 rounded-full text-sm font-medium
                                      ${bookingStatus.className}
                                      transition-all duration-300
                                    `}>
                                      {bookingStatus.text}
                                    </div>

                                    {!isBooked && (
                                      <div className="absolute top-2 right-2">
                                        <div className="status-dot-available" />
                                      </div>
                                    )}
                                  </div>

                                  {isBooked && (
                                    <div className="booked-overlay">
                                      <span className="booked-stamp">
                                        BOOKED
                                      </span>
                                    </div>
                                  )}
                                </button>
                              );
                            })}
                          </div>
                        </div>
                      ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Enhanced Modals */}
        {showBookingForm && selectedCard !== null && slot[selectedCard] && (
          <BookingForm 
            formData={formData}
            handleSubmit={handleSubmit}
            onClose={handleCloseForm}
            selectedTime={`${slot[selectedCard].start_time} - ${slot[selectedCard].end_time}`}
            selectedCourt={(selectedCard % 4) + 1}
          />
        )}

        {/* Success/Error Modals */}
        {isBookingSuccessful && (
          <SuccessModal onClose={() => {
            setIsBookingSuccessful(false);
            setShowBookingForm(false);
          }} />
        )}
        {isBookingFailed && (
          <ErrorModal onClose={() => {
            setIsBookingFailed(false);
            setShowBookingForm(false);
          }} />
        )}
      </div>
    </>
  );
}

const SuccessModal = ({ onClose }) => (
  <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div className="bg-white rounded-xl p-6 max-w-sm w-full">
      <div className="text-center">
        <CheckIcon className="text-green-500 text-5xl mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Booking Successful!
        </h2>
        <p className="text-gray-600 mb-6">
          Your court has been successfully booked.
        </p>
        <button
          onClick={onClose}
          className="w-full bg-green-500 text-white py-3 rounded-lg font-semibold hover:bg-green-600 transition-colors"
        >
          Done
        </button>
      </div>
    </div>
  </div>
);

const ErrorModal = ({ onClose }) => (
  <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div className="bg-white rounded-xl p-6 max-w-sm w-full">
      <div className="text-center">
        <ClearIcon className="text-red-500 text-5xl mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Booking Failed
        </h2>
        <p className="text-gray-600 mb-6">
          Something went wrong. Please try again later.
        </p>
        <button
          onClick={onClose}
          className="w-full bg-red-500 text-white py-3 rounded-lg font-semibold hover:bg-red-600 transition-colors"
        >
          Close
        </button>
      </div>
    </div>
  </div>
);

export default Badminton_Booking;
