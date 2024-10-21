import React from 'react';

function Footer() {
  return (
    <footer className="bg-gray-800 text-white py-6">
      <div className="container mx-auto px-4 text-center">
        <p className="text-lg font-semibold">Sport Complex - Mae Fah Luang University</p>
        <p className="mt-2">Tel: <a href="tel:0-5391-7821" className="text-gray-400 hover:text-white">0-5391-7821</a></p>
        <p>Email: <a href="mailto:sport-complex@mfu.ac.th" className="text-gray-400 hover:text-white">sport-complex@mfu.ac.th</a></p>
        
        <p className="text-gray-500 mt-4">Copyright © 2021 - MFU Library</p>
      </div>
    </footer>
  );
}

export default Footer;