
funciton listen_ticker()
{
    const url = Gdax.SocketUrl();
    const socket = new websocket(url);
    socket.on('open', () => {
        const subscribe = JSON.stringify({
            "type": "subscribe",
            "channels": [
                {
                    "name": "level2",
                    "product_ids": [
                        "ETH-USD"
                    ]
                }
            ]
        });
        socket.send(subscribe);
    });
    socket.on('message', (message) => {
        message = JSON.parse(message);
        console.log(message);
    });
    socket.on('close', () => {
        console.log('closed');
    });
    socket.on('error', (error) => {
        console.error(error);
    });
}

module.exports = listen_ticker;