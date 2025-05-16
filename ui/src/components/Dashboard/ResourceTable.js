import React, { Fragment, useEffect, useState } from "react";
import { connect, useSelector } from "react-redux";
import PropTypes from "prop-types";
import numeral from "numeral";
import MUIDataTable from "mui-datatables";
import TextUtils from "utils/Text";
import TagsDialog from "../Dialog/Tags";
import ReportProblemIcon from "@material-ui/icons/ReportProblem";
import InfoIcon from "@material-ui/icons/Info";
import { getHistory } from "../../utils/History";
import { useTableFilters } from "../../Hooks/TableHooks";
import CustomToolbar from "./CustomToolbar";

import {
  makeStyles,
  Card,
  CardContent,
  LinearProgress,
  Paper,
  Box,
  Typography,
  Fade,
} from "@material-ui/core";

import Moment from "moment";

const useStyles = makeStyles((theme) => ({
  tableWrapper: {
    marginBottom: theme.spacing(3),
    borderRadius: theme.shape.borderRadius,
    overflow: 'hidden',
    boxShadow: '0 2px 10px rgba(0,0,0,0.05)',
    backgroundColor: theme.palette.background.paper,
    '& .MuiPaper-root': {
      boxShadow: 'none',
      borderRadius: 0,
    },
    '& .MuiToolbar-root': {
      backgroundColor: theme.palette.background.paper,
      padding: theme.spacing(2, 3),
    },
    '& .MuiTableCell-head': {
      backgroundColor: theme.palette.background.default,
      fontWeight: 600,
      color: theme.palette.text.primary,
      fontSize: '0.875rem',
      borderBottom: `2px solid ${theme.palette.divider}`,
      padding: theme.spacing(1.5, 2),
    },
    '& .MuiTableCell-body': {
      fontSize: '0.875rem',
      padding: theme.spacing(1.5, 2),
      borderBottom: `1px solid ${theme.palette.divider}`,
    },
    '& .MuiTableRow-root:nth-child(even)': {
      backgroundColor: 'rgba(0, 0, 0, 0.02)',
    },
    '& .MuiTableRow-root:hover': {
      backgroundColor: 'rgba(0, 123, 255, 0.04)',
    },
    '& .MuiTablePagination-root': {
      color: theme.palette.text.primary,
      fontWeight: 500,
      borderTop: `1px solid ${theme.palette.divider}`,
    },
    '& .MuiTableFooter-root': {
      backgroundColor: theme.palette.background.paper,
    },
    '& .MuiChip-root': {
      fontWeight: 500,
      backgroundColor: theme.palette.primary.light,
      color: theme.palette.primary.contrastText,
    },
    '& .MuiIconButton-root': {
      color: theme.palette.primary.main,
    },
  },
  cardError: {
    marginBottom: theme.spacing(3),
    borderRadius: theme.shape.borderRadius,
    boxShadow: '0 2px 10px rgba(0,0,0,0.05)',
    overflow: 'hidden',
  },
  cardContent: {
    padding: theme.spacing(4),
    textAlign: 'center',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
  },
  message: {
    marginTop: theme.spacing(2),
    color: theme.palette.text.secondary,
    maxWidth: '80%',
  },
  alertIcon: {
    fontSize: '3.5rem',
    color: theme.palette.error.main,
    marginBottom: theme.spacing(2),
  },
  infoIcon: {
    fontSize: '3.5rem',
    color: theme.palette.info.main,
    marginBottom: theme.spacing(2),
  },
  progress: {
    width: '50%',
    margin: theme.spacing(4, 0),
    height: 8,
    borderRadius: 4,
  },
  priceValue: {
    fontWeight: 600,
    color: theme.palette.text.primary,
  },
  positiveValue: {
    color: theme.palette.success.main,
    fontWeight: 500,
  },
  negativeValue: {
    color: theme.palette.error.main,
    fontWeight: 500,
  },
  neutralValue: {
    color: theme.palette.text.primary,
    fontWeight: 400,
  },
  dateValue: {
    color: theme.palette.text.secondary,
    fontSize: '0.875rem',
  },
}));

/**
 * @param  {array} {resources  Resources List
 * @param  {string} currentResource  Current Selected Resource
 * @param  {array} currentResourceData  Current Selected Resource data
 * @param  {bool} isResourceTableLoading  isLoading indicator for table}
 */
const ResourceTable = ({
  resources,
  currentResource,
  currentResourceData,
  isResourceTableLoading,
  addFiltersObject,
  removeFiltersObject,
  getFlits,
  getCols,
  checkUncheckColumns,
  getSearchText,
  dispatchSearchText,
}) => {
  const [headers, setHeaders] = useState([]);
  const [errorMessage, setErrorMessage] = useState(false);
  const [hasError, setHasError] = useState(false);
  const classes = useStyles();
  const [setTableFilters] = useTableFilters({});
  const [tableOptions, setTableOptions] = useState({});

  // setting table configuration on first load
  useEffect(() => {
    setTableOptions({
      page: parseInt(getHistory("page", 0)),
      searchText: getHistory("search", ""),
      sortOrder: {
        name: getHistory("sortColumn", ""),
        direction: getHistory("direction", "desc"),
      },
      selectableRows: "none",
      responsive: "standard",
      rowsPerPageOptions: [10, 20, 50, 100],
      rowsPerPage: 20,
      elevation: 0,
    });
  }, []);

  /**
   * format table cell by type
   * @param {string} key TableCell key
   * @returns {func} render function to render cell
   */
  const getRowRender = (key) => {
    let renderr = false;
    switch (key) {
      case "PricePerMonth":
      case "TotalSpendPrice":
        renderr = (data) => (
          <Typography variant="body2" className={classes.priceValue}>
            {numeral(data).format("$ 0,0[.]00")}
          </Typography>
        );
        break;
      case "PricePerHour":
        renderr = (data) => (
          <Typography variant="body2" className={classes.priceValue}>
            {numeral(data).format("$ 0,0[.]000")}
          </Typography>
        );
        break;
      case "Tag":
        renderr = (data) => <TagsDialog tags={data} />;
        break;
      case "LaunchTime":
        renderr = (data) => (
          <Typography variant="body2" className={classes.dateValue}>
            {Moment(data).format("YYYY-MM-DD HH:mm")}
          </Typography>
        );
        break;
      case "Status":
        renderr = (data) => {
          let className = classes.neutralValue;
          if (data === "Active" || data === "Running") className = classes.positiveValue;
          if (data === "Stopped" || data === "Error") className = classes.negativeValue;
          return (
            <Typography variant="body2" className={className}>
              {data}
            </Typography>
          );
        };
        break;
      default:
        renderr = (data) => (
          <Typography variant="body2">{`${data}`}</Typography>
        );
    }
    return renderr;
  };

  /**
   * determines Table header keys
   * @param {object} exampleRow  sample row from data
   * @returns {array} Table header keys
   */
  const getHeaderRow = (row) => {
    const exclude = ["TotalSpendPrice"];
    const keys = Object.keys(row).reduce((filtered, headerKey) => {
      if (exclude.indexOf(headerKey) === -1) {
        const header = {
          name: headerKey,
          label: TextUtils.CamelCaseToTitleCase(headerKey),
          options: {
            customBodyRender: getRowRender(headerKey),
          },
        };
        filtered.push(header);
      }
      return filtered;
    }, []);
    return keys;
  };

  /**
   * Detect resource data changed
   */
  var filterNameArray;
  useEffect(() => {
    let headers = [];
    if (currentResourceData.length) {
      headers = getHeaderRow(currentResourceData[0]);
    }
    filterNameArray = headers && headers.map((obj) => obj.name);
    checkUncheckColumns(filterNameArray);
    setHeaders(headers);
  }, [currentResourceData]);
  /**
   * Detect if we have an error
   */
  useEffect(() => {
    if (!currentResource) {
      return;
    }
    const resourceInfo = resources[currentResource];
    if (resourceInfo && resourceInfo.Status === 1) {
      setHasError(true);
      setErrorMessage(resourceInfo.ErrorMessage);
      return;
    } else {
      setHasError(false);
    }
  }, [currentResource, resources]);
  const [flits, setFlits] = useState({ test: "only" });
  return (
    <Fragment>
      {!hasError && isResourceTableLoading && (
        <Fade in={true}>
          <Card className={classes.cardError}>
            <CardContent className={classes.cardContent}>
              <Typography variant="h6" gutterBottom>
                Loading Resource Data
              </Typography>
              <LinearProgress className={classes.progress} color="primary" />
            </CardContent>
          </Card>
        </Fade>
      )}

      {!isResourceTableLoading && (hasError || !currentResourceData.length) && (
        <Fade in={true}>
          <Card className={classes.cardError}>
            <CardContent className={classes.cardContent}>
              {hasError && (
                <ReportProblemIcon className={classes.alertIcon} />
              )}
              
              {!hasError && !currentResourceData.length && (
                <InfoIcon className={classes.infoIcon} />
              )}

              {hasError && (
                <Typography variant="h6" gutterBottom>
                  Error Scanning Resource
                </Typography>
              )}

              {!isResourceTableLoading &&
                !hasError &&
                !currentResourceData.length && (
                  <Typography variant="h6" gutterBottom>
                    No Data Available
                  </Typography>
                )}

              {(hasError || !currentResourceData.length) && (
                <Typography 
                  variant="body1" 
                  className={classes.message}
                >
                  {hasError 
                    ? "Finala couldn't scan the selected resource. Please check system logs for more information."
                    : "No data was found for the selected resource."}
                </Typography>
              )}

              {errorMessage && (
                <Box mt={2}>
                  <Typography variant="body2" color="error">
                    {errorMessage}
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </Fade>
      )}

      {!hasError && currentResourceData.length > 0 && !isResourceTableLoading && (
        <Fade in={true}>
          <Box className={classes.tableWrapper}>
            <MUIDataTable
              title={
                <Typography variant="h6">
                  {TextUtils.CamelCaseToTitleCase(currentResource)} Resources
                </Typography>
              }
              data={currentResourceData}
              columns={headers}
              options={Object.assign(tableOptions, {
                customSearch: (searchQuery, currentRow, columns) => {
                  return "EMAIL";
                },
                onSearchChange: (searchText) => {
                  dispatchSearchText(searchText);
                  setTableFilters([
                    {
                      key: "search",
                      value: searchText ? searchText : "",
                    },
                  ]);
                },
                onColumnSortChange: (changedColumn, direction) => {
                  setTableFilters([
                    { key: "sortColumn", value: changedColumn },
                    { key: "direction", value: direction },
                  ]);
                },
                onChangePage: (currentPage) => {
                  setTableFilters([{ key: "page", value: currentPage }]);
                },
                onChangeRowsPerPage: (numberOfRows) => {
                  setTableFilters([{ key: "rows", value: numberOfRows }]);
                },
                downloadOptions: {
                  filename: `${currentResource}.csv`,
                },
                customToolbar: () => {
                  return <CustomToolbar />;
                },
                onFilterChipClose: (index, removedFilter, filterList) => {
                  removeFiltersObject({
                    column: "Data." + removedFilter,
                  });
                },
                onFilterChange: (changedColumn, filterList, type) => {
                  if (filterList[changedColumn].length) {
                    addFiltersObject({
                      column: "Data." + changedColumn,
                      filterArr: filterList[changedColumn],
                    });
                  }
                },
              })}
            />
          </Box>
        </Fade>
      )}
    </Fragment>
  );
};

ResourceTable.defaultProps = {};
ResourceTable.propTypes = {
  currentResource: PropTypes.string,
  resources: PropTypes.object,
  currentResourceData: PropTypes.array,
  isResourceTableLoading: PropTypes.bool,
  addFiltersObject: PropTypes.func,
  removeFiltersObject: PropTypes.func,
  dispatchSearchText: PropTypes.func,
  getFlits: PropTypes.object,
  getCols: PropTypes.array,
  checkUncheckColumns: PropTypes.func,
  getSearchText: PropTypes.string,
};
const mapStateToProps = (state) => ({
  resources: state.resources.resources,
  currentResourceData: state.resources.currentResourceData,
  currentResource: state.resources.currentResource,
  isResourceTableLoading: state.resources.isResourceTableLoading,
  getFlits: state.flit,
  getCols: state.cols,
  getSearchText: state.searchMui,
});
const mapDispatchToProps = (dispatch) => ({
  addFiltersObject: (data) => dispatch({ type: "ADD_IN_OBJECT", data }),
  removeFiltersObject: (data) => dispatch({ type: "REMOVE_IN_OBJECT", data }),
  dispatchSearchText: (data) => dispatch({ type: "ON_TEXT_ENTERED", data }),
  checkUncheckColumns: (data) =>
    dispatch({ type: "CHECK_UNCHECK_COLUMNS_CHECKBOX", data }),
});
export default connect(mapStateToProps, mapDispatchToProps)(ResourceTable);
