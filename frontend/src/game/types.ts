import { z } from "zod";

export const ZGame = z.object({
    id: z.string(),
    rounds: z.number(),
    currentRound: z.number(),
    turn: z.union([z.string(), z.null()]),
    maxPlayers: z.number(),
    state: z.number(),
    lastStateChangeTime: z.coerce.date(),
    dateCreated: z.coerce.date()
})

export type Game = z.infer<typeof ZGame>;

const playerActiveStates = ['creating', 'active', 'disconnected'] as const
type PlayerActiveState = typeof playerActiveStates[number];

export function isPlayerActiveState(v: unknown): v is PlayerActiveState {
    return playerActiveStates.includes(v as PlayerActiveState)
}

export const ZPlayer = z.object({
    id: z.string(),
    name: z.string(),
    score: z.number(),
    game: z.string(),
    dateCreated: z.coerce.date(),
    guessedCorrect: z.boolean(),
    activeState:  z.enum(playerActiveStates)
})

export type Player = z.infer<typeof ZPlayer>


export type PlayerEditable = Pick<Player, 'name' | 'activeState' >