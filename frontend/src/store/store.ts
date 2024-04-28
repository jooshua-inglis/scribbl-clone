import gameSlice, { GameState } from '@/game/game.slice';
import { configureStore } from '@reduxjs/toolkit';
import { useRef } from 'react';
import { useDispatch, UseSelector, useSelector, useStore } from 'react-redux';

export const store = configureStore({
    reducer: {
        [gameSlice.reducerPath]: gameSlice.reducer,
    },
});

export type AppStore = typeof store;
export type RootState = ReturnType<AppStore['getState']>;
export type AppDispatch = AppStore['dispatch'];

export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
export const useAppStore = useStore.withTypes<AppStore>();
