import React from 'react';
import { CircularProgress, Box } from '@mui/material';

const LoadingComponent = () => {
  return (
    <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
      <CircularProgress sx={{ color: '#7f1d1d' }} />
    </Box>
  );
};

export default LoadingComponent; 