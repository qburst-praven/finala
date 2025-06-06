import React, { Fragment } from "react";
import { connect } from "react-redux";
import makeStyles from "@mui/styles/makeStyles";
import { useNavigate } from "react-router-dom";
import { setHistory } from "../../utils/History";
import { Button } from "@mui/material";

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
                sx={{
                  color: "#DC143C",
                  borderColor: "#DC143C",
                  "&:hover": {
                    borderColor: "#B01030",
                    backgroundColor: "rgba(220, 20, 60, 0.04)",
                  },
                }}
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
