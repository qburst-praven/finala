import React, { Fragment } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import colors from "./colors.json";
import makeStyles from "@mui/styles/makeStyles";
import { useNavigate } from "react-router-dom";
import { setHistory } from "../../utils/History";

import { Box, Chip } from "@mui/material";
import { titleDirective } from "../../utils/Title";
import { MoneyDirective } from "../../utils/Money";

const useStyles = makeStyles(() => ({
  title: {
    fontFamily: "MuseoModerno",
  },
  resource_chips: {
    fontWeight: "bold",
    fontFamily: "Arial !important",
    margin: "24px 12px",
    borderRadius: "16px",
    backgroundColor: "#FAFAFA",
    borderLeftWidth: "5px",
    borderLeftStyle: "solid",
    borderBottomWidth: "2px",
    borderBottomStyle: "solid",
    fontSize: "14px",
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
  const navigate = useNavigate();
  
  // Separate resources into two categories
  const potentialCostSavingResources = [];
  const unusedResources = [];
  
  Object.values(resources).forEach((resource) => {
    const title = titleDirective(resource.ResourceName);
    const amount = MoneyDirective(resource.TotalSpent);
    
    // Determine category based on pricing data
    const hasPricing = resource.TotalSpent > 0;
    
    if (hasPricing) {
      resource.title = `${title} (${amount} / ${resource.ResourceCount})`;
      resource.display_title = `${title}`;
      potentialCostSavingResources.push(resource);
    } else {
      resource.title = `${title} (${resource.ResourceCount} resources)`;
      resource.display_title = `${title}`;
      unusedResources.push(resource);
    }
  });

  // Sort Potential Cost Saving resources by TotalSpent (descending)
  potentialCostSavingResources.sort((a, b) => {
    if (a.TotalSpent > b.TotalSpent) return -1;
    if (a.TotalSpent < b.TotalSpent) return 1;
    const countA = a.ResourceCount || 0;
    const countB = b.ResourceCount || 0;
    return countB - countA;
  });

  // Sort Unused resources by ResourceCount (descending)
  unusedResources.sort((a, b) => {
    const countA = a.ResourceCount || 0;
    const countB = b.ResourceCount || 0;
    if (countA > countB) return -1;
    if (countA < countB) return 1;
    return 0;
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

    setHistory(navigate, {
      filters: filters,
    });
  };

  return (
    <Fragment>
      {/* Potential Cost Saving Resources Section */}
      {potentialCostSavingResources.length > 0 && (
        <Box mb={3}>
          <h4 className={classes.title} style={{ color: "#ff6b35" }}>
            💰 Potential Cost Saving Resources:
          </h4>
          <p style={{ fontSize: "14px", color: "#666", marginTop: "5px" }}>
            Resources with pricing data that can save you money
          </p>
          {potentialCostSavingResources.map((resource, i) => {
            const chipColor =
              colors[i] && colors[i].hex ? colors[i].hex : "#ff6b35";
            return (
              <Chip
                key={`cost-saving-${i}`}
                className={classes.resource_chips}
                label={resource.title}
                onClick={() => setSelectedResource(resource)}
                style={{
                  borderLeft: `5px solid ${chipColor}`,
                  borderBottom: `2px solid ${chipColor}`,
                }}
              />
            );
          })}
        </Box>
      )}

      {/* Unused Resources Section */}
      {unusedResources.length > 0 && (
        <Box mb={3}>
          <h4 className={classes.title} style={{ color: "#4caf50" }}>
            🗑️ Unused Resources:
          </h4>
          <p style={{ fontSize: "14px", color: "#666", marginTop: "5px" }}>
            Resources without pricing that can be removed to improve security and management
          </p>
          {unusedResources.map((resource, i) => {
            const chipColor = "#4caf50"; // Green color for unused resources
            return (
              <Chip
                key={`unused-${i}`}
                className={classes.resource_chips}
                label={resource.title}
                onClick={() => setSelectedResource(resource)}
                style={{
                  borderLeft: `5px solid ${chipColor}`,
                  borderBottom: `2px solid ${chipColor}`,
                }}
              />
            );
          })}
        </Box>
      )}
    </Fragment>
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
