import { createToast, type ToastOptions } from "mosha-vue-toastify";

export default (message: string, options: ToastOptions = {}) => {
  createToast(message, {
    timeout: 2000,
    transition: "slide",
    ...options,
  });
};
