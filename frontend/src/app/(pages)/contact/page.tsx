import React from 'react'
import NavBar from '../../components/navbar/navbar'

const page = () => {
  return (
    <div className='h-[645px]'>
        <NavBar activePage="contact"/>
        <div className='w-full h-full flex justify-center items-center'>
          <div className='bg-gray-100 rounded-md inline-flex flex-col gap-5 p-10 h-1/2 w-1/2'> 
            <span>6531503042@lamduan.mfu.ac.th</span>
            <span>Nimitsu jung benji tanbooutor</span>
            <span>IG : plscallfrank</span>
            <span>Facebook : Klavivach</span>
            <span>098-123-7654</span>
          </div>
        </div>
    </div>
  )
}

export default page