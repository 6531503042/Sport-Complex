"use client";

import React, { useEffect, useState, Fragment } from "react";
import "./gym.css";
import NavBar from "@/app/components/navbar/navbar";
import CheckIcon from "@mui/icons-material/Check";
import ClearIcon from "@mui/icons-material/Clear";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import GroupIcon from "@mui/icons-material/Group";
import ReportProblemIcon from "@mui/icons-material/ReportProblem";
import Link from "next/link"
import { useRouter } from 'next/navigation';
import FitnessCenterIcon from '@mui/icons-material/FitnessCenter';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import PeopleIcon from '@mui/icons-material/People';
import EventAvailableIcon from '@mui/icons-material/EventAvailable';
import WarningIcon from '@mui/icons-material/Warning';
import { motion, AnimatePresence } from 'framer-motion';
import { Tooltip } from '@mui/material';

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

interface ModalProps {
  onClose: () => void;
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

// Helper functions
const isSlotBooked = (slot: Slot): boolean => {
  return slot.current_bookings >= slot.max_bookings;
};

const getStatusClasses = (slot: Slot): string => {
  const isBooked = isSlotBooked(slot);
  const remainingSpots = slot.max_bookings - slot.current_bookings;

  if (isBooked) return 'bg-red-100 text-red-800';
  if (remainingSpots <= 2) return 'bg-yellow-100 text-yellow-800';
  return 'bg-green-100 text-green-800';
};

const getStatusText = (slot: Slot): string => {
  const isBooked = isSlotBooked(slot);
  const remainingSpots = slot.max_bookings - slot.current_bookings;

  if (isBooked) return 'Fully Booked';
  return `${remainingSpots} ${remainingSpots === 1 ? 'spot' : 'spots'} left`;
};

const BookingForm = ({ formData, handleSubmit, onClose, selectedTime }: BookingFormProps) => (
  <motion.div 
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    exit={{ opacity: 0 }}
    className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50"
  >
    <motion.div 
      initial={{ scale: 0.9, y: 20 }}
      animate={{ scale: 1, y: 0 }}
      exit={{ scale: 0.9, y: 20 }}
      className="bg-white rounded-2xl p-8 max-w-md w-full shadow-2xl"
    >
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 mb-1">Complete Booking</h2>
          <p className="text-sm text-gray-500">Fill in your details to confirm your slot</p>
        </div>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-gray-600 transition-colors p-2 hover:bg-gray-100 rounded-full"
        >
          ✕
        </button>
      </div>

      {selectedTime && (
        <motion.div 
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-6 p-4 bg-purple-50 rounded-xl border border-purple-100"
        >
          <div className="flex items-center gap-3">
            <div className="p-2 bg-purple-100 rounded-lg">
              <AccessTimeIcon className="text-purple-600" />
            </div>
            <div>
              <p className="text-sm text-purple-600 font-medium">Selected Time</p>
              <p className="text-lg font-semibold text-purple-900">{selectedTime}</p>
            </div>
          </div>
        </motion.div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-4">
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
                          focus:ring-2 focus:ring-purple-500 focus:border-transparent
                          transition-all duration-200"
              />
              <Tooltip title="Name cannot be changed" arrow>
                <div className="absolute right-3 top-1/2 -translate-y-1/2">
                  <WarningIcon className="text-gray-400" />
                </div>
              </Tooltip>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              ID Number
            </label>
            <div className="relative">
              <input
                type="text"
                value={formData.id}
                readOnly
                className="text-black w-full px-4 py-3 rounded-xl border border-gray-300 bg-gray-50
                          focus:ring-2 focus:ring-purple-500 focus:border-transparent
                          transition-all duration-200"
              />
              <Tooltip title="ID cannot be changed" arrow>
                <div className="absolute right-3 top-1/2 -translate-y-1/2">
                  <WarningIcon className="text-gray-400" />
                </div>
              </Tooltip>
            </div>
          </div>
        </div>

        <div className="flex gap-4 pt-6">
          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            type="button"
            onClick={onClose}
            className="flex-1 px-4 py-3 text-gray-700 bg-gray-100 rounded-xl
                     hover:bg-gray-200 transition-colors"
          >
            Cancel
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            type="submit"
            className="flex-1 px-4 py-3 bg-gradient-to-r from-purple-600 to-indigo-600
                     text-white rounded-xl font-medium
                     hover:from-purple-700 hover:to-indigo-700
                     transition-colors shadow-lg hover:shadow-xl"
          >
            Confirm Booking
          </motion.button>
        </div>
      </form>
    </motion.div>
  </motion.div>
);

// Success Modal Component
const SuccessModal: React.FC<ModalProps> = ({ onClose }) => (
  <motion.div
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    exit={{ opacity: 0 }}
    className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50"
  >
    <motion.div
      initial={{ scale: 0.9, y: 20 }}
      animate={{ scale: 1, y: 0 }}
      exit={{ scale: 0.9, y: 20 }}
      className="bg-white rounded-xl p-6 max-w-sm w-full"
    >
      <div className="text-center">
        <CheckIcon className="text-green-500 text-5xl mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Booking Successful!
        </h2>
        <p className="text-gray-600 mb-6">
          Your slot has been successfully booked.
        </p>
        <motion.button
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={onClose}
          className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 text-white py-3 rounded-lg font-semibold"
        >
          Done
        </motion.button>
      </div>
    </motion.div>
  </motion.div>
);

// Error Modal Component
const ErrorModal: React.FC<ModalProps> = ({ onClose }) => (
  <motion.div
    initial={{ opacity: 0 }}
    animate={{ opacity: 1 }}
    exit={{ opacity: 0 }}
    className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50"
  >
    <motion.div
      initial={{ scale: 0.9, y: 20 }}
      animate={{ scale: 1, y: 0 }}
      exit={{ scale: 0.9, y: 20 }}
      className="bg-white rounded-xl p-6 max-w-sm w-full"
    >
      <div className="text-center">
        <ClearIcon className="text-red-500 text-5xl mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Booking Failed
        </h2>
        <p className="text-gray-600 mb-6">
          User was Booked this time
        </p>
        <motion.button
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={onClose}
          className="w-full bg-red-500 text-white py-3 rounded-lg font-semibold hover:bg-red-600"
        >
          Close
        </motion.button>
      </div>
    </motion.div>
  </motion.div>
);

function Gym_Booking({ params }: UserDataParams) {
  const router = useRouter();
  const { id } = params;
  const [storedRefreshToken, setStoredRefreshToken] = useState<string | null>(
    null
  );
  const [paymentId, setPaymentId] = useState<string | null>(null); // State for storing payment_id

  const [selectedCard, setSelectedCard] = useState<number | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    id: "",
    phone: "",
  });

  const [isBookingSuccessful, setIsBookingSuccessful] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false);
  const [isBookingFailed, setIsBookingFailed] = useState(false);

  const [slots, setSlots] = useState<Slot[]>([]);
  const [showBookingForm, setShowBookingForm] = useState(false);
  const [handleCloseForm] = useState(() => () => {
    setShowBookingForm(false);
    setSelectedCard(null);
  });

  const handleCardClick = (slot: Slot) => {
    if (!isSlotBooked(slot)) {
      const index = slots.findIndex(s => s._id === slot._id);
      if (index !== -1) {
        setSelectedCard(index);
        setShowBookingForm(true);
      }
    }
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

      const response = await fetch("http://localhost:1326/booking_v1/fitness/booking", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(bookingData),
      });

      if (!response.ok) {
        setIsBookingFailed(true);
        return;
      }

      const result = await response.json();
      
      if (!result.payment_id) {
        setIsBookingFailed(true);
        return;
      }

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
        "http://localhost:1335/facility_v1/fitness/slot_v1/slots",
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
      setSlots([]);
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

    // Fetch slot data on initial render and set up the interval for updating
    getSlot();
    const intervalId = setInterval(getSlot, 1000);

    return () => {
      clearInterval(intervalId);
    };
  }, [id]);

  return (
    <>
      <NavBar activePage="gym" />
      <div className="min-h-screen bg-gradient-to-br from-purple-50 via-indigo-50 to-white py-8 px-4">
        <div className="max-w-7xl mx-auto">
          {/* Enhanced Header Section */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center mb-12"
          >
            <motion.div 
              whileHover={{ scale: 1.1, rotate: 360 }}
              transition={{ duration: 0.5 }}
              className="inline-block p-4 rounded-full bg-gradient-to-br from-purple-100 to-indigo-100 mb-6"
            >
              <FitnessCenterIcon className="text-purple-600 text-4xl" />
            </motion.div>
            <h1 className="text-5xl font-bold text-gray-900 mb-4 bg-gradient-to-r from-purple-600 to-indigo-600 inline-block text-transparent bg-clip-text">
              Fitness Center Booking
            </h1>
            <p className="text-gray-600 text-xl max-w-2xl mx-auto">
              Reserve your spot for a great workout session
            </p>
          </motion.div>

         

          {/* Slots Grid */}
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
          >
            {slots?.map((slot, index) => (
              <motion.div
                key={slot._id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                className={`
                  relative rounded-xl overflow-hidden transition-all duration-300
                  ${isSlotBooked(slot) 
                    ? 'bg-gray-50 cursor-not-allowed' 
                    : 'bg-white hover:bg-purple-50 cursor-pointer transform hover:scale-105'
                  }
                  shadow-lg hover:shadow-xl
                `}
                onClick={() => !isSlotBooked(slot) && handleCardClick(slot)}
              >
                {/* Time Header */}
                <div className="bg-gradient-to-r from-purple-600 to-indigo-600 text-white p-6">
                  <div className="flex items-center justify-center gap-3">
                    <AccessTimeIcon className="text-purple-200" />
                    <h3 className="text-xl font-semibold">
                      {slot.start_time} - {slot.end_time}
                    </h3>
                  </div>
                </div>

                {/* Slot Details */}
                <div className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-2">
                      <PeopleIcon className="text-gray-400" />
                      <span className="text-gray-600 font-medium">
                        {slot.current_bookings} / {slot.max_bookings}
                      </span>
                    </div>
                    <motion.div
                      whileHover={{ scale: 1.05 }}
                      className={`
                        px-4 py-2 rounded-full text-sm font-medium
                        ${getStatusClasses(slot)}
                      `}
                    >
                      {getStatusText(slot)}
                    </motion.div>
                  </div>

                  {!isSlotBooked(slot) && (
                    <motion.div
                      animate={{ scale: [1, 1.2, 1] }}
                      transition={{ repeat: Infinity, duration: 2 }}
                      className="absolute top-2 right-2"
                    >
                      <div className="w-3 h-3 bg-green-400 rounded-full shadow-lg shadow-green-200" />
                    </motion.div>
                  )}
                </div>

                {isSlotBooked(slot) && (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="absolute inset-0 bg-gray-100/50 backdrop-blur-[1px] 
                             flex items-center justify-center"
                  >
                    <motion.span
                      animate={{ rotate: [0, -12, 0] }}
                      transition={{ repeat: Infinity, duration: 2 }}
                      className="text-red-600 font-bold text-xl border-2 border-red-600 
                               px-6 py-3 rounded-lg bg-white/90 shadow-xl"
                    >
                      Maximum
                    </motion.span>
                  </motion.div>
                )}
              </motion.div>
            ))}
          </motion.div>

          {/* Modals */}
          <AnimatePresence>
            {showBookingForm && selectedCard !== null && slots[selectedCard] && (
              <BookingForm 
                formData={formData}
                handleSubmit={handleSubmit}
                onClose={handleCloseForm}
                selectedTime={`${slots[selectedCard].start_time} - ${slots[selectedCard].end_time}`}
              />
            )}

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
          </AnimatePresence>
        </div>
      </div>
    </>
  );
}

export default Gym_Booking;

