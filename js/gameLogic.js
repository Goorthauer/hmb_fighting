import { gameState, draggingCharacter, selectedAbility } from './state.js';
import { cellWidth, cellHeight } from './constants.js';
import { findCharacter } from './utils.js';
import { sendMessage } from './websocket.js';

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
    const currentChar = findCharacter(gameState.teams, gameState.currentTurn);
    if (!currentChar) return false;
    const dist = Math.abs(gridX - currentChar.position[0]) + Math.abs(gridY - currentChar.position[1]);
    return gameState.phase === 'move' &&
        gridX >= 0 && gridX < 16 &&
        gridY >= 0 && gridY < 9 &&
        gameState.board[gridX][gridY] === -1 &&
        dist <= currentChar.stamina;
}

export function canAttackOrUseAbility(gridX, gridY, myTeam) {
    if (gameState.phase !== 'action' || gameState.board[gridX][gridY] === -1) return false;
    const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
    return target && target.team !== myTeam;
}

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
