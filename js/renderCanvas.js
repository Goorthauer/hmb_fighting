import { gameState, draggingCharacter, dragOffsetX, dragOffsetY, movePath, setMovePath, selectedAbility } from './state.js'; // Добавлен selectedAbility
import { ctx, cellWidth, cellHeight } from './constants.js';
import { findCharacter } from './utils.js';

let movingCharacter = null;

export function drawBoard(data) {
    if (!ctx) {
        console.error('Canvas context (ctx) is not initialized');
        return;
    }
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

    // Рисуем сетку
    for (let x = 0; x < 20; x++) {
        for (let y = 0; y < 10; y++) {
            ctx.strokeStyle = '#ccc';
            ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
        }
    }

    if (!data || !data.Teams || !Array.isArray(data.Teams) || !data.Board) {
        console.error('Invalid data in drawBoard:', data);
        return;
    }

    // Подсветка зон при перетаскивании
    if (draggingCharacter) {
        const startX = draggingCharacter.Position[0];
        const startY = draggingCharacter.Position[1];
        const stamina = draggingCharacter.Stamina;
        const isTwoHanded = (draggingCharacter.Weapon === 'two_handed_halberd' || draggingCharacter.Weapon === 'two_handed_sword');
        const weaponRange = isTwoHanded ? 2 : 1;
        const gridX = Math.floor(dragOffsetX / cellWidth);
        const gridY = Math.floor(dragOffsetY / cellHeight);

        if (data.Phase === 'move') {
            for (let x = Math.max(0, startX - stamina); x <= Math.min(19, startX + stamina); x++) {
                for (let y = Math.max(0, startY - stamina); y <= Math.min(9, startY + stamina); y++) {
                    const dist = Math.abs(x - startX) + Math.abs(y - startY);
                    if (dist <= stamina && data.Board[x][y] === -1) {
                        ctx.fillStyle = 'rgba(0, 255, 0, 0.2)';
                        ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                        ctx.strokeStyle = 'green';
                        ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                    }
                }
            }
            if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10 && data.Board[gridX][gridY] === -1) {
                drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'blue');
            }
        } else if (data.Phase === 'action') {
            if (selectedAbility) {
                for (let x = Math.max(0, startX - 1); x <= Math.min(19, startX + 1); x++) {
                    for (let y = Math.max(0, startY - 1); y <= Math.min(9, startY + 1); y++) {
                        const dist = Math.max(Math.abs(x - startX), Math.abs(y - startY));
                        if (dist <= 1 && data.Board[x][y] !== -1) {
                            const target = findCharacter(data.Teams, data.Board[x][y]);
                            if (target && target.Team !== draggingCharacter.Team) {
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
                        const isValidRange = (isTwoHanded && dist === weaponRange) || (!isTwoHanded && dist <= weaponRange);
                        if (isValidRange && data.Board[x][y] !== -1) {
                            const target = findCharacter(data.Teams, data.Board[x][y]);
                            if (target && target.Team !== draggingCharacter.Team) {
                                ctx.fillStyle = 'rgba(255, 0, 0, 0.2)';
                                ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                                ctx.strokeStyle = 'red';
                                ctx.strokeRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
                            }
                        }
                    }
                }
            }
            if (gridX >= 0 && gridX < 20 && gridY >= 0 && gridY < 10 && data.Board[gridX][gridY] !== -1) {
                const target = findCharacter(data.Teams, data.Board[gridX][gridY]);
                if (target && target.Team !== draggingCharacter.Team) {
                    const attackDist = Math.max(Math.abs(gridX - startX), Math.abs(gridY - startY));
                    if ((isTwoHanded && attackDist === weaponRange) || (!isTwoHanded && attackDist <= weaponRange)) {
                        drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'red', true);
                    }
                }
            }
        }
    }

    // Рисуем персонажей
    for (let x = 0; x < 20; x++) {
        for (let y = 0; y < 10; y++) {
            const charId = data.Board[x][y];
            if (charId !== -1 && (!movingCharacter || movingCharacter.ID !== charId)) {
                const char = findCharacter(data.Teams, charId);
                if (!char) {
                    console.warn(`Character with ID ${charId} not found in Teams`);
                    continue;
                }
                ctx.fillStyle = char.Team === 0 ? '#800080' : '#FFD700';
                ctx.fillRect(x * cellWidth + 5, y * cellHeight + 5, cellWidth - 10, cellHeight - 10);
                ctx.fillStyle = '#fff';
                ctx.font = '12px Arial';
                ctx.fillText(`${char.Name} (${char.HP})`, x * cellWidth + 5, y * cellHeight + 20);
                if (char.ID === data.CurrentTurn) {
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
        ctx.fillStyle = char.Team === 0 ? '#800080' : '#FFD700';
        ctx.fillRect(char.currentX + 5, char.currentY + 5, cellWidth - 10, cellHeight - 10);
        ctx.fillStyle = '#fff';
        ctx.fillText(`${char.Name} (${char.HP})`, char.currentX + 5, char.currentY + 20);
    }

    // Рисуем перетаскиваемого персонажа
    if (draggingCharacter) {
        const char = draggingCharacter;
        ctx.fillStyle = char.Team === 0 ? '#800080' : '#FFD700';
        ctx.fillRect(dragOffsetX - cellWidth / 2 + 5, dragOffsetY - cellHeight / 2 + 5, cellWidth - 10, cellHeight - 10);
        ctx.fillStyle = '#fff';
        ctx.fillText(`${char.Name} (${char.HP})`, dragOffsetX - cellWidth / 2 + 5, dragOffsetY - cellHeight / 2 + 20);
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