import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { UserI } from "../../types/types";

const initialUserState: UserI = {
  isLoggedIn: false,
  username: "",
};

const userSlice = createSlice({
  name: "userSlice",
  initialState: initialUserState,
  reducers: {
    loginUser: (state, action: PayloadAction<UserI>) => {
      state.isLoggedIn = true;
      state.username = action.payload.username;
    },
    logoutUser: (state) => {
      state = initialUserState;
    },
  },
});

export const { loginUser, logoutUser } = userSlice.actions;
export default userSlice.reducer;
