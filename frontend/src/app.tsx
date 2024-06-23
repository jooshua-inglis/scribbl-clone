import { apiSelectWord, apiStartGame } from '@/game/api';
import { clearGame, setAuthToken } from '@/game/game.slice';
import {
    GameEventTypes,
    GameStates,
    GameWebsocket,
    useGameWebsocket,
    useInitialiseData,
} from '@/game/gameWebSocket';
import Home from '@/home';
import { useAppDispatch, useAppSelector } from '@/store/store';
import { DrawingCanvas, Line } from '@/ui/drawingCanvas';
import { Guess, Guesser } from '@/ui/guesser';
import { PlayerList } from '@/ui/playerList';
import { useCallback, useEffect, useState } from 'react';
import classes from './app.module.css';
import CountDown from '@/ui/timer';
import { addDuration } from '@/utils/date';

const gameWebsocket = GameWebsocket.getInstance();

export default function App() {
    useGameWebsocket();
    useInitialiseData();

    const dispatch = useAppDispatch();
    const game = useAppSelector((state) => state.game.game);
    if (!game) {
        return <Home />;
    }

    const Page = game.state === GameStates.STATE_WAITING_FOR_PLAYERS ? Lobby : Main;

    return (
        <div>
            <button
                onClick={() => {
                    dispatch(setAuthToken(null));
                    dispatch(clearGame());
                }}>
                Exit Game
            </button>
            <Page />
            <div
                style={{
                    backgroundColor: '#7be395',
                    width: 'fit-content',
                    margin: '10px auto',
                    padding: '5px',
                    borderRadius: '10px',
                }}>
                GameId: {game.id}
            </div>
        </div>
    );
}

function Lobby() {
    const gameId = useAppSelector((state) => state.game.gameId);
    if (!gameId) {
        return null;
    }
    return (
        <>
            <button onClick={() => apiStartGame(gameId)}>Start Game</button>
        </>
    );
}

function Main() {
    const gameState = useAppSelector((state) => state.game);
    const [lines, setLines] = useState<Line[]>([]);
    const [guesses, setGuesses] = useState<Guess[]>([])

    const onDraw = useCallback((lines: Line[]) => {
        const i = lines.length - 1;
        gameWebsocket.sendMessage({
            type: GameEventTypes.DRAWING,
            payload: {
                line: lines[i],
                index: i,
            },
        });
        setLines(lines);
    }, []);

    useEffect(() => {
        if (gameState.dataState !== 'loaded') {
            return;
        }
        if (gameState.game.turn === gameState.playerId) {
            return;
        }
        const drawingCleanup = gameWebsocket.addEventListener(
            GameEventTypes.DRAWING,
            ({ payload }) => {
                setLines((preState) => {
                    const newState = [...preState];
                    newState[payload.index] = payload.line;
                    return newState;
                });
            }
        );
        return () => {
            drawingCleanup();
        };
    }, [gameState.dataState, gameState.game?.turn, gameState.playerId]);

    // Clear screen when the state isn't in the drawing phase.
    useEffect(() => {
        if (gameState.dataState !== 'loaded') {
            return;
        }
        if (gameState.game.state !== GameStates.STATE_DRAWING) {
            setLines([]);
        }
    }, [gameState.dataState, gameState.game?.state]);

    useEffect(() => {
        const guessCleanup = gameWebsocket.addEventListener(GameEventTypes.GUESS_OCCURRED, ({payload}) => {
            setGuesses((prevValue) => {
                const player = gameState.players[payload.playerId]
                if (!player) {
                    return prevValue
                }
                return [...prevValue, {
                    isCorrect: payload.isCorrect,
                    playerName: player.name,
                    guess: payload.guess
                }]
            })
        })

        return () => {
            guessCleanup()
        }
    }, [gameState.players])

    if (gameState.dataState !== 'loaded') {
        return <div>Loading</div>;
    }

    const { game, playerId, players } = gameState;
    console.log(gameState.players)

    return (
        <div className={classes['game-container']}>
            <dialog
                open={game.state === GameStates.STATE_SELECTING_WORD && game.turn === playerId}
                style={{
                    position: 'fixed',
                    top: '50%',
                    transform: 'translate(0,-50%)',
                }}>
                <h1>Pick a word</h1>
                {['a', 'b', 'c'].map((word) => (
                    <button key={word} onClick={() => apiSelectWord(game.id, word)}>
                        {word}
                    </button>
                ))}
            </dialog>
            <div className={classes['banner']}>
                Time remaining: 10s Round: {game?.currentRound} / {game?.rounds}
                <CountDown time={addDuration(gameState.game.lastStateChangeTime, 30000)} />
            </div>
            <div className={classes['leaderboard']}>
                <PlayerList
                    players={Object.values(players).sort(
                        (a, b) => a.dateCreated.getTime() - b.dateCreated.getTime()
                    )}
                    playerId={playerId}
                    drawerId={game.turn}
                />
            </div>
            <div className={classes['game']}>
                <DrawingCanvas
                    lines={lines}
                    width={780}
                    height={780}
                    onDraw={onDraw}
                    drawable={playerId === game.turn && game.state === GameStates.STATE_DRAWING}
                />
            </div>
            <div className={classes['guessing']}>
                <Guesser guesses={guesses} />
            </div>
        </div>
    );
}
