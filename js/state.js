export let gameState = null;
export let draggingCharacter = null;
export let dragOffsetX = 0;
export let dragOffsetY = 0;
export let selectedCharacter = null;
export let selectedAbility = null;
export let isSpectator = false;
export let currentRoom = localStorage.getItem('currentRoom') || null;
export let movePath = [];

export function setGameState(state) {
    gameState = state;
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