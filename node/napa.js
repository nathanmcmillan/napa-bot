const WebSocket = require('ws');

const connection = new WebSocket('wss://ws-feed.gdax.com')

connection.onopen = function () {
    console.log('open')
    connection.send('{"type": "subscribe", "product_ids": ["BTC-USD"], "channels": ["ticker"]}')
}

connection.onclose = function () {
    console.log('close')
}

connection.onerror = function(error) {
    console.error('error ' + error)
}

connection.onmessage = function(event) {
    console.error(event.data)
}