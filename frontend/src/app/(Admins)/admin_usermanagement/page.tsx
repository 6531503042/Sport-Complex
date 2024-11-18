"use client";

import React, { useEffect, useState } from "react";
import Sidebar from "../../components/sidebar_admin/sidebar";
import Logo from "../../assets/Logo.png";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import { IUser } from "../../models/user";
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useAuth } from "../../context/AuthContext";
import { CircularProgress, Fade } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import { useRouter } from 'next/navigation';
import { motion } from "framer-motion";
import { Tooltip } from "@mui/material";
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import { useTheme } from '@mui/material/styles';

interface UserData {
  name: string;
  email: string;
  password: string;
  user_roles: { role_title: string; role_code: number }[];
}

interface UserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (userData: UserData) => void;
  user?: IUser;
  mode: 'add' | 'edit';
}

const UserModal: React.FC<UserModalProps> = ({ isOpen, onClose, onSubmit, user, mode }) => {
  const [formData, setFormData] = useState({
    name: user?.name || '',
    email: user?.email || '',
    password: '',
    user_roles: user?.user_roles || [{ role_title: 'user', role_code: 0 }]
  });

  useEffect(() => {
    if (user) {
      setFormData({
        name: user.name,
        email: user.email,
        password: '',
        user_roles: user.user_roles
      });
    }
  }, [user]);

  const handleRoleChange = (roleTitle: string) => {
    const roleCode = roleTitle === 'admin' ? 1 : 0;
    setFormData({
      ...formData,
      user_roles: [{ role_title: roleTitle, role_code: roleCode }]
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <motion.div 
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.9, opacity: 0 }}
        className="bg-white rounded-lg p-8 w-96 shadow-xl"
      >
        <h2 className="text-2xl font-bold mb-6">{mode === 'add' ? 'Add New User' : 'Edit User'}</h2>
        <form onSubmit={(e) => {
          e.preventDefault();
          onSubmit(formData);
        }}>
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2">Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2">Email</label>
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              className="w-full p-2 border rounded"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2">Password</label>
            <input
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              className="w-full p-2 border rounded"
              required={mode === 'add'}
            />
          </div>
          <div className="mb-6">
            <label className="block text-gray-700 text-sm font-bold mb-2">Role</label>
            <select
              value={formData.user_roles[0]?.role_title || 'user'}
              onChange={(e) => handleRoleChange(e.target.value)}
              className="w-full p-2 border rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div className="flex justify-end gap-4">
            <motion.button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-600 hover:text-gray-800"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              Cancel
            </motion.button>
            <motion.button
              type="submit"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              {mode === 'add' ? 'Add User' : 'Save Changes'}
            </motion.button>
          </div>
        </form>
      </motion.div>
    </div>
  );
};

const API_BASE_URL = 'http://localhost:1325/user_v1';

const UserManagementPage = () => {
  const { user: authUser, logout } = useAuth();
  const [users, setUsers] = useState<IUser[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isInitializing, setIsInitializing] = useState(true);
  const [modalMode, setModalMode] = useState<'add' | 'edit'>('add');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<IUser | undefined>(undefined);
  const [searchTerm, setSearchTerm] = useState('');
  const [isSearchFocused, setIsSearchFocused] = useState(false);
  const router = useRouter();

  const theme = useTheme();

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

<<<<<<< HEAD
  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1
    }
  };
=======
  return (
    <div className="w-screen h-screen flex relative">
      <Sidebar activePage="admin_usermanagement"/>
      <div className="bg-white text-black w-full p-10 flex flex-col overflow-y-auto">
        <div className="inline-flex justify-between w-full items-end">
          <div className="text-lg font-medium">User Management</div>
          <img src={Logo.src} alt="Logo" className="w-7 h-min" />
        </div>
        <br />
        <div className="bg-zinc-500 h-[1px] rounded-lg"></div>
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
>>>>>>> bcfa478ac8b6508d73b9d33cbaa72dc20d79e3e0

  const fetchUsers = async () => {
    try {
      setIsLoading(true);
      const accessToken = localStorage.getItem('access_token');
      if (!accessToken) {
        toast.error('Authentication required');
        return;
      }

      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        }
      });

      if (!response.ok) {
        if (response.status === 401) {
          toast.error('Session expired. Please login again.');
          logout();
          return;
        }
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.data) {
        setUsers(data.data);
      } else if (Array.isArray(data)) {
        setUsers(data);
      } else {
        setUsers([]);
      }
    } catch (error) {
      console.error('Error fetching users:', error);
      toast.error('Failed to fetch users');
      setUsers([]);
    } finally {
      setIsLoading(false);
      setIsInitializing(false);
    }
  };

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('access_token');
      const storedUser = localStorage.getItem('user');

      if (!token || !storedUser) {
        router.push('/login');
        return;
      }

      try {
        // Parse stored user data
        const userData = JSON.parse(storedUser);
        if (!userData) {
          router.push('/login');
          return;
        }

        // If we have both token and user data, initialize the page
        await fetchUsers();
      } catch (error) {
        console.error('Error initializing page:', error);
      } finally {
        setIsInitializing(false);
      }
    };

    checkAuth();
  }, []);

  const handleAddUser = async (userData: UserData) => {
    try {
      const accessToken = localStorage.getItem('access_token');
      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify(userData),
      });

      if (response.ok) {
        toast.success('User added successfully');
        fetchUsers();
        setIsModalOpen(false);
      } else {
        const error = await response.json();
        toast.error(error.message || 'Failed to add user');
      }
    } catch (error) {
      console.error('Error adding user:', error);
      toast.error('Error adding user');
    }
  };

  const handleEditUser = async (userData: UserData) => {
    if (!selectedUser) return;
    
    try {
      const accessToken = localStorage.getItem('access_token');
      const response = await fetch(`${API_BASE_URL}/admin/users/${selectedUser._id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify(userData),
      });

      if (response.ok) {
        toast.success('User updated successfully');
        fetchUsers();
        setIsModalOpen(false);
      } else {
        const error = await response.json();
        toast.error(error.message || 'Failed to update user');
      }
    } catch (error) {
      console.error('Error updating user:', error);
      toast.error('Error updating user');
    }
  };

  const handleDeleteUser = async (userId: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        const accessToken = localStorage.getItem('access_token');
        const response = await fetch(`${API_BASE_URL}/admin/users/${userId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${accessToken}`
          }
        });

        if (response.ok) {
          toast.success('User deleted successfully');
          fetchUsers();
        } else {
          const error = await response.json();
          toast.error(error.message || 'Failed to delete user');
        }
      } catch (error) {
        console.error('Error deleting user:', error);
        toast.error('Error deleting user');
      }
    }
  };

  const filteredUsers = users?.filter(user => 
    user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  // Show loading state while initializing
  if (isInitializing) {
    return (
      <div className="w-screen h-screen flex items-center justify-center">
        <CircularProgress />
      </div>
    );
  }

  // Only show the main content if we have the auth token
  const token = localStorage.getItem('access_token');
  if (!token) {
    return null; // Return null instead of redirecting
  }

  return (
    <div className="w-screen h-screen flex flex-row bg-gradient-to-br from-gray-50 to-gray-100">
      <Sidebar activePage="admin_usermanagement"/>
      <div className="flex-1 p-8 overflow-hidden">
        <motion.div 
          className="max-w-7xl mx-auto"
          initial="hidden"
          animate="visible"
          variants={containerVariants}
        >
          {/* Header Section */}
          <motion.div 
            className="flex justify-between items-center mb-8"
            variants={itemVariants}
          >
            <div>
              <h1 className="text-3xl font-bold text-gray-900 tracking-tight">
                User Management
              </h1>
              <p className="mt-2 text-sm text-gray-600">
                Manage and monitor user accounts and permissions
              </p>
            </div>
            <motion.img 
              src={Logo.src} 
              alt="Logo" 
              className="w-12 h-auto"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
            />
          </motion.div>

          {/* Search and Add User Section */}
          <motion.div 
            className="mb-6 flex justify-between items-center"
            variants={itemVariants}
          >
            <div className="relative">
              <div className={`
                flex items-center transition-all duration-300 bg-white rounded-lg
                ${isSearchFocused 
                  ? 'ring-2 ring-blue-500 shadow-lg transform scale-105' 
                  : 'shadow hover:shadow-md transform hover:scale-102'
                }
              `}>
                <SearchIcon className="text-gray-400 ml-3" />
                <input
                  type="text"
                  placeholder="Search users..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  onFocus={() => setIsSearchFocused(true)}
                  onBlur={() => setIsSearchFocused(false)}
                  className="p-3 pl-2 w-64 border-none rounded-lg focus:outline-none"
                />
              </div>
            </div>
            <Tooltip title="Add new user" arrow>
              <motion.button
                onClick={() => {
                  setModalMode('add');
                  setSelectedUser(undefined);
                  setIsModalOpen(true);
                }}
                className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 
                          flex items-center gap-2 shadow-md"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <PersonAddIcon /> Add New User
              </motion.button>
            </Tooltip>
          </motion.div>

          {/* Users Table */}
          <motion.div 
            className="bg-white rounded-xl shadow-xl overflow-hidden"
            variants={itemVariants}
          >
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    {["Name", "Email", "Role", "Actions"].map((header) => (
                      <th
                        key={header}
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                      >
                        {header}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {isLoading ? (
                    <tr>
                      <td colSpan={4} className="px-6 py-4">
                        <div className="flex justify-center items-center space-x-3">
                          <CircularProgress size={24} />
                          <span className="text-gray-500 animate-pulse">Loading users...</span>
                        </div>
                      </td>
                    </tr>
                  ) : filteredUsers.length === 0 ? (
                    <tr>
                      <td colSpan={4} className="px-6 py-8">
                        <div className="text-center text-gray-500">
                          <motion.div
                            initial={{ scale: 0 }}
                            animate={{ scale: 1 }}
                            transition={{ type: "spring", stiffness: 200, damping: 20 }}
                          >
                            <p className="text-lg font-semibold">No users found</p>
                            <p className="text-sm mt-1">Try adjusting your search or add a new user</p>
                          </motion.div>
                        </div>
                      </td>
                    </tr>
                  ) : (
                    filteredUsers.map((user, index) => (
                      <motion.tr 
                        key={user._id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className="hover:bg-gray-50 transition-all duration-150"
                      >
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-semibold mr-3">
                              {user.name.charAt(0).toUpperCase()}
                            </div>
                            <div className="text-sm font-medium text-gray-900">{user.name}</div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {user.email}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full 
                            ${user.user_roles?.[0]?.role_code === 1 
                              ? 'bg-purple-100 text-purple-800' 
                              : 'bg-green-100 text-green-800'}`}
                          >
                            {user.user_roles?.[0]?.role_title || 'User'}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          <div className="flex space-x-3">
                            <Tooltip title="Edit user" arrow>
                              <motion.button
                                onClick={() => {
                                  setSelectedUser(user);
                                  setModalMode('edit');
                                  setIsModalOpen(true);
                                }}
                                className="text-indigo-600 hover:text-indigo-900"
                                whileHover={{ scale: 1.2 }}
                                whileTap={{ scale: 0.9 }}
                              >
                                <EditIcon />
                              </motion.button>
                            </Tooltip>
                            <Tooltip title="Delete user" arrow>
                              <motion.button
                                onClick={() => handleDeleteUser(user._id)}
                                className="text-red-600 hover:text-red-900"
                                whileHover={{ scale: 1.2 }}
                                whileTap={{ scale: 0.9 }}
                              >
                                <DeleteIcon />
                              </motion.button>
                            </Tooltip>
                          </div>
                        </td>
                      </motion.tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </motion.div>

          {/* Add the Modal here */}
          {isModalOpen && (
            <UserModal
              isOpen={isModalOpen}
              onClose={() => setIsModalOpen(false)}
              onSubmit={modalMode === 'add' ? handleAddUser : handleEditUser}
              user={selectedUser}
              mode={modalMode}
            />
          )}
        </motion.div>
      </div>
      <ToastContainer 
        position="top-right" 
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="colored"
        limit={3}
      />
    </div>
  );
};

export default UserManagementPage;

