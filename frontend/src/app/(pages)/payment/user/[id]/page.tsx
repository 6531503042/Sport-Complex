"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import NavBar from "@/app/components/navbar/navbar"; // อัปเดต path ถ้าจำเป็น
import Link from "next/link"; // นำเข้า Link
import router from "next/router";

const PaymentUserPage: React.FC = () => {
  const [userName, setUserName] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [payments, setPayments] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { id } = useParams();

  useEffect(() => {
    const fetchPayments = async () => {
      try {
        const res = await fetch(`http://localhost:1327/payment_v1/payments/user/${id}`);
        if (!res.ok) {
          throw new Error("Failed to fetch payment data.");
        }
        const data = await res.json();
        setPayments(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
      
    };

    fetchPayments();
    const userData = localStorage.getItem("user");
    if (userData) {
      const user = JSON.parse(userData);
      setUserName(user.name);
      setUserId(user.id);
      const id = user.id?.replace("user:", "") || ""; 
      setUserId(id);  
    } else {
      router.replace("/login");
    }
  }, [id]);

  return (
    <div>
      <NavBar activePage="payment" />
      <div className="container mx-auto p-6">
        <h1 className="text-2xl font-bold mb-6">Payment History for {userName} </h1>
        {loading ? (
          <p>Loading...</p>
        ) : error ? (
          <p className="text-red-500">{error}</p>
        ) : payments.length > 0 ? (
          <table className="table-auto w-full border-collapse border border-gray-300">
            <thead>
              <tr className="bg-gray-100">
                <th className="border border-gray-300 px-4 py-2">Payment ID</th>
                <th className="border border-gray-300 px-4 py-2">Amount</th>
                <th className="border border-gray-300 px-4 py-2">Currency</th>
                <th className="border border-gray-300 px-4 py-2">Payment Method</th>
                <th className="border border-gray-300 px-4 py-2">Facility Name</th>
                <th className="border border-gray-300 px-4 py-2">Status</th>
                <th className="border border-gray-300 px-4 py-2">Date</th>
                <th className="border border-gray-300 px-4 py-2">Action</th>
              </tr>
            </thead>
            <tbody>
              {payments.map((payment) => (
                <tr key={payment.payment_id}>
                  <td className="border border-gray-300 px-4 py-2">{payment._id}</td>
                  <td className="border border-gray-300 px-4 py-2">{payment.amount}</td>
                  <td className="border border-gray-300 px-4 py-2">{payment.currency}</td>
                  <td className="border border-gray-300 px-4 py-2">{payment.payment_method}</td>
                  <td className="border border-gray-300 px-4 py-2">{payment.facility_name}</td>
                  <td className="border border-gray-300 px-4 py-2">{payment.status}</td>
                  <td className="border border-gray-300 px-4 py-2">
                    {new Date(payment.created_at).toLocaleString()}
                  </td>
                  <td className="border border-gray-300 px-4 py-2">
                    <Link
                      href={`/payment/${payment._id}`} // ลิงก์ไปยังหน้ารายละเอียดของ payment โดยใช้ _id
                      className="text-blue-500 hover:underline"
                    >
                      View Payment
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <p>No payment records found.</p>
        )}
      </div>
    </div>
  );
};

export default PaymentUserPage;
