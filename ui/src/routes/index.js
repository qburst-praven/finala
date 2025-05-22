import React from "react";
import { connect } from "react-redux";
import {
  Route,
  Routes,
  Navigate,
  useLocation,
  useNavigate,
} from "react-router-dom";
import PropTypes from "prop-types";
import Dashboard from "../components/Dashboard/Index";
import PageLoader from "../components/PageLoader";
import NotFound from "../components/NotFound";
import NoData from "../components/NoData";
import DataFactory from "../components/DataFactory";

// Auth components
import LoginPage from "../components/auth/LoginPage";
import ProtectedRoute from "../components/auth/ProtectedRoute";

import { CssBaseline, Box } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";

const useStyles = makeStyles(() => ({
  root: {
    background: "#f1f5f9",
    color: "#27303f",
  },
  content: {
    background: "#f1f5f9",
    color: "#27303f",
    minHeight: "100vh",
    display: "flex",
    flexDirection: "column",
  },
  hide: {
    display: "none",
  },
}));

const RouterIndex = ({ isAppLoading, executions }) => {
  const classes = useStyles();
  const location = useLocation();
  const isAuthenticated = !!localStorage.getItem("finalaAuthToken");

  if (isAuthenticated && location.pathname === "/login") {
    return <Navigate to="/" replace />;
  }

  return (
    <Routes>
      <Route
        path="/login"
        element={isAuthenticated ? <Navigate to="/" replace /> : <LoginPage />}
      />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <AppLayout isAppLoading={isAppLoading} executions={executions} />
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
};

const AppLayout = ({ isAppLoading, executions }) => {
  const navigate = useNavigate();

  const classes = useStyles();
  return (
    <div className={classes.root}>
      <CssBaseline />
      <DataFactory />
      <main className={classes.content}>
        <Box
          component="div"
          m={3}
          sx={{ flexGrow: 1, p: { xs: 2, sm: 3 }, pt: 0 }}
        >
          {isAppLoading && <PageLoader />}
          {!isAppLoading && !executions.length && <NoData />}
          {!isAppLoading && executions.length > 0 && (
            <Box component="div">
              <Dashboard />
            </Box>
          )}
        </Box>
      </main>
    </div>
  );
};

AppLayout.propTypes = {
  isAppLoading: PropTypes.bool,
  executions: PropTypes.array,
};

const mapStateToProps = (state) => ({
  executions: state.executions.list,
  isAppLoading: state.executions.isAppLoading,
});

const mapDispatchToProps = () => ({});

RouterIndex.defaultProps = {};
RouterIndex.propTypes = {
  isAppLoading: PropTypes.bool,
  executions: PropTypes.array,
};

export default connect(mapStateToProps, mapDispatchToProps)(RouterIndex);
