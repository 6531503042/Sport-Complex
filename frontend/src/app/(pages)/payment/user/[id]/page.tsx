"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import NavBar from "@/app/components/navbar/navbar";
import Link from "next/link";
import router from "next/router";
import styles from "./PaymentUserPage.module.css"; 

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
        <h1 className="text-2xl font-bold mb-6">Payment History for {userName}</h1>
        {loading ? (
          <p>Loading...</p>
        ) : error ? (
          <p className={styles.textRed500}>{error}</p>
        ) : payments.length > 0 ? (
          <table className={styles.tableAuto}>
            <thead>
              <tr>
                <th>Payment ID</th>
                <th>Amount</th>
                <th>Currency</th>
                <th>Payment Method</th>
                <th>Facility Name</th>
                <th>Status</th>
                <th>Date</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {payments.map((payment) => (
                <tr key={payment.payment_id}>
                  <td>{payment._id}</td>
                  <td>{payment.amount}</td>
                  <td>{payment.currency}</td>
                  <td>{payment.payment_method}</td>
                  <td>{payment.facility_name}</td>
                  <td>{payment.status}</td>
                  <td>{new Date(payment.created_at).toLocaleString()}</td>
                  <td>
                    <Link href={`/payment/${payment._id}`} className="text-blue-500 hover:underline">
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
