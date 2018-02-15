const private = require('../private.json');
const https = require('https');
const crypto = require('crypto');
const sqlite = require('sqlite3');
const sandbox = true;
const actual = false;
const apiRest = sandbox ? 'api-public.sandbox.gdax.com' : actual ? 'api.gdax.com' : '';
const apiSocket = 'wss://ws-feed.gdax.com';
const apiTime = '/time';
const apiProducts = '/products';
const apiAccounts = '/accounts';
const apiOrders = '/orders';

function getProducts(resolve, reject) {
    const options = {
        host: apiRest,
        port: 443,
        path: apiProducts,
        method: 'GET',
        headers: {
            'User-Agent': 'bot',
            'Accept': 'application/json'
        }
    };
    const callback = (response) => {
        response.on('error', (error) => {
            console.log('error ' + error);
            reject(error);
        });
        let body = '';
        response.on('data', (data) => {
            body += data;
        });
        response.on('end', () => {
            body = JSON.parse(body);
            console.log(body);
            resolve();
        });
        console.log('code ' + response.statusCode);
    };
    const get = https.request(options, callback);
    get.end();
}

function getAccounts(resolve, reject) {
    const body = '';
    const restMethod = 'GET';
    const requestPath = apiAccounts;
    const apiKey = private.key;
    const secret = private.secret;
    const phrase = private.phrase;
    const time = Date.now() / 1000;
    const prehash = time + restMethod + requestPath + body;
    const secretKey = Buffer(secret, 'base64');
    const hmac = crypto.createHmac('sha256', secretKey);
    const signature = hmac.update(prehash).digest('base64');
    const options = {
        host: apiRest,
        port: 443,
        path: requestPath,
        method: restMethod,
        headers: {
            'User-Agent': 'bot',
            'Accept': 'application/json',
            'CB-ACCESS-KEY': apiKey,
            'CB-ACCESS-SIGN': signature,
            'CB-ACCESS-TIMESTAMP': time,
            'CB-ACCESS-PASSPHRASE': phrase
        }
    };
    const callback = (response) => {
        response.on('error', (error) => {
            console.log('error ' + error);
            reject(error);
        });
        let body = '';
        response.on('data', (data) => {
            body += data;
        });
        response.on('end', () => {
            body = JSON.parse(body);
            console.log(body);
            resolve();
        });
        console.log('code ' + response.statusCode);
    };
    const get = https.request(options, callback);
    get.end();
}

function getOrders(resolve, reject) {
    const body = '';
    const restMethod = 'GET';
    const requestPath = apiOrders;
    const apiKey = private.key;
    const secret = private.secret;
    const phrase = private.phrase;
    const time = Date.now() / 1000;
    const prehash = time + restMethod + requestPath + body;
    const secretKey = Buffer(secret, 'base64');
    const hmac = crypto.createHmac('sha256', secretKey);
    const signature = hmac.update(prehash).digest('base64');
    const options = {
        host: apiRest,
        port: 443,
        path: requestPath,
        method: restMethod,
        headers: {
            'User-Agent': 'bot',
            'Accept': 'application/json',
            'CB-ACCESS-KEY': apiKey,
            'CB-ACCESS-SIGN': signature,
            'CB-ACCESS-TIMESTAMP': time,
            'CB-ACCESS-PASSPHRASE': phrase
        }
    };
    const callback = (response) => {
        response.on('error', (error) => {
            console.log('error ' + error);
            reject(error);
        });
        let body = '';
        response.on('data', (data) => {
            body += data;
        });
        response.on('end', () => {
            body = JSON.parse(body);
            console.log(body);
            resolve();
        });
        console.log('code ' + response.statusCode);
    };
    const get = https.request(options, callback);
    get.end();
}

function placeLimitOrder(resolve, reject) {
    const body = JSON.stringify({
        "type": "limit", //
        "side": "buy", // [buy] or [sell]
        "product_id": "BTC-USD", // pair
        "stp": "co", // cancel other
        "price": "0.100", // price per coin
        "size": "0.01", // number of coins
        "time_in_force": "GTT", // good till cancelled
        "cancel_after": "hour", // [min], [hour], or [day]
        "post_only": true // prevent any fee
    });
    const restMethod = 'POST';
    const requestPath = apiOrders;
    const apiKey = private.key;
    const secret = private.secret;
    const phrase = private.phrase;
    const time = Date.now() / 1000;
    const prehash = time + restMethod + requestPath + body;
    const secretKey = Buffer(secret, 'base64');
    const hmac = crypto.createHmac('sha256', secretKey);
    const signature = hmac.update(prehash).digest('base64');
    const options = {
        host: apiRest,
        port: 443,
        path: requestPath,
        method: restMethod,
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(body),
            'User-Agent': 'bot',
            'Accept': 'application/json',
            'CB-ACCESS-KEY': apiKey,
            'CB-ACCESS-SIGN': signature,
            'CB-ACCESS-TIMESTAMP': time,
            'CB-ACCESS-PASSPHRASE': phrase
        }
    };
    const callback = (response) => {
        response.on('error', (error) => {
            console.log('error ' + error);
            reject(error);
        });
        let body = '';
        response.on('data', (data) => {
            body += data;
        });
        response.on('end', () => {
            body = JSON.parse(body);
            console.log(body);
            resolve();
        });
        console.log('code ' + response.statusCode);
    };
    const post = https.request(options, callback);
    post.write(body);
    post.end();
}

async function main() {
    console.log('napa bot');
    if (apiRest == '') {
        return;
    }
    console.log(apiRest);
    console.log('getting database');
    
    let db;
    
    await (() => {
        return new Promise((resolve, reject) => {
            db = new sqlite.Database('trade.db', (error) => {
                if (error) {
                    console.error(error.message);
                    reject(error);
                }
                console.log('connected to database');
                resolve();
            });
        });
    })();
    
    /* await (() => {
        return new Promise((resolve, reject) => {
            let query = `insert into trades(price) select '0.01'`;
            db.run(query, (error) => {
                if (error) {
                    console.error(error.message);
                    reject(error);
                }
                console.log('sql insert');
                resolve();
            });
        });
    })(); */
    
    await (() => {
        return new Promise((resolve, reject) => {
            let query = `select * from trades`;
            db.all(query, [], (error, rows) => {
                if (error) {
                    console.error(error.message);
                    reject(error);
                }
                console.log('sql query');
                rows.forEach((row) => {
                    console.log(row.id);
                });
                resolve();
            });
        });
    })();
    
    // getProducts();
    //doSync(getAccounts);
    doSync(placeLimitOrder);
    //doSync(getOrders);
    
    doSync((resolve, reject) => {
        db.close((error) => {
            if (error) {
                console.error(error.message);
                reject(error);
            }
            console.log('closed database');
            resolve();
        });
    });
    
    await sleep(100);
    console.log('done');
}

async function sleep(millis) {
    return new Promise(resolve => setTimeout(resolve, millis));
}

async function doSync(call) {
    await givePromise(call);
}

function givePromise(call) {
    let func = (resolve, reject) => {
        call(resolve, reject);
    };
    return new Promise(func);
}

// main();

const websocket = require('ws');

const matches = [];
const tickerLog = [];
async function otherMain() {
    const socket = new websocket(apiSocket);

    socket.on('message', (data) => {
        const message = JSON.parse(data);
        if (message.type === 'error') {
            console.error(message);
            return;
        }
        console.log(message);
        if (message.best_bid && message.best_ask) {
            tickerLog.push({
                bid: message.best_bid,
                ask: message.best_ask
            });
        }
    });
    socket.on('open', () => {
        console.log('opened');
        const subscribe = JSON.stringify({
            "type": "subscribe",
            "channels": [
                /* {
                    "name": "heartbeat",
                    "product_ids": [
                        "ETH-USD"
                    ]
                } */
                /* {
                    "name": "matches",
                    "product_ids": [
                        "ETH-USD"
                    ]
                } */
                {
                    "name": "ticker",
                    "product_ids": [
                        "ETH-USD"
                    ]
                }
                /* {
                    "name": "level2",
                    "product_ids": [
                        "ETH-USD"
                    ]
                } */
            ]
        });
        socket.send(subscribe);
    });
    socket.on('close', () => {
        console.log('closed');
    });
    socket.on('error', (error) => {
        console.error(error);
    });

    await sleep(1000 * 5);
    socket.close();
    console.log(tickerLog);
};
otherMain();