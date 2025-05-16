import React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import colors from "./colors.json";
import Chart from "react-apexcharts";
import { titleDirective } from "../../utils/Title";
import { MoneyDirective } from "../../utils/Money";
import { setHistory } from "../../utils/History";

import {
  Box,
  Card,
  CardContent,
  LinearProgress,
  makeStyles,
  Typography,
  Paper,
  Fade,
} from "@material-ui/core";
import ReportProblemIcon from "@material-ui/icons/ReportProblem";
import BarChartIcon from "@material-ui/icons/BarChart";
import InfoIcon from "@material-ui/icons/Info";

const useStyles = makeStyles((theme) => ({
  chartContainer: {
    marginBottom: theme.spacing(3),
    borderRadius: theme.shape.borderRadius,
    boxShadow: '0 2px 10px rgba(0,0,0,0.05)',
    backgroundColor: theme.palette.background.paper,
    overflow: 'hidden',
  },
  chartHeader: {
    padding: theme.spacing(2, 3),
    borderBottom: `1px solid ${theme.palette.divider}`,
    display: 'flex',
    alignItems: 'center',
  },
  chartIcon: {
    marginRight: theme.spacing(1.5),
    color: theme.palette.primary.main,
  },
  chartTitle: {
    color: theme.palette.text.primary,
    fontWeight: 600,
  },
  chartContent: {
    padding: theme.spacing(3),
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    position: 'relative',
  },
  noDataContainer: {
    textAlign: "center",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    padding: theme.spacing(5),
  },
  alertIcon: {
    fontSize: "3.5rem",
    color: theme.palette.info.main,
    marginBottom: theme.spacing(2),
  },
  noDataTitle: {
    marginBottom: theme.spacing(1),
    color: theme.palette.text.primary,
  },
  noDataMessage: {
    color: theme.palette.text.secondary,
    maxWidth: '400px',
  },
  loadingProgress: {
    width: '50%',
    height: 8,
    borderRadius: 4,
    marginTop: theme.spacing(4),
    marginBottom: theme.spacing(4),
  },
  tooltipCustom: {
    backgroundColor: theme.palette.background.paper,
    color: theme.palette.text.primary,
    borderRadius: theme.shape.borderRadius,
    boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
    border: 'none',
    padding: theme.spacing(1.5),
    fontSize: '0.875rem',
  },
}));

/**
 * @param  {array} {resources  Resources List
 * @param  {array} filters  Filters List
 * @param  {bool} isResourceListLoading  isLoading state for resources
 * @param  {func} addFilter Add filter to  filters list
 * @param  {func} setResource Update Selected Resource}
 */
const ResourcesChart = ({
  resources,
  filters,
  isResourceListLoading,
  addFilter,
  setResource,
}) => {
  const classes = useStyles();
  const colorList = colors.map((color) => color.hex);
  const sortedResources = Object.values(resources)
    .filter((row) => row.TotalSpent > 0)
    .sort((a, b) => (a.TotalSpent >= b.TotalSpent ? -1 : 1));

  const chartOptions = {
    options: {
      chart: {
        type: "bar",
        width: "100%",
        fontFamily: '"Nunito", "Helvetica", "Arial", sans-serif',
        toolbar: {
          show: true,
          tools: {
            download: true,
            selection: false,
            zoom: false,
            zoomin: false,
            zoomout: false,
            pan: false,
            reset: false,
          },
        },
        events: {
          dataPointSelection: function (event, chartContext, config) {
            const dataPointIndex = config.dataPointIndex;
            const res = sortedResources;
            const selectedResource = res[dataPointIndex];
            setSelectedResource(selectedResource);
          },
        },
        animations: {
          enabled: true,
          easing: 'easeinout',
          speed: 800,
          animateGradually: {
            enabled: true,
            delay: 150
          },
          dynamicAnimation: {
            enabled: true,
            speed: 350
          }
        }
      },
      colors: colorList,
      tooltip: {
        theme: "light",
        custom: function({ series, seriesIndex, dataPointIndex, w }) {
          return `<div class="${classes.tooltipCustom}">
            <div style="font-weight: 600">${w.globals.labels[dataPointIndex]}</div>
            <div style="margin-top: 4px; display: flex; align-items: center;">
              <span style="color: #666">Value:</span>
              <span style="font-weight: 500; margin-left: 8px">${MoneyDirective(series[seriesIndex][dataPointIndex])}</span>
            </div>
            <div style="font-size: 0.75rem; margin-top: 4px; color: #888">Click to filter by this resource</div>
          </div>`;
        }
      },
      dataLabels: {
        enabled: true,
        formatter: function (val, opt) {
          return MoneyDirective(val);
        },
        style: {
          fontSize: '12px',
          fontWeight: 600,
          colors: ['#fff']
        },
        background: {
          enabled: true,
          foreColor: '#000',
          padding: 4,
          borderRadius: 4,
          borderWidth: 0,
          opacity: 0.5,
        },
      },
      grid: {
        borderColor: 'rgba(0,0,0,0.05)',
        strokeDashArray: 5,
      },
      plotOptions: {
        bar: {
          horizontal: true,
          distributed: true,
          startingShape: "rounded",
          endingShape: "rounded",
          columnWidth: "70%",
          barHeight: "70%",
          borderRadius: 6,
          dataLabels: {
            position: 'center',
          }
        },
      },
      xaxis: {
        categories: [],
        labels: {
          style: {
            colors: '#495057',
            fontSize: '12px',
            fontFamily: '"Nunito", sans-serif',
          }
        },
        axisBorder: {
          color: 'rgba(0,0,0,0.05)',
        },
        axisTicks: {
          color: 'rgba(0,0,0,0.05)',
        }
      },
      yaxis: {
        labels: {
          style: {
            colors: '#495057',
            fontSize: '12px',
            fontFamily: '"Nunito", sans-serif',
          }
        }
      },
      states: {
        hover: {
          filter: {
            type: 'darken',
            value: 0.9,
          }
        },
        active: {
          filter: {
            type: 'darken',
            value: 0.85,
          }
        }
      }
    },
    series: [
      {
        name: "Total Spent",
        data: [],
      },
    ],
  };

  /**
   * update chart height according to number of resources
   */
  const getChartHeight = () => {
    if (!sortedResources || !sortedResources.length) {
      return 500;
    }
    return Math.max(500, 120 * sortedResources.length);
  };

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

  /**
   * push resources into chart
   */
  sortedResources.forEach((resource) => {
    const title = titleDirective(resource.ResourceName);
    const amount = MoneyDirective(resource.TotalSpent);
    resource.title = `${title} (${amount})`;
    resource.display_title = `${title}`;

    chartOptions.options.xaxis.categories.push(title);
    chartOptions.series[0].data.push(resource.TotalSpent);
    return resource;
  });

  return (
    <Fade in={true}>
      <Paper className={classes.chartContainer} elevation={0}>
        <Box className={classes.chartHeader}>
          <BarChartIcon className={classes.chartIcon} />
          <Typography variant="h6" className={classes.chartTitle}>
            Resource Analysis
          </Typography>
        </Box>
        
        <Box className={classes.chartContent} style={{ minHeight: getChartHeight() }}>
          {!isResourceListLoading && sortedResources.length > 0 && (
            <Chart
              id="MainChart"
              height={getChartHeight()}
              options={chartOptions.options}
              series={chartOptions.series}
              type="bar"
            />
          )}
          
          {isResourceListLoading && (
            <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center">
              <Typography variant="h6" gutterBottom>
                Loading Resource Data
              </Typography>
              <LinearProgress className={classes.loadingProgress} color="primary" />
            </Box>
          )}
          
          {!isResourceListLoading && !sortedResources.length && (
            <Box className={classes.noDataContainer}>
              <InfoIcon className={classes.alertIcon} />
              <Typography variant="h6" className={classes.noDataTitle}>
                No Resources Found
              </Typography>
              <Typography variant="body2" className={classes.noDataMessage}>
                No resource data is available for the current filters. Try changing your filters or check back later.
              </Typography>
            </Box>
          )}
        </Box>
      </Paper>
    </Fade>
  );
};

ResourcesChart.defaultProps = {};
ResourcesChart.propTypes = {
  resources: PropTypes.object,
  filters: PropTypes.array,
  isResourceListLoading: PropTypes.bool,
  addFilter: PropTypes.func,
  setResource: PropTypes.func,
};

const mapStateToProps = (state) => ({
  resources: state.resources.resources,
  isResourceListLoading: state.resources.isResourceListLoading,
  filters: state.filters.filters,
});

const mapDispatchToProps = (dispatch) => ({
  addFilter: (data) => dispatch({ type: "ADD_FILTER", data }),
  setResource: (data) => dispatch({ type: "SET_RESOURCE", data }),
});

export default connect(mapStateToProps, mapDispatchToProps)(ResourcesChart);

