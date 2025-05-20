import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import App from "./App";
import configureStore from "./configureStore";

const store = configureStore();

const container = document.getElementById("react-root");

const root = createRoot(container);

root.render(
  <Provider store={store}>
    <App />
  </Provider>,
);
