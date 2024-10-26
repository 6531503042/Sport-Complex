'use client'
import NavBar from '../../components/navbar/navbar';
import '@fortawesome/fontawesome-free/css/all.min.css';
import { Button, Card, CardBody, CardHeader, Input } from "@nextui-org/react";
import { useState } from "react";

const Page = () => {
  const [visible, setVisible] = useState(false); // state for modal visibility
  const [selectedCourt, setSelectedCourt] = useState(""); // state for selected court

  const cards = [
    {
      title: "13:00 - 14:00",
      icon: "fa-clock",
      subCards: [
        { title: "Court 1", status: "available" },
        { title: "Court 2", status: "booked" },
        { title: "Court 3", status: "booked" },
        { title: "Court 4", status: "booked" },
      ],
    },
    {
      title: "14:00 - 15:00",
      icon: "fa-clock",
      subCards: [
        { title: "Court 1", status: "booked" },
        { title: "Court 2", status: "available" },
        { title: "Court 3", status: "booked" },
        { title: "Court 4", status: "booked" },
      ],
    },
    {
      title: "15:00 - 16:00",
      icon: "fa-clock",
      subCards: [
        { title: "Court 1", status: "booked" },
        { title: "Court 2", status: "booked" },
        { title: "Court 3", status: "available" },
        { title: "Court 4", status: "booked" },
      ],
    },
    {
      title: "16:00 - 17:00",
      icon: "fa-clock",
      subCards: [
        { title: "Court 1", status: "booked" },
        { title: "Court 2", status: "booked" },
        { title: "Court 3", status: "booked" },
        { title: "Court 4", status: "available" },
      ],
    },
    {
      title: "17:00 - 18:00",
      icon: "fa-clock",
      subCards: [
        { title: "Court 1", status: "booked" },
        { title: "Court 2", status: "booked" },
        { title: "Court 3", status: "booked" },
        { title: "Court 4", status: "available" },
      ],
    },
    {
      title: "18:00 - 19:00",
      icon: "fa-clock",
      subCards: [
        { title: "Court 1", status: "booked" },
        { title: "Court 2", status: "booked" },
        { title: "Court 3", status: "booked" },
        { title: "Court 4", status: "available" },
      ],
    },
  ];

  const [formData, setFormData] = useState({ name: "", id: "", phone: "" }); // useState for form data

  const handleButtonClick = (courtTitle: string) => {
    setSelectedCourt(courtTitle); // set the selected court title
    setVisible(true); // show the modal
  };

  const handleClose = () => {
    setVisible(false);
    setFormData({ name: "", id: "", phone: "" }); // Reset form data when closing
    setSelectedCourt(""); // Reset selected court
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value }); // Update form data
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Handle booking submission logic here
    console.log("Booking submitted:", formData); // Corrected formData reference
    handleClose(); // Close the modal after submission
  };

  return (
    <div className="flex flex-col bg-gray-100 h-screen">
      {/*  NavBar   */}
      <NavBar activePage="gym" /> {/* Change "gym" to the appropriate page based on your routing */}

      {/* Title */}
      <h1 className="flex flex-col items-center text-3xl font-bold mb-8 mt-9">Badminton Booking</h1>
      <p className="flex flex-col text-gray-600 items-center text-sm mb-3">Click on an available court to book</p>

      {/* Render cards */}
      <div className="flex flex-col items-center mb-4 mt-9">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {cards.map((card, index) => (
            <Card key={index} className="py-4 bg-gray-200 relative rounded-lg shadow-lg mb-8 w-96">
              <CardHeader className="pb-0 pt-2 px-4 flex-col items-start">
                <div className="flex items-center">
                  <i className={`fa ${card.icon} mr-2`} aria-hidden="true"></i>
                  <p className="text-tiny uppercase font-bold">{card.title}</p>
                </div>
              </CardHeader>
              {/* Calendar icon in the top right corner */}
              <div className="absolute top-2 right-2 mt-3 mb-10">
                <i className="fas fa-calendar-alt" style={{ color: "#A9A9A9", marginLeft: "-25px" }} aria-hidden="true"></i>
              </div>
              {/* Show 4 subCards */}
              <div className="grid grid-cols-2 gap-3 p-4">
                {card.subCards.map((subCard, smallIndex) => (
                  <Card key={smallIndex} className="bg-gray-300 p-2 rounded-md">
                    <CardBody>
                      <div className="flex items-center">
                        {/* Show check or cross icon based on status */}
                        <span className="mr-2">
                          {subCard.status === "available" ? (
                            <i className="fas fa-check-circle text-green-500" aria-hidden="true"></i>
                          ) : (
                            <i className="fas fa-times-circle text-red-900" aria-hidden="true"></i>
                          )}
                        </span>
                        {/* Booking button */}
                        <Button
                          color={subCard.status === "available" ? "primary" : "default"}
                          className={subCard.status === "booked" ? "text-gray-500 cursor-not-allowed" : "text-green-600"}
                          onClick={() => handleButtonClick(subCard.title)}
                          disabled={subCard.status === "booked"}
                        >
                          {subCard.title}
                        </Button>
                      </div>
                    </CardBody>
                  </Card>
                ))}
              </div>
            </Card>
          ))}
        </div>
      </div>
      {/* Modal for booking confirmation */}
      {visible && (
        <>
          <div
            style={{
              position: "fixed",
              top: 0,
              left: 0,
              width: "100%",
              height: "100%",
              backgroundColor: "rgba(0,0,0,0.5)",
              zIndex: 999,
            }}
            onClick={handleClose}
          />
          <div
            style={{
              position: "fixed",
              top: "50%",
              left: "50%",
              transform: "translate(-50%, -50%)",
              zIndex: 1000,
              backgroundColor: "#f9f9f9",
              padding: "20px",
              borderRadius: "12px",
              boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
              maxWidth: "400px",
              width: "100%",
              margin: "20px auto",
            }}
          >
            <button
              onClick={handleClose}
              style={{
                position: "absolute",
                top: "5px",
                right: "10px",
                background: "transparent",
                border: "none",
                fontSize: "20px",
                cursor: "pointer",
              }}
            >
              &times;
            </button>

            {/* overlay */}
            <h3 style={{
              fontSize: "1.5em",
              fontWeight: "bold",
              color: "#E74C3C",
              textAlign: "center",
              margin: "20px 0",
              textTransform: "uppercase",
              letterSpacing: "1.5px",
              textShadow: "1px 1px 2px rgba(0, 0, 0, 0.2)",
            }}>Booking Confirmation</h3>

            <p style={{ fontSize: "1em", marginBottom: "20px", textAlign: "center" }}>
              You have selected {selectedCourt}. Please provide your information to complete the booking.
            </p>

            {/* Form for user information */}
            <form onSubmit={handleSubmit}>
              <div className="mb-4">
                <Input
                  type="text"
                  name="name"
                  placeholder="Name"
                  value={formData.name}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className="mb-4">
                <Input
                  type="text"
                  name="id"
                  placeholder="ID Card Number"
                  value={formData.id}
                  onChange={handleChange}
                  required
                />
              </div>
              <div className="mb-4">
                <Input
                  type="text"
                  name="phone"
                  placeholder="Phone Number"
                  value={formData.phone}
                  onChange={handleChange}
                  required
                />
              </div>
              <Button type="submit" color="success" className="w-full">Confirm Booking</Button>
            </form>
          </div>
        </>
      )}
    </div>
  );
};

export default Page;
