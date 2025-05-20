// import { createBrowserHistory } from "history"; // Remove
import {
  applyMiddleware,
  compose,
  legacy_createStore as createStore,
} from "redux";
// import { routerMiddleware } from "connected-react-router"; // Remove
import createRootReducer from "./reducers";
import { thunk as thunkMiddleware } from "redux-thunk";

// export const history = createBrowserHistory({ // Remove
//   basename: "/",
// });

export default function configureStore(preloadedState) {
  const composeEnhancer =
    window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
  const store = createStore(
    createRootReducer(), // Call without history
    preloadedState,
    // composeEnhancer(applyMiddleware(thunkMiddleware, routerMiddleware(history))) // Remove routerMiddleware
    composeEnhancer(applyMiddleware(thunkMiddleware)), // Added trailing comma
  );
  return store;
}
