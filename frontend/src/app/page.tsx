"use client";
import React from 'react';
import { useAuth } from '../app/context/AuthContext';
import HomePage from '@/app/(pages)/homepage/page'

const Main = () => {
  const { user, logout } = useAuth();

  return (
    <div>
      {user ? (
        <HomePage/>
      ) : (
        <p>Loading...</p>
      )}
    </div>
  );
};

export default Main;
