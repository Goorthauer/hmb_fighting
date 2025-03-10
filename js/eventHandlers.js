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
import { findPath, getGridPosition, canMove, canAttackOrUseAbility, isWithinAttackRange, killUnplacedCharacters } from './gameLogic.js';
import { findCharacter } from './utils.js';
import { canvas } from './constants.js';

// Настройка слушателей событий
export function setupEventListeners(myTeam) {
    if (myTeam === null || myTeam === undefined) {
        console.error('myTeam is not initialized yet, delaying event listeners setup');
        return;
    }
    canvas.addEventListener('mousemove', handleMouseMove);
    canvas.addEventListener('mouseup', (e) => handleMouseUp(e, myTeam));
    canvas.addEventListener('click', (e) => handleClick(e, myTeam));
    document.getElementById('endTurnBtn').addEventListener('click', () => handleEndTurn(myTeam));
    document.getElementById('startGameBtn').addEventListener('click', () => handleStartGame(myTeam));
    setupCardDragListeners(myTeam);
}

// Настройка слушателей событий для перетаскивания карт
export function setupCardDragListeners(myTeam) {
    const characterCards = document.getElementById('characterCards');
    if (!characterCards) return;

    characterCards.addEventListener('dragstart', (e) => handleDragStart(e, myTeam));
    canvas.addEventListener('dragover', (e) => e.preventDefault());
    canvas.addEventListener('drop', (e) => handleDrop(e, myTeam));
}

// Обработка начала перетаскивания карты
function handleDragStart(event, myTeam) {
    if (gameState.phase !== 'setup') return;
    const card = event.target.closest('.card');
    if (!card || !gameState || !gameState.teams) return;

    const charId = parseInt(card.dataset.id);
    const char = findCharacter(gameState.teams, charId);
    if (char && char.team === myTeam) {
        console.log('Drag started for:', char.name);
        setDraggingCharacter(char);
        event.dataTransfer.setData('text/plain', char.id.toString());
    } else {
        event.preventDefault();
    }
}

// Обработка движения мыши
function handleMouseMove(event) {
    if (!draggingCharacter) return;
    const { x, y } = getGridPosition(event);
    setDragOffsetX(x);
    setDragOffsetY(y);
    drawBoard(gameState);
}

// Обработка отпускания кнопки мыши
function handleMouseUp(event, myTeam) {
    if (!draggingCharacter || !gameState) {
        setDraggingCharacter(null);
        return;
    }
    const { gridX, gridY } = getGridPosition(event);
    handlePlaceOrMove(gridX, gridY, myTeam);
}

// Обработка события drop
function handleDrop(event, myTeam) {
    event.preventDefault();
    if (!draggingCharacter || !gameState) {
        setDraggingCharacter(null);
        return;
    }
    const { gridX, gridY } = getGridPosition(event);
    handlePlaceOrMove(gridX, gridY, myTeam);
}

// Обработка размещения или перемещения персонажа
function handlePlaceOrMove(gridX, gridY, myTeam) {
    const clientID = localStorage.getItem('clientID');
    const draggedChar = draggingCharacter;

    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9) {
        if (gameState.phase === 'setup') {
            handlePlaceCharacter(gridX, gridY, myTeam, clientID, draggedChar);
        } else if (draggedChar.id === gameState.currentTurn && draggedChar.team === myTeam) {
            if (gameState.phase === 'move') {
                handleMoveCharacter(gridX, gridY, myTeam, clientID, draggedChar);
            } else if (gameState.phase === 'action') {
                handleAttackOrUseAbility(gridX, gridY, myTeam, clientID, draggedChar);
            }
        }
    }
    setDraggingCharacter(null);
    drawBoard(gameState);
}

// Обработка размещения персонажа в фазе setup
function handlePlaceCharacter(gridX, gridY, myTeam, clientID, draggedChar) {
    const isValidZone = (myTeam === 0 && gridX < 8) || (myTeam === 1 && gridX >= 8);
    const placedCount = gameState.teams[myTeam].characters.filter(c => c.position[0] !== -1).length;
    if (isValidZone && gameState.board[gridX][gridY] === -1 && placedCount < 5) {
        console.log(`Placing ${draggedChar.name} at (${gridX}, ${gridY})`);
        sendMessage(JSON.stringify({
            type: 'place',
            clientID: clientID,
            characterID: draggedChar.id,
            position: [gridX, gridY]
        }));
    }
}

// Обработка перемещения персонажа в фазе move
function handleMoveCharacter(gridX, gridY, myTeam, clientID, draggedChar) {
    if (canMove(gridX, gridY)) {
        console.log(`Moving ${draggedChar.name} to (${gridX}, ${gridY})`);
        const path = findPath(draggedChar.position[0], draggedChar.position[1], gridX, gridY, draggedChar.stamina);
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
        console.log(`Attacking from move phase ${target.name} at (${gridX}, ${gridY})`);
        const path = findPath(draggedChar.position[0], draggedChar.position[1], gridX, gridY, draggedChar.stamina);
        setMovePath(path);
        animateMove(draggedChar, path, () => {
            if (selectedAbility) {
                handleUseAbility(clientID, draggedChar, target);
            } else {
                handleAttack(clientID, draggedChar, target);
            }
        }, true);
    }
}

// Обработка атаки или использования способности в фазе action
function handleAttackOrUseAbility(gridX, gridY, myTeam, clientID, draggedChar) {
    if (canAttackOrUseAbility(gridX, gridY, myTeam) && isWithinAttackRange(draggedChar, gridX, gridY, gameState.weaponsConfig, selectedAbility)) {
        const target = findCharacter(gameState.teams, gameState.board[gridX][gridY]);
        console.log(`Attacking/Using ability on ${target.name} at (${gridX}, ${gridY})`);
        const path = findPath(draggedChar.position[0], draggedChar.position[1], gridX, gridY, draggedChar.stamina);
        setMovePath(path);
        animateMove(draggedChar, path, () => {
            if (selectedAbility) {
                handleUseAbility(clientID, draggedChar, target);
            } else {
                handleAttack(clientID, draggedChar, target);
            }
        }, true);
    }
}

// Обработка использования способности
function handleUseAbility(clientID, draggedChar, target) {
    sendMessage(JSON.stringify({
        type: 'ability',
        clientID: clientID,
        characterID: draggedChar.id,
        targetID: target.id,
        ability: selectedAbility.name.toLowerCase()
    }));
    draggedChar.abilities = draggedChar.abilities.filter(a => a !== selectedAbility.name);
    setSelectedAbility(null);
}

// Обработка атаки
function handleAttack(clientID, draggedChar, target) {
    sendMessage(JSON.stringify({
        type: 'attack',
        clientID: clientID,
        characterID: draggedChar.id,
        targetID: target.id
    }));
}

// Обработка клика
function handleClick(event, myTeam) {
    if (!gameState || !gameState.teams || draggingCharacter) return;

    const { gridX, gridY } = getGridPosition(event);
    const clientID = localStorage.getItem('clientID');
    const currentChar = findCharacter(gameState.teams, gameState.currentTurn);

    if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && currentChar?.team === myTeam) {
        const charId = gameState.board[gridX][gridY];
        if (charId === gameState.currentTurn && (gameState.phase === 'move' || gameState.phase === 'action')) {
            setDraggingCharacter(currentChar);
            const { x, y } = getGridPosition(event);
            setDragOffsetX(x);
            setDragOffsetY(y);
            drawBoard(gameState);
            return;
        }

        if (gameState.phase === 'move') {
            handleMoveCharacter(gridX, gridY, myTeam, clientID, currentChar);
        } else if (gameState.phase === 'action' && canAttackOrUseAbility(gridX, gridY, myTeam)) {
            handleAttackOrUseAbility(gridX, gridY, myTeam, clientID, currentChar);
        }
    }
}

// Обработка завершения хода
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

// Обработка начала игры
function handleStartGame(myTeam) {
    const clientID = localStorage.getItem('clientID');
    if (gameState.phase === 'setup' && gameState.teams[myTeam].characters.filter(c => c.position[0] !== -1).length >= 5) {
        console.log('Starting game');
        // Убиваем непоставленных персонажей перед началом игры
        killUnplacedCharacters(myTeam);
        sendMessage(JSON.stringify({
            type: 'start',
            clientID: clientID
        }));
    } else {
        console.warn('Cannot start game: not enough characters placed (minimum 5 required)');
    }
}