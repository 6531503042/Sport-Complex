import React from 'react'
import NavBar from '../../components/navbar/navbar'

const page = () => {
  return (
    <div>
        <NavBar activePage="football"/>
        <div className=''>
          Football
        </div>
    </div>
  )
}

export default page