"use client";

import React, { useCallback } from "react";
import useEmblaCarousel from "embla-carousel-react";
import Autoplay from 'embla-carousel-autoplay'
import './football-booking.css'


interface ImageLink {
    imageUrl: string;
    linkUrl: string;
  }

function Football_Booking() {
    const [emblaRef,emblaApi] = useEmblaCarousel({ loop: true }, [Autoplay({delay:5000})]) // 5s

    const scrollPrev = useCallback(() => {
        if (emblaApi) emblaApi.scrollPrev()
      }, [emblaApi])
    
      const scrollNext = useCallback(() => {
        if (emblaApi) emblaApi.scrollNext()
      }, [emblaApi])


      const imageLinks: ImageLink[] = [
        {
          imageUrl: 'https://en.mfu.ac.th/fileadmin/mainsite_news_eng/news/2024/TCUE.jpg',
          linkUrl: 'https://en.mfu.ac.th/fileadmin/mainsite_news_eng/news/2024/TCUE.jpg',
        },
        {
          imageUrl: 'https://www.mfu.ac.th/fileadmin/mainsite_news_thai/Au_Photo_PR_news/2565/10/NEWS-MFU-WEB.png',
          linkUrl: 'https://www.mfu.ac.th/fileadmin/mainsite_news_thai/Au_Photo_PR_news/2565/10/NEWS-MFU-WEB.png',
        },
        {
            imageUrl: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQXqoUfOanL0mvHBzukEfzTyWTzBZk-ss1FsQ&s',
            linkUrl: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQXqoUfOanL0mvHBzukEfzTyWTzBZk-ss1FsQ&s',
          },
        // Add more objects as needed
      ];

  return (
    <div>
      <div className="embla">
        <div className="embla_viewport mx-auto h-96 max-w-full border" ref={emblaRef}>
        <div className="embla_container h-full">
            {imageLinks.map((item, index) => (
              <div className="embla_slide flex items-center justify-center" key={index}>
                <a href={item.linkUrl} target="_blank" rel="noopener noreferrer">
                  <img src={item.imageUrl} alt={`Slide ${index}`} className="w-full h-full object-cover" />
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
      <div>

      </div>
    </div>
  );
}

export default Football_Booking;
