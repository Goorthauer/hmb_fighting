import {
    gameState,
    draggingCharacter,
    selectedCharacter,
    selectedAbility,
    cellWidth,
    cellHeight,
    setDraggingCharacter,
    setMovePath,
    setSelectedAbility
} from './state.js';
import { animateMove } from './render.js';
import { sendMessage } from './websocket.js';

export function handleCanvasMouseDown(event, myTeam, setDraggingCharacter, setSelectedCharacter) {
    if (!gameState || !gameState.Teams || !Array.isArray(gameState.Teams)) return null;
    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const gridX = Math.floor(x / cellWidth);
    const gridY = Math.floor(y / cellHeight);

    if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10) {
        const charId = gameState.Board[gridX][gridY];
        if (charId !== -1) {
            const char = findCharacter(gameState.Teams, charId);
            if (char && char.Team === myTeam && char.ID === gameState.CurrentTurn) {
                setDraggingCharacter(char);
                return { x, y };
            } else {
                setSelectedCharacter(char);
            }
        }
    }
    return null;
}

export function handleCanvasMouseMove(event) {
    if (!draggingCharacter) return null;
    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    return { x, y };
}

export function handleCanvasMouseUp(event, clientID, setDraggingCharacter) {
    if (!draggingCharacter || !gameState) return;
    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const gridX = Math.floor(x / cellWidth);
    const gridY = Math.floor(y / cellHeight);

    if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10) {
        const startX = draggingCharacter.Position[0];
        const startY = draggingCharacter.Position[1];
        const path = calculatePath(startX, startY, gridX, gridY);
        setMovePath(path);

        if (gameState.Phase === 'move' && gameState.Board[gridX][gridY] === -1) {
            const charToMove = { ...draggingCharacter };
            animateMove(charToMove, path, () => {
                sendMessage(JSON.stringify({
                    type: 'move',
                    clientID: clientID,
                    characterID: charToMove.ID,
                    position: [gridX, gridY]
                }));
            });
        }
    }
    setDraggingCharacter(null);
}

function calculatePath(startX, startY, endX, endY) {
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

export function handleCanvasClick(event, myTeam, clientID, setSelectedAbility) {
    if (!gameState || !gameState.Teams || draggingCharacter) return;

    const rect = event.target.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const gridX = Math.floor(x / 50); // Используем 50, как в вашем коде
    const gridY = Math.floor(y / 50);

    if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10) {
        const charId = gameState.Board[gridX][gridY];
        if (charId !== -1) {
            const target = findCharacter(gameState.Teams, charId);
            const currentChar = findCharacter(gameState.Teams, gameState.CurrentTurn);
            if (
                gameState.Phase === 'action' &&
                target &&
                currentChar &&
                target.Team !== myTeam &&
                currentChar.Team === myTeam
            ) {
                const startX = currentChar.Position[0];
                const startY = currentChar.Position[1];
                const path = calculatePath(startX, startY, gridX, gridY);

                const charToAct = { ...currentChar };
                const targetToAct = { ...target };

                if (selectedAbility) {
                    animateMove(charToAct, path, () => {
                        sendMessage(JSON.stringify({
                            type: 'ability',
                            clientID: clientID,
                            characterID: charToAct.ID,
                            targetID: targetToAct.ID,
                            ability: selectedAbility.Name
                        }));
                        currentChar.Abilities = currentChar.Abilities.filter(a => a.Name !== selectedAbility.Name);
                        setSelectedAbility(null);
                    }, true);
                } else {
                    animateMove(charToAct, path, () => {
                        sendMessage(JSON.stringify({
                            type: 'attack',
                            clientID: clientID,
                            characterID: charToAct.ID,
                            targetID: targetToAct.ID,
                            ability: null
                        }));
                        setSelectedAbility(null);
                    }, true);
                }
            }
        }
    }
}

export function handleEndTurn(myTeam, clientID) {
    const currentChar = findCharacter(gameState.Teams, gameState.CurrentTurn);
    if (currentChar && currentChar.Team === myTeam) {
        const action = {
            type: 'end_turn',
            clientID: clientID
        };
        sendMessage(JSON.stringify(action));
    }
}

function findCharacter(teams, id) {
    if (!teams || !Array.isArray(teams)) return null;
    for (let team of teams) {
        for (let char of team.Characters) {
            if (char.ID === id) return char;
        }
    }
    return null;
}