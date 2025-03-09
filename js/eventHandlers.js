import {
    gameState,
    draggingCharacter,
    selectedAbility,
    setDraggingCharacter,
    setSelectedAbility,
    setMovePath,
    setDragOffsetX,
    setDragOffsetY
} from './state.js';
import { sendMessage } from './websocket.js';
import { animateMove, drawBoard } from './renderCanvas.js';
import { calculatePath, getGridPosition, canMove, canAttackOrUseAbility, isWithinAttackRange } from './gameLogic.js';
import { findCharacter } from './utils.js';
import { canvas } from './constants.js';

export function setupEventListeners(myTeam) {
    if (myTeam === null || myTeam === undefined) {
        console.error('myTeam is not initialized yet, delaying event listeners setup');
        return;
    }
    canvas.addEventListener('mousemove', (e) => handleMouseMove(e));
    canvas.addEventListener('mouseup', (e) => handleMouseUp(e, myTeam));
    canvas.addEventListener('click', (e) => handleClick(e, myTeam));
    document.getElementById('endTurnBtn').addEventListener('click', () => handleEndTurn(myTeam));
    document.getElementById('startGameBtn').addEventListener('click', () => handleStartGame(myTeam));
    setupCardDragListeners(myTeam);
}

export function setupCardDragListeners(myTeam) {
    const characterCards = document.getElementById('characterCards');
    if (!characterCards) return;

    // Перетаскивание карточек только в фазе setup
    characterCards.addEventListener('dragstart', (e) => {
        if (gameState.phase !== 'setup') return;
        const card = e.target.closest('.card');
        if (!card || !gameState || !gameState.teams) return;

        const charId = parseInt(card.dataset.id);
        const char = findCharacter(gameState.teams, charId);
        if (char && char.team === myTeam) {
            console.log('Drag started for:', char.name);
            setDraggingCharacter(char);
            e.dataTransfer.setData('text/plain', char.id.toString());
        } else {
            e.preventDefault();
        }
    });

    canvas.addEventListener('dragover', (e) => e.preventDefault());
    canvas.addEventListener('drop', (e) => handleDrop(e, myTeam));
}

function handleMouseMove(event) {
    if (!draggingCharacter) return;
    const { x, y } = getGridPosition(event);
    setDragOffsetX(x);
    setDragOffsetY(y);
    drawBoard(gameState);
}

function handleMouseUp(event, myTeam) {
    if (!draggingCharacter || !gameState) {
        setDraggingCharacter(null);
        return;
    }
    const { gridX, gridY } = getGridPosition(event);
    handlePlaceOrMove(gridX, gridY, myTeam);
}

function handleDrop(event, myTeam) {
    event.preventDefault();
    if (!draggingCharacter || !gameState) {
        setDraggingCharacter(null);
        return;
    }
    const { gridX, gridY } = getGridPosition(event);
    handlePlaceOrMove(gridX, gridY, myTeam);
}

function handlePlaceOrMove(gridX, gridY, myTeam) {
    const clientID = localStorage.getItem('clientID');
    const draggedChar = draggingCharacter;

    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9) {
        if (gameState.phase === 'setup') {
            const isValidZone = (myTeam === 0 && gridX < 8) || (myTeam === 1 && gridX >= 8);
            const placedCount = gameState.teams[myTeam].characters.filter(c => c.position[0] !== -1).length;
            if (isValidZone && gameState.board[gridX][gridY] === -1 && placedCount < 5) {
                console.log(`Placing ${draggedChar.name} at (${gridX}, ${gridY})`);
                const path = calculatePath(draggedChar.position[0], draggedChar.position[1], gridX, gridY);
                animateMove(draggedChar, path, () => {
                    sendMessage(JSON.stringify({
                        type: 'place',
                        clientID: clientID,
                        characterID: draggedChar.id,
                        position: [gridX, gridY]
                    }));
                });
            }
        } else if (gameState.phase === 'move' && draggedChar.id === gameState.currentTurn && draggedChar.team === myTeam) {
            if (canMove(gridX, gridY)) {
                console.log(`Moving ${draggedChar.name} to (${gridX}, ${gridY})`);
                const path = calculatePath(draggedChar.position[0], draggedChar.position[1], gridX, gridY);
                setMovePath(path);
                animateMove(draggedChar, path, () => {
                    sendMessage(JSON.stringify({
                        type: 'move',
                        clientID: clientID,
                        characterID: draggedChar.id,
                        position: [gridX, gridY]
                    }));
                });
            }
        } else if (gameState.phase === 'action' && draggedChar.id === gameState.currentTurn && draggedChar.team === myTeam) {
            if (canAttackOrUseAbility(gridX, gridY, myTeam) && isWithinAttackRange(draggedChar, gridX, gridY, gameState.weaponsConfig, selectedAbility)) {
                const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
                console.log(`Attacking/Using ability on ${target.name} at (${gridX}, ${gridY})`);
                const path = calculatePath(draggedChar.position[0], draggedChar.position[1], gridX, gridY);
                setMovePath(path);
                animateMove(draggedChar, path, () => {
                    if (selectedAbility) {
                        sendMessage(JSON.stringify({
                            type: 'ability',
                            clientID: clientID,
                            characterID: draggedChar.id,
                            targetID: target.id,
                            ability: selectedAbility.name.toLowerCase()
                        }));
                        draggedChar.abilities = draggedChar.abilities.filter(a => a !== selectedAbility.name);
                        setSelectedAbility(null);
                    } else {
                        sendMessage(JSON.stringify({
                            type: 'attack',
                            clientID: clientID,
                            characterID: draggedChar.id,
                            targetID: target.id
                        }));
                    }
                }, true);
            }
        }
    }
    setDraggingCharacter(null);
    drawBoard(gameState);
}

function handleClick(event, myTeam) {
    if (!gameState || !gameState.teams || draggingCharacter) return;

    const { gridX, gridY } = getGridPosition(event);
    const clientID = localStorage.getItem('clientID');
    const currentChar = findCharacter(gameState.teams, gameState.currentTurn);

    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && currentChar?.team === myTeam) {
        // Выбор персонажа на поле для перетаскивания
        const charId = gameState.board[gridX][gridY];
        if (charId === gameState.currentTurn && (gameState.phase === 'move' || gameState.phase === 'action')) {
            setDraggingCharacter(currentChar);
            const { x, y } = getGridPosition(event);
            setDragOffsetX(x);
            setDragOffsetY(y);
            drawBoard(gameState);
            return;
        }

        // Ход или атака
        if (gameState.phase === 'move' && canMove(gridX, gridY)) {
            console.log(`Click to move ${currentChar.name} to (${gridX}, ${gridY})`);
            const path = calculatePath(currentChar.position[0], currentChar.position[1], gridX, gridY);
            setMovePath(path);
            animateMove(currentChar, path, () => {
                sendMessage(JSON.stringify({
                    type: 'move',
                    clientID: clientID,
                    characterID: currentChar.id,
                    position: [gridX, gridY]
                }));
            });
        } else if (gameState.phase === 'action' && canAttackOrUseAbility(gridX, gridY, myTeam)) {
            const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
            if (isWithinAttackRange(currentChar, gridX, gridY, gameState.weaponsConfig, selectedAbility)) {
                console.log(`Click to attack/use ability on ${target.name}`);
                const path = calculatePath(currentChar.position[0], currentChar.position[1], gridX, gridY);
                animateMove(currentChar, path, () => {
                    if (selectedAbility) {
                        sendMessage(JSON.stringify({
                            type: 'ability',
                            clientID: clientID,
                            characterID: currentChar.id,
                            targetID: target.id,
                            ability: selectedAbility.name.toLowerCase()
                        }));
                        currentChar.abilities = currentChar.abilities.filter(a => a !== selectedAbility.name);
                        setSelectedAbility(null);
                    } else {
                        sendMessage(JSON.stringify({
                            type: 'attack',
                            clientID: clientID,
                            characterID: currentChar.id,
                            targetID: target.id
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
    if (currentChar && currentChar.team === myTeam && gameState.phase !== 'setup') {
        console.log('Ending turn');
        sendMessage(JSON.stringify({
            type: 'end_turn',
            clientID: clientID
        }));
    }
}

function handleStartGame(myTeam) {
    const clientID = localStorage.getItem('clientID');
    if (gameState.phase === 'setup' && gameState.teams[myTeam].characters.filter(c => c.position[0] !== -1).length >= 5) {
        console.log('Starting game');
        sendMessage(JSON.stringify({
            type: 'start',
            clientID: clientID
        }));
    }
}