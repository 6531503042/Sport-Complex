"use client";

import React, { useState, useCallback } from "react";
import useEmblaCarousel from "embla-carousel-react";
import Autoplay from "embla-carousel-autoplay";
import "./football-booking.css";
import Footer from "@/app/components/Footer";

interface ImageLink {
  imageUrl: string;
  linkUrl: string;
}

interface TimeSlot {
  time: string;
  isAvailable: boolean;
  member:string;
}

function Football_Booking() {
  const [emblaRef, emblaApi] = useEmblaCarousel({ loop: true }, [
    Autoplay({ delay: 5000 }),
  ]); // 5s
  const [selectedTime, setSelectedTime] = useState<string | null>(null);
  const [showPopup, setShowPopup] = useState(false);

  const timeSlots: TimeSlot[] = [
    { time: "10:00 - 11:00", isAvailable: true ,member:"0/1"},
    { time: "11:00 - 12:00", isAvailable: false ,member:"1/1"},
    { time: "12:00 - 13:00", isAvailable: true ,member:"0/1"},
    { time: "13:00 - 14:00", isAvailable: true ,member:"0/1"},
    { time: "14:00 - 15:00", isAvailable: false ,member:"1/1"},
    { time: "15:00 - 16:00", isAvailable: true,member:"0/1" },
    { time: "16:00 - 17:00", isAvailable: true ,member:"0/1"},
    { time: "17:00 - 18:00", isAvailable: false ,member:"1/1"},
    { time: "18:00 - 19:00", isAvailable: true ,member:"0/1"},
    { time: "19:00 - 20:00", isAvailable: true ,member:"0/1"},
  ];

  const scrollPrev = useCallback(() => {
    if (emblaApi) emblaApi.scrollPrev();
  }, [emblaApi]);

  const scrollNext = useCallback(() => {
    if (emblaApi) emblaApi.scrollNext();
  }, [emblaApi]);

  const handleTimeClick = (time: string, isAvailable: boolean) => {
    if (isAvailable) {
      setSelectedTime(time);
      setShowPopup(true);
    } else {
      alert("This time slot is unavailable.");
    }
  };

  const closePopup = () => {
    setShowPopup(false);
    setSelectedTime(null);
  };

  const imageLinks: ImageLink[] = [
    {
      imageUrl:
        "https://en.mfu.ac.th/fileadmin/mainsite_news_eng/news/2024/TCUE.jpg",
      linkUrl:
        "https://en.mfu.ac.th/fileadmin/mainsite_news_eng/news/2024/TCUE.jpg",
    },
    {
      imageUrl:
        "https://www.mfu.ac.th/fileadmin/mainsite_news_thai/Au_Photo_PR_news/2565/10/NEWS-MFU-WEB.png",
      linkUrl:
        "https://www.mfu.ac.th/fileadmin/mainsite_news_thai/Au_Photo_PR_news/2565/10/NEWS-MFU-WEB.png",
    },
    {
      imageUrl:
        "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQXqoUfOanL0mvHBzukEfzTyWTzBZk-ss1FsQ&s",
      linkUrl:
        "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQXqoUfOanL0mvHBzukEfzTyWTzBZk-ss1FsQ&s",
    },
  ];

  return (
    <div>
      <div className="embla">
        <div
          className="embla_viewport mx-auto h-80 max-w-full border"
          ref={emblaRef}
        >
          <div className="embla_container h-full">
            {imageLinks.map((item, index) => (
              <div
                className="embla_slide flex items-center justify-center"
                key={index}
              >
                <a
                  href={item.linkUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <img
                    src={item.imageUrl}
                    alt={`Slide ${index}`}
                    className="w-full h-full object-cover"
                  />
                </a>
              </div>
            ))}
          </div>
        </div>
        <div>
          <button className="embla_prev" onClick={scrollPrev}>
            Prev
          </button>
          <span> </span>
          <button className="embla_next" onClick={scrollNext}>
            Next
          </button>
        </div>
      </div>

      <hr className="my-10 mx-5" />

      <div className="mx-20 mb-20">
        <h1 className="text-center text-3xl font-bold mb-20">
          FOOTBALL BOOKING
        </h1>
        <div className="grid grid-cols-2">
          <div className="container_list_time h-[700px] bg-[#387F39] rounded-2xl p-10 overflow-y-auto shadow-2xl">
            <div className="grid grid-cols-1 gap-6 ">
              {timeSlots.map((slot, index) => (
                <div
                key={index}
                className={`list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl cursor-pointer h-[114px] ${
                  slot.isAvailable ? "text-black cursor-pointer" : "text-[#AAAAAA] cursor-not-allowed"
                }`}
                onClick={() => handleTimeClick(slot.time, slot.isAvailable)}
              >
                {slot.time} {"  "}
                {slot.isAvailable ? "Available" : "Unavailable"} {"  "}
                {slot.member}
              </div>
              
              ))}
            </div>
          </div>
          <div className="flex justify-center items-center">
            <div className="text-center text-lg">
              Please select a time for booking...
            </div>
          </div>
        </div>
      </div>

      {showPopup && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center">
          <div className="bg-white p-8 rounded-lg shadow-lg w-[600px] h-[450px]">
            <h2 className="text-xl font-semibold mb-4">Booking Time: {selectedTime}</h2>
            <form>
              <label htmlFor="name" className="block mb-2">
                Name:
              </label>
              <input
                type="text"
                id="name"
                className="w-full border p-2 mb-4"
                placeholder="Enter your name"
              />
              <label htmlFor="id" className="block mb-2">
                Lecturer / Staff / Student ID:
              </label>
              <input
                type="text"
                id="id"
                className="w-full border p-2 mb-4"
                placeholder="Enter your ID"
              />
              <label htmlFor="phone" className="block mb-2">
                Phone Number:
              </label>
              <input
                type="text"
                id="phone"
                className="w-full border p-2 mb-4"
                placeholder="Enter your phone number"
              />
              <div className="flex justify-end">
                <button
                  type="button"
                  onClick={closePopup}
                  className="bg-red-500 text-white px-4 py-2 rounded-lg"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="bg-blue-500 text-white px-4 py-2 rounded-lg ml-2"
                >
                  Book Now
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <Footer />
    </div>
  );
}

export default Football_Booking;
