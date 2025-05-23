import React from "react";
import PropTypes from "prop-types";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import Link from "@mui/material/Link";

/**
 *
 * @param {string} tags - tags list
 */
const DialogTags = ({ tags }) => {
  const [open, setOpen] = React.useState(false);

  /**
   * Open dialog
   */
  const handleClickOpen = () => {
    setOpen(true);
  };

  /**
   * Close dialog
   */
  const handleClose = () => {
    setOpen(false);
  };

  return (
    <React.Fragment>
      <Link component="button" variant="body2" onClick={handleClickOpen}>
        Tags
      </Link>

      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="max-width-dialog-title"
      >
        <DialogTitle id="max-width-dialog-title">Tags</DialogTitle>
        <DialogContent>
          <DialogContentText></DialogContentText>
          <pre>{<div>{JSON.stringify(tags, null, 2)}</div>}</pre>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </React.Fragment>
  );
};

DialogTags.propTypes = {
  tags: PropTypes.object,
};

DialogTags.defaultProps = {};

export default DialogTags;
