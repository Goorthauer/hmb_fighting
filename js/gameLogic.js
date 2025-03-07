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
    return gameState.Phase === 'move' && gameState.Board[gridX][gridY] === -1;
}

export function canAttackOrUseAbility(gridX, gridY, myTeam) {
    if (gameState.Phase !== 'action' || gameState.Board[gridX][gridY] === -1) return false;
    const target = findCharacter(gameState.Teams, gameState.Board[gridX][gridY]);
    return target && target.Team !== myTeam;
}

export function isWithinAttackRange(attacker, targetX, targetY) {
    const startX = attacker.Position[0];
    const startY = attacker.Position[1];
    const isTwoHanded = (attacker.Weapon === 'two_handed_halberd' || attacker.Weapon === 'two_handed_sword');
    const weaponRange = isTwoHanded ? 2 : 1;
    const dist = Math.max(Math.abs(targetX - startX), Math.abs(targetY - startY));
    return (isTwoHanded && dist === weaponRange) || (!isTwoHanded && dist <= weaponRange);
}