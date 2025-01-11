"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";
import { useToast } from "../../../components/ui/use-toast";
import { CheckCircle, Download, Timer } from "lucide-react";
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

  const ProgressBar = ({ progress }: { progress: number }) => {
    return (
      <div className="w-full bg-gray-200 rounded-full h-4 flex items-center">
        <div
          className="bg-black h-4 rounded-full transition-all"
          style={{ width: `${progress}%` }}
        ></div>
      </div>
    );
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
            <p className="text-3xl font-bold text-primary">{`฿${paymentData.amount}.00`}</p>
          </div>
        </div>

        <div className="relative">
          <div className="aspect-square bg-white pb-10 pt-10 rounded-lg shadow-lg flex items-center justify-center border-2 border-gray-100">
            <div className="flex flex-col">
                <img 
                src="https://media.discordapp.net/attachments/1120602456473215047/1308337656614096937/image.png?ex=6782caa5&is=67817925&hm=371d94c2ed396b49d7d0fa12d2329fcf7db5d6ec19bcdaf1720745e360f8c1b1&=&format=webp&quality=lossless&width=2160&height=780" 
                alt="Logo" 
                className="w-auto h-auto mt-3 mb-3 ml-8 mr-8 pr-1"/>
              <img
                src={paymentData.qr_code_url}
                alt="QR Code"
                className="w-auto h-auto ml-16 mr-16"
              />
              <div className="text-2xl text-center mt-8">
                SCAN WITH ANY UPI APP
              </div>
            </div>
          </div>
          <div className="absolute -top-2 -right-2">
            <div className="bg-black text-white text-lg px-4 py-1 rounded-full flex items-center gap-2">
              <Timer className="w-4 h-4" />
              <span>{formatTime(timeLeft)}</span>
            </div>
          </div>
        </div>

        <div className="space-y-2">
          <div className="flex justify-between text-sm text-gray-600">
            <span>Time remaining</span>
            <span>{formatTime(timeLeft)}</span>
          </div>
          <ProgressBar progress={progress} />
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
