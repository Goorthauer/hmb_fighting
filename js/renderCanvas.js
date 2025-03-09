import {
    gameState,
    draggingCharacter,
    dragOffsetX,
    dragOffsetY,
    movePath,
    setMovePath,
    selectedAbility
} from './state.js';
import { ctx, cellWidth, cellHeight } from './constants.js';
import { findCharacter } from './utils.js';
import { findPath } from './gameLogic.js';

let movingCharacter = null;
export const imagesCache = {};

const DEFAULT_IMAGES = {
    character: '/static/characters/default.png',
    ability: '/static/abilities/default.jpg',
    weapon: '/static/weapons/default.png',
    shield: '/static/shields/default.png',
    icon: '/static/icons/default.png'
};

// Загрузка изображения с обработкой ошибок
function loadImage(url, defaultUrl) {
    if (!url || url.trim() === '') {
        console.warn(`Image URL is empty or invalid, using default: ${defaultUrl}`);
        url = defaultUrl;
    }
    if (!imagesCache[url]) {
        const img = new Image();
        img.src = url;
        img.onload = () => {
            console.log(`Image loaded: ${url}`);
            if (gameState) drawBoard(gameState);
        };
        img.onerror = () => {
            console.warn(`Failed to load image: ${url}, falling back to ${defaultUrl}`);
            if (url !== defaultUrl) {
                img.src = defaultUrl;
                img.onload = () => console.log(`Default image loaded: ${defaultUrl}`);
            }
        };
        imagesCache[url] = img;
    }
    return imagesCache[url];
}

// Очистка холста
function clearCanvas() {
    if (!ctx) {
        console.error('Canvas context (ctx) is not initialized');
        return;
    }
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
}

// Отрисовка сетки игрового поля
function drawGrid() {
    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            ctx.strokeStyle = '#ccc';
            ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        }
    }
}

// Отрисовка зон для фазы setup
function drawSetupZones(data) {
    const myTeam = data.teamID;
    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            const isMyZone = (myTeam === 0 && x < 8) || (myTeam === 1 && x >= 8);
            ctx.fillStyle = isMyZone ? 'rgba(0, 255, 0, 0.1)' : 'rgba(255, 0, 0, 0.1)';
            ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        }
    }
}

// Отрисовка доступных ходов для фазы move
function drawMoveRange(data, currentChar) {
    const startX = currentChar.position[0];
    const startY = currentChar.position[1];
    const stamina = currentChar.stamina;

    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            if (data.board[x][y] === -1) {
                const path = findPath(startX, startY, x, y, stamina);
                if (path.length > 0) {
                    ctx.fillStyle = 'rgba(0, 255, 0, 0.2)';
                    ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                }
            }
        }
    }
}

// Отрисовка зоны атаки для фазы action
function drawAttackRange(data, currentChar) {
    const startX = currentChar.position[0];
    const startY = currentChar.position[1];
    const range = selectedAbility ? selectedAbility.range : (data.weaponsConfig[currentChar.weapon]?.range || 1);

    for (let x = Math.max(0, startX - range); x <= Math.min(15, startX + range); x++) {
        for (let y = Math.max(0, startY - range); y <= Math.min(8, startY + range); y++) {
            const dx = Math.abs(x - startX);
            const dy = Math.abs(y - startY);
            if (dx <= range && dy <= range && data.board[x][y] !== -1 && data.board[x][y] !== currentChar.id) {
                const target = findCharacter(data.teams, data.board[x][y]);
                if (target && target.team !== data.teamID) {
                    ctx.fillStyle = 'rgba(255, 0, 0, 0.2)';
                    ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                }
            }
        }
    }
}

// Отрисовка персонажа на поле
function drawCharacter(char, x, y, isCurrentTurn = false) {
    const teamIconURL = gameState.teamsConfig[char.team].iconURL;
    const charImage = loadImage(teamIconURL, DEFAULT_IMAGES.icon);

    if (charImage && charImage.complete) {
        const iconSize = 40;
        const offsetX = (cellWidth - iconSize) / 2;
        const offsetY = (cellHeight - 15 - iconSize) / 2;
        ctx.drawImage(charImage, x * cellWidth + offsetX, y * cellHeight + offsetY, iconSize, iconSize);
    }

    ctx.fillStyle = char.team === 0 ? 'rgba(128, 0, 128, 0.8)' : 'rgba(255, 215, 0, 0.8)';
    ctx.fillRect(x * cellWidth, y * cellHeight + cellHeight - 15, cellWidth, 15);
    ctx.fillStyle = '#fff';
    ctx.font = '12px Arial';
    ctx.textAlign = 'center';
    ctx.fillText(char.name, x * cellWidth + cellWidth / 2, y * cellHeight + cellHeight - 7.5);

    if (isCurrentTurn) {
        ctx.strokeStyle = 'yellow';
        ctx.lineWidth = 3;
        ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        ctx.lineWidth = 1;
    }
}

// Отрисовка всех персонажей на поле
function drawCharacters(data) {
    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            const charId = data.board[x][y];
            if (charId !== -1 && (!movingCharacter || movingCharacter.id !== charId)) {
                const char = findCharacter(data.teams, charId);
                if (char) {
                    drawCharacter(char, x, y, char.id === data.currentTurn);
                }
            }
        }
    }
}

// Отрисовка стрелки для перемещения или атаки
function drawArrow(fromX, fromY, toX, toY, color = 'blue') {
    ctx.strokeStyle = color;
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.moveTo(fromX, fromY);
    ctx.lineTo(toX, toY);
    ctx.stroke();

    const angle = Math.atan2(toY - fromY, toX - fromX);
    const arrowSize = 12;
    ctx.beginPath();
    ctx.moveTo(toX, toY);
    ctx.lineTo(toX - arrowSize * Math.cos(angle - Math.PI / 6), toY - arrowSize * Math.sin(angle - Math.PI / 6));
    ctx.lineTo(toX - arrowSize * Math.cos(angle + Math.PI / 6), toY - arrowSize * Math.sin(angle + Math.PI / 6));
    ctx.closePath();
    ctx.fillStyle = color;
    ctx.fill();
}

// Отрисовка перемещаемого персонажа
function drawDraggingCharacter(data) {
    if (!draggingCharacter) return;

    const gridX = Math.floor(dragOffsetX / cellWidth);
    const gridY = Math.floor(dragOffsetY / cellHeight);

    if (data.phase === 'setup') {
        const myTeam = data.teamID;
        const isValidZone = (myTeam === 0 && gridX < 8) || (myTeam === 1 && gridX >= 8);
        if (isValidZone && gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && data.board[gridX][gridY] === -1) {
            ctx.fillStyle = 'rgba(0, 255, 0, 0.3)';
            ctx.fillRect(gridX * cellWidth, gridY * cellHeight, cellWidth, cellHeight);
        }
    } else if (data.phase === 'move' || data.phase === 'action') {
        const startX = draggingCharacter.position[0];
        const startY = draggingCharacter.position[1];
        if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9) {
            const target = data.board[gridX][gridY] !== -1 ? findCharacter(data.teams, data.board[gridX][gridY]) : null;
            const color = target && target.team !== data.teamID ? 'red' : 'blue';
            drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, color);
        }
    }

    drawCharacter(draggingCharacter, dragOffsetX / cellWidth, dragOffsetY / cellHeight);
}

// Отрисовка перемещающегося персонажа
function drawMovingCharacter() {
    if (!movingCharacter) return;
    drawCharacter(movingCharacter, movingCharacter.currentX / cellWidth, movingCharacter.currentY / cellHeight);
}

// Основная функция отрисовки игрового поля
export function drawBoard(data) {
    clearCanvas();
    drawGrid();

    if (data.phase === 'setup') {
        drawSetupZones(data);
    }

    if (!data || !data.teams || !Array.isArray(data.teams) || !data.board || !data.teamsConfig) {
        console.error('Invalid data in drawBoard:', data);
        return;
    }

    const currentChar = findCharacter(data.teams, data.currentTurn);
    if (currentChar && currentChar.team === data.teamID) {
        if (data.phase === 'move') {
            drawMoveRange(data, currentChar);
        } else if (data.phase === 'action') {
            drawAttackRange(data, currentChar);
        }
    }

    drawCharacters(data);
    drawDraggingCharacter(data);
    drawMovingCharacter();
}

// Анимация перемещения персонажа
export function animateMove(character, path, callback, isAttack = false) {
    if (path.length === 0) {
        setMovePath([]);
        callback();
        return;
    }

    movingCharacter = { ...character, currentX: path[0].x * cellWidth, currentY: path[0].y * cellHeight };
    let step = 0;
    const speed = 6;
    let returning = false;

    function animate() {
        if (!isAttack && step >= path.length - 1) {
            movingCharacter = null;
            setMovePath([]);
            callback();
            return;
        }

        if (isAttack && step >= path.length - 1 && !returning) {
            returning = true;
            step = path.length - 1;
        }

        if (isAttack && returning && step <= 0) {
            movingCharacter = null;
            setMovePath([]);
            callback();
            return;
        }

        const nextX = path[returning ? step - 1 : step + 1]?.x * cellWidth || path[step].x * cellWidth;
        const nextY = path[returning ? step - 1 : step + 1]?.y * cellHeight || path[step].y * cellHeight;
        const dx = nextX - movingCharacter.currentX;
        const dy = nextY - movingCharacter.currentY;
        const dist = Math.sqrt(dx * dx + dy * dy);

        if (dist <= speed) {
            movingCharacter.currentX = nextX;
            movingCharacter.currentY = nextY;
            if (returning) step--;
            else step++;
        } else {
            movingCharacter.currentX += (dx / dist) * speed;
            movingCharacter.currentY += (dy / dist) * speed;
        }

        drawBoard(gameState);
        requestAnimationFrame(animate);
    }

    requestAnimationFrame(animate);
}