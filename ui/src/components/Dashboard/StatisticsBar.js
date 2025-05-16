import React from "react";
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
  Divider,
} from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import AttachMoneyIcon from '@material-ui/icons/AttachMoney';
import TrendingDownIcon from '@material-ui/icons/TrendingDown';
import StorageIcon from '@material-ui/icons/Storage';

const useStyles = makeStyles((theme) => ({
  statsCard: {
    borderRadius: theme.shape.borderRadius,
    boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
    overflow: 'hidden',
    marginBottom: theme.spacing(3),
  },
  statItem: {
    padding: theme.spacing(2),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    transition: 'all 0.3s ease',
    position: 'relative',
    '&:hover': {
      backgroundColor: 'rgba(0,0,0,0.01)',
    },
  },
  statIconWrapper: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    width: 60,
    height: 60,
    borderRadius: '50%',
    marginBottom: theme.spacing(2),
  },
  monthlyIcon: {
    backgroundColor: 'rgba(255, 87, 34, 0.1)',
    color: '#FF5722',
  },
  dailyIcon: {
    backgroundColor: 'rgba(156, 39, 176, 0.1)',
    color: '#9C27B0',
  },
  resourceIcon: {
    backgroundColor: 'rgba(76, 175, 80, 0.1)',
    color: '#4CAF50',
  },
  statIcon: {
    fontSize: '2rem',
  },
  monthly: {
    fontSize: '2.2rem',
    fontWeight: 600,
    color: '#FF5722',
    fontFamily: theme.typography.h2.fontFamily,
    minHeight: '3.2rem',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  daily: {
    fontSize: '2.2rem',
    fontWeight: 600,
    color: '#9C27B0',
    fontFamily: theme.typography.h2.fontFamily,
    minHeight: '3.2rem',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  resource: {
    fontSize: '2rem',
    fontWeight: 500,
    color: '#4CAF50',
    fontFamily: theme.typography.h3.fontFamily,
    minHeight: '3.2rem',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    textAlign: 'center',
  },
  statLabel: {
    color: theme.palette.text.secondary,
    fontSize: '0.9rem',
    fontWeight: 500,
    textAlign: 'center',
    marginTop: theme.spacing(1),
  },
  progress: {
    margin: theme.spacing(3),
    height: 6,
    borderRadius: 3,
  },
  verticalDivider: {
    position: 'absolute',
    top: '20%',
    bottom: '20%',
    width: 1,
    backgroundColor: 'rgba(0,0,0,0.08)',
  },
  leftDivider: {
    right: 0,
  },
  rightDivider: {
    left: 0,
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
  const TotalSpent = Object.values(resources).reduce((acc, resource) => {
    let TotalSpent = resource.TotalSpent;

    if (currentResource && currentResource !== resource.ResourceName) {
      TotalSpent = 0;
    }

    if (resource.TotalSpent > HighestResourceValue) {
      HighestResourceValue = resource.TotalSpent;
      HighestResourceName = resource.ResourceName;
    }

    return acc + TotalSpent;
  }, 0);

  const DailySpent = TotalSpent / 30;

  return (
    <Card className={classes.statsCard} elevation={0}>
      <Grid container>
        <Grid item sm={4} xs={12}>
          <Tooltip title="Monthly Unused resources are effected from filters" placement="top">
            <Box className={classes.statItem}>
              <Box className={`${classes.statIconWrapper} ${classes.monthlyIcon}`}>
                <AttachMoneyIcon className={classes.statIcon} />
              </Box>
              {isResourceListLoading ? (
                <LinearProgress className={classes.progress} color="primary" />
              ) : (
                <Typography className={classes.monthly}>
                  {MoneyDirective(TotalSpent)}
                </Typography>
              )}
              <Typography className={classes.statLabel}>
                Monthly Unused Resources
              </Typography>
              <Box className={`${classes.verticalDivider} ${classes.leftDivider}`} />
            </Box>
          </Tooltip>
        </Grid>
        
        <Grid item sm={4} xs={12}>
          <Tooltip title="Daily waste is the amount you pay daily for unused resources and can be saved" placement="top">
            <Box className={classes.statItem}>
              <Box className={`${classes.statIconWrapper} ${classes.dailyIcon}`}>
                <TrendingDownIcon className={classes.statIcon} />
              </Box>
              {isResourceListLoading ? (
                <LinearProgress className={classes.progress} color="secondary" />
              ) : (
                <Typography className={classes.daily}>
                  {MoneyDirective(DailySpent)}
                </Typography>
              )}
              <Typography className={classes.statLabel}>
                Daily Waste
              </Typography>
              <Box className={`${classes.verticalDivider} ${classes.leftDivider}`} />
              <Box className={`${classes.verticalDivider} ${classes.rightDivider}`} />
            </Box>
          </Tooltip>
        </Grid>
        
        <Grid item sm={4} xs={12}>
          <Box className={classes.statItem}>
            <Box className={`${classes.statIconWrapper} ${classes.resourceIcon}`}>
              <StorageIcon className={classes.statIcon} />
            </Box>
            {isResourceListLoading ? (
              <LinearProgress className={classes.progress} style={{ color: '#4CAF50' }} />
            ) : (
              <Typography className={classes.resource}>
                {titleDirective(HighestResourceName).toUpperCase()}
              </Typography>
            )}
            <Typography className={classes.statLabel}>
              Most Unused Resource
            </Typography>
          </Box>
        </Grid>
      </Grid>
    </Card>
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
