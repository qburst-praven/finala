import React, { Fragment, useState, useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import EmailIcon from "@material-ui/icons/Email";
import GetAppIcon from "@material-ui/icons/GetApp";
import ViewColumnIcon from "@material-ui/icons/ViewColumn";
import FilterListIcon from "@material-ui/icons/FilterList";
import { makeStyles } from "@material-ui/core/styles";
import {
  Modal,
  Backdrop,
  Fade,
  Paper,
  Typography,
  FormControl,
  FormLabel,
  TextField,
  Button,
  Hidden,
  Snackbar,
  Box,
  IconButton,
  Tooltip,
  Divider,
} from "@material-ui/core";
import Alert from "@material-ui/lab/Alert";
import { http } from "../../services/request.service";
import CloseIcon from "@material-ui/icons/Close";

const useStyles = makeStyles((theme) => ({
  modalPaper: {
    position: "absolute",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    width: 450,
    backgroundColor: theme.palette.background.paper,
    boxShadow: theme.shadows[5],
    padding: 0,
    outline: "none",
    borderRadius: theme.shape.borderRadius,
    overflow: "hidden",
  },
  modalHeader: {
    padding: theme.spacing(2, 3),
    backgroundColor: theme.palette.primary.main,
    color: theme.palette.primary.contrastText,
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  modalContent: {
    padding: theme.spacing(3),
  },
  modalFooter: {
    padding: theme.spacing(2),
    borderTop: `1px solid ${theme.palette.divider}`,
    display: "flex",
    justifyContent: "flex-end",
  },
  formNote: {
    marginBottom: theme.spacing(2),
    fontSize: "0.875rem",
    color: theme.palette.text.secondary,
  },
  formLabel: {
    marginBottom: theme.spacing(1),
    color: theme.palette.text.primary,
    fontWeight: 500,
  },
  toolbarButton: {
    color: theme.palette.primary.main,
    margin: theme.spacing(0, 0.5),
    padding: theme.spacing(1),
    '&:hover': {
      backgroundColor: 'rgba(0, 123, 255, 0.08)',
    },
  },
  closeButton: {
    color: theme.palette.primary.contrastText,
  },
  buttonProgress: {
    position: 'absolute',
    top: '50%',
    left: '50%',
    marginTop: -12,
    marginLeft: -12,
  },
  buttonWrapper: {
    position: 'relative',
  },
  toolbarContainer: {
    display: 'flex',
    alignItems: 'center',
  },
}));

const CustomToolbar = (props, getFlits) => {
  const classes = useStyles();
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
      }, {});
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
    setOpen(true);
    if (getCookie("toEmails")) {
      setFormData({ ...formData, toEmails: getCookie("toEmails") });
    }
  };
  
  const handleClose = () => {
    setOpen(false);
  };
  
  const [formData, setFormData] = useState({
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
    setCookie("toEmails", formData.toEmails, 7);
    
    const fullUrl = `${http.baseURL}/api/v1/send-report`;
    try {
      fetch(fullUrl, {
        method: "POST",
        body: JSON.stringify(formData),
      })
        .then((response) => response.json())
        .then((data) => {
          setDisabledBtn(false);
          if (data.status === 200) {
            setOpenSnackError(false);
            setOpenSnackSuccess(true);
            setOpen(false);
            setTimeout(() => {
              setOpenSnackSuccess(false);
            }, 6000);
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
        });
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
    <div className={classes.toolbarContainer}>
      <Tooltip title="Email Report">
        <IconButton className={classes.toolbarButton} onClick={handleClick}>
          <EmailIcon />
        </IconButton>
      </Tooltip>
      
      <Modal
        open={open}
        onClose={handleClose}
        closeAfterTransition
        BackdropComponent={Backdrop}
        BackdropProps={{
          timeout: 500,
        }}
      >
        <Fade in={open}>
          <Paper className={classes.modalPaper}>
            <Box className={classes.modalHeader}>
              <Typography variant="h6">Email Report</Typography>
              <IconButton className={classes.closeButton} size="small" onClick={handleClose}>
                <CloseIcon />
              </IconButton>
            </Box>
            
            <Box className={classes.modalContent}>
              <Typography variant="body2" className={classes.formNote}>
                You can enter multiple email addresses separated by commas.
              </Typography>
              
              <form onSubmit={handleSubmit}>
                <FormControl fullWidth>
                  <Typography variant="subtitle2" className={classes.formLabel}>
                    Email Addresses
                  </Typography>
                  <TextField
                    name="toEmails"
                    required
                    size="medium"
                    color="primary"
                    placeholder="example@domain.com, another@domain.com"
                    onChange={handleInputChange}
                    variant="outlined"
                    fullWidth
                    value={formData.toEmails || ""}
                    multiline
                    rows={3}
                  />
                  <input
                    type="hidden"
                    value={formData.executionID}
                    name="executionID"
                    onChange={handleInputChange}
                  />
                </FormControl>
              </form>
            </Box>
            
            <Box className={classes.modalFooter}>
              <Button 
                color="primary" 
                onClick={handleClose}
                style={{ marginRight: 8 }}
              >
                Cancel
              </Button>
              <Button
                color="primary"
                variant="contained"
                onClick={handleSubmit}
                disabled={disabledBtn}
              >
                Send Report
              </Button>
            </Box>
          </Paper>
        </Fade>
      </Modal>
      
      <Snackbar 
        open={openSnackSuccess} 
        autoHideDuration={6000} 
        onClose={handleStackClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleStackClose} severity="success" variant="filled">
          Report will be sent to the specified email addresses shortly.
        </Alert>
      </Snackbar>
      
      <Snackbar 
        open={openSnackError} 
        autoHideDuration={6000} 
        onClose={handleStackClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleStackClose} severity="error" variant="filled">
          Something went wrong. Please try again.
        </Alert>
      </Snackbar>
    </div>
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
