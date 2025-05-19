import React from "react";
import PropTypes from "prop-types";
import { BrowserRouter } from "react-router-dom";
import Routes from "./routes";
import { connect } from "react-redux";
import "./styles/index.scss";
import { createTheme, ThemeProvider } from "@mui/material/styles";

// Create a default theme
const theme = createTheme();

// Main application class
class App extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <ThemeProvider theme={theme}>
        <BrowserRouter basename="">
          <Routes />
        </BrowserRouter>
      </ThemeProvider>
    );
  }
}

function mapStateToProps() {
  return {};
}

App.propTypes = {
  // history: PropTypes.object,
};
export default connect(mapStateToProps)(App);
