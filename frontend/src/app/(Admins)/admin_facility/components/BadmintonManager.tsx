"use client";

import React, { useState, useEffect } from "react";
import { TextField, Button, Dialog, Tab, Tabs, Box } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import SearchIcon from "@mui/icons-material/Search";
import { toast } from 'react-hot-toast';
import { muiStyles } from '../styles/shared';
import { motion } from 'framer-motion';

interface Court {
  _id: string;
  court_number: number;
  status: number;
}

interface Slot {
  _id: string;
  start_time: string;
  end_time: string;
  court_id: string;
  status: number;
  max_bookings: number;
  current_bookings: number;
}

const BadmintonManager = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [courts, setCourts] = useState<Court[]>([]);
  const [slots, setSlots] = useState<Slot[]>([]);
  const [openDialog, setOpenDialog] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedItem, setSelectedItem] = useState<Court | Slot | null>(null);
  const [courtFormData, setCourtFormData] = useState({
    court_number: 1,
    status: 0,
  });
  const [slotFormData, setSlotFormData] = useState({
    start_time: "",
    end_time: "",
    court_id: "",
    status: 0,
    max_bookings: 1,
    current_bookings: 0,
  });

  const fetchCourts = async () => {
    try {
      const url = "http://localhost:1335/facility_v1/badminton_v1/courts";
      console.log('Fetching courts from:', url);
      
      const response = await fetch(url);
      console.log('Courts response:', response);
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
      }
      
      const courts = await response.json();
      console.log('Fetched courts:', courts);
      
      if (Array.isArray(courts)) {
        setCourts(courts);
      } else {
        console.warn('Courts response is not an array:', courts);
        setCourts([]);
      }
    } catch (error: any) {
      console.error('Courts fetch error:', error);
      toast.error(`Error fetching courts: ${error.message}`);
      setCourts([]);
    }
  };

  const fetchSlots = async () => {
    try {
      const url = "http://localhost:1335/facility_v1/badminton_v1/slots";
      console.log('Fetching slots from:', url);
      
      const response = await fetch(url);
      console.log('Slots response:', response);
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
      }
      
      const slots = await response.json();
      console.log('Fetched slots:', slots);
      
      if (Array.isArray(slots)) {
        setSlots(slots);
      } else {
        console.warn('Slots response is not an array:', slots);
        setSlots([]);
      }
    } catch (error: any) {
      console.error('Slots fetch error:', error);
      toast.error(`Error fetching slots: ${error.message}`);
      setSlots([]);
    }
  };

  useEffect(() => {
    console.log('BadmintonManager mounted');
    fetchCourts();
    fetchSlots();
  }, []);

  const handleSubmit = async () => {
    try {
      const isEditMode = selectedItem !== null;
      const isCourt = activeTab === 0;
      
      const url = isCourt
        ? "http://localhost:1335/facility_v1/badminton_v1/court"
        : "http://localhost:1335/facility_v1/badminton_v1/slot";

      const method = isEditMode ? "PUT" : "POST";
      const body = isCourt ? courtFormData : slotFormData;

      const response = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });

      if (response.ok) {
        setOpenDialog(false);
        if (isCourt) {
          fetchCourts();
          resetCourtForm();
        } else {
          fetchSlots();
          resetSlotForm();
        }
      }
    } catch (error) {
      console.error("Error saving:", error);
    }
  };

  const handleDelete = async (id: string, isCourt: boolean) => {
    if (window.confirm("Are you sure you want to delete this item?")) {
      try {
        const url = isCourt
          ? `http://localhost:1335/facility_v1/badminton_v1/court/${id}`
          : `http://localhost:1335/facility_v1/badminton_v1/slot/${id}`;

        await fetch(url, {
          method: "DELETE",
        });

        if (isCourt) {
          fetchCourts();
        } else {
          fetchSlots();
        }
      } catch (error) {
        console.error("Error deleting:", error);
      }
    }
  };

  const resetCourtForm = () => {
    setCourtFormData({
      court_number: 1,
      status: 0,
    });
    setSelectedItem(null);
  };

  const resetSlotForm = () => {
    setSlotFormData({
      start_time: "",
      end_time: "",
      court_id: "",
      status: 0,
      max_bookings: 1,
      current_bookings: 0,
    });
    setSelectedItem(null);
  };

  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="space-y-4"
    >
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs 
          value={activeTab} 
          onChange={(_, newValue) => setActiveTab(newValue)}
          sx={{
            '& .MuiTab-root': {
              minWidth: 120,
              fontWeight: 600,
              color: '#64748b',
              '&.Mui-selected': {
                color: '#7f1d1d',
              },
            },
            '& .MuiTabs-indicator': {
              backgroundColor: '#7f1d1d',
            },
          }}
        >
          <Tab label="Courts" />
          <Tab label="Slots" />
        </Tabs>
      </Box>

      <div className="flex justify-between items-center">
        <div className="relative w-64">
          <SearchIcon className="absolute left-3 top-2.5 text-gray-400" />
          <input
            type="text"
            placeholder={activeTab === 0 ? "Search courts..." : "Search slots..."}
            className="pl-10 pr-4 py-2 w-full border rounded-lg"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => {
            if (activeTab === 0) resetCourtForm();
            else resetSlotForm();
            setOpenDialog(true);
          }}
          sx={{
            backgroundColor: '#7f1d1d !important',
            color: '#ffffff !important',
            '& .MuiButton-startIcon': {
              color: '#ffffff !important',
            },
            '& .MuiSvgIcon-root': {
              color: '#ffffff !important',
              fontSize: '1.5rem',
            },
            '&:hover': {
              backgroundColor: '#991b1b !important',
              transform: 'translateY(-2px)',
              boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
            },
            fontWeight: 600,
            padding: '8px 16px',
            textTransform: 'none',
            fontSize: '1rem',
            borderRadius: '8px',
            transition: 'all 0.2s ease-in-out',
          }}
        >
          Add {activeTab === 0 ? "Court" : "Slot"}
        </Button>
      </div>

      {activeTab === 0 ? (
        // Courts Table
        <div className="bg-white rounded-lg shadow">
          <table className="min-w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Court Number
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {courts.map((court) => (
                <tr key={court._id}>
                  <td className="px-6 py-4 whitespace-nowrap">{court.court_number}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {court.status === 0 ? "Available" : "Occupied"}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap space-x-2">
                    <Button
                      size="small"
                      startIcon={<EditIcon />}
                      onClick={() => {
                        setSelectedItem(court);
                        setCourtFormData({
                          court_number: court.court_number,
                          status: court.status,
                        });
                        setOpenDialog(true);
                      }}
                    >
                      Edit
                    </Button>
                    <Button
                      size="small"
                      startIcon={<DeleteIcon />}
                      color="error"
                      onClick={() => handleDelete(court._id, true)}
                    >
                      Delete
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        // Slots Table
        <div className="bg-white rounded-lg shadow">
          <table className="min-w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Start Time
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  End Time
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Court
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {slots.map((slot) => (
                <tr key={slot._id}>
                  <td className="px-6 py-4 whitespace-nowrap">{slot.start_time}</td>
                  <td className="px-6 py-4 whitespace-nowrap">{slot.end_time}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {courts.find(c => c._id === slot.court_id)?.court_number || "N/A"}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {slot.status === 0 ? "Available" : "Booked"}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap space-x-2">
                    <Button
                      size="small"
                      startIcon={<EditIcon />}
                      onClick={() => {
                        setSelectedItem(slot);
                        setSlotFormData({
                          ...slot,
                        });
                        setOpenDialog(true);
                      }}
                    >
                      Edit
                    </Button>
                    <Button
                      size="small"
                      startIcon={<DeleteIcon />}
                      color="error"
                      onClick={() => handleDelete(slot._id, false)}
                    >
                      Delete
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
        <div className="p-6 space-y-4">
          <h2 className="text-xl font-semibold">
            {selectedItem ? "Edit" : "Add New"} {activeTab === 0 ? "Court" : "Slot"}
          </h2>
          <div className="space-y-4">
            {activeTab === 0 ? (
              // Court Form
              <>
                <TextField
                  fullWidth
                  label="Court Number"
                  type="number"
                  value={courtFormData.court_number}
                  onChange={(e) => setCourtFormData({ ...courtFormData, court_number: parseInt(e.target.value) })}
                />
                <TextField
                  fullWidth
                  select
                  label="Status"
                  value={courtFormData.status}
                  onChange={(e) => setCourtFormData({ ...courtFormData, status: parseInt(e.target.value) })}
                  SelectProps={{ native: true }}
                >
                  <option value={0}>Available</option>
                  <option value={1}>Occupied</option>
                </TextField>
              </>
            ) : (
              // Slot Form
              <>
                <TextField
                  fullWidth
                  label="Start Time"
                  type="time"
                  value={slotFormData.start_time}
                  onChange={(e) => setSlotFormData({ ...slotFormData, start_time: e.target.value })}
                />
                <TextField
                  fullWidth
                  label="End Time"
                  type="time"
                  value={slotFormData.end_time}
                  onChange={(e) => setSlotFormData({ ...slotFormData, end_time: e.target.value })}
                />
                <TextField
                  fullWidth
                  select
                  label="Court"
                  value={slotFormData.court_id}
                  onChange={(e) => setSlotFormData({ ...slotFormData, court_id: e.target.value })}
                  SelectProps={{ native: true }}
                >
                  <option value="">Select Court</option>
                  {courts.map((court) => (
                    <option key={court._id} value={court._id}>
                      Court {court.court_number}
                    </option>
                  ))}
                </TextField>
                <TextField
                  fullWidth
                  label="Max Bookings"
                  type="number"
                  value={slotFormData.max_bookings}
                  onChange={(e) => setSlotFormData({ ...slotFormData, max_bookings: parseInt(e.target.value) })}
                />
              </>
            )}
          </div>
          <div className="flex justify-end space-x-2">
            <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
            <Button
              variant="contained"
              onClick={handleSubmit}
              className="bg-red-900 hover:bg-red-800"
            >
              {selectedItem ? "Update" : "Create"}
            </Button>
          </div>
        </div>
      </Dialog>
    </motion.div>
  );
};

export default BadmintonManager; 