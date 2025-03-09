import { findCharacter } from './utils.js';

export let gameState = null;
export let draggingCharacter = null;
export let dragOffsetX = 0;
export let dragOffsetY = 0;
export let selectedCharacter = null;
export let selectedAbility = null;
export let isSpectator = false;
export let currentRoom = localStorage.getItem('currentRoom') || null;
export let movePath = [];
export let currentChar = null;

export function setGameState(state) {
    console.log('Setting game state:', state);
    gameState = state;
    if (state.currentTurn !== -1) {
        currentChar = findCharacter(state.teams, state.currentTurn);
        console.log('Updated currentChar:', currentChar);
    } else {
        currentChar = null;
    }
}

export function setDraggingCharacter(char) {
    draggingCharacter = char;
}

export function setDragOffsetX(x) {
    dragOffsetX = x;
}

export function setDragOffsetY(y) {
    dragOffsetY = y;
}

export function setSelectedCharacter(char) {
    selectedCharacter = char;
}

export function setSelectedAbility(ability) {
    selectedAbility = ability;
}

export function setIsSpectator(spectator) {
    isSpectator = spectator;
}

export function setCurrentRoom(room) {
    currentRoom = room;
    if (room) {
        localStorage.setItem('currentRoom', room);
    } else {
        localStorage.removeItem('currentRoom');
    }
}

export function setMovePath(path) {
    movePath = path;
}