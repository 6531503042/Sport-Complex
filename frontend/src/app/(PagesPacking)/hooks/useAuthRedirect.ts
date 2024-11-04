import { useEffect } from "react";
import { useRouter } from "next/navigation";

const useAuthRedirect = () => {
  const router = useRouter();

  useEffect(() => {
    const userData = localStorage.getItem("user");
    if (!userData) {
      router.replace("/login");
    }
  }, [router]);
};

export default useAuthRedirect;
