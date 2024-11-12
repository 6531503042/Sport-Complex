"use client";
import React, { useState, useEffect } from 'react';
import { useAuth } from '../app/context/AuthContext';
import HomePage from '@/app/(pages)/homepage/page';
import LoginPage from './(pages)/login/page';
import LoadingScreen from '@/app/components/loading_screen/loading'; // Adjust the import path as needed

const Main = () => {
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Simulate a loading delay or use actual loading state here
    const loadTimeout = setTimeout(() => {
      setLoading(false);
    }, 2000); // 2-second loading simulation

    return () => clearTimeout(loadTimeout);
  }, [user]);

  if (loading) {
    return <LoadingScreen />;
  }

  return (
    <div>
      {user ? <HomePage /> : <HomePage />}
    </div>
  );
};

export default Main;
