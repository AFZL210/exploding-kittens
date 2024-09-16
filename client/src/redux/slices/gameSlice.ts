import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import { GameStateI, openCardI } from "../../types/types";
import apiClient from "../../lib/apiClient";

const initialState: GameStateI = {
  cards: [],
  defuseCards: 0,
  isLost: false,
  isWon: false,
  isLoading: true,
  remainingCards: 0,
};

export const fetchGameState = createAsyncThunk(
  "fetchGameState",
  async (username: string) => {
    const res = await apiClient.get(`/getcards?username=${username}`);
    return res.data;
  }
);

export const openCard = createAsyncThunk(
  "openCard",
  async (payload: openCardI) => {
    const res = await apiClient.post(`/play?username=${payload.username}`, {
      index: payload.index,
    });
    return res.data;
  }
);

const gameSlice = createSlice({
  name: "GameState",
  initialState: initialState,
  reducers: {
    addItem: (_state, action: PayloadAction<GameStateI>) => {
      return action.payload;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(fetchGameState.pending, (state) => {
      state.isLoading = true;
    });
    builder.addCase(
      fetchGameState.fulfilled,
      (_state, action: PayloadAction<GameStateI>) => {
        return { ...action.payload, isLoading: false };
      }
    );

    builder.addCase(openCard.pending, (state) => {
      state.isLoading = true;
    });
    builder.addCase(openCard.fulfilled, (state, action: PayloadAction<GameStateI>) => {
      state.isLoading = false;
      state.isLost = action.payload.isLost;
      state.isWon = action.payload.isWon;
      state.remainingCards = action.payload.remainingCards;
      state.cards = action.payload.cards;
    });
  },
});

export const { addItem } = gameSlice.actions;
export default gameSlice.reducer;
