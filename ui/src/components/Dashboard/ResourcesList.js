import React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import colors from "./colors.json";
import { makeStyles } from "@material-ui/core/styles";
import { setHistory } from "../../utils/History";

import { Box, Chip, Typography, Paper, Grid } from "@material-ui/core";
import { titleDirective } from "../../utils/Title";
import { MoneyDirective } from "../../utils/Money";
import StorageIcon from '@material-ui/icons/Storage';

const useStyles = makeStyles((theme) => ({
  resourcesContainer: {
    padding: theme.spacing(2, 3),
    marginBottom: theme.spacing(3),
    borderRadius: theme.shape.borderRadius,
    backgroundColor: theme.palette.background.paper,
    boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
  },
  headerContainer: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: theme.spacing(2),
  },
  resourceIcon: {
    marginRight: theme.spacing(1),
    color: theme.palette.primary.main,
  },
  title: {
    color: theme.palette.text.primary,
    fontWeight: 600,
    fontFamily: theme.typography.h3.fontFamily,
    flexGrow: 1,
  },
  chipContainer: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: theme.spacing(1),
  },
  resource_chips: {
    margin: theme.spacing(0.5),
    transition: 'all 0.2s ease',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: theme.palette.background.default,
    fontSize: '0.875rem',
    fontWeight: 500,
    boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
    '&:hover': {
      boxShadow: '0 3px 6px rgba(0,0,0,0.15)',
      backgroundColor: theme.palette.background.paper,
    },
    border: 'none',
    '& .MuiChip-label': {
      paddingLeft: theme.spacing(1),
    },
    position: 'relative',
    paddingLeft: theme.spacing(0.5),
  },
  colorIndicator: {
    width: 5,
    position: 'absolute',
    left: 0,
    top: 0,
    bottom: 0,
    borderTopLeftRadius: theme.shape.borderRadius,
    borderBottomLeftRadius: theme.shape.borderRadius,
  },
}));

/**
 * @param  {array} {resources  Resources List
 * @param  {array} filters  Filters List
 * @param  {func} addFilter Add filter to  filters list
 * @param  {func} setResource Update Selected Resource}
 */
const ResourcesList = ({ resources, filters, addFilter, setResource }) => {
  const classes = useStyles();
  const resourcesList = Object.values(resources)
    .sort((a, b) => (a.TotalSpent > b.TotalSpent ? -1 : 1))
    .map((resource) => {
      const title = titleDirective(resource.ResourceName);
      const amount = MoneyDirective(resource.TotalSpent);
      resource.title = `${title} (${amount})`;
      resource.display_title = `${title}`;

      return resource;
    });

  /**
   *
   * @param {object} resource Set selected resource
   */
  const setSelectedResource = (resource) => {
    const filter = {
      title: `Resource:${resource.display_title}`,
      id: `resource:${resource.ResourceName}`,
      type: "resource",
    };
    setResource(resource.ResourceName);
    addFilter(filter);

    setHistory({
      filters: filters,
    });
  };

  if (!resourcesList.length) {
    return null;
  }

  return (
    <Paper className={classes.resourcesContainer} elevation={0}>
      <Box className={classes.headerContainer}>
        <StorageIcon className={classes.resourceIcon} />
        <Typography variant="h6" className={classes.title}>
          Resources
        </Typography>
      </Box>
      <Grid container className={classes.chipContainer}>
        {resourcesList.map((resource, i) => (
          <Grid item key={i}>
            <Chip
              className={classes.resource_chips}
              onClick={() => setSelectedResource(resource)}
              label={resource.title}
              clickable
            />
            <Box 
              className={classes.colorIndicator} 
              style={{ backgroundColor: colors[i % colors.length]?.hex || '#cccccc' }}
            />
          </Grid>
        ))}
      </Grid>
    </Paper>
  );
};

ResourcesList.defaultProps = {};
ResourcesList.propTypes = {
  resources: PropTypes.object,
  filters: PropTypes.array,
  addFilter: PropTypes.func,
  setResource: PropTypes.func,
};

const mapStateToProps = (state) => ({
  resources: state.resources.resources,
  filters: state.filters.filters,
});
const mapDispatchToProps = (dispatch) => ({
  addFilter: (data) => dispatch({ type: "ADD_FILTER", data }),
  setResource: (data) => dispatch({ type: "SET_RESOURCE", data }),
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesList);
