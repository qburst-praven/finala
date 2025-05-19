import React, { Fragment, useState, useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import EmailIcon from "@mui/icons-material/Email";
import makeStyles from "@mui/styles/makeStyles";
import Modal from "@mui/material/Modal";
import Alert from "@mui/material/Alert";
import Snackbar from "@mui/material/Snackbar";
import {
  FormControl,
  FormLabel,
  TextField,
  Button,
  Hidden,
} from "@mui/material";
import { http } from "../../services/request.service";

function rand() {
  return Math.round(Math.random() * 20) - 10;
}
function getModalStyle() {
  const top = 50;
  const left = 50;
  return {
    top: `${top}%`,
    left: `${left}%`,
    transform: `translate(-${top}%, -${left}%)`,
    border: `none`,
  };
}
const useStyles = makeStyles((theme) => ({
  paper: {
    position: "absolute",
    width: 400,
    backgroundColor: theme.palette.background.paper,
    border: "2px solid #000",
    boxShadow: theme.shadows[5],
    padding: theme.spacing(2, 4, 3),
  },
}));
const CustomToolbar = (props, getFlits) => {
  console.log("BASE URL", http.baseURL);
  const classes = useStyles();
  // getModalStyle is not a pure function, we roll the style only on the first render
  const [modalStyle] = useState(getModalStyle);
  const [open, setOpen] = useState(false);
  const [openSnackSuccess, setOpenSnackSuccess] = useState(false);
  const [openSnackError, setOpenSnackError] = useState(false);
  const [executionId, setExecutionId] = useState(null);
  const [disabledBtn, setDisabledBtn] = useState(false);
  useEffect(() => {
    var url = window.location.search;
    url = url
      .replace("?", "")
      .split("&")
      .map((param) => param.split("="))
      .reduce((values, [key, value]) => {
        values[key] = value;
        return values;
      }, {}); // remove the ?
    setFormData({
      ...formData,
      executionID: url.executionId,
      resourceType: url.filters ? url.filters.replace("resource:", "") : "",
      filters: props.getFlits,
    });
  }, []);
  const setCookie = (name, value, days) => {
    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`;
  };
  const getCookie = (name) => {
    const cookieName = `${name}=`;
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(";");
    for (let cookie of cookieArray) {
      while (cookie.charAt(0) === " ") {
        cookie = cookie.substring(1);
      }
      if (cookie.indexOf(cookieName) === 0) {
        return cookie.substring(cookieName.length, cookie.length);
      }
    }
  };
  const deleteCookie = (name) => {
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`;
  };
  const handleStackClose = () => {
    setOpenSnackError(false);
    setOpenSnackSuccess(false);
  };
  const handleClick = () => {
    console.log("clicked on icon!");
    console.log(getCookie("toEmails"));
    setOpen(true);
    if (getCookie("toEmails")) {
      //if cookie set then auto filled in form
      setFormData({ ...formData, toEmails: getCookie("toEmails") });
    }
  };
  const handleOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };
  const [formData, setFormData] = useState({
    // Initialize your form fields here
    toEmails: null,
    columns: props.getCols,
  });
  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };
  const handleSubmit = async (event) => {
    setDisabledBtn(true);
    deleteCookie("toEmails");
    event.preventDefault();
    formData.filters = props.getFlits;
    formData.search = props.getSearchText;
    setCookie("toEmails", formData.toEmails, 7); // Sets a cookie named 'cookieName' with value 'cookieValue' that expires in 7 days
    // const fullUrl = `http://127.0.0.1:8081/api/v1/send-report`;
    const fullUrl = `${http.baseURL}/api/v1/send-report`;
    try {
      fetch(fullUrl, {
        method: "POST", // Specify the HTTP method
        body: JSON.stringify(formData), // Collect form data
      })
        .then((response) => response.json()) // Read response as text
        .then((data) => {
          setDisabledBtn(false);
          if (data.status === 200) {
            setOpenSnackError(false);
            setOpenSnackSuccess(true);
            setTimeout(() => {
              setOpenSnackSuccess(false);
            }, 6000);
            //reset form
            setFormData({
              ...formData,
              toEmails: "",
            });
          } else {
            setOpenSnackError(true);
            setOpenSnackSuccess(false);
            setTimeout(() => {
              setOpenSnackError(false);
            }, 6000);
          }
        }); // Alert the response
    } catch (error) {
      setDisabledBtn(false);
      setOpenSnackError(true);
      setOpenSnackSuccess(false);
      setTimeout(() => {
        setOpenSnackError(false);
      }, 6000);
    }
  };
  return (
    <Fragment>
      <span onClick={handleClick}>
        <button className="MuiButtonBase-root MuiIconButton-root jss430">
          <EmailIcon />
        </button>
      </span>
      <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby="simple-modal-title"
        aria-describedby="simple-modal-description"
      >
        <div style={modalStyle} className={classes.paper}>
          <h2 id="simple-modal-title">Report Send To</h2>
          <hr />
          {/* {"GET FLITES object :-" + JSON.stringify(props.getFlits, null, 2)} */}
          {/* {"GET cols array:-" + JSON.stringify(props.getCols, null, 2)} */}
          <form onSubmit={handleSubmit}>
            <FormControl>
              <h4>Note:You can pass multiple email by comma separated.</h4>
              <FormLabel>Enter Email:</FormLabel>
              <TextField
                name="toEmails"
                required="true"
                size="medium"
                color="primary"
                placeholder="abc@gmail.com,xyz@gmail.com"
                onChange={handleInputChange}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.toEmails}
                multiline={true}
                rows={3}
              ></TextField>
              <input
                type="hidden"
                value={formData.executionID}
                name="executionID"
                onChange={handleInputChange}
              ></input>
              <Button
                color="primary"
                variant="contained"
                type="submit"
                disabled={disabledBtn}
              >
                Submit
              </Button>
            </FormControl>
          </form>
        </div>
      </Modal>
      <Snackbar open={openSnackSuccess} autoHideDuration={600}>
        <Alert onClose={handleStackClose} severity="success">
          Success! Report will send on entered email id(s) soon.
        </Alert>
      </Snackbar>
      <Snackbar open={openSnackError} autoHideDuration={600}>
        <Alert onClose={handleStackClose} severity="error">
          Something went wrong! Please try again.
        </Alert>
      </Snackbar>
    </Fragment>
  );
};
CustomToolbar.propTypes = {
  dbFilter: PropTypes.object,
  getFlits: PropTypes.object,
  getCols: PropTypes.array,
  getSearchText: PropTypes.string,
};

const mapStateToProps = (state) => ({
  getFlits: state.flit,
  getCols: state.cols,
  getSearchText: state.searchMui,
});

export default connect(mapStateToProps)(CustomToolbar);
