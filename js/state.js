import { findCharacter } from './utils.js';

// Инициализация состояния
let gameState = null;
let draggingCharacter = null;
let dragOffsetX = 0;
let dragOffsetY = 0;
let selectedCharacter = null;
let selectedAbility = null;
let isSpectator = false;
let currentRoom = localStorage.getItem('currentRoom') || null;
let movePath = [];
let currentChar = null;

// Установка состояния игры
export function setGameState(state) {
    console.log('Setting game state:', state);
    gameState = state;
    updateCurrentChar(state);
}

// Обновление текущего персонажа
function updateCurrentChar(state) {
    if (state.currentTurn !== -1) {
        currentChar = findCharacter(state.teams, state.currentTurn);
        console.log('Updated currentChar:', currentChar);
    } else {
        currentChar = null;
    }
}

// Установка перетаскиваемого персонажа
export function setDraggingCharacter(char) {
    draggingCharacter = char;
}

// Установка смещения по оси X
export function setDragOffsetX(x) {
    dragOffsetX = x;
}

// Установка смещения по оси Y
export function setDragOffsetY(y) {
    dragOffsetY = y;
}

// Установка выбранного персонажа
export function setSelectedCharacter(char) {
    selectedCharacter = char;
}

// Установка выбранной способности
export function setSelectedAbility(ability) {
    selectedAbility = ability;
}

// Установка режима наблюдателя
export function setIsSpectator(spectator) {
    isSpectator = spectator;
}

// Установка текущей комнаты
export function setCurrentRoom(room) {
    currentRoom = room;
    if (room) {
        localStorage.setItem('currentRoom', room);
    } else {
        localStorage.removeItem('currentRoom');
    }
}

// Установка пути перемещения
export function setMovePath(path) {
    movePath = path;
}

// Экспорт состояния
export {
    gameState,
    draggingCharacter,
    dragOffsetX,
    dragOffsetY,
    selectedCharacter,
    selectedAbility,
    isSpectator,
    currentRoom,
    movePath,
    currentChar
};