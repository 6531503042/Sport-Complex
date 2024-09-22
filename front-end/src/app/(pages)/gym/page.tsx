import React from 'react'
import NavBar from '../../components/navbar/navbar'

const page = () => {
  return (
    <div className=''>
        <NavBar activePage="gym"/>
        <div className='w-full h-screen inline-flex flex-row'>
          <div className='bg-red-500 flex-none w-2/3'>1</div>
          <div className='bg-black flex-none w-1/3'>2</div>
        </div>
    </div>
  )
}

export default page