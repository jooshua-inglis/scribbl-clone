import { decodeAuthPayload } from '@/auth/token';
import { Game, Player } from '@/game/types';
import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

type GameStatePre = {
    authToken: string | null;
    playerId: string | null;
    gameId: string | null;
    game: null;
    players: Record<string, Player>;
    dataState: 'pending' | 'loading';
};

type GameStateLoaded = {
    authToken: string;
    playerId: string;
    gameId: string;
    game: Game;
    players: Record<string, Player>;
    dataState: 'loaded';
};

type GameState = GameStateLoaded | GameStatePre;

function initialState(): GameState {
    const authToken = localStorage.getItem('authToken');
    return {
        game: null,
        players: {},
        authToken: authToken,
        playerId: null,
        gameId: null,
        ...decodeAuthPayload(authToken),
        dataState: 'pending',
    };
}

export const gameSlice = createSlice({
    name: 'game',
    initialState,
    reducers: {
        setAuthToken: (state, { payload }: PayloadAction<GameState['authToken']>) => {
            if (payload === state.authToken) {
                console.log('they are the same');
                return;
            }
            if (!payload) {
                localStorage.removeItem('authToken');
                return;
            }
            state = { ...state, ...decodeAuthPayload(payload) };
            localStorage.setItem('authToken', payload);
            return state;
        },
        clearGame: () => {
            return initialState();
        },
        addPlayers: (state, { payload }: PayloadAction<Player[] | Player>) => {
            if (!Array.isArray(payload)) {
                payload = [payload];
            }
            for (const player of payload) {
                if (player.id in state.players) {
                    console.warn(
                        'Tried to add a player that already exists, use edit player instead'
                    );
                    continue;
                }
                state.players[player.id] = player;
            }
            return state;
        },
        removePlayers: (state, { payload }: PayloadAction<string[] | string>) => {
            if (typeof payload === 'string') {
                delete state.players[payload];
                return;
            }
            for (const playerId of payload) {
                delete state.players[playerId];
            }
            return state;
        },
        editPlayer: (
            state,
            { payload }: PayloadAction<{ id: string; updateSet: Partial<Player> }>
        ) => {
            if (!(payload.id in state.players)) {
                return;
            }
            state.players[payload.id] = { ...state.players[payload.id], ...payload.updateSet };
            return state;
        },
        setGame: (state, { payload }: PayloadAction<Game>) => {
            state.game = payload;
            return state;
        },
        setLoadingState: (state, { payload }: PayloadAction<GameState['dataState']>) => {
            state.dataState = payload;
            return state;
        },
        updateGame: (state, { payload }: PayloadAction<Partial<Game>>) => {
            if (!state.game) {
                return;
            }
            console.log("updating game")
            console.log(payload)
            state.game = { ...state.game, ...payload };
            return state;
        },
    },
});

export const {
    setAuthToken,
    addPlayers,
    editPlayer,
    removePlayers,
    setGame,
    clearGame,
    setLoadingState,
    updateGame,
} = gameSlice.actions;

export default gameSlice;
