import React from 'react'
import NavBar from '../../components/navbar/navbar'
import Banner1 from '../../assets/banner_1.jpg'

const page = () => {
  return (
    <div>
      <NavBar />
      <div
        className="flex items-center justify-center h-[400px] text-white bg-cover bg-center"
        style={{
          backgroundImage: `url(${Banner1.src})`,
        }}
      >
        <div className="bg-black bg-opacity-50 p-4 rounded-md">
          <h1 className="text-3xl font-bold">Welcome to the Sport Complex</h1>
          <p className="mt-2">
            Discover our new features and start your fitness journey today!
          </p>
        </div>
      </div>
      <div className="p-4">
        homepage content
      </div>
    </div>
  )
}

export default page