import React, { Fragment } from "react";
import { connect } from "react-redux";
import makeStyles from "@mui/styles/makeStyles";
import { useNavigate } from "react-router-dom";
import { setHistory } from "../../utils/History";
import { Button } from "@mui/material";
import { Card, CardContent, Typography } from "@mui/material";

import PropTypes from "prop-types";
import FilterBar from "./FilterBar";
import StatisticsBar from "./StatisticsBar";
import ResourceScanning from "./ResourceScanning";
import ResourcesChart from "./ResourcesChart";
import ResourcesList from "./ResourcesList";
import ResourceTable from "./ResourceTable";
import ExecutionIndex from "../Executions/Index";
import Logo from "../Logo";
import { Grid, Box } from "@mui/material";

const useStyles = makeStyles(() => ({
  root: {
    width: "100%",
  },
  title: {
    fontFamily: "MuseoModerno",
  },
  logoGrid: {
    textAlign: "left",
  },
  selectorGrid: {
    textAlign: "right",
  },
  overviewCard: {
    backgroundColor: "#fbfcfd",
    border: "1px solid #e1e5e9",
    borderRadius: "8px",
    boxShadow: "0 1px 3px rgba(0, 0, 0, 0.05)",
    marginBottom: "24px",
    transition: "box-shadow 0.2s ease",
    "&:hover": {
      boxShadow: "0 4px 12px rgba(0, 0, 0, 0.08)",
    },
  },
  overviewTitle: {
    fontFamily: "MuseoModerno",
    fontWeight: "600",
    fontSize: "1.4rem",
    color: "#1a202c",
    marginBottom: "12px",
    display: "flex",
    alignItems: "center",
    gap: "8px",
  },
  overviewDescription: {
    color: "#4a5568",
    lineHeight: "1.6",
    fontSize: "0.95rem",
  },
  categoryLabel: {
    fontWeight: "600",
    fontSize: "0.9rem",
  },
  costSavingLabel: {
    color: "#d69e2e",
  },
  unusedLabel: {
    color: "#38a169",
  },
  logoutButton: {
    backgroundColor: "#fff",
    color: "#e53e3e",
    border: "1px solid #e53e3e",
    fontWeight: "500",
    borderRadius: "6px",
    padding: "8px 16px",
    textTransform: "none",
    fontSize: "0.9rem",
    transition: "all 0.2s ease",
    "&:hover": {
      backgroundColor: "#e53e3e",
      color: "#fff",
      transform: "translateY(-1px)",
      boxShadow: "0 2px 8px rgba(229, 62, 62, 0.25)",
    },
  },
}));

/**
 * @param  {string} {currentResource  Current Selected Resource
 * @param  {func} setFilters  Update Filters
 * @param  {func} setResource  Update Selected Resource
 * @param  {array} filters   Filters List } */
const DashboardIndex = ({
  currentResource,
  setFilters,
  setResource,
  filters,
}) => {
  const classes = useStyles();
  const navigate = useNavigate();
  /**
   * Will clear selected filter and show main page
   */
  const gotoHome = () => {
    const updatedFilters = filters.filter((filter) => {
      return filter.id.substr(0, 8) !== "resource";
    });
    setResource(null);
    setFilters(updatedFilters);
    setHistory(navigate, { filters: updatedFilters });
  };

  return (
    <Fragment>
      <Box mb={2}>
        <Grid container className={classes.root} spacing={0}>
          <Grid item sm={8} xs={12} className={classes.logoGrid}>
            <a href="javascript:void(0)" onClick={gotoHome}>
              <Logo />
            </a>
            <ResourceScanning />
          </Grid>
          <Grid item sm={4} xs={12} className={classes.selectorGrid}>
            <Box sx={{ display: "flex", justifyContent: "flex-end", mb: 1 }}>
              <Button
                variant="outlined"
                onClick={() => {
                  localStorage.removeItem("finalaAuthToken");
                  navigate("/login");
                }}
                className={classes.logoutButton}
              >
                Logout
              </Button>
            </Box>
            <ExecutionIndex />
          </Grid>
        </Grid>
      </Box>



      <FilterBar />
      <StatisticsBar />
      <ResourcesList />
      {currentResource ? <ResourceTable /> : <ResourcesChart />}
    </Fragment>
  );
};

DashboardIndex.defaultProps = {};
DashboardIndex.propTypes = {
  currentResource: PropTypes.string,
  filters: PropTypes.array,
  setFilters: PropTypes.func,
  setResource: PropTypes.func,
};

const mapStateToProps = (state) => ({
  currentResource: state.resources.currentResource,
  filters: state.filters.filters,
});

const mapDispatchToProps = (dispatch) => ({
  setFilters: (data) => dispatch({ type: "SET_FILTERS", data }),
  setResource: (data) => dispatch({ type: "SET_RESOURCE", data }),
});

export default connect(mapStateToProps, mapDispatchToProps)(DashboardIndex);
