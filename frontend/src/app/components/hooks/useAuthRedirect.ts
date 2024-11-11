import { useEffect } from "react";
import { useRouter } from "next/navigation";

const useAuthRedirect = () => {
  const router = useRouter();

  useEffect(() => {
    // Check if there is user data stored in localStorage
    const userData = localStorage.getItem("user");

    // If there's no user data (i.e., not logged in), redirect to login page
    if (!userData) {
      router.replace("/login");
    }
  }, [router]);
};

export default useAuthRedirect;
