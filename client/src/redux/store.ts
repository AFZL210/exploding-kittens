import { configureStore } from "@reduxjs/toolkit";
import cardSlice from "./slices/cardSlice";
import userSlice from "./slices/userSlice";

export const store = configureStore({
    reducer: {
        cards: cardSlice,
        user: userSlice
    }
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;