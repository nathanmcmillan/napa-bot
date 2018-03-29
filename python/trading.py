import time
import http.client
import hmac
import hashlib
import time
import base64
import json
import gdax
import os
import printing


def process(auth, product, orders, orders_file, funds, funds_file, signal):
    if signal == 'wait':
        return
    updates = False
    ticker, status = gdax.get_ticker(product)
    if status != 200:
        print('could not get ticker', status)
        return
    if signal == 'buy':
        for existing_order in orders:
            coin_price = existing_order.coin_price()
            if percent_change(coin_price, ticker.price) < 0.05:
                print('existing order', existing_order.id, 'bought at $', coin_price, 'within range of ticker at $', ticker.price)
                return
        accounts, status = gdax.get_accounts(auth)
        product_fund = funds[product]
        available_usd = accounts['USD'].available
        if product_fund > available_usd and product_fund > 20.0:
            buy_size = product_fund / 2.0
            pending_order, status = buy(auth, product, str(buy_size))
            if status == 200:
                settled_order = wait_til_settled(auth, pending_order.id)
                cost = settled_order.executed_value + settled_order.fill_fees
                funds[product] = funds[product] - cost
                orders.append(settled_order)
                updates = True
                printing.log('bought {} cost ${}'.format(settled_order.id, cost))
            else:
                printing.log('{} failed to buy {}'.format(status, pending_order))
        else:
            printing.log('not enough funds ${} {} / ${} available'.format(product_fund, product, available_usd))
    elif signal == 'sell':
        if len(orders) == 0:
            print('nothing to sell')
            return
        for order_to_sell in orders[:]:
            min_price = order_to_sell.profit_price()
            print(product, '|', ticker.price, '>', min_price, '?')
            if ticker.price > min_price:
                pending_order, status = sell(auth, order_to_sell)
                if status == 200:
                    settled_order = wait_til_settled(auth, pending_order.id)
                    profit = settled_order.executed_value - order_to_sell.executed_value - settled_order.fill_fees
                    funds[product] = funds[product] + profit * 0.85
                    orders.remove(order_to_sell)
                    updates = True
                    printing.log('sold {} for profit ${}'.format(settled_order.id, profit))
                else:
                    printing.log('{} failed to sell {}'.format(status, pending_order))
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