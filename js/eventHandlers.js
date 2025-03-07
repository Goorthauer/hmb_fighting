import { gameState, draggingCharacter, selectedCharacter, selectedAbility, setDraggingCharacter, setSelectedCharacter, setSelectedAbility, setMovePath, setDragOffsetX, setDragOffsetY } from './state.js';
import { connectWebSocket, sendMessage } from './websocket.js';
import { animateMove, drawBoard } from './renderCanvas.js';
import { calculatePath, getGridPosition, canMove, canAttackOrUseAbility, isWithinAttackRange } from './gameLogic.js';
import { findCharacter } from './utils.js';
import { canvas } from './constants.js';

export function setupEventListeners(myTeam) {
    if (!myTeam && myTeam !== 0) {
        console.error('myTeam is not initialized yet, delaying event listeners setup');
        return;
    }
    console.log('Setting up event listeners with myTeam:', myTeam);

    canvas.addEventListener('mousedown', (e) => handleMouseDown(e, myTeam));
    canvas.addEventListener('mousemove', (e) => handleMouseMove(e));
    canvas.addEventListener('mouseup', (e) => handleMouseUp(e, myTeam));
    canvas.addEventListener('click', (e) => handleClick(e, myTeam));
    document.getElementById('endTurnBtn').addEventListener('click', () => handleEndTurn(myTeam));
}

function handleMouseDown(event, myTeam) {
    if (!gameState || !gameState.Teams || !Array.isArray(gameState.Teams)) {
        console.log('Game state not fully initialized yet');
        return;
    }

    const { gridX, gridY, x, y } = getGridPosition(event);
    if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10) {
        const charId = gameState.Board[gridX][gridY];
        if (charId !== -1) {
            const char = findCharacter(gameState.Teams, charId);
            if (char && char.Team === myTeam && char.ID === gameState.CurrentTurn) {
                console.log('MouseDown: Setting draggingCharacter:', char);
                setDraggingCharacter(char);
                setDragOffsetX(x);
                setDragOffsetY(y);
            } else {
                setSelectedCharacter(char);
            }
        }
    }
}

function handleMouseMove(event) {
    if (!draggingCharacter) return;
    const { x, y } = getGridPosition(event);
    setDragOffsetX(x);
    setDragOffsetY(y);
    drawBoard(gameState);
}

function handleMouseUp(event, myTeam) {
    console.log('MouseUp: draggingCharacter:', draggingCharacter, 'gameState:', gameState);
    if (!draggingCharacter || !gameState) {
        console.log('No dragging character or game state not fully initialized');
        setDraggingCharacter(null);
        return;
    }

    const { gridX, gridY } = getGridPosition(event);
    const clientID = localStorage.getItem('clientID');
    const draggedChar = { ...draggingCharacter };

    if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10) {
        const path = calculatePath(draggedChar.Position[0], draggedChar.Position[1], gridX, gridY);
        setMovePath(path);

        if (canMove(gridX, gridY)) {
            animateMove(draggedChar, path, () => {
                sendMessage(JSON.stringify({
                    type: 'move',
                    clientID: clientID,
                    characterID: draggedChar.ID,
                    position: [gridX, gridY]
                }));
            });
        } else if (canAttackOrUseAbility(gridX, gridY, myTeam) && isWithinAttackRange(draggedChar, gridX, gridY)) {
            const target = findCharacter(gameState.Teams, gameState.Board[gridX][gridY]);
            animateMove(draggedChar, path, () => {
                if (selectedAbility) {
                    sendMessage(JSON.stringify({
                        type: 'ability',
                        clientID: clientID,
                        characterID: draggedChar.ID,
                        targetID: target.ID,
                        ability: selectedAbility.Name
                    }));
                } else {
                    sendMessage(JSON.stringify({
                        type: 'attack',
                        clientID: clientID,
                        characterID: draggedChar.ID,
                        targetID: target.ID
                    }));
                }
            }, true);
        }
    }
    setDraggingCharacter(null);
}

function handleClick(event, myTeam) {
    if (!gameState || !gameState.Teams || draggingCharacter) return;

    const { gridX, gridY } = getGridPosition(event);
    const clientID = localStorage.getItem('clientID');
    if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10) {
        const charId = gameState.Board[gridX][gridY];
        if (charId !== -1 && canAttackOrUseAbility(gridX, gridY, myTeam)) {
            const target = findCharacter(gameState.Teams, charId);
            const currentChar = findCharacter(gameState.Teams, gameState.CurrentTurn);
            if (currentChar && currentChar.Team === myTeam && isWithinAttackRange(currentChar, gridX, gridY)) {
                const path = calculatePath(currentChar.Position[0], currentChar.Position[1], gridX, gridY);
                const charToAct = { ...currentChar };
                const targetToAct = { ...target };

                animateMove(charToAct, path, () => {
                    if (selectedAbility) {
                        sendMessage(JSON.stringify({
                            type: 'ability',
                            clientID: clientID,
                            characterID: charToAct.ID,
                            targetID: targetToAct.ID,
                            ability: selectedAbility.Name
                        }));
                        currentChar.Abilities = currentChar.Abilities.filter(a => a.Name !== selectedAbility.Name);
                        setSelectedAbility(null);
                    } else {
                        sendMessage(JSON.stringify({
                            type: 'attack',
                            clientID: clientID,
                            characterID: charToAct.ID,
                            targetID: targetToAct.ID
                        }));
                    }
                }, true);
            }
        }
    }
}

function handleEndTurn(myTeam) {
    const clientID = localStorage.getItem('clientID');
    const currentChar = findCharacter(gameState.Teams, gameState.CurrentTurn);
    if (currentChar && currentChar.Team === myTeam) {
        sendMessage(JSON.stringify({
            type: 'end_turn',
            clientID: clientID
        }));
    }
}