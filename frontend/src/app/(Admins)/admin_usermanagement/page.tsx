"use client";

import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Paper,
  IconButton,
  TextField,
  InputAdornment,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  Drawer,
  useTheme,
  useMediaQuery
} from '@mui/material';
import { 
  Search, 
  Edit, 
  Delete, 
  VerifiedUser,
  AdminPanelSettings,
  Person,
  Mail,
  Phone,
  CalendarToday,
  Close,
  Lock
} from '@mui/icons-material';
import { useAuth } from '@/app/context/AuthContext';
import { toast, Toaster } from 'react-hot-toast';
import Logo from "@/app/assets/Logo.png";
import { useRouter } from 'next/navigation';
import Sidebar from '@/app/components/sidebar_admin/sidebar';

interface IUser {
  _id: string;
  name: string;
  email: string;
  phone: string;
  created_at: string;
  user_roles: Array<{ role_code: number; role_title: string }>;
}

interface UserFormData {
  name: string;
  email: string;
  password?: string;
  role_code: number;
}

const UserManagementPage = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const [users, setUsers] = useState<IUser[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredUsers, setFilteredUsers] = useState<IUser[]>([]);
  const [openDialog, setOpenDialog] = useState(false);
  const [selectedUser, setSelectedUser] = useState<IUser | null>(null);
  const [formData, setFormData] = useState<UserFormData>({
    name: '',
    email: '',
    password: '',
    role_code: 0
  });
  const router = useRouter();
  const [mobileOpen, setMobileOpen] = useState(false);

  useEffect(() => {
    fetchUsers();
  }, []);

  useEffect(() => {
    const filtered = users.filter(user => 
      user.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      user.email.toLowerCase().includes(searchTerm.toLowerCase())
    );
    setFilteredUsers(filtered);
  }, [searchTerm, users]);

  const fetchUsers = async () => {
    setIsLoading(true);
    try {
      const accessToken = localStorage.getItem('access_token');
      
      if (!accessToken) {
        throw new Error('No access token found');
      }

      const response = await fetch('http://localhost:1325/user_v1/admin/users', {
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error('Failed to fetch users');
      }

      const data = await response.json();
      
      if (Array.isArray(data)) {
        setUsers(data);
        setFilteredUsers(data);
      } else {
        console.warn('Unexpected data format:', data);
        setUsers([]);
        setFilteredUsers([]);
      }
    } catch (error) {
      toast.error('Failed to fetch users');
      console.error('Error:', error);
      
      if (error.message.includes('401')) {
        router.push('/login');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const getRoleIcon = (roleCode: number) => {
    switch (roleCode) {
      case 1:
        return <AdminPanelSettings className="text-red-600" />;
      case 0:
        return <VerifiedUser className="text-green-600" />;
      default:
        return <Person className="text-blue-600" />;
    }
  };

  const getRoleLabel = (roleCode: number) => {
    switch (roleCode) {
      case 1:
        return { label: 'Admin', color: 'error' };
      case 0:
        return { label: 'User', color: 'success' };
      default:
        return { label: 'Guest', color: 'info' };
    }
  };

  const handleAddUser = async () => {
    try {
      const accessToken = localStorage.getItem('access_token');
      if (!accessToken) throw new Error('No access token found');

      const response = await fetch('http://localhost:1325/user_v1/admin/users', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });

      if (!response.ok) throw new Error('Failed to add user');

      const newUser = await response.json();
      setUsers(prev => [...prev, newUser]);
      toast.success('User added successfully');
      setOpenDialog(false);
      resetForm();
    } catch (error) {
      toast.error('Failed to add user');
      console.error('Error:', error);
    }
  };

  const handleEditUser = async () => {
    if (!selectedUser) return;

    try {
      const accessToken = localStorage.getItem('access_token');
      if (!accessToken) throw new Error('No access token found');

      const response = await fetch(`http://localhost:1325/user_v1/admin/users/${selectedUser._id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });

      if (!response.ok) throw new Error('Failed to update user');

      const updatedUser = await response.json();
      setUsers(prev => prev.map(user => 
        user._id === selectedUser._id ? updatedUser : user
      ));
      toast.success('User updated successfully');
      setOpenDialog(false);
      resetForm();
    } catch (error) {
      toast.error('Failed to update user');
      console.error('Error:', error);
    }
  };

  const handleDeleteUser = async (userId: string) => {
    if (!window.confirm('Are you sure you want to delete this user?')) return;

    try {
      const accessToken = localStorage.getItem('access_token');
      if (!accessToken) throw new Error('No access token found');

      const response = await fetch(`http://localhost:1325/user_v1/admin/users/${userId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
        }
      });

      if (!response.ok) throw new Error('Failed to delete user');

      setUsers(prev => prev.filter(user => user._id !== userId));
      toast.success('User deleted successfully');
    } catch (error) {
      toast.error('Failed to delete user');
      console.error('Error:', error);
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      email: '',
      password: '',
      role_code: 0
    });
    setSelectedUser(null);
  };

  const handleSubmit = () => {
    if (selectedUser) {
      handleEditUser();
    } else {
      handleAddUser();
    }
  };

  const handleAddNewUser = () => {
    resetForm();
    setOpenDialog(true);
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Desktop Sidebar */}
      {!isMobile && (
        <div className="sticky top-0 h-screen flex-shrink-0">
          <Sidebar activePage="users" />
        </div>
      )}

      {/* Mobile Drawer */}
      <Drawer
        variant="temporary"
        anchor="left"
        open={mobileOpen}
        onClose={() => setMobileOpen(false)}
        ModalProps={{
          keepMounted: true,
        }}
        sx={{
          display: { xs: 'block', md: 'none' },
          '& .MuiDrawer-paper': { 
            width: 280,
            boxSizing: 'border-box',
            backgroundColor: '#7f1d1d',
          },
        }}
      >
        <Sidebar activePage="users" />
      </Drawer>

      {/* Main Content */}
      <div className="flex-1 min-w-0">
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="p-4 md:p-6 lg:p-8"
        >
          <Toaster position="top-right" />
          
          {/* Header */}
          <div className="flex justify-between items-center mb-6">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">User Management</h1>
              <p className="text-gray-500 mt-1">Manage and monitor user accounts</p>
            </div>
            <img src={Logo.src} alt="Logo" className="w-7 h-min" />
          </div>

          {/* Search and Actions */}
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 mb-6">
            <TextField
              placeholder="Search users..."
              variant="outlined"
              size="small"
              fullWidth
              className="md:max-w-xs"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search className="text-gray-400" />
                  </InputAdornment>
                ),
              }}
            />
            
            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={handleAddNewUser}
              className="w-full md:w-auto px-4 py-2 bg-red-600 text-white rounded-lg
                hover:bg-red-700 transition-all duration-200 shadow-lg hover:shadow-xl
                flex items-center justify-center gap-2"
            >
              <Person className="w-5 h-5" />
              <span>Add New User</span>
            </motion.button>
          </div>

          {/* User Form Dialog */}
          <Dialog 
            open={openDialog} 
            onClose={() => {
              setOpenDialog(false);
              resetForm();
            }}
            maxWidth="sm"
            fullWidth
            PaperProps={{
              sx: {
                borderRadius: '16px',
                boxShadow: '0 25px 50px -12px rgba(0,0,0,0.25)',
              }
            }}
          >
            <motion.div 
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="p-8"
            >
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-2xl font-bold text-gray-900">
                  {selectedUser ? 'Edit User' : 'Add New User'}
                </h2>
                <motion.div
                  whileHover={{ rotate: 90 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <IconButton
                    onClick={() => {
                      setOpenDialog(false);
                      resetForm();
                    }}
                  >
                    <Close />
                  </IconButton>
                </motion.div>
              </div>

              <div className="space-y-6">
                <TextField
                  fullWidth
                  label="Name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <Person className="text-gray-400" />
                      </InputAdornment>
                    ),
                  }}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      borderRadius: '12px',
                      '&:hover': {
                        '& .MuiOutlinedInput-notchedOutline': {
                          borderColor: '#ef4444',
                        }
                      },
                      '&.Mui-focused': {
                        '& .MuiOutlinedInput-notchedOutline': {
                          borderColor: '#ef4444',
                          borderWidth: '2px',
                        }
                      }
                    }
                  }}
                />

                <TextField
                  fullWidth
                  label="Email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <Mail className="text-gray-400" />
                      </InputAdornment>
                    ),
                  }}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      borderRadius: '12px',
                      '&:hover': {
                        '& .MuiOutlinedInput-notchedOutline': {
                          borderColor: '#ef4444',
                        }
                      },
                      '&.Mui-focused': {
                        '& .MuiOutlinedInput-notchedOutline': {
                          borderColor: '#ef4444',
                          borderWidth: '2px',
                        }
                      }
                    }
                  }}
                />

                {!selectedUser && (
                  <TextField
                    fullWidth
                    label="Password"
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                    InputProps={{
                      startAdornment: (
                        <InputAdornment position="start">
                          <Lock className="text-gray-400" />
                        </InputAdornment>
                      ),
                    }}
                    sx={{
                      '& .MuiOutlinedInput-root': {
                        borderRadius: '12px',
                        '&:hover': {
                          '& .MuiOutlinedInput-notchedOutline': {
                            borderColor: '#ef4444',
                          }
                        },
                        '&.Mui-focused': {
                          '& .MuiOutlinedInput-notchedOutline': {
                            borderColor: '#ef4444',
                            borderWidth: '2px',
                          }
                        }
                      }
                    }}
                  />
                )}

                <FormControl fullWidth>
                  <InputLabel>Role</InputLabel>
                  <Select
                    value={formData.role_code}
                    onChange={(e) => setFormData({ ...formData, role_code: Number(e.target.value) })}
                    label="Role"
                    sx={{
                      borderRadius: '12px',
                      '& .MuiOutlinedInput-notchedOutline': {
                        borderColor: '#e5e7eb',
                      },
                      '&:hover': {
                        '& .MuiOutlinedInput-notchedOutline': {
                          borderColor: '#ef4444',
                        }
                      },
                      '&.Mui-focused': {
                        '& .MuiOutlinedInput-notchedOutline': {
                          borderColor: '#ef4444',
                          borderWidth: '2px',
                        }
                      }
                    }}
                  >
                    <MenuItem value={0}>
                      <div className="flex items-center gap-2">
                        <VerifiedUser className="text-green-600" />
                        <span>User</span>
                      </div>
                    </MenuItem>
                    <MenuItem value={1}>
                      <div className="flex items-center gap-2">
                        <AdminPanelSettings className="text-red-600" />
                        <span>Admin</span>
                      </div>
                    </MenuItem>
                  </Select>
                </FormControl>
              </div>

              <div className="flex justify-end gap-3 mt-8">
                <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                  <Button 
                    onClick={() => {
                      setOpenDialog(false);
                      resetForm();
                    }}
                    variant="outlined"
                    sx={{
                      borderRadius: '10px',
                      borderColor: '#e5e7eb',
                      color: '#6b7280',
                      '&:hover': {
                        borderColor: '#d1d5db',
                        backgroundColor: '#f9fafb',
                      }
                    }}
                  >
                    Cancel
                  </Button>
                </motion.div>
                
                <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                  <Button
                    variant="contained"
                    onClick={handleSubmit}
                    sx={{
                      borderRadius: '10px',
                      backgroundColor: '#ef4444',
                      '&:hover': {
                        backgroundColor: '#dc2626',
                      },
                      boxShadow: '0 4px 6px -1px rgba(239, 68, 68, 0.1), 0 2px 4px -1px rgba(239, 68, 68, 0.06)',
                    }}
                  >
                    {selectedUser ? 'Update' : 'Add'} User
                  </Button>
                </motion.div>
              </div>
            </motion.div>
          </Dialog>

          {/* Users Table */}
          <TableContainer component={Paper} elevation={0} className="border rounded-xl">
            <Table>
              <TableHead className="bg-gray-50">
                <TableRow>
                  <TableCell>User</TableCell>
                  <TableCell>Contact</TableCell>
                  <TableCell>Role</TableCell>
                  <TableCell>Joined Date</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center" className="py-8">
                      <motion.div
                        animate={{ rotate: 360 }}
                        transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                        className="w-6 h-6 border-2 border-red-600 border-t-transparent rounded-full mx-auto"
                      />
                    </TableCell>
                  </TableRow>
                ) : filteredUsers.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center" className="py-8">
                      <Person className="w-12 h-12 text-gray-300 mx-auto mb-2" />
                      <p className="text-gray-500">No users found</p>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredUsers.map((user, index) => (
                    <motion.tr
                      key={user._id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="hover:bg-gray-50"
                    >
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center">
                            <Person className="text-gray-500" />
                          </div>
                          <div>
                            <p className="font-medium text-gray-900">{user.name}</p>
                            <p className="text-sm text-gray-500">{user.email}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-col gap-1">
                          <div className="flex items-center gap-2">
                            <Mail className="w-4 h-4 text-gray-400" />
                            <span className="text-sm">{user.email}</span>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        {user.user_roles?.map((role, idx) => (
                          <Chip
                            key={idx}
                            icon={getRoleIcon(role.role_code)}
                            label={role.role_title}
                            color={getRoleLabel(role.role_code).color as any}
                            size="small"
                            className="font-medium"
                          />
                        ))}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <CalendarToday className="w-4 h-4 text-gray-400" />
                          <span className="text-sm">
                            {new Date(user.created_at).toLocaleDateString()}
                          </span>
                        </div>
                      </TableCell>
                      <TableCell align="right">
                        <IconButton 
                          size="small"
                          className="text-blue-600 hover:bg-blue-50"
                          onClick={() => {
                            setSelectedUser(user);
                            setFormData({
                              name: user.name,
                              email: user.email,
                              role_code: user.user_roles[0]?.role_code || 0
                            });
                            setOpenDialog(true);
                          }}
                        >
                          <Edit />
                        </IconButton>
                        <IconButton 
                          size="small"
                          className="text-red-600 hover:bg-red-50"
                          onClick={() => handleDeleteUser(user._id)}
                        >
                          <Delete />
                        </IconButton>
                      </TableCell>
                    </motion.tr>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </motion.div>
      </div>
    </div>
  );
};

export default UserManagementPage;

