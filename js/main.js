import { connectWebSocket } from './websocket.js';
import { drawBoard } from './renderCanvas.js';
import { updateCharacterCards, updateAbilityCards, updatePhaseAndProgress, updateBattleLog } from './renderUI.js';
import { setGameState, setCurrentRoom, setIsSpectator, setSelectedCharacter, setSelectedAbility } from './state.js';
import { setupEventListeners } from './eventHandlers.js';

let myTeam = null;

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

    connectWebSocket(room, false, async (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log('Received WebSocket data:', data);
            setGameState(data);
            myTeam = data.TeamID;

            // Обновляем UI
            updateCharacterCards(data);
            updateAbilityCards(myTeam, data);
            drawBoard(data);
            updatePhaseAndProgress(data);
            updateBattleLog(data);

            // Устанавливаем обработчики событий только после получения первого состояния
            if (!document.getElementById('gameCanvas').hasAttribute('data-listeners-set')) {
                setupEventListeners(myTeam);
                document.getElementById('gameCanvas').setAttribute('data-listeners-set', 'true');
            }
        } catch (error) {
            console.error('Error processing WebSocket message:', error);
        }
    });
});