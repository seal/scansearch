import axios from "axios";

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "/api",
});

const token = localStorage.getItem("token");

if (token) {
  instance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
}

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    // whatever you want to do with the error
    if (axios.isAxiosError(error)) {
      if (
        error.response?.status === 401 &&
        !error.request.responseURL.includes("login")
      ) {
        localStorage.clear();
        location.reload();
      }
    }

    throw error;
  }
);

export default instance;
