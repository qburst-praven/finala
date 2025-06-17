import React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { setHistory } from "../../utils/History";
import { useNavigate } from "react-router-dom";
import { MoneyDirective } from "../../utils/Money";
import {
  Card,
  CardContent,
  Chip,
  Grid,
  Typography,
  Box,
  Divider,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";

const useStyles = makeStyles(() => ({
  card: {
    marginBottom: "24px",
    backgroundColor: "#ffffff",
    border: "1px solid #e2e8f0",
    borderRadius: "8px",
    boxShadow: "0 1px 3px rgba(0, 0, 0, 0.06)",
    transition: "box-shadow 0.2s ease",
    "&:hover": {
      boxShadow: "0 4px 12px rgba(0, 0, 0, 0.08)",
    },
  },
  cardContent: {
    padding: "24px 24px 8px 24px !important",
  },
  sectionTitle: {
    fontFamily: "MuseoModerno",
    fontWeight: "700",
    fontSize: "1.2rem",
    marginBottom: "16px",
    display: "flex",
    alignItems: "center",
    gap: "8px",
  },
  costSavingTitle: {
    color: "#d69e2e",
  },
  unusedTitle: {
    color: "#38a169",
  },
  resourceChip: {
    margin: "4px",
    borderRadius: "6px",
    fontWeight: "500",
    fontSize: "0.85rem",
    border: "1px solid transparent",
    transition: "all 0.2s ease",
    cursor: "pointer",
    "&:hover": {
      transform: "translateY(-1px)",
      boxShadow: "0 4px 8px rgba(0, 0, 0, 0.12)",
    },
  },
  costSavingChip: {
    "&:hover": {
      borderColor: "rgba(214, 158, 46, 0.3)",
      backgroundColor: "rgba(214, 158, 46, 0.08)",
    },
  },
  unusedChip: {
    "&:hover": {
      borderColor: "rgba(56, 161, 105, 0.3)",
      backgroundColor: "rgba(56, 161, 105, 0.08)",
    },
  },
  emptyState: {
    textAlign: "center",
    padding: "40px 20px",
    color: "#718096",
    fontStyle: "italic",
  },
  divider: {
    margin: "20px 0 16px 0",
    backgroundColor: "#e2e8f0",
  },
  resourceGrid: {
    minHeight: "60px",
    display: "flex",
    alignItems: "flex-start",
    flexWrap: "wrap",
    gap: "8px",
    marginBottom: "0px",
  },
}));

const colors = [
  "#3f51b5", "#f44336", "#ff9800", "#4caf50", "#9c27b0",
  "#e91e63", "#00bcd4", "#795548", "#607d8b", "#ff5722"
];

const ResourcesList = ({ resources, setResource, setFilters, filters }) => {
  const classes = useStyles();
  const navigate = useNavigate();

  // Separate resources by category
  const costSavingResources = Object.values(resources || {})
    .filter((resource) => 
      resource.Category === "potential_cost_saving" || 
      (resource.TotalSpent && resource.TotalSpent > 0)
    )
    .sort((a, b) => (b.TotalSpent || 0) - (a.TotalSpent || 0));

  const unusedResources = Object.values(resources || {})
    .filter((resource) => 
      (resource.Category === "unused_resource" || 
       (!resource.TotalSpent || resource.TotalSpent === 0)) &&
      (resource.ResourceCount && resource.ResourceCount > 0)
    )
    .sort((a, b) => (b.ResourceCount || 0) - (a.ResourceCount || 0));

  const setResourceFilter = (resourceName) => {
    const newFilters = filters.filter((filter) => 
      filter.id.substr(0, 8) !== "resource"
    );
    newFilters.push({
      title: `Resource:${resourceName}`,
      id: `resource:${resourceName}`,
      value: resourceName,
      type: "resource",
    });
    setResource(resourceName);
    setFilters(newFilters);
    setHistory(navigate, { filters: newFilters });
  };

  const renderResourceChips = (resourceList, isUnused = false) => {
    if (!resourceList.length) {
      return (
        <Typography className={classes.emptyState}>
          No {isUnused ? "unused" : "cost-saving"} resources found
        </Typography>
      );
    }

    return (
      <Box className={classes.resourceGrid}>
        {resourceList.map((resource, index) => {
          const colorIndex = index % colors.length;
          const backgroundColor = colors[colorIndex];
          
          return (
            <Chip
              key={resource.ResourceName}
              label={
                isUnused
                  ? `${resource.ResourceName} (${resource.ResourceCount || 0})`
                  : `${resource.ResourceName} (${MoneyDirective(resource.TotalSpent || 0)})`
              }
              onClick={() => setResourceFilter(resource.ResourceName)}
              className={`${classes.resourceChip} ${
                isUnused ? classes.unusedChip : classes.costSavingChip
              }`}
              style={{
                backgroundColor,
                color: "#ffffff",
              }}
            />
          );
        })}
      </Box>
    );
  };

  return (
    <Card className={classes.card}>
      <CardContent className={classes.cardContent}>
        {/* Cost Saving Resources Section */}
        <Typography className={`${classes.sectionTitle} ${classes.costSavingTitle}`}>
          💰 Potential Cost Saving Resources
        </Typography>
        {renderResourceChips(costSavingResources, false)}

        <Divider className={classes.divider} />

        {/* Unused Resources Section */}
        <Typography className={`${classes.sectionTitle} ${classes.unusedTitle}`}>
          🗑️ Unused Resources
        </Typography>
        {renderResourceChips(unusedResources, true)}
      </CardContent>
    </Card>
  );
};

ResourcesList.defaultProps = {
  resources: {},
};

ResourcesList.propTypes = {
  resources: PropTypes.object,
  setResource: PropTypes.func.isRequired,
  setFilters: PropTypes.func.isRequired,
  filters: PropTypes.array.isRequired,
};

const mapStateToProps = (state) => ({
  resources: state.resources.resources,
  filters: state.filters.filters,
});

const mapDispatchToProps = (dispatch) => ({
  setResource: (data) => dispatch({ type: "SET_RESOURCE", data }),
  setFilters: (data) => dispatch({ type: "SET_FILTERS", data }),
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesList);
