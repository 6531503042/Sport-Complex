// Example of correct booking request format
const bookingRequest = {
  slot_id: selectedSlot._id,
  user_id: userId,
  facility_type: "fitness", // Make sure this matches the backend expectation
  booking_date: selectedDate,
  status: 1, // 1 for active booking
  price: price,
  payment_status: "pending"
};

try {
  const response = await fetch('http://localhost:1326/booking_v1/fitness/booking', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(bookingRequest)
  });

  if (!response.ok) {
    const errorData = await response.json();
    console.error('Booking error:', errorData);
    throw new Error(errorData.message || 'Failed to create booking');
  }

  const data = await response.json();
  console.log('Booking successful:', data);
} catch (error) {
  console.error('Error creating booking:', error);
  toast.error('Failed to create booking. Please try again.');
} 