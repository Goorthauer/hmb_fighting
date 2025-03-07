export const canvas = document.getElementById('gameCanvas');
export const ctx = canvas.getContext('2d');
export const cellWidth = 50;
export const cellHeight = 50;

export let gameState = null;
export let previousState = null;
export let draggingCharacter = null;
export let dragOffsetX = 0;
export let dragOffsetY = 0;
export let selectedCharacter = null;
export let selectedAbility = null;
export let isSpectator = false;
export let currentRoom = localStorage.getItem('currentRoom') || null;
export let movePath = []; // Добавляем movePath как изменяемую переменную состояния

export function setGameState(state) {
    gameState = state;
}

export function setDraggingCharacter(char) {
    draggingCharacter = char;
}

export function setDragOffset(x, y) {
    dragOffsetX = x;
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
    movePath = path; // Добавляем сеттер для movePath
}