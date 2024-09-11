"use client";

import React, { useCallback } from "react";
import useEmblaCarousel from "embla-carousel-react";
import Autoplay from "embla-carousel-autoplay";
import "./football-booking.css";
import Footer from "@/app/components/Footer";

interface ImageLink {
  imageUrl: string;
  linkUrl: string;
}

function Football_Booking() {
  const [emblaRef, emblaApi] = useEmblaCarousel({ loop: true }, [
    Autoplay({ delay: 5000 }),
  ]); // 5s

  const scrollPrev = useCallback(() => {
    if (emblaApi) emblaApi.scrollPrev();
  }, [emblaApi]);

  const scrollNext = useCallback(() => {
    if (emblaApi) emblaApi.scrollNext();
  }, [emblaApi]);

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
    // Add more objects as needed
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

      <hr className="my-10 mx-5 " />

      <div className="mx-20 mb-20">
        <h1 className="text-center text-3xl font-bold mb-20">
          FOOTBALL BOOKING
        </h1>
        <div className="grid grid-cols-2">
          <div className="container_list_time w-[1000px] h-[700px] bg-[#387F39] rounded-2xl p-10 overflow-y-auto shadow-2xl">
            <div className="grid grid-cols-1 gap-6">
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                1
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                2
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                3
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                4
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                5
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                6
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                7
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                8
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                9
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                10
              </div>
              <div className="list-time bg-white p-4 rounded-lg text-center font-semibold text-lg shadow-2xl">
                11
              </div>
            </div>
          </div>
          <div className="flex justify-center items-center">
            <div className="text-center text-lg">
              Please select a time for booking...
            </div>
          </div>
        </div>
      </div>
        <Footer />
    </div>
  );
}

export default Football_Booking;

{
  /* <div>
            <form action="">
              <label htmlFor="name">Name</label>
              <input type="text" name="name" placeholder="Enter your name"/>
              <label htmlFor="id">Lecturer / Staff / Student ID</label>
              <input type="text" name="id" placeholder="Enter your ID"/>
              <label htmlFor="phone-number">Phone Number</label>
              <input type="number" name="phone-number"  placeholder="Enter your Phone number"/>
            </form>
           </div> */
}
