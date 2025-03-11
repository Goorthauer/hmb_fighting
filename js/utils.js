import { gameState } from './state.js';

export function findCharacter(teams, id) {
    if (!teams) {
        console.warn('Teams data is missing or undefined');
        return null;
    }

    // Если teams — массив (стандартный случай для большинства фаз)
    if (Array.isArray(teams)) {
        for (let team of teams) {
            if (team && Array.isArray(team.characters)) {
                for (let char of team.characters) {
                    if (char.id === id) return char;
                }
            } else {
                console.warn('Team is missing characters or characters is not an array:', team);
            }
        }
    }
    // Если teams — объект (возможный случай для pick_team)
    else if (typeof teams === 'object') {
        for (let teamId in teams) {
            const team = teams[teamId];
            if (team && Array.isArray(team.characters)) {
                for (let char of team.characters) {
                    if (char.id === id) return char;
                }
            } else {
                console.warn('Team is missing characters or characters is not an array:', team);
            }
        }
    } else {
        console.warn('Unexpected teams format:', teams);
    }

    return null;
}

export function addLogEntry(message) {
    const logEntries = document.getElementById('logEntries');
    if (!logEntries) return;
    const entry = document.createElement('div');
    entry.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
    logEntries.appendChild(entry);
    logEntries.scrollTop = logEntries.scrollHeight;
    console.log(entry)
    console.log(`[${new Date().toLocaleTimeString()}] ${message}`)
}