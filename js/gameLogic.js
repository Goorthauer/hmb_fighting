import { gameState, draggingCharacter, selectedAbility } from './state.js';
import { cellWidth, cellHeight } from './constants.js';
import { findCharacter } from './utils.js';

export function calculatePath(startX, startY, endX, endY) {
    const path = [];
    let currentX = startX;
    let currentY = startY;

    while (currentX !== endX || currentY !== endY) {
        path.push({ x: currentX, y: currentY });
        if (currentX < endX) currentX++;
        else if (currentX > endX) currentX--;
        if (currentY < endY) currentY++;
        else if (currentY > endY) currentY--;
    }
    path.push({ x: endX, y: endY });
    return path;
}

export function getGridPosition(event) {
    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const gridX = Math.floor(x / cellWidth);
    const gridY = Math.floor(y / cellHeight);
    return { gridX, gridY, x, y };
}

export function canMove(gridX, gridY) {
    return gameState.phase === 'move' && gameState.board[gridX][gridY] === -1;
}

export function canAttackOrUseAbility(gridX, gridY, myTeam) {
    if (gameState.phase !== 'action' || gameState.board[gridX][gridY] === -1) return false;
    const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
    return target && target.team !== myTeam;
}

export function isWithinAttackRange(attacker, targetX, targetY, weaponsConfig, ability = null) {
    const startX = attacker.position[0];
    const startY = attacker.position[1];

    // Если используется способность, проверяем её диапазон
    if (ability && ability.range !== undefined) {
        const dist = Math.max(Math.abs(targetX - startX), Math.abs(targetY - startY));
        return dist <= ability.range; // Проверяем, что цель в пределах диапазона способности
    }

    // Если способность не используется, проверяем диапазон оружия
    const weapon = weaponsConfig[attacker.weapon];
    const weaponRange = weapon ? weapon.range : 1;
    const isTwoHanded = weapon ? weapon.isTwoHanded : false;
    const dist = Math.max(Math.abs(targetX - startX), Math.abs(targetY - startY));

    // Проверяем диапазон оружия
    return (isTwoHanded && dist === weaponRange) || (!isTwoHanded && dist <= weaponRange);
}