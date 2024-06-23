import { WEBSOCKET_HOST } from '@/constants';
import { PubSub } from '@/eventListener/eventListener';
import { getGame, getPlayers } from '@/game/api';
import {
    addPlayers,
    clearGame,
    editPlayer,
    setGame,
    setLoadingState,
    updateGame,
} from '@/game/game.slice';
import { Game, isPlayerActiveState, Player, ZGame, ZPlayer } from '@/game/types';
import { useAppDispatch, useAppSelector } from '@/store/store';
import { Line } from '@/ui/drawingCanvas';
import { useEffect } from 'react';

export const GameEventTypes = {
    GUESS_OCCURRED: 0,
    STATE_CHANGE: 1,
    SCORE_UPDATE: 2,
    GAME_UPDATE: 3,
    PLAYER_UPDATE: 4,
    PLAYER_ADDED: 5,
    DRAWING: 6,
} as const;

export type GameEventType = (typeof GameEventTypes)[keyof typeof GameEventTypes];

export const GameStates = {
    STATE_WAITING_FOR_PLAYERS: 0,
    STATE_END: 1,
    STATE_DRAWING: 2,
    STATE_SELECTING_WORD: 3,
} as const;

export type GameStateType = (typeof GameStates)[keyof typeof GameStates];

// ====== EVENT TYPES =========
type UpdateGameEvent = {
    type: (typeof GameEventTypes)['GAME_UPDATE'];
    payload: Partial<Game>;
};

type UpdatePlayerEvent = {
    type: (typeof GameEventTypes)['PLAYER_UPDATE'];
    payload: {
        PlayerId: string;
        Updates: Partial<
            {
                // TODO Scribbl-1 remove caps from schema
                Name: string;
                Score: number;
                GuessedCorrect: boolean;
                ActiveState: 'creating' | 'active' | 'disconnected';
            } & Player
        >;
    };
};

type PlayerAddedEvent = {
    type: (typeof GameEventTypes)['PLAYER_ADDED'];
    payload: Player;
};

type GameUpdateEvent = {
    type: (typeof GameEventTypes)['SCORE_UPDATE'];
    payload: Record<string, number>;
};

type DrawingEvent = {
    type: (typeof GameEventTypes)['DRAWING'];
    payload: {
        line: Line;
        index: number;
    };
};

type GuessEvent = {
    type: (typeof GameEventTypes)['GUESS_OCCURRED'];
    payload: {
        guess: string;
        playerId: string;
        isCorrect: boolean;
    };
};

type GameWebsocketEvent =
    | UpdatePlayerEvent
    | UpdateGameEvent
    | PlayerAddedEvent
    | GameUpdateEvent
    | DrawingEvent
    | GuessEvent;

function createConnection(url: string) {
    const ws = new WebSocket(url);
    return new Promise<WebSocket>((resolve, reject) => {
        ws.onerror = (error) => {
            ws.onerror = () => {};
            ws.onopen = () => {};
            reject(error);
        };
        ws.onopen = () => {
            ws.onerror = () => {};
            ws.onopen = () => {};
            resolve(ws);
        };
    });
}

export class GameWebsocket {
    private static instance: GameWebsocket;
    private pubSub: PubSub<GameWebsocketEvent> = new PubSub();
    private gameId: string | undefined = undefined;
    private ws: WebSocket | undefined = undefined;

    addEventListener = this.pubSub.addEventListener.bind(this.pubSub);
    dispatch = this.pubSub.dispatchEvent.bind(this.pubSub);

    private constructor() {}

    public sendMessage(event: GameWebsocketEvent) {
        const payload = {
            EventType: event.type,
            EventPayload: event.payload,
        } as const;
        this.ws?.send(JSON.stringify(payload));
    }

    public async setPlayerId(gameId: string) {
        if (!gameId || gameId === this.gameId) {
            return;
        }
        this.gameId = gameId;

        this.ws?.close();

        this.ws = await createConnection(`${WEBSOCKET_HOST}/game_connection/${this.gameId}`);

        this.ws.onerror = () => {
            console.error('Websockets encountered an error');
        };

        this.ws.onmessage = (ev) => {
            const { data } = ev;
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const parsedData: any = JSON.parse(data);

            if (parsedData.EventType === GameEventTypes.PLAYER_ADDED) {
                const parsedPlayer = ZPlayer.safeParse(parsedData.EventPayload);
                if (!parsedPlayer.success) {
                    console.trace('Invalid player');
                    return;
                }
                parsedData.EventPayload = parsedPlayer.data;
            }
            if (parsedData.EventType === GameEventTypes.GAME_UPDATE) {
                console.log(parsedData);
                const parsedGame = ZGame.partial().safeParse(parsedData.EventPayload);
                if (!parsedGame.success) {
                    console.trace('Invalid game');
                    return;
                }
                parsedData.EventPayload = parsedGame.data;
            }
            this.pubSub.dispatchEvent({
                type: parsedData.EventType,
                payload: parsedData.EventPayload,
            });
        };
    }

    public static getInstance(): GameWebsocket {
        if (!GameWebsocket.instance) {
            GameWebsocket.instance = new GameWebsocket();
        }

        return GameWebsocket.instance;
    }
}
const gameWebSocket = GameWebsocket.getInstance();

export function useInitialiseData() {
    const playerId = useAppSelector((state) => state.game.playerId);
    const gameId = useAppSelector((state) => state.game.gameId);
    const dispatch = useAppDispatch();

    useEffect(() => {
        if (!playerId || !gameId) {
            return;
        }
        dispatch(setLoadingState('loading'));

        gameWebSocket.setPlayerId(playerId);
        dispatch(clearGame());
        Promise.all([
            getGame(gameId).then((game) => dispatch(setGame(game))),
            getPlayers(gameId).then((players) => dispatch(addPlayers(players))),
        ]).then(() => {
            dispatch(setLoadingState('loaded'));
        });
    }, [playerId, gameId, dispatch]);
}

export function useGameWebsocket() {
    const dispatch = useAppDispatch();
    useEffect(() => {
        const updatedPlayerCleanup = gameWebSocket.addEventListener(
            GameEventTypes.PLAYER_UPDATE,
            (event) => {
                const updateSet: Partial<Player> = {};
                if (isPlayerActiveState(event.payload.Updates.ActiveState)) {
                    updateSet.activeState = event.payload.Updates.ActiveState;
                }
                // TODO Scribbl-1 remove caps from schema
                if (event.payload.Updates.Name) {
                    updateSet.name = event.payload.Updates.Name;
                }
                if (event.payload.Updates.Score) {
                    updateSet.score = event.payload.Updates.Score;
                }
                if ('GuessedCorrect' in event.payload.Updates) {
                    updateSet.guessedCorrect = event.payload.Updates.GuessedCorrect;
                }
                console.log(updateSet);
                dispatch(editPlayer({ id: event.payload.PlayerId, updateSet }));
            }
        );

        const addedPlayerCleanup = gameWebSocket.addEventListener(
            GameEventTypes.PLAYER_ADDED,
            (event) => {
                dispatch(addPlayers(event.payload));
            }
        );

        const gameUpdatedCleanup = gameWebSocket.addEventListener(
            GameEventTypes.GAME_UPDATE,
            (event) => {
                dispatch(updateGame(event.payload));
            }
        );

        const scoreUpdateCleanup = gameWebSocket.addEventListener(
            GameEventTypes.SCORE_UPDATE,
            (event) => {
                for (const [key, score] of Object.entries(event.payload)) {
                    dispatch(editPlayer({ id: key, updateSet: { score } }));
                }
            }
        );

        return () => {
            updatedPlayerCleanup();
            addedPlayerCleanup();
            gameUpdatedCleanup();
            scoreUpdateCleanup();
        };
    }, [dispatch]);
}
