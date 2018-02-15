
const private = require('../private.json');

const gdax = require('gdax');

const https = require('https');
const crypto = require('crypto');
const sqlite = require('sqlite3');
const websocket = require('ws');

const sandbox = true;
const actual = false;
const apiRest = sandbox ? 'api-public.sandbox.gdax.com' : actual ? 'api.gdax.com' : '';
const apiSocket = 'wss://ws-feed.gdax.com';
const apiTime = '/time';
const apiProducts = '/products';
const apiAccounts = '/accounts';
const apiOrders = '/orders';

function app_run()
{
    gdax.listen_ticker();
}

app_run();