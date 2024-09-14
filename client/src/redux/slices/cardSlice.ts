import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { CardI } from "../../types/types";

const cardSlice = createSlice({
  name: "Cards",
  initialState: [] as CardI[],
  reducers: {
    addItem: (state, action: PayloadAction<CardI>) => {
      state.push(action.payload);
    },
  },
});

export const {addItem} = cardSlice.actions;
export default cardSlice.reducer;
