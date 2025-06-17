import React, { Fragment } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { titleDirective } from "../../utils/Title";
import { MoneyDirective } from "../../utils/Money";
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
  LinearProgress,
  Tooltip,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";

const useStyles = makeStyles(() => ({
  unused: {
    fontSize: "42px",
    color: "orangered",
    fontFamily: "MuseoModerno",
    minHeight: "63px",
  },
  unused_daily: {
    fontSize: "42px",
    color: "purple",
    fontFamily: "MuseoModerno",
    minHeight: "63px",
  },
  unused_resource: {
    fontSize: "42px",
    color: "darkgreen",
    fontFamily: "Nunito",
    fontWeight: "400",
    minHeight: "63px",
  },
  middleGrid: {
    textAlign: "center",
    borderLeft: "1px dashed #c1c1c1",
    borderRight: "1px dashed #c1c1c1",
  },
  grid: {
    textAlign: "center",
  },
  progress: {
    margin: "30px",
  },
}));

/**
 * @param  {array} {resources  Resources List
 * @param  {bool} isResourceListLoading  isLoading indicator for resources
 * @param  {func} currentResource  Current Selected Resource
 */
const StatisticsBar = ({
  resources,
  isResourceListLoading,
  currentResource,
}) => {
  const classes = useStyles();

  let HighestResourceName = "";
  let HighestResourceValue = 0;
  let TotalUnusedResourcesCount = 0;
  
  const TotalSpent = Object.values(resources).reduce((acc, resource) => {
    let TotalSpent = resource.TotalSpent;

    if (currentResource && currentResource !== resource.ResourceName) {
      TotalSpent = 0;
    }

    if (resource.TotalSpent > HighestResourceValue) {
      HighestResourceValue = resource.TotalSpent;
      HighestResourceName = resource.ResourceName;
    }

    // Count unused resources (resources without pricing)
    if (resource.TotalSpent === 0) {
      TotalUnusedResourcesCount += resource.ResourceCount || 0;
    }

    return acc + TotalSpent;
  }, 0);

  const DailySpent = TotalSpent / 30;

  return (
    <Fragment>
      <Box mb={3}>
        <Card>
          <CardContent>
            <Grid container className={classes.root} spacing={2}>
              <Grid item sm={3} xs={12} className={classes.grid}>
                <Tooltip title="Monthly potential cost savings from unused resources with pricing">
                  <div>
                    {isResourceListLoading && (
                      <LinearProgress className={classes.progress} />
                    )}
                    {!isResourceListLoading && (
                      <Typography
                        variant="inherit"
                        sx={{
                          fontSize: "42px",
                          color: "orangered",
                          fontFamily: "MuseoModerno",
                          minHeight: "63px",
                        }}
                      >
                        {MoneyDirective(TotalSpent)}
                      </Typography>
                    )}
                    <Typography>💰 Monthly potential savings</Typography>
                  </div>
                </Tooltip>
              </Grid>
              <Grid item sm={3} xs={12} className={classes.middleGrid}>
                <Tooltip title="Daily cost savings you can achieve by optimizing unused resources">
                  <div>
                    {isResourceListLoading && (
                      <LinearProgress className={classes.progress} />
                    )}
                    {!isResourceListLoading && (
                      <Typography
                        variant="inherit"
                        sx={{
                          fontSize: "42px",
                          color: "purple",
                          fontFamily: "MuseoModerno",
                          minHeight: "63px",
                        }}
                      >
                        {MoneyDirective(DailySpent)}
                      </Typography>
                    )}
                    <Typography>📅 Daily potential savings</Typography>
                  </div>
                </Tooltip>
              </Grid>
              <Grid item sm={3} xs={12} className={classes.middleGrid}>
                <Tooltip title="Number of unused resources that can be removed (no cost impact but improves security)">
                  <div>
                    {isResourceListLoading && (
                      <LinearProgress className={classes.progress} />
                    )}
                    {!isResourceListLoading && (
                      <Typography
                        variant="inherit"
                        sx={{
                          fontSize: "42px",
                          color: "#4caf50",
                          fontFamily: "MuseoModerno",
                          minHeight: "63px",
                        }}
                      >
                        {TotalUnusedResourcesCount}
                      </Typography>
                    )}
                    <Typography>🗑️ Unused resources</Typography>
                  </div>
                </Tooltip>
              </Grid>
              <Grid item sm={3} xs={12} className={classes.grid}>
                {isResourceListLoading && (
                  <LinearProgress className={classes.progress} />
                )}
                {!isResourceListLoading && (
                  <Typography
                    variant="inherit"
                    sx={{
                      fontSize: "42px",
                      color: "darkgreen",
                      fontFamily: "Nunito",
                      fontWeight: "400",
                      minHeight: "63px",
                    }}
                  >
                    {titleDirective(HighestResourceName).toUpperCase()}
                  </Typography>
                )}
                <Typography>🎯 Highest cost resource</Typography>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </Box>
    </Fragment>
  );
};

StatisticsBar.defaultProps = {};
StatisticsBar.propTypes = {
  isScanning: PropTypes.bool,
  isResourceListLoading: PropTypes.bool,
  currentResource: PropTypes.string,
  resources: PropTypes.object,
  filters: PropTypes.array,
};

const mapStateToProps = (state) => ({
  resources: state.resources.resources,
  filters: state.filters.filters,
  isResourceListLoading: state.resources.isResourceListLoading,
  currentResource: state.resources.currentResource,
});
const mapDispatchToProps = () => ({});

export default connect(mapStateToProps, mapDispatchToProps)(StatisticsBar);
