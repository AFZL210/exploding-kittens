import axios from "axios";
import Cookies from "js-cookie";

const SERVER_URL = import.meta.env.VITE_SERVER_URL;

const apiClient = axios.create({
  baseURL: SERVER_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

apiClient.interceptors.request.use(
  (config) => {
    const token = Cookies.get("token");

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('username');
      window.location.href = '/';
      Cookies.remove('token');
    }

    return Promise.reject(error);
  }
);

export default apiClient;
