export const muiStyles = {
  button: {
    primary: {
      backgroundColor: '#7f1d1d !important',
      color: '#ffffff !important',
      padding: '8px 16px',
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
      textTransform: 'none',
      fontSize: '1rem',
      borderRadius: '8px',
      transition: 'all 0.2s ease-in-out',
    },
    iconButton: {
      color: '#7f1d1d',
      backgroundColor: 'white',
      border: '2px solid #7f1d1d',
      '&:hover': {
        backgroundColor: '#7f1d1d',
        color: 'white',
        transform: 'scale(1.05)',
      },
      transition: 'all 0.2s ease',
    },
    deleteButton: {
      color: '#dc2626',
      backgroundColor: 'white',
      border: '2px solid #dc2626',
      '&:hover': {
        backgroundColor: '#dc2626',
        color: 'white',
        transform: 'scale(1.05)',
      },
      transition: 'all 0.2s ease',
    }
  },
  table: {
    header: {
      backgroundColor: '#f8fafc',
      '& th': {
        fontWeight: 600,
        color: '#475569',
        fontSize: '0.875rem',
        textTransform: 'uppercase',
        letterSpacing: '0.05em',
        padding: '16px',
        borderBottom: '2px solid #e2e8f0',
      }
    },
    row: {
      '&:hover': {
        backgroundColor: '#fef2f2 !important',
        transform: 'scale(1.002)',
        boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
      },
      transition: 'all 0.2s ease',
      cursor: 'pointer',
      '& td': {
        padding: '16px',
        fontSize: '0.95rem',
      }
    },
    actionButton: {
      padding: '8px',
      minWidth: 'unset',
      borderRadius: '6px',
      marginLeft: '8px',
    }
  },
  dialog: {
    paper: {
      borderRadius: '12px',
      padding: '24px',
      boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
      '& .MuiDialogTitle-root': {
        fontSize: '1.5rem',
        fontWeight: 600,
        color: '#1f2937',
        paddingBottom: '16px',
        borderBottom: '1px solid #e5e7eb',
      }
    }
  },
  chip: {
    available: {
      backgroundColor: '#dcfce7',
      color: '#166534',
      fontWeight: 600,
      padding: '4px 12px',
      borderRadius: '9999px',
      '& .MuiChip-label': {
        padding: '0',
      },
      border: '1px solid #86efac',
    },
    full: {
      backgroundColor: '#fee2e2',
      color: '#991b1b',
      fontWeight: 600,
      padding: '4px 12px',
      borderRadius: '9999px',
      '& .MuiChip-label': {
        padding: '0',
      },
      border: '1px solid #fca5a5',
    }
  },
  searchField: {
    '& .MuiOutlinedInput-root': {
      borderRadius: '8px',
      backgroundColor: 'white',
      transition: 'all 0.2s ease',
      '&:hover': {
        boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
      },
      '&.Mui-focused': {
        boxShadow: '0 4px 6px rgba(0,0,0,0.05)',
      }
    },
    '& .MuiOutlinedInput-notchedOutline': {
      borderColor: '#e5e7eb',
    },
    '& .MuiInputAdornment-root': {
      color: '#9ca3af',
    }
  },
  addButton: {
    container: {
      display: 'inline-flex',
      position: 'relative',
      overflow: 'hidden',
      borderRadius: '8px',
      '&::after': {
        content: '""',
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        background: 'linear-gradient(45deg, rgba(255,255,255,0.1), rgba(255,255,255,0))',
        transform: 'translateX(-100%)',
        transition: 'transform 0.3s ease',
      },
      '&:hover::after': {
        transform: 'translateX(100%)',
      }
    }
  },
  responsive: {
    table: {
      container: {
        overflowX: 'auto',
        '&::-webkit-scrollbar': {
          height: '6px',
        },
        '&::-webkit-scrollbar-track': {
          background: '#f1f1f1',
        },
        '&::-webkit-scrollbar-thumb': {
          background: '#888',
          borderRadius: '3px',
        },
      },
      cell: {
        '@media (max-width: 600px)': {
          padding: '12px 8px',
          fontSize: '0.875rem',
        },
      },
      headerCell: {
        '@media (max-width: 600px)': {
          padding: '12px 8px',
          fontSize: '0.75rem',
        },
      }
    },
    dialog: {
      content: {
        '@media (max-width: 600px)': {
          margin: '16px',
          width: 'calc(100% - 32px)',
          maxWidth: 'none',
        },
      }
    },
    searchField: {
      '@media (max-width: 600px)': {
        width: '100%',
        marginBottom: '16px',
      }
    },
    actionButtons: {
      container: {
        '@media (max-width: 600px)': {
          display: 'flex',
          flexDirection: 'column',
          gap: '8px',
        }
      }
    }
  }
}; 