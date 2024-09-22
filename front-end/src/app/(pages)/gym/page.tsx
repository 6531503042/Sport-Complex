import React from 'react'
import NavBar from '../../components/navbar/navbar'

const page = () => {
  return (
    <div className='bg-red-500'>
        <NavBar activePage="gym"/>
        <div className='w-full h-screen inline-flex flex-row'>
          <div className='bg-red-500 flex-none w-2/3'>
          <div>
      <h1>Material Icon Football</h1>
      <i className="material-icons">sports_soccer</i>
    </div>
          </div>
          <div className='bg-black flex-none w-1/3'>2</div>
        </div>
    </div>
  )
}

export default page