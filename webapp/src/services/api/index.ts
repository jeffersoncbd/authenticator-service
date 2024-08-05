import axios from "axios";
import { Credentials, LoginResponse } from "./interfaces";

const service = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

export const apiService = {
  login: async (credentials: Credentials) => {
    const response = await service.post<LoginResponse>("/login", {
      ...credentials,
      application: process.env.NEXT_PUBLIC_APP_ID,
    });
    return response.data;
  },
};
