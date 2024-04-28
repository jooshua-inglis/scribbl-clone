import appAxios from '@/auth/appAxios';
import { Game, Player, ZGame, ZPlayer } from '@/game/types';

type SuccessResponse = {
    message: "success"
}

export async function apiJoinGame(
    gameId: string,
    name: string,
): Promise<{ player: Player; token: string }> {
    const response = await appAxios.post<{ player: Player; token: string }>(
        `/game/${gameId}/join`, { name },
    );

    return {
        player: response.data.player,
        token: response.data.token,
    };
}

export async function apiCreateGame(): Promise<Game> {
    const response = await appAxios.post<Game>('/game/new');
    return response.data;
}

export async function apiUpdatePlayer(
    playerId: string,
    updateSet: Partial<Pick<Player, 'activeState' | 'name'>>,
): Promise<Player> {
    // TODO SCRIBBL-1
    const body = {
        ActiveState: updateSet.activeState,
        Name: updateSet.name,
    };
    const response = await appAxios.post<Player>(`/player/${playerId}`, body);
    return response.data;
}

export async function getGame(gameId: string) {
    const response = await appAxios.get<Game>(`/game/${gameId}`);
    return ZGame.parse(response.data)
}

export async function getPlayers(gameId: string) {
    const response = await appAxios.get<Player[]>(`/game/${gameId}/players`)
    return response.data.map(player => (ZPlayer.parse(player)))
}

export async function apiMakeGuess(gameId: string, { guess }: { guess: string }) {
    const response = await appAxios.post<SuccessResponse>(
        `/game/${gameId}/guess`,
        { guess }
    )
    return response.data
}

export async function apiStartGame(gameId: string) {
    const response = await appAxios.post<SuccessResponse>(`/game/${gameId}/start`)
    return response.data
}

export async function apiSelectWord(gameId: string, word: string) {
    const response = await appAxios.post<SuccessResponse>(`/game/${gameId}/select_word`, {
        word
    })
    return response.data
}