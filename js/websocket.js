export let ws = null;
export let clientID = localStorage.getItem('clientID') || null;

export function connectWebSocket(room, isSpectator, onMessage) {
    if (!room || room === 'null') {
        console.error('Room parameter is missing or invalid, cannot connect to WebSocket');
        return;
    }

    const url = `ws://localhost:8080/ws?room=${encodeURIComponent(room)}&spectator=${isSpectator}${clientID ? `&clientID=${encodeURIComponent(clientID)}` : ''}`;
    console.log('Connecting to WebSocket with URL:', url);
    ws = new WebSocket(url);

    ws.onopen = () => {
        console.log(`Connected to ${room} as ${isSpectator ? 'spectator' : 'player'} with clientID: ${clientID}`);
    };

    ws.onmessage = (event) => {
        onMessage(event);
    };

    ws.onerror = (err) => {
        console.error('WebSocket error:', err);
    };

    ws.onclose = () => {
        console.log('Disconnected from server, attempting to reconnect...');
        ws = null;
        setTimeout(() => connectWebSocket(room, isSpectator, onMessage), 1000);
    };
}

export function sendMessage(message) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        console.log('Sending WebSocket message:', message);
        ws.send(message);
    } else {
        console.error('WebSocket is not connected');
    }
}

export function setClientID(id) {
    clientID = id;
    localStorage.setItem('clientID', id);
    console.log('Updated clientID in websocket.js:', clientID);
}