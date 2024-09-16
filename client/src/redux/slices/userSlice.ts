import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { UserI } from "../../types/types";

const initialUserState: UserI = {
  isLoggedIn: false,
  username: "",
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
