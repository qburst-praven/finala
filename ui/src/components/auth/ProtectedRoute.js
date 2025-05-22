import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import PropTypes from "prop-types";

// This component can be enhanced to check token validity, roles, etc.
const ProtectedRoute = ({ children }) => {
  const isAuthenticated = !!localStorage.getItem("finalaAuthToken");
  const location = useLocation();

  if (!isAuthenticated) {
    // Redirect them to the /login page, but save the current location they were
    // trying to go to so we can send them along after they login.
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return children;
};

ProtectedRoute.propTypes = {
  children: PropTypes.node.isRequired,
};

export default ProtectedRoute;
