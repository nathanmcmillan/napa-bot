
const gdax = require('gdax');
const apiRest = 'https://api.gdax.com';
const apiFeed = 'wss://ws-feed.gdax.com';
const sandboxRest = 'https://api-public.sandbox.gdax.com';
const sandboxFeed = 'wss://ws-feed-public.sandbox.gdax.com';
const ethToUsd = 'ETH-USD';
const btcToUsd = 'BTC-USD';
const passphrase = 'phrase';
// const apiKey = 'apiKey';
const base64secret = 'secret';
const publicClient = new gdax.PublicClient();
// const authClient = new gdax.AuthenticatedClient(apiKey, base64secret, passphrase, apiUri);
const buyOrders = [];
const sellOrders = [];

function getProducts() {
  const description = 'get products';
  const callback = (error, response, data) => {
    console.log(description);
    if (error) {
      console.log(error);
      return;
    }
    console.log(data);
  };
  publicClient.getProducts(callback);
}

function getCurrencies() {
  const description = 'get currencies';
  const callback = (error, response, data) => {
    console.log(description);
    if (error) {
      console.log(error);
      return;
    }
    console.log(data);
  };
  publicClient.getCurrencies(callback);
}

function getAccounts() {
  const description = 'get accounts';
  const callback = (error, response, data) => {
    console.log(description);
    if (error) {
      console.log(error);
      return;
    }
    console.log(data);
  };
  authClient.getAccounts(callback);
}

function getOrderbook(pair, api, feed) {
  const book = new gdax.OrderbookSync([pair], api, feed);
  console.log(book.books[pair].state());
}

function listenExchange(pairs, uri) {
  const socket = new gdax.WebsocketClient(pairs, uri);
  const onMessage = (data) => {
    if (data.type !== 'done' || data.reason !== 'filled') {
      return;
    }
    console.log(data);
  };
  const onError = (error) => {
    console.log('socket error');
    console.log(error);
  };
  const onClose = () => {
    console.log('socket closed');
  };
  socket.on('message', onMessage);
  socket.on('error', onError);
  socket.on('close', onClose);
}

function buy(price, size, product) {
  const parameters = {
    'price': price,
    'size': size,
    'product_id': product
  };
  const callback = () => {
    console.log('todo');
  };
  const buyId = authClient.buy(parameters, callback);
  buyOrders.push(buyId);
}

function sell(price, size, product) {
  const parameters = {
    'price': price,
    'size': size,
    'product_id': product
  };
  const callback = () => {
    console.log('todo');
  };
  const sellId = authClient.sell(parameters, callback);
  sellOrders.push(sellId);
}

function getOrders() {
  const callback = () => {
    console.log('todo');
  };
  authClient.getOrders(callback);
}

function checkOrder(orderId) {
  const callback = () => {
    console.log('todo');
  };
  authClient.getOrder(orderId, callback);
}

// getProducts();
// getCurrencies();
// getAccounts();
// listenExchange([ethToUsd], api);
// getOrderbook(ethToUsd, sandboxRest, sandboxFeed);
// getOrderbook(ethToUsd, apiRest, apiFeed);
