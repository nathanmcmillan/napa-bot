const private = require('../private.json');
const https = require('https');
const crypto = require('crypto');
const sqlite = require('sqlite3');
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
    console.log('cryptocurrency trade bot');
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
    
    await (() => {
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
    })();
    
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
    
    // getAccounts();
    
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
    
    /* await (() => {
        return new Promise((resolve, reject) => {
            db.close((error) => {
                if (error) {
                    console.error(error.message);
                    reject(error);
                }
                console.log('closed database');
                resolve();
            });
        });
    })(); */
    
    await sleep(100);
    console.log('done');
}

async function sleep(millis) {
    return new Promise(resolve => setTimeout(resolve, millis));
}

function doSync(call) {
    await (() => {
        return new Promise((resolve, reject) => {
            call(resolve, reject);
        });
    })();
}

main();
