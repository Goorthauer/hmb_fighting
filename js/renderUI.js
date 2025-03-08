import { gameState, selectedCharacter, selectedAbility, setSelectedCharacter, setSelectedAbility } from './state.js';
import { findCharacter, addLogEntry } from './utils.js';

let previousState = null;
let turnOrder = [];
let roundStarted = false;
let currentGameSessionId = null;
let isUpdatingCards = false;
let isUpdatingAbilities = false;

function saveTurnState(gameSessionId) {
    const state = { turnOrder, roundStarted, previousTurn: previousState?.currentTurn };
    localStorage.setItem(`turnState_${gameSessionId}`, JSON.stringify(state));
}

function loadTurnState(gameSessionId) {
    const savedState = localStorage.getItem(`turnState_${gameSessionId}`);
    if (savedState) {
        const state = JSON.parse(savedState);
        turnOrder = state.turnOrder || [];
        roundStarted = state.roundStarted || false;
        previousState = state.previousTurn ? { currentTurn: state.previousTurn } : null;
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
            console.log('Loaded turn state:', { turnOrder, previousTurn: previousState?.currentTurn });
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
    chars.forEach(char => {
        const card = document.createElement('div');
        card.classList.add('card', `team${char.team}`);
        card.dataset.id = char.id;
        if (selectedCharacter && selectedCharacter.id === char.id) card.classList.add('selected');
        if (char.id === data.currentTurn) card.classList.add('current');
        if (char.hp <= 0) card.classList.add('dead');

        card.innerHTML = `
            <div class="image" style="background-image: url('${char.imageURL || 'default-image.png'}');"></div>
            <div class="name">${char.name}</div>
            <div class="info">
                <div class="stat"><i class="fas fa-heartbeat"></i> <span class="label">Sta:</span> ${char.stamina}</div>
                <div class="stat"><i class="fas fa-skull"></i> <span class="label">Atk:</span> ${char.attackMin}-${char.attackMax}</div>
                <div class="stat"><i class="fas fa-shield-alt"></i> <span class="label">Def:</span> ${char.defense}</div>
                <div class="stat"><i class="fas fa-tachometer-alt"></i> <span class="label">Ini:</span> ${char.initiative}</div>
                <div class="stat"><i class="fas fa-gavel"></i> <span class="label"></span> ${data.weaponsConfig[char.weapon]?.name || 'None'}</div>
                <div class="stat"><i class="fas fa-shield"></i> <span class="label"></span> ${data.shieldsConfig[char.shield]?.name || 'None'}</div>
                <div class="hp-container"><div class="hp-diamond"><div class="hp">${char.hp}</div></div></div>
            </div>
        `;
        card.addEventListener('click', () => {
            if (char.hp > 0) setSelectedCharacter(char); // Выбор только живых персонажей
            updateCharacterCards(data);
        });
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
    if (currentChar && currentChar.team === myTeam && currentChar.abilities?.length) {
        currentChar.abilities.forEach(abilityID => {
            console.log(data)
            const ability = data.abilitiesConfig[abilityID]; // Получаем способность из конфигурации
            if (!ability) {
                console.warn(`Ability with ID ${abilityID} not found in AbilitiesConfig`);
                return;
            }

            const card = document.createElement('div');
            card.classList.add('ability-card');
            if (selectedAbility && selectedAbility.name === ability.name) card.classList.add('selected');
            card.innerHTML = `
                <div class="image" style="background-image: url('${ability.imageURL}');"></div>
                <div class="info"><strong>${ability.name}</strong><br>${ability.description || 'No description'}</div>
            `;
            card.addEventListener('click', () => {
                setSelectedAbility(ability);
                updateAbilityCards(myTeam, data);
            });
            abilityCards.appendChild(card);
        });
    }
    isUpdatingAbilities = false;
}
export function updateBattleLog(data) {
    if (!data || !data.teams || !Array.isArray(data.teams)) {
        console.warn('Invalid data in updateBattleLog:', data);
        return;
    }

    if (!previousState || !previousState.teams || !Array.isArray(previousState.teams)) {
        previousState = { ...data, teams: data.teams.map(team => ({ ...team, characters: [...team.characters] })) };
        return;
    }

    const prevChar = findCharacter(previousState.teams, previousState.currentTurn);
    const currChar = findCharacter(data.teams, data.currentTurn);

    if (previousState.phase !== data.phase) {
        addLogEntry(`Phase changed to ${data.phase} for ${currChar ? currChar.name : 'Unknown'}`);
    }

    for (let teamIdx = 0; teamIdx < 2; teamIdx++) {
        const prevTeam = previousState.teams[teamIdx]?.characters;
        const currTeam = data.teams[teamIdx]?.characters;

        if (!prevTeam || !currTeam) {
            console.warn(`Team ${teamIdx} is missing in previousState or data:`, { prevTeam, currTeam });
            continue;
        }

        for (let i = 0; i < Math.min(prevTeam.length, currTeam.length); i++) {
            const prev = prevTeam[i];
            const curr = currTeam[i];
            if (prev.position[0] !== curr.position[0] || prev.position[1] !== curr.position[1]) {
                addLogEntry(`${curr.name} moved from (${prev.position[0]}, ${prev.position[1]}) to (${curr.position[0]}, ${curr.position[1]})`);
            }
            if (prev.hp !== curr.hp && curr.hp > 0) {
                addLogEntry(`${curr.name} took ${prev.hp - curr.hp} damage (HP: ${curr.hp})`);
            }
            if (prev.hp > 0 && curr.hp <= 0) {
                addLogEntry(`${curr.name} was defeated`);
            }
        }
    }
    previousState = { ...data, teams: data.teams.map(team => ({ ...team, characters: [...team.characters] })) };
}