"use client";

import React, { useState, useEffect } from "react";
import { 
  TextField, 
  Button, 
  Dialog,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  InputAdornment,
  CircularProgress
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import SearchIcon from "@mui/icons-material/Search";
import { toast } from 'react-hot-toast';

interface Slot {
  _id: string;
  start_time: string;
  end_time: string;
  status: number;
  max_bookings: number;
  current_bookings: number;
}

interface FacilitySlotManagerProps {
  facilityName: string;
}

const FacilitySlotManager: React.FC<FacilitySlotManagerProps> = ({ facilityName }) => {
  const [slots, setSlots] = useState<Slot[]>([]);
  const [openDialog, setOpenDialog] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    start_time: "",
    end_time: "",
    max_bookings: 1,
    current_bookings: 0,
    facility_type: facilityName
  });

  useEffect(() => {
    fetchSlots();
  }, [facilityName]);

  const fetchSlots = async () => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots`);
      const data = await response.json();
      setSlots(data.data || []);
    } catch (error) {
      toast.error(`Error fetching ${facilityName} slots`);
      console.error("Error fetching slots:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    try {
      const url = selectedSlot 
        ? `http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots/${selectedSlot._id}`
        : `http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots`;

      const method = selectedSlot ? "PUT" : "POST";

      const response = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        toast.success(selectedSlot ? "Slot updated successfully" : "Slot created successfully");
        setOpenDialog(false);
        fetchSlots();
        resetForm();
      }
    } catch (error) {
      toast.error("Error saving slot");
      console.error("Error saving slot:", error);
    }
  };

  const handleDelete = async (slotId: string) => {
    if (window.confirm("Are you sure you want to delete this slot?")) {
      try {
        await fetch(`http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots/${slotId}`, {
          method: "DELETE",
        });
        toast.success("Slot deleted successfully");
        fetchSlots();
      } catch (error) {
        toast.error("Error deleting slot");
        console.error("Error deleting slot:", error);
      }
    }
  };

  const resetForm = () => {
    setFormData({
      start_time: "",
      end_time: "",
      max_bookings: 1,
      current_bookings: 0,
      facility_type: facilityName
    });
    setSelectedSlot(null);
  };

  const filteredSlots = slots.filter(slot => 
    slot.start_time.toLowerCase().includes(searchTerm.toLowerCase()) ||
    slot.end_time.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-6">
        <TextField
          placeholder="Search slots..."
          variant="outlined"
          size="small"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon className="text-gray-400" />
              </InputAdornment>
            ),
          }}
          sx={{ width: 250 }}
        />
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => {
            resetForm();
            setOpenDialog(true);
          }}
          sx={{
            backgroundColor: '#7f1d1d',
            '&:hover': {
              backgroundColor: '#991b1b',
            },
          }}
        >
          Add Slot
        </Button>
      </div>

      <TableContainer component={Paper} elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Start Time</TableCell>
              <TableCell>End Time</TableCell>
              <TableCell>Max Bookings</TableCell>
              <TableCell>Current Bookings</TableCell>
              <TableCell>Status</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 3 }}>
                  <CircularProgress size={40} sx={{ color: '#7f1d1d' }} />
                </TableCell>
              </TableRow>
            ) : filteredSlots.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 3 }}>
                  No slots found
                </TableCell>
              </TableRow>
            ) : (
              filteredSlots.map((slot) => (
                <TableRow key={slot._id} hover>
                  <TableCell>{slot.start_time}</TableCell>
                  <TableCell>{slot.end_time}</TableCell>
                  <TableCell>{slot.max_bookings}</TableCell>
                  <TableCell>{slot.current_bookings}</TableCell>
                  <TableCell>
                    <Chip 
                      label={slot.current_bookings >= slot.max_bookings ? "Full" : "Available"}
                      color={slot.current_bookings >= slot.max_bookings ? "error" : "success"}
                      size="small"
                    />
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => {
                        setSelectedSlot(slot);
                        setFormData({
                          ...formData,
                          start_time: slot.start_time,
                          end_time: slot.end_time,
                          max_bookings: slot.max_bookings,
                        });
                        setOpenDialog(true);
                      }}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                    <IconButton
                      size="small"
                      color="error"
                      onClick={() => handleDelete(slot._id)}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog 
        open={openDialog} 
        onClose={() => setOpenDialog(false)}
        PaperProps={{
          sx: { borderRadius: 2, p: 2 }
        }}
      >
        <div className="p-6 space-y-4 min-w-[400px]">
          <h2 className="text-xl font-semibold mb-4">
            {selectedSlot ? "Edit Slot" : "Add New Slot"}
          </h2>
          <div className="space-y-4">
            <TextField
              fullWidth
              label="Start Time"
              type="time"
              value={formData.start_time}
              onChange={(e) => setFormData({ ...formData, start_time: e.target.value })}
              InputLabelProps={{ shrink: true }}
            />
            <TextField
              fullWidth
              label="End Time"
              type="time"
              value={formData.end_time}
              onChange={(e) => setFormData({ ...formData, end_time: e.target.value })}
              InputLabelProps={{ shrink: true }}
            />
            <TextField
              fullWidth
              label="Max Bookings"
              type="number"
              value={formData.max_bookings}
              onChange={(e) => setFormData({ ...formData, max_bookings: parseInt(e.target.value) })}
              InputProps={{ inputProps: { min: 1 } }}
            />
          </div>
          <div className="flex justify-end space-x-2 mt-6">
            <Button onClick={() => setOpenDialog(false)}>
              Cancel
            </Button>
            <Button
              variant="contained"
              onClick={handleSubmit}
              sx={{
                backgroundColor: '#7f1d1d',
                '&:hover': {
                  backgroundColor: '#991b1b',
                },
              }}
            >
              {selectedSlot ? "Update" : "Create"}
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
};

export default FacilitySlotManager; 