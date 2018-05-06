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
intervals = 26

todo = []

strat = Strategy('trend 7 intervals')
strat.buy = strategy.continue_trend
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('green maru')
strat.buy = strategy.green_maru
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('green hammer')
strat.buy = strategy.green_hammer
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('green star')
strat.buy = strategy.green_star
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('macd buy')
strat.buy = strategy.macd_buy
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('rsi buy')
strat.buy = strategy.rsi_buy
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('liquidation drop')
strat.buy = strategy.liquidation_drop
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('resistance breakout')
strat.buy = strategy.resistance_breakout
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('trend & maru')
strat.buy = strategy.trend_and_maru
strat.stop_limit = strategy.simple_trail
todo.append(strat)

strat = Strategy('trend & maru + chandelier')
strat.buy = strategy.trend_and_maru
strat.stop_limit = strategy.chandelier
todo.append(strat)

strat = Strategy('resistance breakout + chandelier')
strat.buy = strategy.resistance_breakout
strat.stop_limit = strategy.chandelier
todo.append(strat)

strat = Strategy('trend 7 intervals + large trail')
strat.buy = strategy.continue_trend
strat.stop_limit = strategy.large_trail
todo.append(strat)

strat = Strategy('resistance breakout + large trail')
strat.buy = strategy.resistance_breakout
strat.stop_limit = strategy.large_trail
todo.append(strat)

strat = Strategy('trend 7 intervals + no loss')
strat.buy = strategy.continue_trend
strat.stop_limit = strategy.no_loss
todo.append(strat)

strat = Strategy('resistance breakout + no loss')
strat.buy = strategy.resistance_breakout
strat.stop_limit = strategy.no_loss
todo.append(strat)

strat = Strategy('trend & maru + no loss')
strat.buy = strategy.trend_and_maru
strat.stop_limit = strategy.no_loss
todo.append(strat)

strat = Strategy('trend 7 intervals + sell no trend')
strat.buy = strategy.continue_trend
strat.sell = strategy.not_continue_trend
strat.stop_limit = strategy.no_loss
todo.append(strat)

strat = Strategy('derive + no loss')
strat.buy = strategy.derivative
strat.stop_limit = strategy.no_loss
todo.append(strat)

strat = Strategy('velocity reversal + no loss')
strat.buy = strategy.velocity_reversal
strat.stop_limit = strategy.no_loss
todo.append(strat)

ls = []
for interval, values in candles.items():
    for test in todo:
        print('testing...', end=' ', flush=True)
        data = simulation.run(values, intervals, funds, fees, test, False)
        data.insert(0, interval)
        data.insert(0, test)
        ls.append(data)

ls.sort(key=itemgetter(2), reverse=True)

top = set()
for algo in ls[:]:
    name = algo[0].name
    if name in top:
        ls.remove(algo)
    else:
        top.add(name)

print('----------------------------------------')
single = simulation.SimOrder(candles['5 minute'][0].closing, None, funds)
print('buy & hold ${:,.2f} - high ${:,.2f}'.format(single.size * candles['5 minute'][-1].closing, single.size * 19500.0))

for index in range(min(10, len(ls))):
    print('----------------------------------------')
    top = ls[index]
    print('top', index + 1, ':', top[0].name, ':', top[1])
    print('total ${:,.2f} - coins {:,.3f} - low ${:,.2f} - high ${:,.2f} - buys {:,} - sells {:,} - gains {:,} - losses {:,} - draw down {:,.2f}%'.format(top[2], top[3], top[4], top[5], top[6], top[7], top[8], top[9], top[10]))