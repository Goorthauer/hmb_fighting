import { gameState } from './state.js';

export function findCharacter(teams, id) {
    if (!teams || !Array.isArray(teams)) return null;
    for (let team of teams) {
        for (let char of team.characters) {
            if (char.id === id) return char;
        }
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
}