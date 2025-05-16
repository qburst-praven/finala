import { createMuiTheme } from '@material-ui/core/styles';
import { red } from '@material-ui/core/colors'; // Example for error color

// Define a modern color palette
const PALETTE = {
  primary: {
    main: '#007BFF', // A vibrant blue
    contrastText: '#FFFFFF',
  },
  secondary: {
    main: '#6C757D', // A neutral grey
    contrastText: '#FFFFFF',
  },
  error: {
    main: red.A400, // Standard Material Design error red
  },
  background: {
    default: '#F4F6F8', // A light, clean background
    paper: '#FFFFFF',   // Background for elements like Cards
  },
  text: {
    primary: '#212529', // Dark grey for primary text
    secondary: '#495057', // Lighter grey for secondary text
  },
};

// Define typography
const TYPOGRAPHY = {
  fontFamily: [
    '"Nunito"', // Already imported in index.html
    'Roboto',
    '"Helvetica Neue"',
    'Arial',
    'sans-serif'
  ].join(','),
  h1: {
    fontFamily: '"MuseoModerno", "Nunito", sans-serif', // Using MuseoModerno for headings
    fontWeight: 500,
    fontSize: '2.5rem', // Example size
  },
  h2: {
    fontFamily: '"MuseoModerno", "Nunito", sans-serif',
    fontWeight: 500,
    fontSize: '2rem',
  },
  h3: {
    fontFamily: '"MuseoModerno", "Nunito", sans-serif',
    fontWeight: 500,
    fontSize: '1.75rem',
  },
  body1: {
    fontWeight: 400,
    fontSize: '1rem',
    lineHeight: 1.5,
  },
  button: {
    textTransform: 'none', // Buttons without ALL CAPS
    fontWeight: 600,
  }
};

// Define component overrides
const OVERRIDES = {
  MuiButton: {
    root: {
      borderRadius: 8, // Slightly more rounded buttons
    },
    containedPrimary: {
      '&:hover': {
        backgroundColor: '#0069D9', // Darker shade on hover for primary button
      },
    },
  },
  MuiCard: {
    root: {
      borderRadius: 12, // More rounded cards
      boxShadow: '0px 4px 12px rgba(0,0,0,0.08)', // A softer, modern shadow
    }
  },
  MuiAppBar: { // Example: If you use AppBar
    colorDefault: {
      backgroundColor: PALETTE.background.paper,
      color: PALETTE.text.primary,
    }
  },
  // Add more component overrides as needed for Paper, TextField, Table, etc.
};

const theme = createMuiTheme({
  palette: PALETTE,
  typography: TYPOGRAPHY,
  overrides: OVERRIDES,
  shape: {
    borderRadius: 8, // Global border radius, can be overridden per component
  },
  // You can also define props defaults
  // props: {
  //   MuiButtonBase: {
  //     disableRipple: true, // Example: disable ripple effect globally
  //   },
  // },
});

export default theme; 