import { Player } from '@/game/types';
import classes from './playerList.module.css';

type PlayerListItemProps = {
    player: Player;
    playerId: string;
    drawerId: string | null;
    className?: string;
};
function PlayerListItem({ player, playerId, drawerId, className = '' }: PlayerListItemProps) {
    return (
        <li className={`${classes.playerListItem} ${className}`}>
            <div style={{ fontWeight: playerId === player.id ? 'bold' : 'inherit' }}>
                {player.name}
                {player.id === drawerId && '✏️'}
                {player.guessedCorrect && '✨'}
            </div>
            <div>Score: {player.score}</div>
        </li>
    );
}

type PlayerListProps = {
    players: Player[];
    playerId: string;
    drawerId: string | null;
};
export function PlayerList({ players, playerId, drawerId }: PlayerListProps) {
    return (
        <ul className={`${classes.playerList}`}>
            {players.map((p, i) => (
                <PlayerListItem
                    key={p.id}
                    playerId={playerId}
                    player={p}
                    drawerId={drawerId}
                    className={i % 2 === 0 ? classes.even : classes.odd}
                />
            ))}
        </ul>
    );
}
