import { connectWebSocket, sendMessage, setClientID } from './websocket.js';
import {
    drawBoard,
    updateCharacterCards,
    updateAbilityCards,
    updatePhaseAndProgress,
    updateBattleLog,
    animateMove
} from './render.js';
import {
    gameState,
    draggingCharacter,
    dragOffsetX,
    dragOffsetY,
    selectedCharacter,
    selectedAbility,
    isSpectator,
    currentRoom,
    setGameState,
    setSelectedCharacter,
    setDraggingCharacter,
    setDragOffset,
    setSelectedAbility,
    setIsSpectator,
    setCurrentRoom
} from './state.js';

let myTeam = null;
let previousState = null;

document.getElementById('joinRoomBtn').addEventListener('click', () => {
    const currentClientID = localStorage.getItem('clientID');
    if (!currentClientID) {
        console.error('Cannot join room: clientID is not set. Please register first.');
        alert('Please register before joining a room.');
        return;
    }

    const room = document.getElementById('roomSelect').value;
    console.log('Joining room:', room, 'with clientID:', currentClientID);
    document.getElementById('roomSelection').classList.add('hidden');
    document.getElementById('mainContainer').classList.remove('hidden');
    document.getElementById('wrestleCards').classList.remove('hidden');

    setCurrentRoom(room);
    setIsSpectator(false);

    connectWebSocket(room, false, (event) => {
        const data = JSON.parse(event.data);
        console.log('Received WebSocket data:', data);
        setGameState(data);
        myTeam = data.TeamID;

        updateCharacterCards(setSelectedCharacter, data);
        updateAbilityCards(myTeam, setSelectedAbility, data);
        drawBoard(data);
        updatePhaseAndProgress(data);
        updateBattleLog(data, previousState);

        previousState = data;
    });
});

const canvas = document.getElementById('gameCanvas');
canvas.addEventListener('mousedown', (e) => {
    if (!gameState || !gameState.Board || !gameState.Teams) {
        console.log('Game state not fully initialized yet');
        return;
    }

    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const gridX = Math.floor(x / 50);
    const gridY = Math.floor(y / 50);

    const charId = gameState.Board[gridX][gridY];
    if (charId !== -1) {
        const char = findCharacter(gameState.Teams, charId);
        if (char && char.Team === myTeam && gameState.CurrentTurn && char.ID === gameState.CurrentTurn) {
            setDraggingCharacter(char);
            setDragOffset(x,y);
        }
    }
});

canvas.addEventListener('mousemove', (e) => {
    if (draggingCharacter) {
        const rect = canvas.getBoundingClientRect();
        setDragOffset(e.clientX - rect.left,e.clientY - rect.top);
        drawBoard(gameState);
    }
});

canvas.addEventListener('mouseup', (e) => {
    if (!draggingCharacter || !gameState || !gameState.Board || !gameState.Teams) {
        console.log('No dragging character or game state not fully initialized');
        setDraggingCharacter(null);
        return;
    }

    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const gridX = Math.floor(x / 50);
    const gridY = Math.floor(y / 50);
    const currentClientID = localStorage.getItem('clientID');
    const draggedChar = { ...draggingCharacter }; // Сохраняем копию draggingCharacter

    if (gameState.Phase === 'move' && gameState.Board[gridX][gridY] === -1) {
        const path = [{ x: draggedChar.Position[0], y: draggedChar.Position[1] }, { x: gridX, y: gridY }];
        animateMove(draggedChar, path, () => {
            sendMessage(JSON.stringify({
                type: 'move',
                characterID: draggedChar.ID,
                position: [gridX, gridY],
                clientID: currentClientID
            }));
        });
    } else if (gameState.Phase === 'action' && gameState.Board[gridX][gridY] !== -1) {
        const target = findCharacter(gameState.Teams, gameState.Board[gridX][gridY]);
        if (target && target.Team !== draggedChar.Team) {
            const path = [
                { x: draggedChar.Position[0], y: draggedChar.Position[1] },
                { x: gridX, y: gridY }
            ];
            animateMove(draggedChar, path, () => {
                if (selectedAbility) {
                    sendMessage(JSON.stringify({
                        type: 'ability',
                        characterID: draggedChar.ID,
                        targetID: target.ID,
                        ability: selectedAbility.Name,
                        clientID: currentClientID
                    }));
                } else {
                    sendMessage(JSON.stringify({
                        type: 'attack',
                        characterID: draggedChar.ID,
                        targetID: target.ID,
                        clientID: currentClientID
                    }));
                }
            }, true);
        }
    }

    setDraggingCharacter(null);
});

document.getElementById('endTurnBtn').addEventListener('click', () => {
    const currentClientID = localStorage.getItem('clientID');
    sendMessage(JSON.stringify({
        type: 'end_turn',
        clientID: currentClientID
    }));
});

function findCharacter(teams, id) {
    if (!teams) return null;
    for (let team of teams) {
        for (let char of team.Characters) {
            if (char.ID === id) return char;
        }
    }
    return null;
}