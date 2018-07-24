class Socket {
    constructor() {
        this.src = new WebSocket('ws://localhost:80/websocket');
    }
}