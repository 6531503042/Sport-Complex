"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "../app/context/AuthContext"; // Assuming this is where user state is managed
import Homepage from "@/app/(pages)/homepage/page";
import LoginPage from "@/app/(pages)/login/page";
import LoadingScreen from "@/app/components/loading_screen/loading"; // Import loading screen component

const Main: React.FC = () => {
  const { user, setUser } = useAuth(); // Get user state from AuthContext
  const [loading, setLoading] = useState(true); // Manage loading state
  const router = useRouter();

  useEffect(() => {
    // Simulate authentication check (check localStorage if no user in context)
    const checkUserAuth = () => {
      const storedUser = localStorage.getItem("user");

      if (storedUser) {
        setUser(JSON.parse(storedUser)); // Set user from localStorage if it exists
      }

      setLoading(false); // Stop loading after the check
    };

    checkUserAuth(); // Run the check for user auth on page load
  }, [setUser]);

  useEffect(() => {
    // Redirect user based on authentication status
    if (!loading) {
      if (!user) {
        // Redirect to login if the user is not authenticated
        router.replace("/login");
      } else {
        // Redirect to homepage if the user is authenticated
        router.replace("/homepage");
      }
    }
  }, [loading, user, router]); // Dependency array includes loading, user, and router

  if (loading) {
    return <LoadingScreen />; // Show loading screen while checking user state
  }

  // If the loading is done and user is authenticated, show Homepage or LoginPage
  return <>{user ? <Homepage /> : <LoginPage />}</>;
};

export default Main;
