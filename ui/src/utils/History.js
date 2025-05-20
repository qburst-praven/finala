// import { history } from "configureStore"; // Remove

const windowParams = new window.URLSearchParams(window.location.search);
let savedFilters = windowParams.get("filters");
let savedExecutionId = windowParams.get("executionId");

/**
 *
 * @param {array} filters filters list
 * @returns filters params for request
 */
export const transformFilters = (filters) => {
  // return filters;
  const params = {};
  const list = [];

  filters.forEach((filter) => {
    const [key, value] = filter.id.split(":");
    if (value) {
      if (!params[key]) {
        params[key] = [];
      }
      params[key].push(value);
    }
  });
  for (const key in params) {
    const paramsAsString = params[key].join(",");
    list.push(`${key}:${paramsAsString}`);
  }
  return list.join(";");
};

/**
 *
 * @param {function} navigate - The navigate function from useNavigate() hook
 * @param {object} historyParams - {filters, executionId}
 *  Will set State into url search params, will save the params so we can set only part of the params
 */
export const setHistory = (navigate, historyParams = {}) => {
  // Add navigate parameter
  if (typeof navigate !== "function") {
    console.error("setHistory requires a navigate function from useNavigate.");
    return;
  }
  savedFilters = Object.prototype.hasOwnProperty.call(historyParams, "filters")
    ? transformFilters(historyParams.filters)
    : savedFilters;

  savedExecutionId = Object.prototype.hasOwnProperty.call(
    historyParams,
    "executionId",
  )
    ? historyParams.executionId
    : savedExecutionId;
  const params = {};

  if (savedFilters && savedFilters.length) {
    params.filters = savedFilters;
  }
  if (savedExecutionId) {
    params.executionId = savedExecutionId;
  }

  const searchParams = new window.URLSearchParams(params);
  // history.push({  // Old way
  //   pathname: "/",
  //   search: decodeURIComponent(`?${searchParams.toString()}`),
  // });
  navigate(`/?${decodeURIComponent(searchParams.toString())}`, {
    replace: true,
  }); // New way
};

/**
 *
 * @param {string} query params name from url
 * @param {any} defaultValue return default value in case the query not exists
 * @returns {string} Param value from url
 */
export const getHistory = (query, defaultValue = null) => {
  const searchParams = new window.URLSearchParams(window.location.search);
  const searchQuery = searchParams.get(query);
  if (searchQuery) {
    return searchQuery;
  } else {
    return defaultValue;
  }
};
