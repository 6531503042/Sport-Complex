import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

const useAuthRedirect = () => {
  const [loading, setLoading] = useState(true); // Track if loading is in progress
  const [user, setUser] = useState<any | null>(null); // Store user data if authenticated
  const router = useRouter();

  useEffect(() => {
    // Check authentication on component mount
    const checkAuth = () => {
      const userData = localStorage.getItem("user"); // Check if user data is in localStorage

      if (userData) {
        setUser(JSON.parse(userData)); // Set the user if authenticated
      }

      setLoading(false); // Stop loading once the check is complete
    };

    if (typeof window !== "undefined") {
      checkAuth(); // Run the authentication check
    }
  }, []);

  useEffect(() => {
    // If loading is complete and the user is not authenticated, redirect to login
    if (!loading && !user) {
      router.replace("/login");
    }
  }, [loading, user, router]);

  return { loading, user }; // Return loading and user data
};

export default useAuthRedirect;
