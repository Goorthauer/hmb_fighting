import { gameState, selectedCharacter, selectedAbility, setSelectedCharacter, setSelectedAbility } from './state.js';
import { findCharacter, addLogEntry } from './utils.js';

let previousState = null;
let turnOrder = [];
let roundStarted = false;
let currentGameSessionId = null;
let isUpdatingCards = false;
let isUpdatingAbilities = false;

function saveTurnState(gameSessionId) {
    const state = { turnOrder, roundStarted, previousTurn: previousState?.CurrentTurn };
    localStorage.setItem(`turnState_${gameSessionId}`, JSON.stringify(state));
}

function loadTurnState(gameSessionId) {
    const savedState = localStorage.getItem(`turnState_${gameSessionId}`);
    if (savedState) {
        const state = JSON.parse(savedState);
        turnOrder = state.turnOrder || [];
        roundStarted = state.roundStarted || false;
        previousState = state.previousTurn ? { CurrentTurn: state.previousTurn } : null;
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
    if (!phaseContainer) return;

    const team0Alive = data.Teams?.[0]?.Characters.filter(c => c.HP > 0).length || 0;
    const team1Alive = data.Teams?.[1]?.Characters.filter(c => c.HP > 0).length || 0;
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
    const gameSessionId = data.GameSessionId || 'default_session';

    if (currentGameSessionId && currentGameSessionId !== gameSessionId) {
        resetTurnState();
    }

    let allChars = [];
    if (data.Teams && Array.isArray(data.Teams)) {
        data.Teams.forEach(team => {
            team.Characters.forEach(char => {
                if (char.HP > 0) allChars.push(char);
            });
        });
    }

    if (!roundStarted || turnOrder.length === 0) {
        if (loadTurnState(gameSessionId)) {
            console.log('Loaded turn state:', { turnOrder, previousTurn: previousState?.CurrentTurn });
        } else {
            allChars.sort((a, b) => b.Initiative - a.Initiative);
            turnOrder = allChars.map(char => char.ID);
            roundStarted = true;
            currentGameSessionId = gameSessionId;
        }
    }

    renderCards(characterCards, allChars, data);
    saveTurnState(gameSessionId);
    isUpdatingCards = false;
}

function renderCards(container, chars, data) {
    chars.forEach(char => {
        const card = document.createElement('div');
        card.classList.add('card', `team${char.Team}`);
        card.dataset.id = char.ID;
        if (selectedCharacter && selectedCharacter.ID === char.ID) card.classList.add('selected');
        if (char.ID === data.CurrentTurn) card.classList.add('current');

        card.innerHTML = `
            <div class="image" style="background-image: url('./image/char_${char.ID % 5 || 5}.png');">
                <div class="name">${char.Name}</div>
                <div class="hp-container"><div class="hp-diamond"><div class="hp">${char.HP}</div></div></div>
            </div>
            <div class="info">
                <div class="stat"><span class="label">Stamina:</span> ${char.Stamina}</div>
                <div class="stat"><span class="label">Attack:</span> ${char.AttackMin}-${char.AttackMax}</div>
                <div class="stat"><span class="label">Defense:</span> ${char.Defense}</div>
                <div class="stat"><span class="label">Init:</span> ${char.Initiative}</div>
                <div class="stat"><span class="label">Weapon:</span> ${char.Weapon || 'None'}</div>
                <div class="stat"><span class="label">Shield:</span> ${char.Shield || 'None'}</div>
            </div>
        `;
        card.addEventListener('click', () => {
            if (char.HP > 0) setSelectedCharacter(char);
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

    const currentChar = findCharacter(data.Teams, data.CurrentTurn);
    if (currentChar && currentChar.Team === myTeam && currentChar.Abilities?.length) {
        currentChar.Abilities.forEach(ability => {
            const card = document.createElement('div');
            card.classList.add('ability-card');
            if (selectedAbility && selectedAbility.Name === ability.Name) card.classList.add('selected');
            card.innerHTML = `
                <div class="image" style="background-image: url('./image/ability_${ability.Name.toLowerCase()}.png');"></div>
                <div class="info"><strong>${ability.Name}</strong><br>${ability.Description || 'No description'}</div>
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
    if (!data || !data.Teams || !Array.isArray(data.Teams)) {
        console.warn('Invalid data in updateBattleLog:', data);
        return;
    }

    if (!previousState || !previousState.Teams || !Array.isArray(previousState.Teams)) {
        console.log('Initializing previousState with current data');
        previousState = { ...data, Teams: data.Teams.map(team => ({ ...team, Characters: [...team.Characters] })) };
        return;
    }

    console.log('updateBattleLog - previousState:', previousState, 'data:', data);

    const prevChar = findCharacter(previousState.Teams, previousState.CurrentTurn);
    const currChar = findCharacter(data.Teams, data.CurrentTurn);

    if (previousState.Phase !== data.Phase) {
        addLogEntry(`Phase changed to ${data.Phase} for ${currChar ? currChar.Name : 'Unknown'}`);
    }

    for (let teamIdx = 0; teamIdx < 2; teamIdx++) {
        const prevTeam = previousState.Teams[teamIdx]?.Characters;
        const currTeam = data.Teams[teamIdx]?.Characters;

        if (!prevTeam || !currTeam) {
            console.warn(`Team ${teamIdx} is missing in previousState or data:`, { prevTeam, currTeam });
            continue;
        }

        for (let i = 0; i < Math.min(prevTeam.length, currTeam.length); i++) {
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
    previousState = { ...data, Teams: data.Teams.map(team => ({ ...team, Characters: [...team.Characters] })) };
}