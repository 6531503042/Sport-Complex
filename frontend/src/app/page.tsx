"use client";
import React from 'react';
import { useAuth } from '../app/context/AuthContext';
import HomePage from '@/app/(PagesPacking)/homepage/page'
import LoginPage from './(PagesPacking)/login/page';

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
