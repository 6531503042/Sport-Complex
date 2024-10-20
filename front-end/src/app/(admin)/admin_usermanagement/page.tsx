"use client";

import React, { useEffect, useState } from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";
import AddIcon from "@mui/icons-material/Add";
import Modal from "../../components/popup/popup_addnew_user"; 
import { IUser } from "../../models/user"; 

const UserManagementPage = () => {
  const [users, setUsers] = useState<IUser[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false); 

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await fetch('/api/users');
        const data = await res.json();
        setUsers(data.data);
      } catch (error) {
        console.error('Error fetching users:', error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchUsers();
  }, []);

  const handleAddUser = async (userData: { name: string; email: string; password: string }) => {
    try {
      const res = await fetch('/api/users', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });
      const data = await res.json();
      
      if (res.ok) {
        setUsers((prevUsers) => [...prevUsers, data]); 
      } else {
        console.error('Failed to add user:', data);
      }
    } catch (error) {
      console.error('Error adding user:', error);
    }
  };

  return (
    <div className="w-[1920px] h-[945px] flex flex-row">
      <Sidebar activePage="admin_usermanagement"/>
      <div className="bg-white text-black w-full p-10 flex flex-col">
        <div className="inline-flex justify-between w-full items-end">
          <div className="text-lg font-medium">User Management</div>
          <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        </div>
        <br />
        <div className="border-b rounded-lg border-black"></div>
        <br />
        <div className="w-full flex justify-center">
          <div className="w-2/3 gap-10 inline-flex flex-row justify-center">
            <button className="inline-flex flex-row items-center text-white text-lg bg-red-600 p-2 justify-center rounded-md hover:bg-red-700 hover:drop-shadow-2xl drop-shadow-sm transition-all duration-200">
              <p>EDIT</p>
              <AddIcon />
            </button>
            <button
              className="inline-flex flex-row items-center text-white text-lg bg-green-500 p-2 justify-center rounded-md hover:bg-green-600 hover:drop-shadow-2xl drop-shadow-sm transition-all duration-200"
              onClick={() => setIsModalOpen(true)} 
            >
              <p>NEW</p>
              <AddIcon />
            </button>
          </div>
        </div>
        <br />
        <div className="w-full flex justify-center">
          <div className="w-5/6 flex rounded-md justify-between border p-5">
            <ul className="inline-flex flex-row w-full text-base font-medium items-center text-center">
              <li className="w-1/5">ID</li>
              <li className="w-1/5">Name</li>
              <li className="w-1/5">Email</li>
              <li className="w-1/5">Password</li>
              <li className="w-1/5">Role</li>
            </ul>
          </div>
        </div>
        <br />
        <div className="w-full flex justify-center">
        <div className="w-full flex justify-center">
  <div className="w-5/6 flex flex-col rounded-md border p-5">
    {isLoading ? (
      <p>Loading users...</p>
    ) : users.length === 0 ? (
      <p className="text-red-600">No users found.</p>
    ) : (
      <ul className="flex flex-col w-full text-base font-medium items-center text-center">
        {users.map((user) => (
          <li key={String(user._id)} className="flex flex-row w-full justify-between">
            <span className="w-1/5">{user._id}</span>
            <span className="w-1/5">{user.name}</span>
            <span className="w-1/5">{user.email}</span>
            <span className="w-1/5">********</span>
            <span className="w-1/5">{user.role}</span>
          </li>
        ))}
      </ul>
    )}
  </div>
</div>

        </div>
      </div>
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleAddUser}
      />
    </div>
  );
};

export default UserManagementPage;
function setError(arg0: string) {
  throw new Error("Function not implemented.");
}

