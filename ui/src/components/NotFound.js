import React from "react";
import { connect } from "react-redux";
// import { history } from "configureStore"; // Remove history import

@connect()
/**
 * Route not found
 */
export default class NotFound extends React.Component {
  /**
   * When component mount redirect to root route
   */
  componentDidMount() {
    // history.push("/"); // Remove history.push
    // TODO: If this component is meant to be a "Not Found" page, add UI here.
    // If all unknown routes should redirect to home, handle that in routes/index.js
  }

  /**
   * Component render
   */

  render() {
    // TODO: Return actual "Not Found" UI if this is a page.
    return "Page Not Found"; // Placeholder content
  }
}
