import {gameState, selectedCharacter, selectedAbility, setSelectedCharacter, setSelectedAbility} from './state.js';
import {findCharacter, addLogEntry} from './utils.js';

let previousState = null;
let turnOrder = [];
let roundStarted = false;
let currentGameSessionId = null;
let isUpdatingCards = false;
let isUpdatingAbilities = false;

function saveTurnState(gameSessionId) {
    const state = {turnOrder, roundStarted, previousTurn: previousState?.currentTurn};
    localStorage.setItem(`turnState_${gameSessionId}`, JSON.stringify(state));
}

function loadTurnState(gameSessionId) {
    const savedState = localStorage.getItem(`turnState_${gameSessionId}`);
    if (savedState) {
        const state = JSON.parse(savedState);
        turnOrder = state.turnOrder || [];
        roundStarted = state.roundStarted || false;
        previousState = state.previousTurn ? {currentTurn: state.previousTurn} : null;
        currentGameSessionId = gameSessionId;
        return true;
    }
    return false;
}

function resetTurnState() {
    turnOrder = [];
    roundStarted = false;
    currentGameSessionId = null;
}

export function updatePhaseAndProgress(data) {
    const phaseContainer = document.getElementById('phaseContainer');
    if (!phaseContainer) {
        console.error('Phase container not found');
        return;
    }

    const team0Alive = data.teams?.[0]?.characters.filter(c => c.hp > 0).length || 0;
    const team1Alive = data.teams?.[1]?.characters.filter(c => c.hp > 0).length || 0;
    const currentChar = findCharacter(data.teams, data.currentTurn);
    const phaseText = data.phase === 'move' ? 'Move' : 'Action';

    console.log('Updating phase and progress:', {
        currentTurn: data.currentTurn,
        currentChar,
        phaseText,
        team0Alive,
        team1Alive
    });

    phaseContainer.querySelector('.phase').textContent = `${phaseText} - ${currentChar ? currentChar.name + ' team ' + currentChar.team : 'Unknown'}`;
    phaseContainer.querySelector('.team0').textContent = `⚔️${team0Alive}`;
    phaseContainer.querySelector('.team1').textContent = `${team1Alive}🛡️`;

    phaseContainer.classList.remove('team0-turn', 'team1-turn');
    if (currentChar) {
        phaseContainer.classList.add(`team${currentChar.team}-turn`);
        phaseContainer.style.display = 'none';
        setTimeout(() => phaseContainer.style.display = 'flex', 0);
    }
}

export function updateCharacterCards(data) {
    if (isUpdatingCards) return;
    isUpdatingCards = true;

    const characterCards = document.getElementById('characterCards');
    if (!characterCards) {
        isUpdatingCards = false;
        return;
    }

    characterCards.innerHTML = '';
    const gameSessionId = data.gameSessionId || 'default_session';

    if (currentGameSessionId && currentGameSessionId !== gameSessionId) {
        resetTurnState();
    }

    let allChars = [];
    if (data.teams && Array.isArray(data.teams)) {
        data.teams.forEach(team => {
            team.characters.forEach(char => {
                allChars.push(char); // Убираем фильтр hp > 0, добавляем всех персонажей
            });
        });
    }

    if (!roundStarted || turnOrder.length === 0) {
        if (loadTurnState(gameSessionId)) {
            console.log('Loaded turn state:', {turnOrder, previousTurn: previousState?.currentTurn});
        } else {
            allChars.sort((a, b) => b.initiative - a.initiative);
            turnOrder = allChars.map(char => char.id);
            roundStarted = true;
            currentGameSessionId = gameSessionId;
        }
    }

    renderCards(characterCards, allChars, data);
    saveTurnState(gameSessionId);
    isUpdatingCards = false;
}

function renderCards(container, chars, data) {
    container.innerHTML = ''; // Очищаем контейнер

    // Сортируем персонажей: живые в начале, мёртвые в конце
    const sortedChars = chars.sort((a, b) => {
        if (a.hp <= 0 && b.hp > 0) return 1; // Мёртвые в конец
        if (a.hp > 0 && b.hp <= 0) return -1; // Живые в начало
        return 0; // Остальные не меняют порядок
    });

    sortedChars.forEach(char => {
        const card = document.createElement('div');
        card.classList.add('card', `team${char.team}`);
        card.dataset.id = char.id;
        if (char.id === data.currentTurn) card.classList.add('current');
        if (char.hp <= 0) card.classList.add('dead');

        card.innerHTML = `
            <div class="image" style="background-image: url('${char.imageURL || 'default-image.png'}');"></div>
            <div class="name">${char.name}</div>
            <div class="info">
                <div class="stat"><i class="fas fa-wheelchair-move"></i> <span class="label">Скорость:</span> ${char.stamina}</div>
                <div class="stat"><i class="fas fa-skull"></i> <span class="label">Атака:</span> ${char.attackMin}-${char.attackMax}</div>
                <div class="stat"><i class="fas fa-shield-alt"></i> <span class="label">Защита:</span> ${char.defense}</div>
                <div class="stat"><i class="fas fa-rocket"></i> <span class="label">Инициатива:</span> ${char.initiative}</div>
                <div class="stat"><i class="fas fa-gavel"></i> <span class="label"></span> ${data.weaponsConfig[char.weapon]?.display_name || 'None'}</div>
                <div class="stat"><i class="fas fa-shield"></i> <span class="label"></span> ${data.shieldsConfig[char.shield]?.display_name || 'None'}</div>
                <div class="hp-container"><div class="hp-diamond"><div class="hp">${char.hp}</div></div></div>
            </div>
        `;
        container.appendChild(card);
    });
}
export function updateAbilityCards(myTeam, data) {
    if (isUpdatingAbilities) return;
    isUpdatingAbilities = true;

    const abilityCards = document.getElementById('abilityCards');
    if (!abilityCards) {
        isUpdatingAbilities = false;
        return;
    }
    abilityCards.innerHTML = '';

    const currentChar = findCharacter(data.teams, data.currentTurn);
    if (currentChar) {
        if (currentChar.abilities?.length) {
            // Отрисовываем способности как полноценные карточки
            currentChar.abilities.forEach((abilityID) => {
                const ability = data.abilitiesConfig[abilityID];
                if (!ability) {
                    console.warn(`Ability with ID ${abilityID} not found in AbilitiesConfig`);
                    return;
                }

                const card = document.createElement('div');
                card.classList.add('ability-card');
                if (selectedAbility && selectedAbility.name === ability.name) card.classList.add('selected');

                // Добавляем класс locked, если ход не нашей команды
                if (currentChar.team !== myTeam) {
                    card.classList.add('locked');
                }

                card.innerHTML = `
                    <div class="image" style="background-image: url('${ability.imageURL}');"></div>
                    <div class="info">
                        <strong>${ability.display_name}</strong><br>
                        ${ability.description || 'No description'}
                    </div>
                `;

                // Добавляем обработчик клика только если ход нашей команды
                if (currentChar.team === myTeam) {
                    card.addEventListener('click', () => {
                        setSelectedAbility(ability);
                        updateAbilityCards(myTeam, data);
                    });
                }

                abilityCards.appendChild(card);
            });

            // Отрисовываем рубашки веером чуть ниже
            const stack = document.createElement('div');
            stack.classList.add('no-abilities-stack');

            // Создаём 3 рубашки
            for (let i = 0; i < 3; i++) {
                const card = document.createElement('div');
                card.classList.add('no-abilities-card');
                card.innerHTML = `
                    <div class="image"></div>
                `;
                stack.appendChild(card);
            }

            abilityCards.appendChild(stack);
        } else {
            // Если способностей нет, отрисовываем только рубашки
            const stack = document.createElement('div');
            stack.classList.add('no-abilities-stack');

            // Создаём 3 рубашки
            for (let i = 0; i < 3; i++) {
                const card = document.createElement('div');
                card.classList.add('no-abilities-card');
                card.innerHTML = `
                    <div class="image"></div>
                `;
                stack.appendChild(card);
            }

            abilityCards.appendChild(stack);
        }
    }
    isUpdatingAbilities = false;
}

export function updateBattleLog(data) {
    if (!data || !data.teams || !Array.isArray(data.teams)) {
        console.warn('Invalid data in updateBattleLog:', data);
        return;
    }

    if (!previousState || !previousState.teams || !Array.isArray(previousState.teams)) {
        previousState = {...data, teams: data.teams.map(team => ({...team, characters: [...team.characters]}))};
        return;
    }


    for (let teamIdx = 0; teamIdx < 2; teamIdx++) {
        const prevTeam = previousState.teams[teamIdx]?.characters;
        const currTeam = data.teams[teamIdx]?.characters;

        if (!prevTeam || !currTeam) {
            console.warn(`Team ${teamIdx} is missing in previousState or data:`, {prevTeam, currTeam});
            continue;
        }

        for (let i = 0; i < Math.min(prevTeam.length, currTeam.length); i++) {
            const prev = prevTeam[i];
            const curr = currTeam[i];
            if (prev.position[0] !== curr.position[0] || prev.position[1] !== curr.position[1]) {
                addLogEntry(`${curr.name} походил из (${prev.position[0]}, ${prev.position[1]}) в (${curr.position[0]}, ${curr.position[1]})`);
            }
            if (prev.hp !== curr.hp && curr.hp > 0) {
                addLogEntry(`${curr.name} получил ${prev.hp - curr.hp} урона (Оставшееся здоровье: ${curr.hp})`);
            }
            if (prev.hp > 0 && curr.hp <= 0) {
                addLogEntry(`${curr.name} был побежден`);
            }
        }
    }
    previousState = {...data, teams: data.teams.map(team => ({...team, characters: [...team.characters]}))};
}

export function updateTurnHeader(myTeam, data) {
    const turnHeader = document.getElementById('turnText');
    if (!turnHeader) return;

    const currentChar = findCharacter(data.teams, data.currentTurn);
    if (currentChar) {
        if (currentChar.team === myTeam) {
            turnHeader.textContent = 'ВАШ ХОД';
            turnHeader.style.color = '#4dabf7'; // Синий цвет для вашего хода
            turnHeader.style.textShadow = '0 0 10px rgba(77, 171, 247, 0.7)';
        } else {
            turnHeader.textContent = 'ХОД ПРОТИВНИКА';
            turnHeader.style.color = '#ff6b6b'; // Красный цвет для хода противника
            turnHeader.style.textShadow = '0 0 10px rgba(255, 107, 107, 0.7)';
        }
    }
}

export function updateEndTurnButton(myTeam, data) {
    const endTurnBtn = document.getElementById('endTurnBtn');
    if (!endTurnBtn) return;

    const currentChar = findCharacter(data.teams, data.currentTurn);
    if (currentChar) {
        if (currentChar.team === myTeam) {
            endTurnBtn.textContent = 'Завершить ход';
            endTurnBtn.disabled = false;
            endTurnBtn.classList.remove('disabled');
        } else {
            endTurnBtn.textContent = 'Завершить ход';
            endTurnBtn.disabled = true;
            endTurnBtn.classList.add('disabled');
        }
    }
}