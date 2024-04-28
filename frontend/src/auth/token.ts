import { jwtDecode } from "jwt-decode";


export type PlayerClaim = {
    playerId: string;
    gameId: string;
}

export function decodeAuthPayload(authToken: string | null | undefined): PlayerClaim | null {
    if (!authToken) {
        return null;
    }
    try {
        const claim = jwtDecode<{ playerId: string; gameId: string }>(authToken);
        if (typeof claim === 'string' || claim === null) {
            console.error('Invalid claim type');
            return null;
        }
        return {
            playerId: claim.playerId,
            gameId: claim.gameId,
        };
    } catch (err) {
        console.error(err);
        return null;
    }
}
