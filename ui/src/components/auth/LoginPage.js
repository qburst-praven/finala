import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Container,
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Alert,
} from "@mui/material";
// Assuming logo.png is the one to use. Adjust path if necessary, or use icon.png
// Path will be /icons/logo.png assuming image is moved to ui/public/icons/
// import FinalaLogo from '../../styles/icons/logo.png'; // Keep this commented or remove

// Add the import for the Logo component
import Logo from "../Logo";

const LoginPage = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [apiBaseUrl, setApiBaseUrl] = useState("");
  const [configLoading, setConfigLoading] = useState(true);
  const [configError, setConfigError] = useState("");

  useEffect(() => {
    const fetchApiConfig = async () => {
      try {
        // Fetch from the UI server, which might proxy or serve settings directly
        const response = await fetch("/api/v1/settings");
        if (!response.ok) {
          throw new Error(`Failed to fetch API settings: ${response.status}`);
        }
        const configData = await response.json();
        if (configData && configData.api_endpoint) {
          setApiBaseUrl(configData.api_endpoint);
        } else {
          throw new Error("API endpoint not found in settings");
        }
      } catch (err) {
        /* eslint-disable no-console */
        console.error("Error fetching API config:", err);
        /* eslint-enable no-console */
        setConfigError(
          `Failed to load API configuration: ${err.message}. Using default.`,
        );
        // Fallback or default if settings fetch fails - crucial for direct backend calls if proxy fails
        setApiBaseUrl("http://localhost:8089");
      } finally {
        setConfigLoading(false);
      }
    };

    fetchApiConfig();
  }, []);

  const handleLogin = async (event) => {
    event.preventDefault();

    if (!apiBaseUrl) {
      setError("API configuration is not loaded yet. Please wait.");
      return;
    }

    setLoading(true);
    setError("");

    if (!username || !password) {
      setError("Username and password are required.");
      setLoading(false);
      return;
    }

    try {
      // Construct full URL using fetched or default apiBaseUrl
      const response = await fetch(`${apiBaseUrl}/api/v1/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      // Try to read the response as text first to aid debugging
      const responseText = await response.text();

      if (!response.ok) {
        let apiError = "Login failed. Please try again.";
        try {
          const errorData = JSON.parse(responseText);
          if (errorData && errorData.error) {
            apiError = errorData.error;
          }
        } catch (e) {
          /* eslint-disable no-console */
          console.error(
            "Could not parse error response as JSON:",
            responseText,
          );
          /* eslint-enable no-console */
          if (responseText.length < 200) {
            apiError = `Server error: ${responseText}`;
          } else {
            apiError = `Server error (status ${response.status}). Check console for details.`;
          }
        }
        setError(apiError);
        setLoading(false);
        return;
      }

      try {
        const data = JSON.parse(responseText);
        if (data.token) {
          localStorage.setItem("finalaAuthToken", data.token);
          setUsername("");
          setPassword("");
          navigate("/");
        } else {
          setError(data.error || "Login successful but no token received.");
        }
      } catch (e) {
        /* eslint-disable no-console */
        console.error(
          "Failed to parse successful response JSON:",
          responseText,
          e,
        );
        /* eslint-enable no-console */
        setError("Received an invalid response from the server.");
      }
    } catch (err) {
      /* eslint-disable no-console */
      console.error("Login API error:", err);
      /* eslint-enable no-console */
      // Check if it's a CORS issue or network failure when apiBaseUrl is directly used
      if (err instanceof TypeError && apiBaseUrl.startsWith("http")) {
        // Likely network error (CORS, server down)
        setError(
          `Network error or CORS issue trying to reach ${apiBaseUrl}. Please check backend API and CORS settings.`,
        );
      } else {
        setError("An unexpected error occurred. Please try again later.");
      }
    } finally {
      setLoading(false);
    }
  };

  // Render logic to show loading/error for config
  if (configLoading) {
    return (
      <Container
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          minHeight: "100vh",
        }}
      >
        <CircularProgress />
        <Typography sx={{ ml: 2 }}>Loading API configuration...</Typography>
      </Container>
    );
  }

  return (
    <Container
      component="main"
      maxWidth="xs"
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "100vh",
      }}
    >
      <Paper
        elevation={3}
        sx={{
          padding: 4,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          width: "100%",
        }}
      >
        <Box sx={{ mb: 3 }}>
          <Logo width={228} height={80} />
        </Box>
        <Typography component="h1" variant="h5" sx={{ mb: 2 }}>
          Welcome to Finala
        </Typography>
        <Typography component="p" variant="subtitle1" sx={{ mb: 3 }}>
          Please log in to continue
        </Typography>
        {configError && (
          <Alert severity="warning" sx={{ width: "100%", mb: 2 }}>
            {configError}
          </Alert>
        )}
        <Box component="form" onSubmit={handleLogin} sx={{ width: "100%" }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Username"
            name="username"
            autoComplete="username"
            autoFocus
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            disabled={loading || !!configError}
            sx={{
              "& .MuiOutlinedInput-root": {
                "&.Mui-focused fieldset": {
                  borderColor: "#DC143C",
                },
              },
              "& .MuiInputLabel-root.Mui-focused": {
                color: "#DC143C",
              },
            }}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={loading || !!configError}
            sx={{
              "& .MuiOutlinedInput-root": {
                "&.Mui-focused fieldset": {
                  borderColor: "#DC143C",
                },
              },
              "& .MuiInputLabel-root.Mui-focused": {
                color: "#DC143C",
              },
            }}
          />
          {error && (
            <Alert severity="error" sx={{ width: "100%", mt: 2, mb: 1 }}>
              {error}
            </Alert>
          )}
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{
              mt: 3,
              mb: 2,
              bgcolor: "#DC143C",
              "&:hover": {
                bgcolor: "#B01030",
              },
            }}
            disabled={loading || !!configError}
          >
            {loading ? (
              <CircularProgress size={24} color="inherit" />
            ) : (
              "Log In"
            )}
          </Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default LoginPage;
