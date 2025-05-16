import React from "react";
import ReactDOM from "react-dom";
import { AppContainer } from "react-hot-loader";
import { Provider } from "react-redux";
import { ThemeProvider } from "@material-ui/core/styles";
import App from "./App";
import configureStore, { history } from "./configureStore";
import theme from "./theme";

const store = configureStore();
const render = () => {
  ReactDOM.render(
    <AppContainer>
      <Provider store={store}>
        <ThemeProvider theme={theme}>
          <App history={history} />
        </ThemeProvider>
      </Provider>
    </AppContainer>,
    document.getElementById("react-root")
  );
};

render();
