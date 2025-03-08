import {
    gameState,
    draggingCharacter,
    selectedCharacter,
    selectedAbility,
    setDraggingCharacter,
    setSelectedCharacter,
    setSelectedAbility,
    setMovePath,
    setDragOffsetX,
    setDragOffsetY
} from './state.js';
import {connectWebSocket, sendMessage} from './websocket.js';
import {animateMove, drawBoard} from './renderCanvas.js';
import {calculatePath, getGridPosition, canMove, canAttackOrUseAbility, isWithinAttackRange} from './gameLogic.js';
import {findCharacter} from './utils.js';
import {canvas} from './constants.js';

export function setupEventListeners(myTeam) {
    if (myTeam === null || myTeam === undefined) {
        console.error('myTeam is not initialized yet, delaying event listeners setup');
        return;
    }
    canvas.addEventListener('mousedown', (e) => handleMouseDown(e, myTeam));
    canvas.addEventListener('mousemove', (e) => handleMouseMove(e));
    canvas.addEventListener('mouseup', (e) => handleMouseUp(e, myTeam));
    canvas.addEventListener('click', (e) => handleClick(e, myTeam));
    document.getElementById('endTurnBtn').addEventListener('click', () => handleEndTurn(myTeam));
}

function handleMouseDown(event, myTeam) {
    if (!gameState || !gameState.teams || !Array.isArray(gameState.teams)) {
        return;
    }

    const {gridX, gridY, x, y} = getGridPosition(event);
    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9) { // Обновлено с 20x10 на 16x9
        const charId = gameState.board[gridX][gridY];
        if (charId !== -1) {
            const char = findCharacter(gameState.teams, charId);
            if (char && char.team === myTeam && char.id === gameState.currentTurn) {
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
    const {x, y} = getGridPosition(event);
    setDragOffsetX(x);
    setDragOffsetY(y);
    drawBoard(gameState);
}

function handleMouseUp(event, myTeam) {
    if (!draggingCharacter || !gameState) {
        setDraggingCharacter(null);
        return;
    }

    const {gridX, gridY} = getGridPosition(event);
    const clientID = localStorage.getItem('clientID');
    const draggedChar = {...draggingCharacter};

    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9) { // Обновлено с 20x10 на 16x9
        const path = calculatePath(draggedChar.position[0], draggedChar.position[1], gridX, gridY);

        if (canMove(gridX, gridY)) {
            setMovePath(path);
            animateMove(draggedChar, path, () => {
                sendMessage(JSON.stringify({
                    type: 'move',
                    clientID: clientID,
                    characterID: draggedChar.id,
                    position: [gridX, gridY]
                }));
            });
        } else if (canAttackOrUseAbility(gridX, gridY, myTeam) && isWithinAttackRange(draggedChar, gridX, gridY, gameState.weaponsConfig, selectedAbility)) {
            const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
            setMovePath(path);
            animateMove(draggedChar, path, () => {
                if (selectedAbility) {
                    sendMessage(JSON.stringify({
                        type: 'ability',
                        clientID: clientID,
                        characterID: draggedChar.id,
                        targetID: target.id,
                        ability: selectedAbility.name
                    }));
                } else {
                }
            }, true);
        } else {
            setDraggingCharacter(null);
            drawBoard(gameState);
            return;
        }
    } else {
        setDraggingCharacter(null);
        drawBoard(gameState);
        return;
    }
    setDraggingCharacter(null);
}

function handleClick(event, myTeam) {
    if (!gameState || !gameState.teams || draggingCharacter) return;

    const {gridX, gridY} = getGridPosition(event);
    const clientID = localStorage.getItem('clientID');
    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9) { // Обновлено с 20x10 на 16x9
        const charId = gameState.board[gridX][gridY];
        if (charId !== -1 && canAttackOrUseAbility(gridX, gridY, myTeam)) {
            const target = findCharacter(gameState.teams, charId);
            const currentChar = findCharacter(gameState.teams, gameState.currentTurn);
            if (currentChar && currentChar.team === myTeam && isWithinAttackRange(currentChar, gridX, gridY, gameState.weaponsConfig, selectedAbility)) {
                const path = calculatePath(currentChar.position[0], currentChar.position[1], gridX, gridY);
                const charToAct = {...currentChar};
                const targetToAct = {...target};

                animateMove(charToAct, path, () => {
                    if (selectedAbility) {
                        sendMessage(JSON.stringify({
                            type: 'ability',
                            clientID: clientID,
                            characterID: charToAct.id,
                            targetID: targetToAct.id,
                            ability: selectedAbility.name
                        }));
                        currentChar.abilities = currentChar.abilities.filter(a => a.name !== selectedAbility.name);
                        setSelectedAbility(null);
                    } else {
                        sendMessage(JSON.stringify({
                            type: 'attack',
                            clientID: clientID,
                            characterID: charToAct.id,
                            targetID: targetToAct.id
                        }));
                    }
                }, true);
            }
        }
    }
}

function handleEndTurn(myTeam) {
    const clientID = localStorage.getItem('clientID');
    const currentChar = findCharacter(gameState.teams, gameState.currentTurn);
    if (currentChar && currentChar.team === myTeam) {
        sendMessage(JSON.stringify({
            type: 'end_turn',
            clientID: clientID
        }));
    }
}