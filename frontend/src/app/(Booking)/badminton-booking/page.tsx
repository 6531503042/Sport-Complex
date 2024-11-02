"use client";
import "@fortawesome/fontawesome-free/css/all.min.css";
import { Button, Input } from "@nextui-org/react";
import { useState, ChangeEvent, FormEvent } from "react";
import NavBar from "../../components/navbar/navbar";
import "./badminton.css";
import { AccessTime, Person, Badge, Close } from "@mui/icons-material";

interface FormData {
  name: string;
  id: string;
}

const timeSlots = ["10:00-11:00", "11:00-12:00", "14:00-15:00", "15:00-16:00"];

const Badminton_Booking: React.FC = () => {
  const [visible, setVisible] = useState<boolean>(false); // state for modal visibility
  const [selectedCourt, setSelectedCourt] = useState<string>(""); // state for selected court
  const [selectedTimeSlot, setSelectedTimeSlot] = useState<string>(""); // state for selected time slot
  const [formData, setFormData] = useState<FormData>({ name: "", id: "" });

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    console.log("Booking confirmed:", {
      ...formData,
      court: selectedCourt,
      timeSlot: selectedTimeSlot,
    });
    setVisible(false); // Close the modal after submission
  };

  const handleClose = () => {
    setVisible(false);
    setSelectedCourt("");
    setSelectedTimeSlot("");
  };

  const handleCourtSlotSelect = (courtName: string, slot: string) => {
    setSelectedCourt(courtName);
    setSelectedTimeSlot(slot);
    setVisible(true); // Open modal
  };

  const courts: string[] = [
    "Court 1",
    "Court 2",
    "Court 3",
    "Court 4",
    "Court 5",
    "Court 6",
  ];

  return (
    <div className="flex flex-col bg-gray-100 min-h-screen">
      <NavBar activePage="gym" />
      <h1 className="text-center text-3xl font-bold mt-9 text-gray-800">
        Badminton Booking
      </h1>
      <p className="text-center text-gray-600 text-base mt-2 mb-8">
        Select a court and time slot to book
      </p>

      <div className="flex justify-center mb-10">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 px-4">
          {courts.map((court) => (
            <div
              key={court}
              className="p-6 bg-white shadow-xl rounded-lg hover:shadow-2xl transition transform hover:scale-105 cursor-pointer"
            >
              <h2 className="text-center text-lg font-semibold mb-4 text-gray-700">
                {court}
              </h2>
              <div className="grid grid-cols-2 gap-3">
                {timeSlots.map((slot) => (
                  <button
                    key={slot}
                    type="button"
                    className={`p-3 rounded-lg text-sm font-medium transition ${
                      selectedCourt === court && selectedTimeSlot === slot
                        ? "bg-green-500 text-white" // Selected slot style
                        : "bg-gray-200 text-gray-700 hover:bg-green-100"
                    }`}
                    onClick={() => handleCourtSlotSelect(court, slot)}
                  >
                    <AccessTime className="mr-1" /> {slot}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Modal for booking confirmation */}
      {visible && (
        <div className="fixed inset-0 flex justify-center items-center bg-black bg-opacity-50">
          <div className="modal-content relative bg-white p-8 rounded-lg shadow-2xl max-w-lg mx-auto transition transform scale-95">
            <button
              className="absolute top-4 right-4 text-gray-500 text-2xl hover:text-red-600 transition"
              onClick={handleClose}
            >
              <Close />
            </button>
            <h3 className="text-xl font-semibold mb-4 text-gray-800">
              Booking Confirmation
            </h3>
            <p className="text-gray-600 mb-6 text-base">
              You have selected{" "}
              <strong className="text-gray-800">{selectedCourt}</strong> at{" "}
              <strong className="text-gray-800">{selectedTimeSlot}</strong>.
              Please enter your information to confirm the booking.
            </p>

            <form onSubmit={handleSubmit}>
              <div className="modal-input flex items-center mb-3">
                <Person className="mr-2 text-gray-500" /> {/* Name icon */}
                <Input
                  type="text"
                  name="name"
                  placeholder="Name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className="modal-input flex items-center mb-3">
                <Badge className="mr-2 text-gray-500" /> {/* ID icon */}
                <Input
                  type="text"
                  name="id"
                  placeholder="ID Card Number"
                  value={formData.id}
                  onChange={handleChange}
                  required
                />
              </div>
              <Button
                type="submit"
                color="success"
                className="w-full py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition"
              >
                Confirm Booking
              </Button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Badminton_Booking;
