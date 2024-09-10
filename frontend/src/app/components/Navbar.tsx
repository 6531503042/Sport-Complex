import React from 'react'
import Link from 'next/link'

function Navbar() {
  return (
    <div >
      <ul>
        <li><Link href="/">Home</Link></li>
        <li><Link href="/football-booking">Football</Link></li>
        <li><Link href="/swimming-booking">Swimming</Link></li>
      </ul>
    </div>
  )
}

export default Navbar
