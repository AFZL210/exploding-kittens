import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { UserI } from "../../types/types";

const localUsername = localStorage.getItem("username");

const initialUserState: UserI = {
  isLoggedIn: localUsername === null ? false : true,
  username: localUsername ?? "",
  loading: false,
};

const userSlice = createSlice({
  name: "userSlice",
  initialState: initialUserState,
  reducers: {
    loginUser: (state, action: PayloadAction<UserI>) => {
      state.isLoggedIn = true;
      state.username = action.payload.username;
      state.loading = false;
    },
  },
});

export const { loginUser } = userSlice.actions;
export default userSlice.reducer;
