import { gameState, previousState } from './state.js';

export function findCharacter(id) {
    if (!gameState) return null;
    for (let team of [gameState.Teams[0].characters, gameState.Teams[1].characters]) {
        for (let char of team) {
            if (char.id === id) return char;
        }
    }
    return null;
}

export function addLogEntry(message) {
    const logEntries = document.getElementById('logEntries');
    const entry = document.createElement('div');
    entry.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
    logEntries.appendChild(entry);
    logEntries.scrollTop = logEntries.scrollHeight;
}

export function animateCell(x, y, animationClass) {
    const tempDiv = document.createElement('div');
    tempDiv.classList.add(animationClass);
    tempDiv.style.position = 'absolute';
    tempDiv.style.left = `${x * 50 + document.getElementById('gameCanvas').offsetLeft}px`;
    tempDiv.style.top = `${y * 50 + document.getElementById('gameCanvas').offsetTop}px`;
    tempDiv.style.width = `50px`;
    tempDiv.style.height = `50px`;
    document.getElementById('gameContainer').appendChild(tempDiv);
    setTimeout(() => tempDiv.remove(), 1000);
}

export function updateBattleLog() {
    if (!previousState || !gameState) return;

    const prevChar = findCharacter(previousState.currentTurn);
    const currChar = findCharacter(gameState.currentTurn);

    if (previousState.phase !== gameState.phase) {
        addLogEntry(`Phase changed to ${gameState.phase} for ${currChar.name}`);
    }

    for (let teamIdx = 0; teamIdx < 2; teamIdx++) {
        const prevTeam = previousState.teams[teamIdx].characters;
        const currTeam = gameState.Teams[teamIdx].characters;
        for (let i = 0; i < prevTeam.length; i++) {
            const prev = prevTeam[i];
            const curr = currTeam[i];
            if (prev.position[0] !== curr.position[0] || prev.position[1] !== curr.position[1]) {
                addLogEntry(`${curr.name} moved from (${prev.position[0]}, ${prev.position[1]}) to (${curr.position[0]}, ${curr.position[1]})`);
                animateCell(curr.position[0], curr.position[1], 'blink-move');
            }
            if (prev.hp !== curr.hp && curr.hp > 0) {
                addLogEntry(`${curr.name} took ${prev.hp - curr.hp} damage (HP: ${curr.hp})`);
                animateCell(curr.position[0], curr.position[1], 'blink-attack');
            }
            if (prev.hp > 0 && curr.hp <= 0) {
                addLogEntry(`${curr.name} was defeated`);
                animateCell(curr.position[0], curr.position[1], 'blink-attack');
            }
            if (prev.is_defending !== curr.is_defending && curr.is_defending) {
                addLogEntry(`${curr.name} is now defending`);
            }
        }
    }

    const team0Alive = gameState.Teams[0].characters.filter(c => c.hp > 0).length;
    const team1Alive = gameState.Teams[1].characters.filter(c => c.hp > 0).length;
    if (team0Alive === 0 || team1Alive === 0) {
        const winner = team0Alive > 0 ? 'Team 0' : 'Team 1';
        const victoryMessage = document.createElement('div');
        victoryMessage.id = 'victoryMessage';
        victoryMessage.textContent = `${winner} won!`;
        document.body.appendChild(victoryMessage);

        const fireworksContainer = document.createElement('div');
        fireworksContainer.id = 'fireworks';
        document.body.appendChild(fireworksContainer);

        for (let i = 0; i < 50; i++) {
            const firework = document.createElement('div');
            firework.classList.add('firework');
            const angle = Math.random() * 2 * Math.PI;
            const distance = Math.random() * 300 + 100;
            const x = Math.cos(angle) * distance;
            const y = Math.sin(angle) * distance;
            firework.style.backgroundColor = `hsl(${Math.random() * 360}, 100%, 50%)`;
            firework.style.left = '50%';
            firework.style.top = '50%';
            firework.style.setProperty('--x', `${x}px`);
            firework.style.setProperty('--y', `${y}px`);
            fireworksContainer.appendChild(firework);
        }
    }
}