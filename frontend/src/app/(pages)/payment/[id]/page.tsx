"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Button } from "../../../components/ui/button";
import { Progress } from "../../../components/ui/progress";
import { Card } from "../../../components/ui/card";
import { useToast } from "../../../components/ui/use-toast";
import { CheckCircle, Download } from "lucide-react";
import saveAs from "file-saver";

const Payment = () => {
  const router = useRouter();
  const { id } = useParams(); // ดึง id จาก path
  const [timeLeft, setTimeLeft] = useState(600);
  const [progress, setProgress] = useState(100);
  const { toast } = useToast();
  const [showDialog, setShowDialog] = useState(false);
  const [paymentData, setPaymentData] = useState<any>(null);

  // Log id when it changes
  useEffect(() => {
    console.log("Payment ID from URL:", id); // This log should show the ID
    if (!id) {
      console.error("Payment ID is missing"); // Log when id is missing
      return;
    }

    const fetchPaymentData = async () => {
      try {
        const response = await fetch(
          `http://localhost:1327/payment_v1/payments/${id}`
        );
        if (!response.ok) {
          throw new Error(`Failed to fetch data: ${response.statusText}`); // Throw error if fetch fails
        }
        const data = await response.json();
        console.log("Payment data fetched:", data); // Log the fetched payment data
        setPaymentData(data);
      } catch (error) {
        console.error("Error fetching payment data:", error);
        toast({
          title: "Error",
          description: "Failed to load payment data. Please try again later.",
          duration: 3000,
        });
      }
    };

    fetchPaymentData();
  }, [id]); // Re-run effect if 'id' changes

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 0) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });

      setProgress((timeLeft / 600) * 100); // ใช้ timeLeft โดยตรง
    }, 1000);

    return () => clearInterval(timer);
  }, [timeLeft]);

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
  };

  const handlePaymentConfirmation = () => {
    toast({
      title: "Payment Verification",
      description: "Verifying your payment. Please wait...",
      duration: 3000,
    });
    setShowDialog(true); // แสดง dialog
  };

  const handleCloseDialog = () => {
    setShowDialog(false);
  };


  const handleGoHome = () => {
    router.push("/");
  };

  const handleSaveQRCode = () => {
    if (paymentData?.qr_code_url) {
      fetch(paymentData.qr_code_url)
        .then((response) => response.blob())
        .then((blob) => {
          saveAs(blob, `QRCode_${paymentData.booking_id}.png`);
          toast({
            title: "Downloading QR Code",
            description: "Your QR code is being downloaded...",
            duration: 3000,
          });
        })
        .catch((error) => {
          console.error("Error downloading QR code:", error);
        });
    }
  };

  if (!paymentData) {
    return <p>Loading...</p>;
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <Card className="w-full max-w-md p-6 space-y-6 bg-white shadow-lg rounded-xl">
        <div className="text-center space-y-2">
          <h1 className="text-2xl font-bold text-gray-900">QR Payment</h1>
          <div className="space-y-1">
            <h2 className="text-lg font-semibold text-gray-700">
              {paymentData.facility_name}
            </h2>
            <p className="text-black text-3xl font-bold text-primary">{`฿${paymentData.amount}.00`}</p>
          </div>
        </div>

        <div className="relative">
          <div className="aspect-square bg-gray-100 rounded-lg flex items-center justify-center border-2 border-gray-200">
            <img
              src={paymentData.qr_code_url}
              alt="QR Code"
              className="w-48 h-48"
            />
          </div>
        </div>

        <div className="space-y-2">
          <div className="flex justify-between text-sm text-gray-600">
            <span>Time remaining</span>
            <span>{formatTime(timeLeft)}</span>
          </div>
          <Progress value={progress} className="h-2" />
        </div>

        <div className="space-y-3">
          <Button
            onClick={handleSaveQRCode}
            className="w-full bg-blue-50 hover:bg-blue-100 text-blue-600 border-blue-200 flex items-center justify-center gap-2 h-12"
          >
            <Download className="w-5 h-5" />
            Save QR Code
          </Button>

          <Button
            onClick={handlePaymentConfirmation}
            className="w-full bg-green-600 hover:bg-green-700 text-white flex items-center justify-center gap-2 h-12"
          >
            <CheckCircle className="w-5 h-5" />
            I&apos;ve Already Paid
          </Button>

          <Button
            onClick={handleGoHome}
            className="w-full bg-gray-50 hover:bg-gray-100 text-gray-600 border-gray-200 flex items-center justify-center gap-2 h-12"
          >
            Go to Home
          </Button>
        </div>

        <div className="text-center text-sm text-gray-500">
          <p>Having trouble? Contact support</p>
        </div>
        {showDialog && (
          <div className="fixed bottom-4 right-4 bg-white p-4 shadow-lg rounded-lg">
            <p className="text-gray-800">Payment verification in progress...</p>
            <Button onClick={handleCloseDialog} className="mt-2 text-blue-600">
              Close
            </Button>
          </div>
        )}
      </Card>
    </div>
  );
};

export default Payment;
