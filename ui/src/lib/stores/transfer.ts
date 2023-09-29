import { writable } from "svelte/store";

export const from = writable('Checking');
export const to = writable('Savings');
export const amount = writable('0.00');