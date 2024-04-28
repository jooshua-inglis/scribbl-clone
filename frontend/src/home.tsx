import { apiCreateGame, apiJoinGame } from '@/game/api';
import { setAuthToken } from '@/game/game.slice';
import { useCallback, useState } from 'react';
import { useDispatch } from 'react-redux';

export default function Home() {
    const [gameId, setGameId] = useState('');
    const [playerName, setPlayerName] = useState('');
    const dispatch = useDispatch();

    const joinGameHandler = useCallback(async () => {
        if(playerName.length === 0) {
            return
        }
        try {
            const { token } = await apiJoinGame(gameId, playerName);
            dispatch(setAuthToken(token));
        } catch (e) {
            console.error(e);
        }
    }, [gameId, dispatch, playerName]);

    const createGameHandler = useCallback(async () => {
        if(playerName.length === 0) {
            return
        }
        try {
            const game = await apiCreateGame();
            const { token } = await apiJoinGame(game.id, playerName);
            dispatch(setAuthToken(token));
        } catch (e) {
            console.error(e);
        }
    }, [dispatch, playerName]);

    return (
        <div className="p-2">
            <div>
                <label htmlFor="nameInput">Name</label>
                <input id='nameInput' value={playerName} onInput={e => setPlayerName(e.currentTarget.value)}/>
            </div>
            <div>
                <button onClick={createGameHandler}>Create Game</button>
            </div>
            <div>
                <input onInput={(e) => setGameId(e.currentTarget.value)} type="text" />
                <button onClick={joinGameHandler}>Join Game</button>
            </div>
        </div>
    );
}
