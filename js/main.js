import { connectWebSocket } from './websocket.js';
import { drawBoard, preloadImages, imagesCache } from './renderCanvas.js';
import { updateCharacterCards, updateAbilityCards, updatePhaseAndProgress, updateBattleLog } from './renderUI.js';
import { setGameState, setCurrentRoom, setIsSpectator, setSelectedCharacter, setSelectedAbility } from './state.js';
import { setupEventListeners } from './eventHandlers.js';

let myTeam = null;

// Функция для ожидания загрузки всех изображений из кэша
function waitForImages() {
    const images = Object.values(imagesCache);
    return Promise.all(
        images.map((img) => {
            if (img.complete && img.naturalWidth !== 0) {
                return Promise.resolve(); // Изображение уже загружено и валидно
            }
            return new Promise((resolve) => {
                const timeout = setTimeout(() => {
                    console.warn(`Image loading timed out: ${img.src}`);
                    resolve(); // Продолжаем, даже если изображение не загрузилось
                }, 500); // Таймаут 500 ms

                img.onload = () => {
                    clearTimeout(timeout);
                    if (img.naturalWidth !== 0) resolve();
                    else resolve(); // Даже если загружено с ошибкой, продолжаем
                };
                img.onerror = () => {
                    clearTimeout(timeout);
                    console.warn(`Image failed to load: ${img.src}`);
                    resolve(); // Ошибка загрузки, продолжаем
                };
            });
        })
    );
}

// Обработка клика по кнопке "Join Room"
function handleJoinRoom() {
    const currentClientID = localStorage.getItem('clientID');
    if (!currentClientID) {
        console.error('Cannot join room: clientID is not set. Please register first.');
        alert('Please register before joining a room.');
        return;
    }

    const room = document.getElementById('roomSelect').value;
    document.getElementById('roomSelection').classList.add('hidden');
    document.getElementById('mainContainer').classList.remove('hidden');
    document.getElementById('wrestleCards').classList.remove('hidden');

    setCurrentRoom(room);
    setIsSpectator(false);

    connectWebSocket(room, false, handleWebSocketMessage);
}

// Обработка сообщений от WebSocket
async function handleWebSocketMessage(event) {
    try {
        const data = JSON.parse(event.data);
        console.log('Received WebSocket data:', data);
        setGameState(data);
        myTeam = data.teamID;

        // Предварительно загружаем изображения
        preloadImages(data);

        // Ждём, пока все изображения загрузятся
        await waitForImages().catch((error) => {
            console.warn('Image loading failed or timed out, continuing without images:', error);
        });

        // Обновляем интерфейс
        updateUI(data);

        // Устанавливаем обработчики событий только после получения первого состояния
        setupEventListenersIfNeeded(myTeam);
    } catch (error) {
        console.error('Error processing WebSocket message:', error);
    }
}

// Обновление интерфейса
function updateUI(data) {
    updateCharacterCards(data);
    updateAbilityCards(myTeam, data); // Передаём myTeam и data
    drawBoard(data);
    updatePhaseAndProgress(data);
    updateBattleLog(data);
}

// Установка обработчиков событий, если они еще не установлены
function setupEventListenersIfNeeded(myTeam) {
    if (!document.getElementById('gameCanvas').hasAttribute('data-listeners-set')) {
        setupEventListeners(myTeam);
        document.getElementById('gameCanvas').setAttribute('data-listeners-set', 'true');
    }
}

// Инициализация
function initialize() {
    document.getElementById('joinRoomBtn').addEventListener('click', handleJoinRoom);
}

// Запуск приложения
initialize();