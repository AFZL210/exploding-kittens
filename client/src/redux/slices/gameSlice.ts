import { createSlice, PayloadAction, createAsyncThunk } from "@reduxjs/toolkit";
import { GameStateI, openCardI } from "../../types/types";
import apiClient from "../../lib/apiClient";

const initialState: GameStateI = {
  cards: [],
  defuseCards: 0,
  gameOver: false,
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
    await apiClient.post(`/play?username=${payload.username}`, {
      index: payload.index,
    });
    return payload.index;
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
    builder.addCase(openCard.fulfilled, (state, action) => {
      state.isLoading = false;
      state.cards[action.payload].isFlipped = true;
      console.log(action.payload);
    });
  },
});

export const { addItem } = gameSlice.actions;
export default gameSlice.reducer;
