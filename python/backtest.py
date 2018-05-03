import strategy
import simulation
from strategy import Strategy
from gdax import Candle
from operator import itemgetter

print('----------------------------------------')
print('|              napa test               |')
print('----------------------------------------')

bear = False
file_in = '../BTC-USD-300.txt'
candles = {}
candles['5 minute'] = []
candles['30 minute'] = []
candles['1 hour'] = []
candles['6 hour'] = []
candles['1 day'] = []
candles['7 day'] = []
candles['30 day'] = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        if candle.time < 1513515600 and bear:
            continue
        candles['5 minute'].append(candle)
        if candle.time % 1800 == 0:
            candles['30 minute'].append(candle)
        if candle.time % 3600 == 0:
            candles['1 hour'].append(candle)
        if candle.time % 21600 == 0:
            candles['6 hour'].append(candle)
        if candle.time % 86400 == 0:
            candles['1 day'].append(candle)
        if candle.time % 604800 == 0:
            candles['7 day'].append(candle)
        if candle.time % 2592000 == 0:
            candles['30 day'].append(candle)

fees = 0.003
funds = 1000.0
intervals = 22

todo = []

strat = Strategy('green maru', 0.1)
strat.buy.append(strategy.green_maru)
todo.append(strat)

ls = []
for interval, values in candles.items():
    for test in todo:
        print('testing...', end=' ', flush=True)
        data = simulation.run(values, intervals, funds, fees, test, False)
        data.insert(0, test)
        ls.append(data)
        data.append(interval)

ls.sort(key=itemgetter(1), reverse=True)

for index in range(min(5, len(ls))):
    print('----------------------------------------')
    top = ls[index]
    print('top', index + 1, top[0].name, top[7])
    print('total ${:,.2f} - coins {:,.3f} - low ${:,.2f} - high ${:,.2f} - buys {:,} - sells {:,}'.format(top[1], top[2], top[3], top[4], top[5], top[6]))