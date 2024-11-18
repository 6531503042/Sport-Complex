"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'react-hot-toast';

interface BookingFormProps {
  selectedSlot: {
    _id: string;
    start_time: string;
    end_time: string;
    max_bookings: number;
    current_bookings: number;
  };
  userId: string;
  facilityType: string;
  price: number;
}

const BookingForm: React.FC<BookingFormProps> = ({ selectedSlot, userId, facilityType, price }) => {
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleBooking = async () => {
    try {
      setIsLoading(true);
      
      // Validate inputs
      if (!selectedSlot || !userId) {
        toast.error('Missing required booking information');
        return;
      }

      // Format the request to match backend expectations
      const bookingRequest = {
        slot_id: selectedSlot._id,
        user_id: userId,
        facility_type: facilityType,
        booking_date: new Date(), // Current date
        status: 1,
        price: price,
        payment_status: "pending",
        slot_type: "normal", // Add this for regular facilities
        badminton_slot_id: null // Add this for regular facilities
      };

      console.log('Sending booking request:', bookingRequest);

      const response = await fetch(`http://localhost:1326/booking_v1/${facilityType}/booking`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify(bookingRequest)
      });

      if (!response.ok) {
        const errorData = await response.json();
        console.error('Booking error response:', errorData);
        throw new Error(errorData.message || `Booking failed with status: ${response.status}`);
      }

      const data = await response.json();
      console.log('Booking successful:', data);

      toast.success('Booking created successfully!');
      router.push('/booking-confirmation');
    } catch (error: any) {
      console.error('Booking error:', error);
      toast.error(error.message || 'Failed to create booking');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="p-4 bg-white rounded-lg shadow">
      <h3 className="text-lg font-semibold mb-4">Booking Details</h3>
      <div className="space-y-3">
        <div>
          <p className="text-sm text-gray-600">Time Slot</p>
          <p className="font-medium">{selectedSlot.start_time} - {selectedSlot.end_time}</p>
        </div>
        <div>
          <p className="text-sm text-gray-600">Facility</p>
          <p className="font-medium capitalize">{facilityType}</p>
        </div>
        <div>
          <p className="text-sm text-gray-600">Price</p>
          <p className="font-medium">{price} THB</p>
        </div>
        <button
          onClick={handleBooking}
          disabled={isLoading}
          className={`w-full py-2 px-4 bg-red-900 text-white rounded-lg hover:bg-red-800 transition-colors
            ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
        >
          {isLoading ? 'Processing...' : 'Confirm Booking'}
        </button>
      </div>
    </div>
  );
};

export default BookingForm; 