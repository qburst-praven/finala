import React from "react";
import { connect } from "react-redux";
import { makeStyles } from "@material-ui/core/styles";
import { setHistory } from "../../utils/History";

import PropTypes from "prop-types";
import FilterBar from "./FilterBar";
import StatisticsBar from "./StatisticsBar";
import ResourceScanning from "./ResourceScanning";
import ResourcesChart from "./ResourcesChart";
import ResourcesList from "./ResourcesList";
import ResourceTable from "./ResourceTable";
import ExecutionIndex from "../Executions/Index";
import Logo from "../Logo";
import { Grid, Box, Container, Paper, Typography } from "@material-ui/core";
import DashboardIcon from '@material-ui/icons/Dashboard';

const useStyles = makeStyles((theme) => ({
  dashboardWrapper: {
    padding: theme.spacing(2, 0),
  },
  headerContainer: {
    borderRadius: theme.shape.borderRadius,
    padding: theme.spacing(2, 3),
    marginBottom: theme.spacing(3),
    backgroundColor: theme.palette.background.paper,
    boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
    display: 'flex',
    alignItems: 'center',
  },
  headerLeft: {
    display: 'flex',
    alignItems: 'center',
  },
  headerIcon: {
    color: theme.palette.primary.main,
    marginRight: theme.spacing(2),
    fontSize: '1.5rem',
  },
  title: {
    color: theme.palette.text.primary,
    fontFamily: theme.typography.h3.fontFamily,
    fontWeight: 600,
    display: 'flex',
    alignItems: 'center',
  },
  logoGrid: {
    display: 'flex',
    alignItems: 'center',
  },
  logo: {
    cursor: 'pointer',
  },
  scanningText: {
    marginLeft: theme.spacing(2),
  },
  selectorGrid: {
    display: 'flex',
    justifyContent: 'flex-end',
    alignItems: 'center',
  },
  contentArea: {
    padding: theme.spacing(0),
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
  /**
   * Will clear selected filter and show main page
   */
  const gotoHome = () => {
    const updatedFilters = filters.filter(
      (filter) => filter.id.substr(0, 8) !== "resource"
    );
    setResource(null);
    setFilters(updatedFilters);
    setHistory({ filters: updatedFilters });
  };

  return (
    <Container maxWidth="xl" className={classes.dashboardWrapper}>
      <Paper className={classes.headerContainer} elevation={0}>
        <Grid container>
          <Grid item sm={8} xs={12} className={classes.logoGrid}>
            <Box className={classes.headerLeft}>
              <DashboardIcon className={classes.headerIcon} />
              <Typography variant="h5" className={classes.title}>
                Dashboard
              </Typography>
            </Box>
          </Grid>
          <Grid item sm={4} xs={12} className={classes.selectorGrid}>
            <ExecutionIndex />
          </Grid>
        </Grid>
      </Paper>

      <Box mb={2} className={classes.logo} onClick={gotoHome}>
        <Logo />
      </Box>
      
      <ResourceScanning />
      
      <Box className={classes.contentArea}>
        <FilterBar />
        <StatisticsBar />
        <ResourcesList />
        {currentResource ? <ResourceTable /> : <ResourcesChart />}
      </Box>
    </Container>
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
