import { gameState, draggingCharacter, selectedAbility } from './state.js';
import { cellWidth, cellHeight } from './constants.js';
import { findCharacter, addLogEntry } from './utils.js';
import { sendMessage } from './websocket.js';

// Узел для алгоритма A*
class Node {
    constructor(x, y, g, h, parent = null) {
        this.x = x;
        this.y = y;
        this.g = g; // Стоимость пути от старта
        this.h = h; // Эвристика до цели
        this.f = g + h; // Общая стоимость
        this.parent = parent;
    }
}

// Эвристика для A*
function heuristic(x1, y1, x2, y2) {
    return Math.abs(x1 - x2) + Math.abs(y1 - y2);
}

// Поиск пути с использованием алгоритма A*
export function findPath(startX, startY, endX, endY, stamina) {
    const openList = [];
    const closedList = new Set();
    const startNode = new Node(startX, startY, 0, heuristic(startX, startY, endX, endY));
    openList.push(startNode);

    while (openList.length > 0) {
        openList.sort((a, b) => a.f - b.f);
        const current = openList.shift();
        const key = `${current.x},${current.y}`;
        if (closedList.has(key)) continue;

        closedList.add(key);

        if (current.x === endX && current.y === endY) {
            return reconstructPath(current, stamina);
        }

        const neighbors = getNeighbors(current, endX, endY, stamina);
        for (const neighbor of neighbors) {
            const neighborKey = `${neighbor.x},${neighbor.y}`;
            if (!closedList.has(neighborKey)) {
                openList.push(neighbor);
            }
        }
    }
    return []; // Путь не найден
}

// Восстановление пути
function reconstructPath(node, stamina) {
    const path = [];
    let current = node;
    while (current) {
        path.unshift({ x: current.x, y: current.y });
        current = current.parent;
    }
    return path.length - 1 <= stamina ? path : [];
}

// Получение соседних узлов
function getNeighbors(current, endX, endY, stamina) {
    const neighbors = [
        { dx: 0, dy: -1 }, { dx: 0, dy: 1 }, { dx: -1, dy: 0 }, { dx: 1, dy: 0 }
    ];
    const result = [];

    for (const { dx, dy } of neighbors) {
        const newX = current.x + dx;
        const newY = current.y + dy;

        if (isValidPosition(newX, newY, endX, endY)) {
            const g = current.g + 1;
            if (g <= stamina) {
                const h = heuristic(newX, newY, endX, endY);
                result.push(new Node(newX, newY, g, h, current));
            }
        }
    }
    return result;
}

// Проверка валидности позиции
function isValidPosition(x, y, endX, endY) {
    if (x < 0 || x >= 16 || y < 0 || y >= 9) return false;
    return gameState.board[x][y] === -1 || (x === endX && y === endY);
}

// Получение позиции на сетке
export function getGridPosition(event) {
    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const gridX = Math.floor(x / cellWidth);
    const gridY = Math.floor(y / cellHeight);
    return { gridX, gridY, x, y };
}

// Проверка возможности перемещения
export function canMove(gridX, gridY) {
    const currentChar = findCharacter(gameState.teams, gameState.currentTurn);
    if (!currentChar) return false;

    const startX = currentChar.position[0];
    const startY = currentChar.position[1];

    if (gameState.phase !== 'move' ||
        gridX < 0 || gridX >= 16 ||
        gridY < 0 || gridY >= 9 ||
        gameState.board[gridX][gridY] !== -1) {
        return false;
    }

    const path = findPath(startX, startY, gridX, gridY, currentChar.stamina);
    return path.length > 0;
}

// Проверка возможности атаки или использования способности
export function canAttackOrUseAbility(gridX, gridY, myTeam) {
    if (gameState.board[gridX][gridY] === -1) return false;
    const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
    return target && target.team !== myTeam;
}

// Проверка нахождения цели в зоне атаки
export function isWithinAttackRange(attacker, targetX, targetY, weaponsConfig, ability = null) {
    const startX = attacker.position[0];
    const startY = attacker.position[1];
    const dx = Math.abs(targetX - startX);
    const dy = Math.abs(targetY - startY);

    if (ability && ability.range !== undefined) {
        return dx <= ability.range && dy <= ability.range;
    }

    const weapon = weaponsConfig[attacker.weapon];
    const weaponRange = weapon ? weapon.range : 1;
    return dx <= weaponRange && dy <= weaponRange;
}

// Убийство непоставленных персонажей после фазы setup
export function killUnplacedCharacters(myTeam) {
    if (!gameState || !gameState.teams || gameState.phase !== 'setup') return;

    const team = gameState.teams[myTeam];
    if (!team || !team.characters) return;

    team.characters.forEach(char => {
        if (char.position[0] === -1 && char.position[1] === -1) {
            char.hp = 0; // Устанавливаем здоровье в 0
            addLogEntry(`${char.name} не был размещён, поэтому убран из пула`);
        }
    });
}