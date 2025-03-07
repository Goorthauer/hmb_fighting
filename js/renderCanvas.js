import { gameState, draggingCharacter, dragOffsetX, dragOffsetY, movePath, setMovePath, selectedAbility } from './state.js';
import { ctx, cellWidth, cellHeight } from './constants.js';
import { findCharacter } from './utils.js';

let movingCharacter = null;
export const imagesCache = {}; // Кэш для изображений, теперь экспортируем

// Пути к изображениям по умолчанию
const DEFAULT_IMAGES = {
    character: '/static/characters/default.png',
    ability: '/static/abilities/default.jpg',
    weapon: '/static/weapons/default.png',
    shield: '/static/shields/default.png',
    icon: '/static/icons/default.png'
};

// Загрузка изображения с кэшированием и запасным вариантом
function loadImage(url, defaultUrl) {
    if (!url || url.trim() === '') {
        console.warn(`Image URL is empty or invalid, using default: ${defaultUrl}`);
        url = defaultUrl;
    }
    if (!imagesCache[url]) {
        const img = new Image();
        img.src = url;
        img.onerror = () => {
            console.warn(`Failed to load image: ${url}, falling back to ${defaultUrl}`);
            img.src = defaultUrl;
        };
        imagesCache[url] = img;
    }
    return imagesCache[url];
}

// Предварительная загрузка всех изображений из gameState
export function preloadImages(data) {
    if (!data || !data.teams || !data.weaponsConfig || !data.shieldsConfig || !data.teamsConfig) {
        console.error('Invalid data for preloading images:', data);
        return;
    }

    data.teamsConfig.forEach((teamConfig, index) => {
        loadImage(teamConfig.iconURL, DEFAULT_IMAGES.icon);
    });

    data.teams.forEach((team, teamIndex) => {
        team.characters.forEach((char, charIndex) => {
            loadImage(char.imageURL, DEFAULT_IMAGES.character);
            char.abilities.forEach((ability, abilityIndex) => {
                loadImage(ability.imageURL, DEFAULT_IMAGES.ability);
            });
        });
    });

    Object.entries(data.weaponsConfig).forEach(([key, weapon]) => {
        loadImage(weapon.imageURL, DEFAULT_IMAGES.weapon);
    });

    Object.entries(data.shieldsConfig).forEach(([key, shield]) => {
        loadImage(shield.imageURL, DEFAULT_IMAGES.shield);
    });

}

export function drawBoard(data) {
    if (!ctx) {
        console.error('Canvas context (ctx) is not initialized');
        return;
    }
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    for (let x = 0; x < 20; x++) {
        for (let y = 0; y < 10; y++) {
            ctx.strokeStyle = '#ccc';
            ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        }
    }

    if (!data || !data.teams || !Array.isArray(data.teams) || !data.board) {
        console.error('Invalid data in drawBoard:', data);
        return;
    }

    // Подсветка зон при перетаскивании (без изменений)
    if (draggingCharacter) {
        const startX = draggingCharacter.position[0];
        const startY = draggingCharacter.position[1];
        const stamina = draggingCharacter.stamina;
        const weapon = data.weaponsConfig[draggingCharacter.weapon];
        const weaponRange = weapon ? weapon.range : 1;
        const gridX = Math.floor(dragOffsetX / cellWidth);
        const gridY = Math.floor(dragOffsetY / cellHeight);

        if (data.phase === 'move') {
            for (let x = Math.max(0, startX - stamina); x <= Math.min(19, startX + stamina); x++) {
                for (let y = Math.max(0, startY - stamina); y <= Math.min(9, startY + stamina); y++) {
                    const dist = Math.abs(x - startX) + Math.abs(y - startY);
                    if (dist <= stamina && data.board[x][y] === -1) {
                        ctx.fillStyle = 'rgba(0, 255, 0, 0.2)';
                        ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                        ctx.strokeStyle = 'green';
                        ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    }
                }
            }
            if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10 && data.board[gridX][gridY] === -1) {
                drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'blue');
            }
        } else if (data.phase === 'action') {
            if (selectedAbility) {
                const abilityRange = selectedAbility.range || 1;
                for (let x = Math.max(0, startX - abilityRange); x <= Math.min(19, startX + abilityRange); x++) {
                    for (let y = Math.max(0, startY - abilityRange); y <= Math.min(9, startY + abilityRange); y++) {
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
                for (let x = Math.max(0, startX - weaponRange); x <= Math.min(19, startX + weaponRange); x++) {
                    for (let y = Math.max(0, startY - weaponRange); y <= Math.min(9, startY + weaponRange); y++) {
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
            if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10 && data.board[gridX][gridY] !== -1) {
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
    // Рисуем персонажей
    for (let x = 0; x < 20; x++) {
        for (let y = 0; y < 10; y++) {
            const charId = data.board[x][y];
            if (charId !== -1 && (!movingCharacter || movingCharacter.id !== charId)) {
                const char = findCharacter(data.teams, charId);
                if (!char) {
                    console.warn(`Character with ID ${charId} not found in teams`);
                    continue;
                }

                // Запасной вариант: рисуем прямоугольник
                ctx.fillStyle = char.team === 0 ? '#800080' : '#FFD700';
                ctx.fillRect(x * cellWidth + 5, y * cellHeight + 5, cellWidth - 10, cellHeight - 10);

                // Иконка команды
                const teamIcon = imagesCache[data.teamsConfig[char.team].iconURL] || imagesCache[DEFAULT_IMAGES.icon];
                if (teamIcon && teamIcon.complete && teamIcon.naturalWidth !== 0) {
                    ctx.drawImage(teamIcon, x * cellWidth + 2, y * cellHeight + 2, cellWidth - 4, cellHeight - 4);
                } else {
                    console.warn(`Team icon for ${char.name} is not valid:`, teamIcon);
                }

                // Изображение персонажа
                const charImage = imagesCache[char.imageURL] || imagesCache[DEFAULT_IMAGES.character];
                if (charImage && charImage.complete && charImage.naturalWidth !== 0) {
                    ctx.drawImage(charImage, x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                } else {
                    console.warn(`Character image for ${char.name} is not valid:`, charImage);
                }

                // Оружие и щит
                const weapon = data.weaponsConfig[char.weapon];
                if (weapon && weapon.imageURL) {
                    const weaponImage = imagesCache[weapon.imageURL] || imagesCache[DEFAULT_IMAGES.weapon];
                    if (weaponImage && weaponImage.complete && weaponImage.naturalWidth !== 0) {
                        ctx.drawImage(weaponImage, x * cellWidth + cellWidth - 15, y * cellHeight, 15, 15);
                    }
                }
                const shield = data.shieldsConfig[char.shield];
                if (shield && shield.imageURL) {
                    const shieldImage = imagesCache[shield.imageURL] || imagesCache[DEFAULT_IMAGES.shield];
                    if (shieldImage && shieldImage.complete && shieldImage.naturalWidth !== 0) {
                        ctx.drawImage(shieldImage, x * cellWidth, y * cellHeight + cellHeight - 15, 15, 15);
                    }
                }

                ctx.fillStyle = '#fff';
                ctx.font = '12px Arial';
                ctx.fillText(`${char.name} (${char.hp})`, x * cellWidth + 5, y * cellHeight + 20);

                if (char.id === data.currentTurn) {
                    ctx.strokeStyle = 'yellow';
                    ctx.lineWidth = 3;
                    ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    ctx.lineWidth = 1;
                }
            }
        }
    }

    // Рисуем движущегося персонажа
    if (movingCharacter) {
        const char = movingCharacter;
        ctx.fillStyle = char.team === 0 ? '#800080' : '#FFD700';
        ctx.fillRect(char.currentX + 5, char.currentY + 5, cellWidth - 10, cellHeight - 10);

        const teamIcon = imagesCache[data.teamsConfig[char.team].iconURL] || imagesCache[DEFAULT_IMAGES.icon];
        if (teamIcon && teamIcon.complete) {
            ctx.drawImage(teamIcon, char.currentX + 2, char.currentY + 2, cellWidth - 4, cellHeight - 4);
        }
        const charImage = imagesCache[char.imageURL] || imagesCache[DEFAULT_IMAGES.character];
        if (charImage && charImage.complete) {
            ctx.drawImage(charImage, char.currentX, char.currentY, cellWidth, cellHeight);
        }
        ctx.fillStyle = '#fff';
        ctx.fillText(`${char.name} (${char.hp})`, char.currentX + 5, char.currentY + 20);
    }

    // Рисуем перетаскиваемого персонажа
    if (draggingCharacter) {
        const char = draggingCharacter;
        ctx.fillStyle = char.team === 0 ? '#800080' : '#FFD700';
        ctx.fillRect(dragOffsetX - cellWidth / 2 + 5, dragOffsetY - cellHeight / 2 + 5, cellWidth - 10, cellHeight - 10);

        const teamIcon = imagesCache[data.teamsConfig[char.team].iconURL] || imagesCache[DEFAULT_IMAGES.icon];
        if (teamIcon && teamIcon.complete) {
            ctx.drawImage(teamIcon, dragOffsetX - cellWidth / 2 + 2, dragOffsetY - cellHeight / 2 + 2, cellWidth - 4, cellHeight - 4);
        }
        const charImage = imagesCache[char.imageURL] || imagesCache[DEFAULT_IMAGES.character];
        if (charImage && charImage.complete) {
            ctx.drawImage(charImage, dragOffsetX - cellWidth / 2, dragOffsetY - cellHeight / 2, cellWidth, cellHeight);
        }
        ctx.fillStyle = '#fff';
        ctx.fillText(`${char.name} (${char.hp})`, dragOffsetX - cellWidth / 2 + 5, dragOffsetY - cellHeight / 2 + 20);
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
    const arrowSize = 10;
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
    const speed = 5;
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