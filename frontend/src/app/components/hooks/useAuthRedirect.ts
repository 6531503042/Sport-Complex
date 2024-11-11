import { useEffect } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import { getToken } from "@/app/utils/auth";

const useAuthRedirect = () => {
  const router = useRouter();

  useEffect(() => {
    const checkAuth = async () => {
      const token = getToken();

      if (!token) {
        router.replace("/login");
        return;
      }

      try {
        const response = await axios.get("/api/auth/check-auth", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          router.replace("/login");
        }
      } catch (error) {
        console.error("Authorization check failed:", error);
        router.replace("/login");
      }
    };

    checkAuth();
  }, [router]);
};

export default useAuthRedirect;
