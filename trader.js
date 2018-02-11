const private = require('../private.json');
const https = require('https');
const crypto = require('crypto');
const apiRest = 'api.gdax.com';
const apiTime = '/time';
const apiProducts = '/products';
const apiAccounts = '/accounts';

function getProducts() {
    const options = {
        host: apiRest,
        port: 443,
        path: apiProducts,
        method: 'GET',
        headers: {
            'accept': 'application/json',
            'user-agent': 'bot'
        }
    };
    const callback = (response) => {
        response.on('error', (error) => {
            console.log('error ' + error);
        });
        let body = '';
        response.on('data', (data) => {
            body += data;
        });
        response.on('end', () => {
            // body = JSON.parse(body);
            console.log('end ' + body);
        });
        console.log('code ' + response.statusCode);
    };
    const get = https.request(options, callback);
    get.end();
}

function getAccounts() {
    const body = JSON.stringify({
        price: '1.0'
    });
    const restMethod = 'GET';
    const requestPath = apiAccounts;
    const apiKey = private.apiKey;
    const secret = private.secret;
    const timeStamp = Date.now() / 1000;
    const passPhrase = private.passPhrase;
    const prehash = timeStamp + restMethod + requestPath + body;
    const secretKey = Buffer(secret, 'base64');
    const hmac = crypto.createHmac('sha256', secretKey);
    const signature = hmac.update(prehash).digest('base64');
    const options = {
        host: apiRest,
        port: 443,
        path: requestPath,
        method: restMethod,
        headers: {
            'accept': 'application/json',
            'user-agent': 'bot',
            'CB-ACCESS-KEY': apiKey,
            'CB-ACCESS-SIGN': signature,
            'CB-ACCESS-TIMESTAMP': timeStamp,
            'CB-ACCESS-PASSPHRASE': passPhrase
        }
    };
    const callback = (response) => {
        response.on('error', (error) => {
            console.log('error ' + error);
        });
        let body = '';
        response.on('data', (data) => {
            body += data;
        });
        response.on('end', () => {
            // body = JSON.parse(body);
            console.log('end ' + body);
        });
        console.log('code ' + response.statusCode);
    };
    const get = https.request(options, callback);
    get.end();
}

async function main() {
    console.log('trading');
    getAccounts();
    await sleep(100);
    console.log('done');
}

async function sleep(millis) {
    return new Promise(resolve => setTimeout(resolve, millis));
}

main();
