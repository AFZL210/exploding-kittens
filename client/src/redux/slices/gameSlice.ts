import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import { GameStateI } from "../../types/types";
import apiClient from "../../lib/apiClient";

const initialState: GameStateI = {
  cards: [],
  defuseCards: 0,
  gameOver: false,
  isWon: false,
  isLoading: true,
  remainingCards: 0
}

export const fetchGameState = createAsyncThunk('fetchGameState', async () => {
  const res = await apiClient.get(`/gamestate?username=a`);
  return res.data;
});

const gameSlice = createSlice({
  name: "GameState",
  initialState: initialState,
  reducers: {
    addItem: (state, action: PayloadAction<GameStateI>) => {
      state = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchGameState.pending, (state) => {
      state.isLoading = true;
    })
    builder.addCase(fetchGameState.fulfilled, (state, action: PayloadAction<GameStateI>) => {
      state = { ...action.payload, isLoading: false };
    })
  }
});

export const { addItem } = gameSlice.actions;
export default gameSlice.reducer;
