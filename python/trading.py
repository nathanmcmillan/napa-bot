import time
import http.client
import hmac
import hashlib
import time
import base64
import json
import gdax
import os


def process(auth, product, orders, orders_file, funds, funds_file, macd):
    if macd.signal == 'wait':
        return
    updates = False
    ticker, status = gdax.get_ticker(product)
    if macd.signal == 'buy':
        for existing_order in orders:
            coin_price = existing_order.coin_price()
            if percent_change(coin_price, ticker.price) < 0.05:
                return
        accounts, status = gdax.get_accounts(auth)
        product_fund = funds[product]
        available_usd = accounts['USD'].available
        if product_fund > available_usd and product_fund > 20.0:
            buy_size = product_fund / 2.0
            pending_order, status = buy(auth, product, str(buy_size))
            if status == 200:
                settled_order = wait_til_settled(auth, pending_order.id)
                funds[product] = funds[product] - settled_order.executed_value - settled_order.fill_fees
                orders.append(settled_order)
                updates = True
    elif macd.signal == 'sell':
        for order_to_sell in orders[:]:
            min_price = order_to_sell.profit_price()
            if ticker.price > min_price:
                pending_order, status = sell(auth, order_to_sell)
                if status == 200:
                    settled_order = wait_til_settled(auth, pending_order.id)
                    profits = settled_order.executed_value - order_to_sell.executed_value - settled_order.fill_fees
                    funds[product] = funds[product] + profits * 0.85
                    orders.remove(order_to_sell)
                    updates = True
    if updates:
        update_orders_file(orders_file, orders)
        update_funds_file(funds_file, funds)


def update_orders_file(orders_file, order_list):
    data = ''
    for current_order in order_list:
        data += current_order.id + '\n'
    orders_file.write(data)


def update_funds_file(funds_file, fund_map):
    data = ''
    for currency, amount in fund_map.items():
        data += currency + ' ' + amount + '\n'
    funds_file.write(data)


def wait_til_settled(auth, order_id):
    while True:
        time.sleep(1)
        order_update, status = gdax.get_order(auth, order_id)
        if status == 200 and order_update.settled:
            return order_update
        print('waiting for order to settle')


def buy(auth, product_id, funds):
    js_map = {'type': 'market', 'side': 'buy', 'product_id': product_id, 'funds': funds}
    js = json.dumps(js_map)
    print(js)
    return gdax.place_order(auth, js)


def sell(auth, order):
    js_map = {'type': 'market', 'side': 'sell', 'product_id': order.product_id, 'sell': str(order.filled_size)}
    js = json.dumps(js_map)
    print(js)
    return gdax.place_order(auth, js)


def percent_change(a, b):
    return abs(a - b) / b