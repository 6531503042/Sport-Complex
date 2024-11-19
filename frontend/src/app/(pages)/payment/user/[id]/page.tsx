"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import NavBar from "@/app/components/navbar/navbar";
import Link from "next/link";
import { motion } from "framer-motion";
import { 
  Receipt, 
  CreditCard, 
  Calendar, 
  Building, 
  CheckCircle,
  XCircle,
  Loader2
} from "lucide-react";
import styles from "./PaymentUserPage.module.css";

interface Payment {
  _id: string;
  payment_id: string;
  amount: number;
  currency: string;
  payment_method: string;
  facility_name: string;
  status: string;
  created_at: string;
}

const PaymentUserPage: React.FC = () => {
  const [userName, setUserName] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [payments, setPayments] = useState<Payment[]>([]);
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
        setPayments(Array.isArray(data) ? data : []);
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
      const id = user.id?.replace("user:", "") || "";
      setUserId(id);
    }
  }, [id]);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'paid':
        return 'bg-green-100 text-green-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'paid':
        return <CheckCircle className="w-5 h-5 text-green-600" />;
      case 'pending':
        return <Loader2 className="w-5 h-5 text-yellow-600 animate-spin" />;
      case 'failed':
        return <XCircle className="w-5 h-5 text-red-600" />;
      default:
        return null;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar activePage="payment" />
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="container mx-auto p-6 max-w-7xl"
      >
        <div className="bg-white rounded-2xl shadow-lg p-6 md:p-8">
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex items-center justify-between mb-8"
          >
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Payment History</h1>
              <p className="text-gray-600 mt-1">Track all your facility bookings and payments</p>
            </div>
            <div className="text-right">
              <p className="text-lg font-medium text-gray-900">{userName}</p>
              <p className="text-sm text-gray-500">User ID: {userId}</p>
            </div>
          </motion.div>

          {loading ? (
            <div className="flex items-center justify-center h-64">
              <Loader2 className="w-8 h-8 text-blue-500 animate-spin" />
            </div>
          ) : error ? (
            <motion.div 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="bg-red-50 text-red-600 p-4 rounded-xl text-center"
            >
              {error}
            </motion.div>
          ) : payments && payments.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="bg-gray-50">
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Payment ID</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Amount</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Method</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Facility</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Status</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Date</th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600">Action</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {payments.map((payment, index) => (
                    <motion.tr 
                      key={payment._id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="hover:bg-gray-50 transition-colors"
                    >
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <Receipt className="w-5 h-5 text-gray-400" />
                          <span className="text-sm text-gray-900">{payment._id.slice(0, 8)}...</span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <CreditCard className="w-5 h-5 text-gray-400" />
                          <span className="text-sm font-medium text-gray-900">
                            ฿{payment.amount.toFixed(2)}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <span className="text-sm text-gray-600">{payment.payment_method}</span>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <Building className="w-5 h-5 text-gray-400" />
                          <span className="text-sm text-gray-600">{payment.facility_name}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          {getStatusIcon(payment.status)}
                          <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(payment.status)}`}>
                            {payment.status}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <Calendar className="w-5 h-5 text-gray-400" />
                          <span className="text-sm text-gray-600">
                            {new Date(payment.created_at).toLocaleDateString()}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <Link 
                          href={`/payment/${payment._id}`}
                          className="inline-flex items-center px-3 py-1 rounded-lg text-sm font-medium
                            text-blue-600 hover:text-blue-700 hover:bg-blue-50
                            transition-colors duration-200"
                        >
                          View Details
                        </Link>
                      </td>
                    </motion.tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <motion.div 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="text-center py-12"
            >
              <Receipt className="w-16 h-16 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900">No Payments Found</h3>
              <p className="text-gray-500 mt-1">You haven't made any payments yet.</p>
            </motion.div>
          )}
        </div>
      </motion.div>
    </div>
  );
};

export default PaymentUserPage;
