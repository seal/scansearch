import axiosInstance from "@/axiosInstance";
import displayAlert from "@/helpers/displayAlert";
import handleAxiosError from "@/helpers/handleAxiosError";
import type { ILoginResponse, IUserDetails } from "@/types/responses.model";
import { useStorage } from "@vueuse/core";
import { startCase } from "lodash";
import { acceptHMRUpdate, defineStore } from "pinia";

export const useUserStore = defineStore("user", {
  state: () => ({
    user: useStorage("user", {} as IUserDetails | null),
    token: useStorage("token", null as string | null),
  }),
  getters: {
    loggedIn: (state) => !!state.token,
    name: ({ user }) =>
      user?.firstname
        ? startCase(`${user.firstname} ${user.lastname}`)
        : user?.email,
  },
  actions: {
    async login(userDetails: { email: string; password: string }) {
      try {
        const {
          data: { user, token },
        } = await axiosInstance.post<ILoginResponse>("/login", userDetails);

        this.user = user;
        this.token = token;

        axiosInstance.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${token}`;
      } catch (err) {
        handleAxiosError(err);
      }
    },
    async register(userDetails: {
      firstName: string;
      lastName: string;
      username: string;
      email: string;
      password: string;
    }) {
      try {
        const { data } = await axiosInstance.post("/register", userDetails);
        displayAlert(data.message, { type: "success" });
      } catch (err) {
        handleAxiosError(err);
      }
    },
    logout() {
      this.user = null;
      this.token = null;

      location.reload();
    },
    async getUser() {
      const { data } = await axiosInstance.get<IUserDetails>("/user");
      this.user = data;
    },
    async updateUser(userDetails: {
      firstname: string;
      lastname: string;
      username: string;
      email: string;
    }) {
      const { firstname, lastname, email, username } = userDetails;

      try {
        await axiosInstance.put("/user", {
          firstname,
          lastname,
          email,
          username,
        });

        await this.getUser();
        displayAlert("User Updated Successfully", {
          type: "success",
        });
      } catch (err) {
        handleAxiosError(err);
      }
    },
    async updatePassword(oldPassword: string, newPassword: string) {
      try {
        await axiosInstance.put("/user", {
          oldPassword,
          newPassword,
        });
      } catch (err) {
        handleAxiosError(err);
      }
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot));
}
