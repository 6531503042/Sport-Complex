"use client";
import React from 'react';
import { useAuth } from '../app/context/AuthContext';
import HomePage from '@/app/(pages)/homepage/page'
import LoginPage from './(pages)/login/page';

const Main = () => {
  const { user, logout } = useAuth();

  return (
    <div>
      {user ? (
        <HomePage/>
      ) : (
        <LoginPage/>
      )}
    </div>
  );
};

export default Main;
