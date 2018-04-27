from rest import request, private_request
import websocket

SITE = 'api.gdax.com'
FEED_SITE = 'wss://ws-feed.gdax.com'


class Ticker:
    def __init__(self, ticker_data):
        self.trade_id = int(ticker_data.get('trade_id') or 0)
        self.price = float(ticker_data.get('price') or 0)
        self.size = float(ticker_data.get('size') or 0)
        self.bid = float(ticker_data.get('bid') or 0)
        self.ask = float(ticker_data.get('ask') or 0)
        self.volume = float(ticker_data.get('volume') or 0)
        self.time = ticker_data.get('time')


class NewOrder:
    def __init__(self, order_data):
        self.id = order_data.get('id')
        self.price = float(order_data.get('price') or 0)
        self.size = float(order_data.get('size') or 0)
        self.product_id = order_data.get('product_id')
        self.side = order_data.get('side')
        self.stp = order_data.get('stp')
        self.type = order_data.get('type')
        self.time_in_force = order_data.get('time_in_force')
        self.post_only = order_data.get('post_only')
        self.created_at = order_data.get('created_at')
        self.fill_fees = float(order_data.get('fiil_fees') or 0)
        self.filled_size = float(order_data.get('filled_size') or 0)
        self.executed_value = float(order_data.get('executed_value') or 0)
        self.status = order_data.get('status')
        self.settled = order_data.get('settled')


class Order:
    def __init__(self, order_data):
        self.id = order_data.get('id')
        self.size = float(order_data.get('size') or 0)
        self.product_id = order_data.get('product_id')
        self.side = order_data.get('side')
        self.stp = order_data.get('stp')
        self.funds = float(order_data.get('funds') or 0)
        self.specified_funds = float(order_data.get('specified_funds') or 0)
        self.type = order_data.get('type')
        self.post_only = order_data.get('post_only')
        self.created_at = order_data.get('created_at')
        self.done_at = order_data.get('done_at')
        self.done_reason = order_data.get('done_reason')
        self.fill_fees = float(order_data.get('fiil_fees') or 0)
        self.filled_size = float(order_data.get('filled_size') or 0)
        self.executed_value = float(order_data.get('executed_value') or 0)
        self.status = order_data.get('status')
        self.settled = order_data.get('settled')

    def coin_price(self):
        return self.executed_value / self.filled_size

    def profit_price(self):
        margin = (self.specified_funds / self.funds - 1.0) * 2.0 + 1.0
        return self.coin_price() * margin


class Account:
    def __init__(self, account_data):
        self.id = account_data['id']
        self.currency = account_data['currency']
        self.balance = float(account_data['balance'])
        self.available = float(account_data['available'])
        self.hold = float(account_data['hold'])
        self.profile_id = account_data['profile_id']


class Candle:
    def __init__(self, candle_data):
        self.time = int(candle_data[0])
        self.low = float(candle_data[1])
        self.high = float(candle_data[2])
        self.open = float(candle_data[3])
        self.closing = float(candle_data[4])
        self.volume = float(candle_data[5])

    def typical_price(self):
        return (self.high + self.low + self.closing) / 3


def place_order(auth, post):
    read, status = private_request(auth, 'POST', SITE, '/orders', post)
    if status != 200 or not isinstance(read, dict):
        return read, status
    return NewOrder(read), status


def get_order(auth, id):
    read, status = private_request(auth, 'GET', SITE, '/orders/' + id, '')
    if status != 200 or not isinstance(read, dict):
        return read, status
    return Order(read), status


def get_accounts(auth):
    read, status = private_request(auth, 'GET', SITE, '/accounts', '')
    if status != 200 or not isinstance(read, list):
        return read, status
    accounts = {}
    for read_account in read:
        new_account = Account(read_account)
        accounts[new_account.currency] = new_account
    return accounts, status


def get_candles(product, start, end, granularity):
    read, status = request('GET', SITE, '/products/' + product + '/candles?start=' + start + '&end=' + end + '&granularity=' + granularity, '')
    if status != 200 or not isinstance(read, list):
        return read, status
    candles = []
    for read_candle in read:
        candles.append(Candle(read_candle))
    candles.sort(key=lambda c: c.time, reverse=False)
    return candles, status


def get_ticker(product):
    read, status = request('GET', SITE, '/products/' + product + '/ticker', '')
    if status != 200 or not isinstance(read, dict):
        return read, status
    return Ticker(read), status


def get_feed(products, channels):
    feed = websocket.create_connection(FEED_SITE)
    params = {'type': 'subscribe', 'product_ids': products, 'channels': channels}
    feed.send(json.dumps(params))