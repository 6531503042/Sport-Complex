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
  CircularProgress,
  useTheme,
  useMediaQuery
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
import SearchIcon from "@mui/icons-material/Search";
import { toast } from 'react-hot-toast';
import { muiStyles } from '../styles/shared';
import { motion } from 'framer-motion';

interface FacilitySlotManagerProps {
  facilityName: string;
}

interface Slot {
  _id: string;
  start_time: string;
  end_time: string;
  status: number;
  max_bookings: number;
  current_bookings: number;
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
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const isTablet = useMediaQuery(theme.breakpoints.down('md'));

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const url = `http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots`;
        console.log('Fetching from URL:', url);
        
        const response = await fetch(url);
        console.log('Response:', response);
        
        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
        }
        
        const result = await response.json();
        console.log('Raw response data:', result);
        
        const slotsData = result.data || result;
        console.log('Slots data to be used:', slotsData);
        
        if (Array.isArray(slotsData)) {
          setSlots(slotsData);
        } else {
          console.warn('Data is not in expected format:', slotsData);
          setSlots([]);
        }
      } catch (error: any) {
        console.error('Fetch error:', error);
        toast.error(`Error: ${error.message}`);
        setSlots([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [facilityName]);

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
        const fetchData = async () => {
          const response = await fetch(`http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots`);
          const data = await response.json();
          setSlots(data.data || []);
        };
        fetchData();
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
        const fetchData = async () => {
          const response = await fetch(`http://localhost:1335/facility_v1/${facilityName}/slot_v1/slots`);
          const data = await response.json();
          setSlots(data.data || []);
        };
        fetchData();
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
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="space-y-4"
    >
      <div className={`flex ${isMobile ? 'flex-col' : 'justify-between'} items-center gap-4 mb-6`}>
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
          sx={{ 
            ...muiStyles.searchField,
            width: isMobile ? '100%' : 250,
          }}
        />
        <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              resetForm();
              setOpenDialog(true);
            }}
            fullWidth={isMobile}
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
            Add Slot
          </Button>
        </motion.div>
      </div>

      <TableContainer 
        component={Paper} 
        elevation={2}
        sx={{ 
          ...muiStyles.responsive.table.container,
          borderRadius: 2,
          overflow: 'hidden',
        }}
      >
        <Table>
          <TableHead>
            <TableRow>
              <TableCell sx={muiStyles.responsive.table.headerCell}>Start Time</TableCell>
              <TableCell sx={muiStyles.responsive.table.headerCell}>End Time</TableCell>
              {!isMobile && (
                <>
                  <TableCell sx={muiStyles.responsive.table.headerCell}>Max Bookings</TableCell>
                  <TableCell sx={muiStyles.responsive.table.headerCell}>Current Bookings</TableCell>
                </>
              )}
              <TableCell sx={muiStyles.responsive.table.headerCell}>Status</TableCell>
              <TableCell align="right" sx={muiStyles.responsive.table.headerCell}>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 6 }}>
                  <motion.div
                    animate={{ rotate: 360 }}
                    transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                  >
                    <CircularProgress size={40} sx={{ color: '#7f1d1d' }} />
                  </motion.div>
                </TableCell>
              </TableRow>
            ) : filteredSlots.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 6 }}>
                  <div className="text-gray-500 flex flex-col items-center">
                    <SearchIcon sx={{ fontSize: 48, color: '#cbd5e1' }} />
                    <p className="mt-2">No slots found</p>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              filteredSlots.map((slot) => (
                <TableRow 
                  key={slot._id} 
                  sx={muiStyles.table.row}
                  component={motion.tr}
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ duration: 0.3 }}
                >
                  <TableCell>{slot.start_time}</TableCell>
                  <TableCell>{slot.end_time}</TableCell>
                  {!isMobile && (
                    <>
                      <TableCell>{slot.max_bookings}</TableCell>
                      <TableCell>{slot.current_bookings}</TableCell>
                    </>
                  )}
                  <TableCell>
                    <Chip 
                      label={slot.current_bookings >= slot.max_bookings ? "Full" : "Available"}
                      sx={slot.current_bookings >= slot.max_bookings ? muiStyles.chip.full : muiStyles.chip.available}
                      size="small"
                    />
                  </TableCell>
                  <TableCell align="right">
                    <motion.div whileHover={{ scale: 1.1 }}>
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
                        sx={muiStyles.table.actionButton}
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                    </motion.div>
                    <motion.div whileHover={{ scale: 1.1 }}>
                      <IconButton
                        size="small"
                        color="error"
                        onClick={() => handleDelete(slot._id)}
                        sx={muiStyles.table.actionButton}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </motion.div>
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
        fullScreen={isMobile}
        PaperProps={{
          sx: {
            ...muiStyles.dialog.paper,
            width: isMobile ? '100%' : 'auto',
            margin: isMobile ? 0 : 24,
            maxHeight: isMobile ? '100%' : '90vh',
          }
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
              sx={muiStyles.button.primary}
            >
              {selectedSlot ? "Update" : "Create"}
            </Button>
          </div>
        </div>
      </Dialog>
    </motion.div>
  );
};

export default FacilitySlotManager;