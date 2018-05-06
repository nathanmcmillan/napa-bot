const fs = require('fs')
const https = require('https')
const WebSocket = require('ws')
/*
const stream = '/stream?streams=btcusdt@ticker'
const connection = new WebSocket('wss://stream.binance.com:9443' + stream)

connection.onopen = function () {
    console.log('open')
    // connection.send('{"type": "subscribe", "product_ids": ["BTC-USD"], "channels": ["ticker"]}')
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
*/

const url = 'api.binance.com'
const path = '/api/v3/ticker/price'
const headers = {
    'Content-Type': 'application/json'
}

let options = {
    method: 'GET',
    host: url,
    path: path,
    port: 443,
    headers: headers
}

let req = https.request(options, function(res) {
    console.log('status ' + res.statusCode)
    res.setEncoding('utf8')
    res.on('data', function(data) {
        console.log('data ' + data)
    })
})

req.on('error', function(error) {
    console.error(error);
})

req.end()

console.log('sent https request')