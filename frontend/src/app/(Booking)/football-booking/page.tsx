"use client";

import React, { useEffect, useState, Fragment } from "react";
import "./football.css";
import NavBar from "@/app/components/navbar/navbar";
import CheckIcon from "@mui/icons-material/Check";
import ClearIcon from "@mui/icons-material/Clear";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ReportProblemIcon from "@mui/icons-material/ReportProblem";
import { useRouter } from 'next/navigation';
import SportsSoccerIcon from '@mui/icons-material/SportsSoccer';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import EventAvailableIcon from '@mui/icons-material/EventAvailable';
import WarningIcon from '@mui/icons-material/Warning';
import GroupIcon from '@mui/icons-material/Group';
import { Tooltip } from "@mui/material";

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

interface BookingFormProps {
  formData: {
    name: string;
    id: string;
    phone: string;
  };
  handleSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  onClose: () => void;
  selectedTime?: string;
}

const BookingForm = ({ formData, handleSubmit, onClose, selectedTime }: BookingFormProps) => (
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

      {selectedTime && (
        <div className="mb-6 p-4 bg-green-50 rounded-lg">
          <div className="flex items-center gap-2 text-green-800">
            <AccessTimeIcon className="text-green-600" />
            <span className="text-lg font-semibold">{selectedTime}</span>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Name
          </label>
          <div className="relative">
          <input
            type="text"
            value={formData.name}
            readOnly
            className="text-black w-full px-4 py-3 rounded-xl border border-gray-300 bg-gray-50
                          focus:ring-2 focus:ring-cyan-500 focus:border-transparent
                          transition-all duration-200"
          />
          <Tooltip title="ID cannot be changed" arrow>
                <div className="absolute right-3 top-1/2 -translate-y-1/2">
                  <WarningIcon className="text-gray-400" />
                </div>
              </Tooltip>
          </div>
        </div>
        
        <div >
          <label className="block text-sm font-medium text-gray-700 mb-2">
            ID Number
          </label>
          <div className="relative">
          <input
            type="text"
            value={formData.id}
            readOnly
            className="text-black w-full px-4 py-3 rounded-xl border border-gray-300 bg-gray-50
                          focus:ring-2 focus:ring-cyan-500 focus:border-transparent
                          transition-all duration-200"
          />
           <Tooltip title="Name cannot be changed" arrow>
                <div className="absolute right-3 top-1/2 -translate-y-1/2">
                  <WarningIcon className="text-gray-400" />
                </div>
              </Tooltip>
         </div>
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
            className="flex-1 px-4 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
          >
            Confirm Booking
          </button>
        </div>
      </form>
    </div>
  </div>
);

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

function Football_Booking({ params }: UserDataParams) {
  const { id } = params;
  const [storedRefreshToken, setStoredRefreshToken] = useState<string | null>(null);
  const router = useRouter();
  const [selectedCard, setSelectedCard] = useState<number | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    id: "",
    phone: "",
  });
  const [isBookingSuccessful, setIsBookingSuccessful] = useState(false);
  const [isBookingFailed, setIsBookingFailed] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false);
  const [showBookingForm, setShowBookingForm] = useState(false);
  const [slots, setSlots] = useState<Slot[]>([]);

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

  const handleCardClick = (slot: Slot) => {
    if (!isSlotBooked(slot)) {
      const index = slots.findIndex(s => s._id === slot._id);
      setSelectedCard(index);
      setShowBookingForm(true);
    }
  };

  const handleCloseForm = () => {
    setShowBookingForm(false);
    setSelectedCard(null);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      if (selectedCard === null || !slots[selectedCard]) {
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
        slot_id: slots[selectedCard]._id,
        status: 1,
        slot_type: "normal",
        badminton_slot_id: null,
      };

      const response = await fetch("http://localhost:1326/booking_v1/football/booking", {
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

      if (!result.payment_id) {
        console.error("No payment_id in response");
        setIsBookingFailed(true);
        return;
      }

      // Store payment information in localStorage
      const paymentInfo = {
        payment_id: result.payment_id,
        booking_id: result.booking_id,
        status: result.status
      };
      localStorage.setItem('currentPaymentInfo', JSON.stringify(paymentInfo));

      setIsBookingSuccessful(true);
      setSelectedCard(null);

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
      setSlots(Array.isArray(slotData) && slotData.length ? slotData : []);
    } catch (error) {
      console.error("Error fetching slot data:", error);
      setSlots([]); // Set to empty array if there's an error fetching
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
  const storedRefreshToken = localStorage.getItem('refresh_token');
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
      <NavBar activePage="football" />
      <div className="min-h-screen bg-gradient-to-br from-green-50 via-emerald-50 to-white py-8 px-4">
        <div className="max-w-7xl mx-auto">
          {/* Enhanced Header Section */}
          <div className="text-center mb-12">
            <div className="inline-block p-3 rounded-full bg-green-100 mb-4">
              <SportsSoccerIcon className="text-green-600 text-4xl" />
            </div>
            <h1 className="text-4xl font-bold text-gray-900 mb-3 bg-gradient-to-r from-green-600 to-emerald-600 inline-block text-transparent bg-clip-text">
              Football Booking
            </h1>
            <p className="text-gray-600 text-lg max-w-2xl mx-auto">
              Reserve your preferred time slot for an amazing football experience !
            </p>
          </div>

         

          <div className="bg-white rounded-2xl shadow-xl p-6 md:p-8 backdrop-blur-lg bg-opacity-90">
            {slots && slots.length === 0 ? (
              <div className="slot-unavailable-card">
                <ReportProblemIcon className="slot-unavailable-icon" />
                <h2 className="text-2xl font-bold text-gray-800 mb-2">
                  No Available Slots
                </h2>
                <p className="text-gray-600">
                  All slots are Unavailable now.
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {slots?.map((lot, index) => {
                  const isSlotFull = isSlotBooked(lot);
                  const remainingSpots = lot.max_bookings - lot.current_bookings;

                  return (
                    <div
                      key={lot._id}
                      className={`
                        relative rounded-xl overflow-hidden transition-all duration-300
                        ${isSlotFull 
                          ? 'bg-gray-50 cursor-not-allowed' 
                          : 'bg-white hover:bg-green-50 cursor-pointer transform hover:scale-105'
                        }
                        shadow-md hover:shadow-xl
                      `}
                      onClick={() => !isSlotFull && handleCardClick(lot)}
                    >
                      {/* Time Header */}
                      <div className="bg-gradient-to-r from-green-600 to-emerald-600 text-white p-4">
                        <div className="flex items-center justify-center gap-2">
                          <AccessTimeIcon className="text-green-200" />
                          <h3 className="text-lg font-semibold">
                            {lot.start_time} - {lot.end_time}
                          </h3>
                        </div>
                      </div>

                      {/* Slot Details */}
                      <div className="p-4">
                        <div className="flex items-center justify-between mb-3">
                          <div className="flex items-center gap-2">
                            <GroupIcon className="text-gray-400" />
                            <span className="text-gray-600">
                              {lot.current_bookings} / {lot.max_bookings}
                            </span>
                          </div>
                          <div className={`
                            px-3 py-1 rounded-full text-sm font-medium
                            ${isSlotFull 
                              ? 'bg-red-100 text-red-800' 
                              : remainingSpots <= 2
                              ? 'bg-yellow-100 text-yellow-800'
                              : 'bg-green-100 text-green-800'
                            }
                          `}>
                            {isSlotFull 
                              ? 'Fully Booked' 
                              : `${remainingSpots} spots left`
                            }
                          </div>
                        </div>

                        {!isSlotFull && (
                          <div className="absolute top-2 right-2">
                            <div className="w-3 h-3 bg-green-400 rounded-full animate-pulse" />
                          </div>
                        )}
                      </div>

                      {isSlotFull && (
                        <div className="absolute inset-0 bg-gray-100/50 backdrop-blur-[1px] flex items-center justify-center">
                          <span className="text-red-600 font-bold text-lg border-2 border-red-600 px-4 py-2 rounded bg-white/80 transform -rotate-12">
                          Reserved
                          </span>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>

        {/* Booking Form Modal */}
        {showBookingForm && selectedCard !== null && slots[selectedCard] && (
          <BookingForm 
            formData={formData}
            handleSubmit={handleSubmit}
            onClose={handleCloseForm}
            selectedTime={`${slots[selectedCard].start_time} - ${slots[selectedCard].end_time}`}
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

interface ModalProps {
  onClose: () => void;
}

const SuccessModal: React.FC<ModalProps> = ({ onClose }) => (
  <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div className="bg-white rounded-xl p-6 max-w-sm w-full">
      <div className="text-center">
        <CheckIcon className="text-green-500 text-5xl mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Booking Successful!
        </h2>
        <p className="text-gray-600 mb-6">
          Your slot has been successfully booked.
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

const ErrorModal: React.FC<ModalProps> = ({ onClose }) => (
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

export default Football_Booking;
