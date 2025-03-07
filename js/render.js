import { gameState, draggingCharacter, dragOffsetX, dragOffsetY, selectedCharacter, selectedAbility, cellWidth, cellHeight, canvas, ctx, movePath, setMovePath, setSelectedCharacter } from './state.js';

let movingCharacter = null;
let previousTurn = null;
let turnOrder = [];
let roundStarted = false;
let currentGameSessionId = null;
let isUpdatingCards = false; // Флаг для предотвращения дублирующих обновлений карточек
let isUpdatingAbilities = false; // Флаг для способностей

function saveTurnState(gameSessionId) {
    const state = {
        turnOrder: turnOrder,
        roundStarted: roundStarted,
        previousTurn: previousTurn
    };
    localStorage.setItem(`turnState_${gameSessionId}`, JSON.stringify(state));
}

function loadTurnState(gameSessionId) {
    const savedState = localStorage.getItem(`turnState_${gameSessionId}`);
    if (savedState) {
        const state = JSON.parse(savedState);
        turnOrder = state.turnOrder || [];
        roundStarted = state.roundStarted || false;
        previousTurn = state.previousTurn || null;
        currentGameSessionId = gameSessionId;
        return true;
    }
    return false;
}

function resetTurnState() {
    turnOrder = [];
    roundStarted = false;
    previousTurn = null;
    currentGameSessionId = null;
}

export function drawBoard(data) {
    if (!ctx) {
        console.error('Canvas context (ctx) is not initialized');
        return;
    }
    ctx.clearRect(0, 0, canvas.width, canvas.height);

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
                    const abilityDist = Math.max(Math.abs(gridX - startX), Math.abs(gridY - startY));
                    if (selectedAbility) {
                        if (abilityDist <= 1) {
                            drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'gold', true);
                        }
                    } else {
                        const isValidRange = (isTwoHanded && attackDist === weaponRange) || (!isTwoHanded && attackDist <= weaponRange);
                        if (isValidRange) {
                            drawArrow(startX * cellWidth + cellWidth / 2, startY * cellHeight + cellHeight / 2, gridX * cellWidth + cellWidth / 2, gridY * cellHeight + cellHeight / 2, 'red', true);
                        }
                    }
                }
            }
        }
    }

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

    if (movingCharacter) {
        const char = movingCharacter;
        ctx.fillStyle = char.Team === 0 ? '#800080' : '#FFD700';
        ctx.fillRect(char.currentX + 5, char.currentY + 5, cellWidth - 10, cellHeight - 10);
        ctx.fillStyle = '#fff';
        ctx.fillText(`${char.Name} (${char.HP})`, char.currentX + 5, char.currentY + 20);
    }

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

export function updatePhaseAndProgress(data) {
    const phaseContainer = document.getElementById('phaseContainer');
    if (!phaseContainer) return;

    const team0Alive = data.Teams[0].Characters.filter(c => c.HP > 0).length;
    const team1Alive = data.Teams[1].Characters.filter(c => c.HP > 0).length;
    const currentChar = findCharacter(data.Teams, data.CurrentTurn);
    const phaseText = data.Phase === 'move' ? 'Move' : 'Action';

    phaseContainer.querySelector('.phase').textContent = `${phaseText} - ${currentChar ? currentChar.Name + ' team ' + currentChar.Team : 'Unknown'}`;
    phaseContainer.querySelector('.team0').textContent = `⚔️${team0Alive}`;
    phaseContainer.querySelector('.team1').textContent = `${team1Alive}🛡️`;

    phaseContainer.classList.remove('team0-turn', 'team1-turn');
    if (currentChar) {
        phaseContainer.classList.add(`team${currentChar.Team}-turn`);
        phaseContainer.style.display = 'none';
        setTimeout(() => phaseContainer.style.display = 'flex', 0);
    } else {
        console.warn('No current character found, unable to set team class');
    }
}

export function updateCharacterCards(setSelectedCharacter, data) {
    if (isUpdatingCards) {
        console.log('Skipping updateCharacterCards - already in progress');
        return;
    }
    isUpdatingCards = true;

    const characterCards = document.getElementById('characterCards');
    if (!characterCards) {
        console.error('characterCards element not found');
        isUpdatingCards = false;
        return;
    }

    characterCards.innerHTML = '';
    const gameSessionId = data.GameSessionId || 'default_session';

    if (currentGameSessionId && currentGameSessionId !== gameSessionId) {
        resetTurnState();
    }

    let allChars = [];
    data.Teams.forEach(team => {
        team.Characters.forEach(char => {
            if (char.HP > 0) {
                allChars.push(char);
            }
        });
    });

    if (!roundStarted || turnOrder.length === 0) {
        if (loadTurnState(gameSessionId)) {
            console.log('Loaded turn state from localStorage:', { turnOrder, previousTurn });
        } else {
            allChars.sort((a, b) => b.Initiative - a.Initiative);
            turnOrder = allChars.map(char => char.ID);
            roundStarted = true;
            currentGameSessionId = gameSessionId; // Исправлено с currentGameSessionیبId
        }
    }

    const currentTurnIndex = turnOrder.indexOf(data.CurrentTurn);
    if (currentTurnIndex !== -1) {
        turnOrder.splice(currentTurnIndex, 1);
        turnOrder.unshift(data.CurrentTurn);
    }

    if (previousTurn !== null && previousTurn !== data.CurrentTurn) {
        const prevTurnIndex = turnOrder.indexOf(previousTurn);
        if (prevTurnIndex !== -1) {
            turnOrder.splice(prevTurnIndex, 1);
            turnOrder.push(previousTurn);
        }
    }

    if (turnOrder[0] === allChars[0].ID && previousTurn !== null && turnOrder.length === allChars.length) {
        allChars.sort((a, b) => b.Initiative - a.Initiative);
        turnOrder = allChars.map(char => char.ID);
    }

    let orderedChars = [];
    const walked = [];
    const notWalked = [];

    allChars.forEach(char => {
        const indexInTurnOrder = turnOrder.indexOf(char.ID);
        if (char.ID === data.CurrentTurn) {
            orderedChars[0] = char;
        } else if (indexInTurnOrder > 0 && indexInTurnOrder < turnOrder.length - walked.length) {
            notWalked.push(char);
        } else if (char.ID !== data.CurrentTurn) {
            walked.push(char);
        }
    });

    notWalked.sort((a, b) => b.Initiative - a.Initiative);
    walked.sort((a, b) => b.Initiative - a.Initiative);

    orderedChars = [orderedChars[0], ...notWalked, ...walked];

    if (previousTurn !== null && previousTurn !== data.CurrentTurn) {
        const prevTurnIndex = orderedChars.findIndex(char => char.ID === previousTurn);
        if (prevTurnIndex !== -1) {
            animateInitiativeShift(characterCards, orderedChars, prevTurnIndex, data, setSelectedCharacter);
        } else {
            renderCards(characterCards, orderedChars, data, setSelectedCharacter);
        }
    } else {
        renderCards(characterCards, orderedChars, data, setSelectedCharacter);
    }

    previousTurn = data.CurrentTurn;
    saveTurnState(gameSessionId);

    isUpdatingCards = false;
}

function renderCards(container, chars, data, setSelectedCharacter) {
    chars.forEach((char, index) => {
        const card = document.createElement('div');
        card.classList.add('card');
        card.classList.add(`team${char.Team}`);
        card.dataset.id = char.ID;
        if (selectedCharacter && selectedCharacter.ID === char.ID) {
            card.classList.add('selected');
        }
        if (char.ID === data.CurrentTurn) {
            card.classList.add('current');
        }
        if (char.HP <= 0) {
            card.classList.add('dead');
        }

        const image = document.createElement('div');
        image.classList.add('image');
        image.style.backgroundImage = `url('./image/char_${char.ID % 5 || 5}.png')`;
        card.appendChild(image);

        const name = document.createElement('div');
        name.classList.add('name');
        name.textContent = char.Name;
        image.appendChild(name);

        const hpContainer = document.createElement('div');
        hpContainer.classList.add('hp-container');
        const hpDiamond = document.createElement('div');
        hpDiamond.classList.add('hp-diamond');
        const hp = document.createElement('div');
        hp.classList.add('hp');
        hp.textContent = char.HP;
        hpDiamond.appendChild(hp);
        hpContainer.appendChild(hpDiamond);
        image.appendChild(hpContainer);

        const info = document.createElement('div');
        info.classList.add('info');
        info.innerHTML = `
            <div class="stat"><span class="label">Stamina:</span> ${char.Stamina}</div>
            <div class="stat"><span class="label">Attack:</span> ${char.AttackMin}-${char.AttackMax}</div>
            <div class="stat"><span class="label">Defense:</span> ${char.Defense}</div>
            <div class="stat"><span class="label">Init:</span> ${char.Initiative}</div>
            <div class="stat"><span class="label">Weapon:</span> ${char.Weapon || 'None'}</div>
            <div class="stat"><span class="label">Shield:</span> ${char.Shield || 'None'}</div>
        `;
        card.appendChild(info);

        card.addEventListener('click', () => {
            if (char.HP > 0) {
                setSelectedCharacter(char);
                updateCharacterCards(setSelectedCharacter, data);
            }
        });

        container.appendChild(card);
    });
}

function animateInitiativeShift(container, chars, prevTurnIndex, data, setSelectedCharacter) {
    container.innerHTML = '';
    renderCards(container, chars, data, setSelectedCharacter);
    const cards = container.children;
    const cardWidth = 140;

    for (let i = 0; i < cards.length; i++) {
        cards[i].style.transition = 'none';
        cards[i].style.transform = `translateX(${i * cardWidth}px)`;
    }

    requestAnimationFrame(() => {
        const prevCard = cards[prevTurnIndex];
        prevCard.style.transition = 'transform 0.5s ease';
        prevCard.style.transform = `translateX(${(cards.length - 1) * cardWidth}px)`;

        prevCard.addEventListener('transitionend', () => {
            container.innerHTML = '';
            const updatedChars = Array.from(chars);
            const [movedChar] = updatedChars.splice(prevTurnIndex, 1);
            updatedChars.push(movedChar);
            renderCards(container, updatedChars, data, setSelectedCharacter);
        }, { once: true });
    });
}

export function updateAbilityCards(myTeam, setSelectedAbility, data) {
    if (isUpdatingAbilities) {
        console.log('Skipping updateAbilityCards - already in progress');
        return;
    }
    isUpdatingAbilities = true;

    const abilityCards = document.getElementById('abilityCards');
    if (!abilityCards) {
        console.error('abilityCards element not found');
        isUpdatingAbilities = false;
        return;
    }
    abilityCards.innerHTML = '';

    if (!data || !data.Teams || !Array.isArray(data.Teams)) {
        isUpdatingAbilities = false;
        return;
    }

    const currentChar = findCharacter(data.Teams, data.CurrentTurn);
    let relevantChar = null;
    if (currentChar && currentChar.Team === myTeam) {
        relevantChar = currentChar;
    } else {
        const myTeamChars = data.Teams[myTeam].Characters.filter(c => c.HP > 0);
        myTeamChars.sort((a, b) => b.Initiative - a.Initiative);
        const currentIndex = myTeamChars.findIndex(c => c.ID === data.CurrentTurn);
        relevantChar = myTeamChars[(currentIndex + 1) % myTeamChars.length] || myTeamChars[0];
    }

    if (relevantChar && Array.isArray(relevantChar.Abilities) && relevantChar.Abilities.length > 0) {
        relevantChar.Abilities.forEach(ability => {
            const card = document.createElement('div');
            card.classList.add('ability-card');
            if (selectedAbility && selectedAbility.Name === ability.Name) {
                card.classList.add('selected');
            }

            const image = document.createElement('div');
            image.classList.add('image');
            image.style.backgroundImage = `url('./image/ability_${ability.Name.toLowerCase()}.png')`;
            card.appendChild(image);

            const info = document.createElement('div');
            info.classList.add('info');
            info.innerHTML = `
                <strong>${ability.Name}</strong><br>
                ${ability.Description || 'No description available'}
            `;
            card.appendChild(info);

            card.addEventListener('click', () => {
                setSelectedAbility(ability);
                updateAbilityCards(myTeam, setSelectedAbility, data);
            });

            abilityCards.appendChild(card);
        });
    } else {
        for (let i = 0; i < 3; i++) {
            const card = document.createElement('div');
            card.classList.add('ability-card', 'placeholder');

            const image = document.createElement('div');
            image.classList.add('image');
            card.appendChild(image);

            const info = document.createElement('div');
            info.classList.add('info');
            info.innerHTML = `<strong>No Ability</strong>`;
            card.appendChild(info);

            abilityCards.appendChild(card);
        }
    }

    isUpdatingAbilities = false;
}

export function updateUI(myTeam) {
    const uiElements = document.querySelectorAll('.team-specific');
    uiElements.forEach(el => {
        el.style.display = (parseInt(el.dataset.team) === myTeam) ? 'block' : 'none';
    });
}

export function updateBattleLog(data, previousState) {
    const logEntries = document.getElementById('logEntries');
    if (!logEntries || !data || !data.Teams || !Array.isArray(data.Teams)) return;

    if (!previousState || !previousState.Teams || !Array.isArray(previousState.Teams)) return;

    const prevChar = findCharacter(previousState.Teams, previousState.CurrentTurn);
    const currChar = findCharacter(data.Teams, data.CurrentTurn);

    if (previousState.Phase !== data.Phase) {
        addLogEntry(`Phase changed to ${data.Phase} for ${currChar ? currChar.Name : 'Unknown'}`);
    }

    for (let teamIdx = 0; teamIdx < 2; teamIdx++) {
        const prevTeam = previousState.Teams[teamIdx].Characters;
        const currTeam = data.Teams[teamIdx].Characters;
        for (let i = 0; i < prevTeam.length; i++) {
            const prev = prevTeam[i];
            const curr = currTeam[i];
            if (prev.Position[0] !== curr.Position[0] || prev.Position[1] !== curr.Position[1]) {
                addLogEntry(`${curr.Name} moved from (${prev.Position[0]}, ${prev.Position[1]}) to (${curr.Position[0]}, ${curr.Position[1]})`);
            }
            if (prev.HP !== curr.HP && curr.HP > 0) {
                addLogEntry(`${curr.Name} took ${prev.HP - curr.HP} damage (HP: ${curr.HP})`);
            }
            if (prev.HP > 0 && curr.HP <= 0) {
                addLogEntry(`${curr.Name} was defeated`);
            }
        }
    }
}

export function addLogEntry(message) {
    const logEntries = document.getElementById('logEntries');
    if (!logEntries) return;
    const entry = document.createElement('div');
    entry.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
    logEntries.appendChild(entry);
    logEntries.scrollTop = logEntries.scrollHeight;
}

export function findCharacter(teams, id) {
    if (!teams || !Array.isArray(teams)) return null;
    for (let team of teams) {
        for (let char of team.Characters) {
            if (char.ID === id) return char;
        }
    }
    return null;
}