import { gameState, draggingCharacter, dragOffsetX, dragOffsetY, movePath, setMovePath, selectedAbility } from './state.js';
import { ctx, cellWidth, cellHeight } from './constants.js';
import { findCharacter } from './utils.js';

let movingCharacter = null;
export const imagesCache = {};

const DEFAULT_IMAGES = {
    character: '/static/characters/default.png',
    ability: '/static/abilities/default.jpg',
    weapon: '/static/weapons/default.png',
    shield: '/static/shields/default.png',
    icon: '/static/icons/default.png'
};

// Функция для динамической загрузки изображения
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
            // Перерисовываем доску после загрузки, если игра уже началась
            if (gameState) {
                drawBoard(gameState);
            }
        };
        img.onerror = () => {
            console.warn(`Failed to load image: ${url}, falling back to ${defaultUrl}`);
            if (url !== defaultUrl) {
                img.src = defaultUrl;
                img.onload = () => console.log(`Default image loaded: ${defaultUrl}`);
                img.onerror = () => console.error(`Failed to load default image: ${defaultUrl}`);
            } else {
                console.error(`Failed to load image: ${url}`);
            }
        };
        imagesCache[url] = img;
    }
    return imagesCache[url];
}

export function drawBoard(data) {
    if (!ctx) {
        console.error('Canvas context (ctx) is not initialized');
        return;
    }

    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    // Отрисовка сетки
    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            ctx.strokeStyle = '#ccc';
            ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        }
    }

    if (!data || !data.teams || !Array.isArray(data.teams) || !data.board || !data.teamsConfig) {
        console.error('Invalid data in drawBoard:', data);
        return;
    }

    // Логика перетаскивания персонажа
    if (draggingCharacter) {
        const startX = draggingCharacter.position[0];
        const startY = draggingCharacter.position[1];
        const stamina = draggingCharacter.stamina;
        const weapon = data.weaponsConfig[draggingCharacter.weapon];
        const weaponRange = weapon ? weapon.range : 1;
        const gridX = Math.floor(dragOffsetX / cellWidth);
        const gridY = Math.floor(dragOffsetY / cellHeight);

        if (data.phase === 'move') {
            for (let x = Math.max(0, startX - stamina); x <= Math.min(15, startX + stamina); x++) {
                for (let y = Math.max(0, startY - stamina); y <= Math.min(8, startY + stamina); y++) {
                    const dist = Math.abs(x - startX) + Math.abs(y - startY);
                    if (dist <= stamina && data.board[x][y] === -1) {
                        ctx.fillStyle = 'rgba(0, 255, 0, 0.2)';
                        ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                        ctx.strokeStyle = 'green';
                        ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    }
                }
            }
            if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && data.board[gridX][gridY] === -1) {
                drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'blue');
            }
        } else if (data.phase === 'action') {
            if (selectedAbility) {
                const abilityRange = selectedAbility.range || 1;
                for (let x = Math.max(0, startX - abilityRange); x <= Math.min(15, startX + abilityRange); x++) {
                    for (let y = Math.max(0, startY - abilityRange); y <= Math.min(8, startY + abilityRange); y++) {
                        const dist = Math.max(Math.abs(x - startX), Math.abs(y - startY));
                        if (dist <= abilityRange && data.board[x][y] !== -1) {
                            const target = findCharacter(data.teams, data.board[x][y]);
                            if (target && target.team !== draggingCharacter.team) {
                                ctx.fillStyle = 'rgba(255, 215, 0, 0.2)';
                                ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                                ctx.strokeStyle = 'gold';
                                ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                            }
                        }
                    }
                }
            } else {
                for (let x = Math.max(0, startX - weaponRange); x <= Math.min(15, startX + weaponRange); x++) {
                    for (let y = Math.max(0, startY - weaponRange); y <= Math.min(8, startY + weaponRange); y++) {
                        const dist = Math.max(Math.abs(x - startX), Math.abs(y - startY));
                        const isValidRange = (weapon && weapon.isTwoHanded && dist === weaponRange) || (!weapon || !weapon.isTwoHanded && dist <= weaponRange);
                        if (isValidRange && data.board[x][y] !== -1) {
                            const target = findCharacter(data.teams, data.board[x][y]);
                            if (target && target.team !== draggingCharacter.team) {
                                ctx.fillStyle = 'rgba(255, 0, 0, 0.2)';
                                ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                                ctx.strokeStyle = 'red';
                                ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                            }
                        }
                    }
                }
            }
            if (gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && data.board[gridX][gridY] !== -1) {
                const target = findCharacter(data.teams, data.board[gridX][gridY]);
                if (target && target.team !== draggingCharacter.team) {
                    const attackDist = Math.max(Math.abs(gridX - startX), Math.abs(gridY - startY));
                    const isValidRange = (weapon && weapon.isTwoHanded && attackDist === weaponRange) || (!weapon || !weapon.isTwoHanded && attackDist <= weaponRange);
                    if (isValidRange) {
                        drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'red', true);
                    }
                }
            }
        }
    }

    // Отрисовка персонажей на поле
    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            const charId = data.board[x][y];
            if (charId !== -1 && (!movingCharacter || movingCharacter.id !== charId)) {
                const char = findCharacter(data.teams, charId);
                if (!char) {
                    console.warn(`Character with ID ${charId} not found in teams`);
                    continue;
                }
                const teamIconURL = data.teamsConfig[char.team].iconURL;
                console.log(`Drawing ${char.name} with team iconURL: ${teamIconURL}`);
                const charImage = loadImage(teamIconURL, DEFAULT_IMAGES.icon);
                if (charImage && charImage.complete && charImage.naturalWidth !== 0) {
                    const iconSize = 40;
                    const offsetX = (cellWidth - iconSize) / 2;
                    const offsetY = (cellHeight - 15 - iconSize) / 2;
                    ctx.drawImage(charImage, x * cellWidth + offsetX, y * cellHeight + offsetY, iconSize, iconSize);
                } else {
                    console.warn(`Team icon for ${char.name} is not yet loaded or invalid:`, charImage);
                }
                ctx.fillStyle = char.team === 0 ? 'rgba(128, 0, 128, 0.8)' : 'rgba(255, 215, 0, 0.8)';
                ctx.fillRect(x * cellWidth, y * cellHeight + cellHeight - 15, cellWidth, 15);
                ctx.strokeStyle = char.team === 0 ? '#800080' : '#ffd700';
                ctx.lineWidth = 1;
                ctx.strokeRect(x * cellWidth, y * cellHeight + cellHeight - 15, cellWidth, 15);

                ctx.fillStyle = '#fff';
                ctx.font = '12px Arial';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillText(char.name, x * cellWidth + cellWidth / 2, y * cellHeight + cellHeight - 7.5);

                if (char.id === data.currentTurn) {
                    ctx.strokeStyle = 'yellow';
                    ctx.lineWidth = 3;
                    ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    ctx.lineWidth = 1;
                }
            }
        }
    }

    // Отрисовка движущегося персонажа
    if (movingCharacter) {
        const char = movingCharacter;
        const teamIconURL = data.teamsConfig[char.team].iconURL;
        const charImage = loadImage(teamIconURL, DEFAULT_IMAGES.icon);
        if (charImage && charImage.complete && charImage.naturalWidth !== 0) {
            const iconSize = 40;
            const offsetX = (cellWidth - iconSize) / 2;
            const offsetY = (cellHeight - 15 - iconSize) / 2;
            ctx.drawImage(charImage, char.currentX + offsetX, char.currentY + offsetY, iconSize, iconSize);
        }
        ctx.fillStyle = char.team === 0 ? 'rgba(128, 0, 128, 0.8)' : 'rgba(255, 215, 0, 0.8)';
        ctx.fillRect(char.currentX, char.currentY + cellHeight - 15, cellWidth, 15);
        ctx.strokeStyle = char.team === 0 ? '#800080' : '#ffd700';
        ctx.strokeRect(char.currentX, char.currentY + cellHeight - 15, cellWidth, 15);
        ctx.fillStyle = '#fff';
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(char.name, char.currentX + cellWidth / 2, char.currentY + cellHeight - 7.5);
    }

    // Отрисовка перетаскиваемого персонажа
    if (draggingCharacter) {
        const char = draggingCharacter;
        const teamIconURL = data.teamsConfig[char.team].iconURL;
        const charImage = loadImage(teamIconURL, DEFAULT_IMAGES.icon);
        if (charImage && charImage.complete && charImage.naturalWidth !== 0) {
            const iconSize = 40;
            const offsetX = (cellWidth - iconSize) / 2;
            const offsetY = (cellHeight - 15 - iconSize) / 2;
            ctx.drawImage(charImage, dragOffsetX - cellWidth / 2 + offsetX, dragOffsetY - cellHeight / 2 + offsetY, iconSize, iconSize);
        }
        ctx.fillStyle = char.team === 0 ? 'rgba(128, 0, 128, 0.8)' : 'rgba(255, 215, 0, 0.8)';
        ctx.fillRect(dragOffsetX - cellWidth / 2, dragOffsetY - cellHeight / 2 + cellHeight - 15, cellWidth, 15);
        ctx.strokeStyle = char.team === 0 ? '#800080' : '#ffd700';
        ctx.strokeRect(dragOffsetX - cellWidth / 2, dragOffsetY - cellHeight / 2 + cellHeight - 15, cellWidth, 15);
        ctx.fillStyle = '#fff';
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(char.name, dragOffsetX, dragOffsetY - cellHeight / 2 + cellHeight - 7.5);
    }
}

function drawArrow(fromX, fromY, toX, toY, color = 'blue', animate = false) {
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

    if (animate) {
        const pulse = Math.sin(Date.now() / 200) * 0.2 + 0.8;
        ctx.lineWidth = 2 * pulse;
        ctx.strokeStyle = `${color === 'blue' ? 'rgba(0, 0, 255, 0.7)' : color === 'red' ? 'rgba(255, 0, 0, 0.7)' : 'rgba(255, 215, 0, 0.7)'}`;
        ctx.beginPath();
        ctx.moveTo(fromX, fromY);
        ctx.lineTo(toX, toY);
        ctx.stroke();
    }
}

export function animateMove(character, path, callback, isAttack = false) {
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