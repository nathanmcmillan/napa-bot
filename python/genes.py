import sys
import signal
import time
import json
import os.path
import patterns
import genetics
import random
from genetics import GetTrend
from gdax import Candle
from trends import ConvergeDiverge
from genetics import Genetics
from operator import itemgetter

class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price


def round(candles, intervals, funds, fees, algorithm, conditions, print_trades):
    candle_count = len(candles)
    orders = []
    limits = []
    low = funds
    high = funds
    coins = 0.0
    buys = 0
    sells = 0
    index = intervals
    while index < candle_count:
        ticker_price = candles[index].closing
        signal = algorithm(candles, index)
        if signal == 'buy':
            usd = funds * conditions['fund_percent']
            if usd > 10.0:
                orders.append(SimOrder(ticker_price, None, usd))
                usd *= (1.0 + fees)
                funds -= usd
                coins += orders[-1].size
                buys += 1
                total = funds + coins * ticker_price
                if total < low:
                    low = total
                if print_trades:
                    print('time - {} - ticker ${:,.2f} - spent ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, usd, funds, coins))
        elif signal == 'sell':
            for order_to_sell in orders[:]:
                change = (ticker_price - order_to_sell.coin_price) / order_to_sell.coin_price
                if change > conditions['min_sell']:
                    orders.remove(order_to_sell)
                    usd = (ticker_price * order_to_sell.size) * (1.0 - fees)
                    funds += usd
                    coins -= order_to_sell.size
                    sells += 1
                    total = funds + coins * ticker_price
                    if total > high:
                        high = total
                    if print_trades:
                        profit = usd - order_to_sell.usd * (1.0 + fees)
                        print('time - {} - ticker ${:,.2f} - profit ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, profit, funds, coins))
        index += 1
    total = 0.0
    coins = 0.0
    end_price = candles[-1].closing
    for order in orders:
        total += order.size * end_price
        coins += order.size
    total += funds
    print('total ${:,.2f} - coins {:,.3f}'.format(total, coins))
    return [total, coins, low, high, buys, sells]


print('----------------------------------------')
print('|       napa genetic simulation        |')
print('----------------------------------------')

file_in = '../BTC-USD-3600.txt'
candles_bull = []
candles_bear = []
candles_all = []
candles_hours = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        candles_all.append(candle)
        if candle.time < 1513515600:
            candles_bull.append(candle)
        else:
            candles_bear.append(candle)
        if candle.time % 86400 == 0:
            candles_hours.append(candle)
candles = candles_bear

fees = 0.003
funds = 1000.0
intervals = 22

epochs = 2
random_limit = 20
top_mix_limit = 10
cooldown = 2.0

bear_list = []
bull_list = []
genetic_list = []
for epoch in range(epochs):

    genetic_random = []
    for _ in range(random_limit):
        genes = Genetics()
        genes.randomize()
        genetic_random.append(genes)

    todo = []
    todo.extend(genetic_random)
    top_len = min(top_mix_limit, len(genetic_list))
    for random_gene in genetic_random:
        for index in range(0, top_len):
            todo.extend(genetics.permutate(random.choice(genetic_list)[0], random_gene))

    todo_len = len(todo)
    index = 0
    while index < todo_len:
        jindex = index + 1
        while jindex < todo_len:
            if genetics.equals(todo[index], todo[jindex]):
                del todo[jindex]
                todo_len -= 1
                continue
            jindex += 1
        for existing_gene in genetic_list:
            if genetics.equals(existing_gene[0], todo[index]):
                del todo[index]
                index -= 1
                todo_len -= 1
                break
        index += 1

    for index in range(0, top_len):
        todo.extend(genetics.mutate(genetic_list[index][0]))

    print('testing', len(todo), 'combinations (', epoch, '/', epochs, ')', flush=True)
    time.sleep(cooldown)
    for genes in todo:
        result = round(candles, intervals, funds, fees, genes.signal, genes.conditions, False)
        if result[5] == 0:
            genes.sell.clear()
        result.insert(0, genes)
        genetic_list.append(result)

    genetic_len = len(genetic_list)
    index = 0
    while index < genetic_len:
        jindex = index + 1
        while jindex < genetic_len:
            if genetics.equals(genetic_list[index][0], genetic_list[jindex][0]):
                genetic_len -= 1
                if genetic_list[index][1] > genetic_list[jindex][1]:
                    del genetic_list[jindex]
                    continue
                else:
                    del genetic_list[index]
                    index -= 1
                    break
            jindex += 1
        index += 1
    genetic_list.sort(key=itemgetter(1), reverse=True)

print('----------------------------------------')
genes = genetic_list[0][0]
result = round(candles_all, intervals, funds, fees, genes.signal, genes.conditions, True)

for index in range(5):
    print('----------------------------------------')
    print('top', index + 1)
    top = genetic_list[index]
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

print('----------------------------------------')
print('candle count {:,}'.format(len(candles)))
