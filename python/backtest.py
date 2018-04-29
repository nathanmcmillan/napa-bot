import sys
import signal
import time
import json
import os.path
import patterns
import genetics
import random
import simulation
from genetics import GetTrend
from gdax import Candle
from trends import ConvergeDiverge
from genetics import Genetics
from operator import itemgetter

print('----------------------------------------')
print('|              napa test               |')
print('----------------------------------------')

file_in = '../BTC-USD-300.txt'
candles_bull = []
candles_bear = []
candles_all = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        candles_all.append(candle)
        if candle.time < 1513515600:
            candles_bull.append(candle)
        else:
            candles_bear.append(candle)
candles = candles_bear

fees = 0.003
funds = 1000.0
intervals = 22

ls = []
for test in todo:
    data = simulation.go(candles, intervals, funds, fees, genes.signal, genes.conditions, False)
    data.insert(0, test)
    ls.append(data)

ls.sort(key=itemgetter(1), reverse=True)

for index in range(min(5, len(ls))):
    print('----------------------------------------')
    print('top', index + 1)
    top = ls[index]
    print('buy: ', end='')
    for _, criteria in top[0].buy.items():
        print(criteria.to_string(), sep='', end=' ')
    print()
    print('sell: ', end='')
    for _, criteria in top[0].sell.items():
        print(criteria.to_string(), sep='', end=' ')
    print()
    print('conditions:', top[0].conditions)
    print('total ${:,.2f} - coins {:,.3f} - low ${:,.2f} - high ${:,.2f} - buys {:,} - sells {:,}'.format(top[1], top[2], top[3], top[4], top[5], top[6]))
    print('entire run - ', end='')
    round(candles_all, intervals, funds, fees, genes.signal, genes.conditions, False)
