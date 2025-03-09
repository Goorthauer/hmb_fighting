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

export function drawBoard(data) {
    if (!ctx) {
        console.error('Canvas context (ctx) is not initialized');
        return;
    }

    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    // Отрисовка сетки
    for (let x = 0; x < 16; x++) {
        for (let y = 0; y < 9; y++) {
            if (data.phase === 'setup') {
                const myTeam = data.teamID;
                if ((myTeam === 0 && x < 8) || (myTeam === 1 && x >= 8)) {
                    ctx.fillStyle = 'rgba(0, 255, 0, 0.1)';
                    ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                } else {
                    ctx.fillStyle = 'rgba(255, 0, 0, 0.1)';
                    ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                }
            }
            ctx.strokeStyle = '#ccc';
            ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        }
    }

    if (!data || !data.teams || !Array.isArray(data.teams) || !data.board || !data.teamsConfig) {
        console.error('Invalid data in drawBoard:', data);
        return;
    }

    // Подсветка доступных ходов/атак
    const currentChar = findCharacter(data.teams, data.currentTurn);
    if (currentChar && currentChar.team === data.teamID) {
        const startX = currentChar.position[0];
        const startY = currentChar.position[1];

        if (data.phase === 'move') {
            const stamina = currentChar.stamina;
            for (let x = Math.max(0, startX - stamina); x <= Math.min(15, startX + stamina); x++) {
                for (let y = Math.max(0, startY - stamina); y <= Math.min(8, startY + stamina); y++) {
                    const dist = Math.abs(x - startX) + Math.abs(y - startY);
                    if (dist <= stamina && data.board[x][y] === -1) {
                        ctx.fillStyle = 'rgba(0, 255, 0, 0.2)';
                        ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    }
                }
            }
        } else if (data.phase === 'action') {
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
    }

    // Отрисовка персонажей
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

                if (char.id === data.currentTurn) {
                    ctx.strokeStyle = 'yellow';
                    ctx.lineWidth = 3;
                    ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    ctx.lineWidth = 1;
                }
            }
        }
    }

    // Отрисовка перетаскиваемого персонажа и стрелки
    if (draggingCharacter) {
        const startX = draggingCharacter.position[0];
        const startY = draggingCharacter.position[1];
        const gridX = Math.floor(dragOffsetX / cellWidth);
        const gridY = Math.floor(dragOffsetY / cellHeight);

        if (data.phase === 'setup') {
            const myTeam = data.teamID;
            const isValidZone = (myTeam === 0 && gridX < 8) || (myTeam === 1 && gridX >= 8);
            if (isValidZone && gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && data.board[gridX][gridY] === -1) {
                ctx.fillStyle = 'rgba(0, 255, 0, 0.3)';
                ctx.fillRect(gridX * cellWidth, gridY * cellHeight, cellWidth, cellHeight);
            }
        } else if (data.phase === 'move' && gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && data.board[gridX][gridY] === -1) {
            drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'blue');
        } else if (data.phase === 'action' && gridX >= 0 && gridX < 16 && gridY >= 0 && gridY < 9 && data.board[gridX][gridY] !== -1) {
            const target = findCharacter(data.teams, data.board[gridX][gridY]);
            if (target && target.team !== data.teamID) {
                drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'red');
            }
        }

        const char = draggingCharacter;
        const teamIconURL = data.teamsConfig[char.team].iconURL;
        const charImage = loadImage(teamIconURL, DEFAULT_IMAGES.icon);
        if (charImage && charImage.complete) {
            const iconSize = 40;
            const offsetX = (cellWidth - iconSize) / 2;
            const offsetY = (cellHeight - 15 - iconSize) / 2;
            ctx.drawImage(charImage, dragOffsetX - cellWidth / 2 + offsetX, dragOffsetY - cellHeight / 2 + offsetY, iconSize, iconSize);
        }
        ctx.fillStyle = char.team === 0 ? 'rgba(128, 0, 128, 0.8)' : 'rgba(255, 215, 0, 0.8)';
        ctx.fillRect(dragOffsetX - cellWidth / 2, dragOffsetY - cellHeight / 2 + cellHeight - 15, cellWidth, 15);
        ctx.fillStyle = '#fff';
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(char.name, dragOffsetX, dragOffsetY - cellHeight / 2 + cellHeight - 7.5);
    }

    if (movingCharacter) {
        const char = movingCharacter;
        const teamIconURL = data.teamsConfig[char.team].iconURL;
        const charImage = loadImage(teamIconURL, DEFAULT_IMAGES.icon);
        if (charImage && charImage.complete) {
            const iconSize = 40;
            const offsetX = (cellWidth - iconSize) / 2;
            const offsetY = (cellHeight - 15 - iconSize) / 2;
            ctx.drawImage(charImage, char.currentX + offsetX, char.currentY + offsetY, iconSize, iconSize);
        }
        ctx.fillStyle = char.team === 0 ? 'rgba(128, 0, 128, 0.8)' : 'rgba(255, 215, 0, 0.8)';
        ctx.fillRect(char.currentX, char.currentY + cellHeight - 15, cellWidth, 15);
        ctx.fillStyle = '#fff';
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(char.name, char.currentX + cellWidth / 2, char.currentY + cellHeight - 7.5);
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
        ctx.strokeStyle = `${color === 'blue' ? 'rgba(0, 0, 255, 0.7)' : 'rgba(255, 0, 0, 0.7)'}`;
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