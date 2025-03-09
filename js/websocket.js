export let ws = null;

export function connectWebSocket(room, isSpectator, onMessage) {
    const accessToken = localStorage.getItem('accessToken');
    const clientID = localStorage.getItem('clientID');
    if (!accessToken || !clientID) {
        console.error('No access token or clientID found, redirecting to login');
        window.location.href = 'index.html';
        return;
    }

    const url = `ws://localhost:8080/ws?room=${encodeURIComponent(room)}&accessToken=${encodeURIComponent(accessToken)}`;
    console.log('Connecting to WebSocket with URL:', url);
    ws = new WebSocket(url);

    ws.onopen = () => {
        console.log(`Connected to ${room} as ${isSpectator ? 'spectator' : 'player'}`);
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('Received WebSocket message:', data);
        if (data.error === 'Invalid token') {
            console.log('Invalid token detected, attempting to refresh...');
            refreshToken(() => connectWebSocket(room, isSpectator, onMessage));
        } else {
            onMessage(event);
        }
    };

    ws.onerror = (err) => {
        console.error('WebSocket error:', err);
    };

    ws.onclose = (event) => {
        console.log('Disconnected from server:', { code: event.code, reason: event.reason });
        localStorage.clear(); // Полная очистка
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

function refreshToken(callback) {
    const refreshTokenVal = localStorage.getItem('refreshToken');
    if (!refreshTokenVal) {
        console.error('No refresh token found, redirecting to login');
        window.location.href = 'index.html';
        return;
    }

    fetch('http://localhost:8080/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken: refreshTokenVal })
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Token refresh failed with status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            localStorage.setItem('accessToken', data.accessToken);
            localStorage.setItem('refreshToken', data.refreshToken);
            localStorage.setItem('clientID', data.clientID);
            console.log('Token refreshed successfully');
            callback();
        })
        .catch(err => {
            console.error('Token refresh failed:', err);
            localStorage.clear();
            window.location.href = 'index.html';
        });
}