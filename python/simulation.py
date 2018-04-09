import sys
import signal
import time
import json
import os.path
from gdax import Candle
from trends import ConvergeDiverge, AverageDirectional
from momentum import RelativeStrength

fees = 0.03


class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price

    def profit_price(self):
        return self.coin_price * (1.0 + fees)


def read_map(path):
    map = {}
    with open(path, 'r') as open_file:
        for line in open_file:
            (key, value) = line.split()
            map[key] = value
    return map


print('----------------------------------------')
print('|           napa simulation            |')
print('----------------------------------------')

file_in = '../candles-btc-usd.txt'
historical_candles = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        historical_candles.append(candle)

ema_short = 12
ema_long = 26
orders = []
historial_candle_count = len(historical_candles)
start = 0
end = ema_long
funds = 1000.0
print('starting funds ${:.2f}'.format(funds))
while end < historial_candle_count:
    candles = historical_candles[start:end]
    candle_count = len(candles)
    macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
    directional_index = AverageDirectional(ema_short)
    directional_index.update(candles)
    relative_strength_index = RelativeStrength(ema_short)
    relative_strength_index.update(candles)
    for index in range(1, candle_count):
        current_candle = candles[index]
        macd.update(current_candle.closing)
    if macd.signal == 'buy':
        # if relative_strength_index.current < 0.3 and directional_index.current > 0.4:
        if funds > 20.0:
            buy_size = funds / 2.0
            funds -= buy_size
            orders.append(SimOrder(candles[-1].closing, None, buy_size))
            print('buying | coin price ${:.2f} using ${:.2f}'.format(candles[-1].closing, buy_size))
        else:
            print('not enough funds ${:.2f}'.format(funds))
    elif macd.signal == 'sell':
        # elif relative_strength_index.current > 0.7 and directional_index.current > 0.4:
        ticker_price = candles[-1].closing
        for order_to_sell in orders[:]:
            min_price = order_to_sell.profit_price()
            if ticker_price > min_price:
                settled_order = SimOrder(ticker_price, order_to_sell.size, None)
                usd = ticker_price * order_to_sell.size
                profit = usd - order_to_sell.usd
                funds += usd + profit * 0.85
                orders.remove(order_to_sell)
                print('selling | ${:.2f} -> ${:.2f} | profit ${:.2f} | funds ${:.2f}'.format(order_to_sell.coin_price, ticker_price, profit, funds))
    start += 1
    end += 1
print('ending funds ${:.2f}'.format(funds))
for order in orders:
    print('coin price ${:.2f} size {:.4f} ${:.2f}'.format(order.coin_price, order.size, order.usd))