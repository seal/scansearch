import axios from "axios";
import displayAlert from "./displayAlert";

export default (err: unknown) => {
  if (axios.isAxiosError(err)) {
    displayAlert(
      err.response?.status === 500
        ? "An unknown error occurred"
        : err.response?.data.message || err.message,
      {
        type: "danger",
      }
    );
  } else {
    console.error(err);
  }
};
