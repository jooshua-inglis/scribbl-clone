import { apiMakeGuess } from '@/game/api';
import { useAppSelector } from '@/store/store';
import { KeyboardEvent, useCallback, useState } from 'react';

export type Guess = {
    playerName: string;
    guess?: string;
    isCorrect: boolean;
}

type GuesserProps = {
    guesses: Guess[]
};
export function Guesser({ guesses }: GuesserProps) {
    const gameId = useAppSelector((state) => state.game.gameId);
    const [guess, setGuess] = useState('');

    const onMakeGuess = useCallback(() => {
        if (!gameId) {
            console.trace('GameId not defined');
            return;
        }
        apiMakeGuess(gameId, { guess });
    }, [guess, gameId]);

    const onEnter = (e: KeyboardEvent) => {
        if (e.key === 'Enter') {
            onMakeGuess();
        }
    };

    return (
        <div
            style={{
                height: '100%',
                display: 'flex',
                gap: '10px',
                flexDirection: 'column-reverse',
                overflow: 'scroll',
            }}>
            <div style={{ display: 'flex' }}>
                <input
                    value={guess}
                    onInput={(e) => setGuess(e.currentTarget.value)}
                    onKeyDown={onEnter}
                    style={{
                        backgroundColor: '#EEFBFA',
                        margin: '10px',
                        borderColor: '#FF9F1C',
                        height: '2rem',
                        borderRadius: '1rem',
                        width: '100%',
                        borderWidth: '4px',
                        borderStyle: 'solid',
                        paddingLeft: '1rem',
                    }}
                />
            </div>
            {[...guesses].reverse().map(guess => {
                return <Bubble {...guess} />
            })}
        </div>
    );
}

type BubbleProps = {
    playerName: string;
    isCorrect: boolean;
    guess?: string;
};
function Bubble({ playerName, isCorrect, guess }: BubbleProps) {
    return (
        <div
            style={{
                backgroundColor: isCorrect ? '#7be395' : '#FFBF69',
                marginInline: '10px',
                padding: '10px',
                borderRadius: '10px',
            }}>
            <span style={{ fontWeight: 'bold' }}>{playerName}: </span>{isCorrect ? "Guessed Correct" : guess ?? ""}
        </div>
    );
}
