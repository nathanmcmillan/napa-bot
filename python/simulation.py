import sys
import signal
import time
import json
import os.path
import patterns
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


def trade(candles, signal):
    global funds
    global orders
    ticker_price = candles[-1].closing
    if signal == 'buy':
        if funds > 20.0:
            '''for existing_order in orders:
                if abs(existing_order.coin_price - ticker_price) / ticker_price < 0.05:
                    print('not buying due to existing order bought at ${}'.format(ticker_price))
                    return'''
            buy_size = funds * 0.5
            funds -= buy_size
            orders.append(SimOrder(candles[-1].closing, None, buy_size))
            print('buy | {} | coin price ${:.2f} using ${:.2f}'.format(candles[-1].time, candles[-1].closing, buy_size))
        else:
            print('not enough funds ${:.2f}'.format(funds))
    elif signal == 'sell':
        for order_to_sell in orders[:]:
            min_price = order_to_sell.profit_price()
            if ticker_price > min_price:
                funds += ticker_price * order_to_sell.size
                orders.remove(order_to_sell)
                print('sell | {} | ${:.2f} -> ${:.2f} | funds ${:.2f}'.format(candles[-1].time, order_to_sell.coin_price, ticker_price, funds))


print('funds ${:.2f}'.format(funds))
while end < historial_candle_count:
    if historical_candles[start].time < 1513504800:
        start += 1
        end += 1
        continue
    candles = historical_candles[start:end]
    candle_count = len(candles)
    macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
    directional_index = AverageDirectional(14)
    directional_index.update(candles)
    relative_strength_index = RelativeStrength(14)
    relative_strength_index.update(candles)
    today_color = patterns.color(candles[-1])
    yesterday_color = patterns.color(candles[-2])
    yesterday_maru = patterns.marubozu(candles[-2])
    hammer = patterns.hammer(candles[-1])
    star = patterns.shooting_star(candles[-1])
    maru = patterns.marubozu(candles[-1])
    trend = patterns.trend(candles, 3)
    for index in range(1, candle_count):
        current_candle = candles[index]
        macd.update(current_candle.closing)
    if trend == 'down' and maru == 'buy':
        trade(candles, 'buy')
    elif today_color == 'red' and yesterday_color == 'red':
        trade(candles, 'sell')
    start += 1
    end += 1
print('funds ${:.2f}'.format(funds))
liquidate = 0.0
for order in orders:
    liquidate += order.size * historical_candles[-1].closing
    print('price ${:.2f} size {:.4f} value ${:.2f}'.format(order.coin_price, order.size, order.usd))
print('total ${:.2f}'.format(funds + liquidate))
'''
$ 26,728.09
buy_funds = funds * 0.5
trend = patterns.trend(candles, 3)
if trend == 'up' and maru == 'buy':
elif trend == 'down' and maru == 'sell':

$ 34,940.68
buy_funds = funds
trend = patterns.trend(candles, 3)
if trend == 'up' and maru == 'buy':
elif trend == 'down' and maru == 'sell':
'''