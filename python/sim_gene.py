import sys
import signal
import time
import json
import os.path
import patterns
import genetics
import simulation
from gdax import Candle
from trends import ConvergeDiverge
from genetics import Genetics
from operator import itemgetter
from simulation import SimOrder
'''
bear
buy: {trend, period: 2, signal: green}
sell:
conditions: {'fund_percent': 0.9577636181452933, 'min_sell': 0.35525307067742246}
'''

print('----------------------------------------')
print('|       napa genetic simulation        |')
print('----------------------------------------')

file_in = '../BTC-USD-3600.txt'
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

fees = 0.005
funds = 1000.0
intervals = 22

epochs = 10
random_limit = 10
top_mix_limit = 10
cooldown = 2.0

bear_list = []
bull_list = []
genetic_list = []
for epoch in range(epochs):

    if epoch % 2 == 0:
        candles = candles_bear
        top_list = bear_list
    else:
        candles = candles_bull
        top_list = bull_list

    genetic_random = []
    for _ in range(random_limit):
        genes = Genetics()
        genes.randomize()
        genetic_random.append(genes)

    todo = []
    todo.extend(genetic_random)
    top_len = min(top_mix_limit, len(genetic_list))
    for index in range(0, top_len):
        top_gene = genetic_list[index][0]
        for random_gene in genetic_random:
            todo.extend(genetics.permutate(top_gene, random_gene))

    todo_len = len(todo)
    print('looking at', todo_len, 'todo')
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
        result = simulation.round(candles, intervals, funds, fees, genes.signal, genes.conditions, False)
        if result[4] > 0:
            if result[5] == 0:
                genes.sell.clear()
            result.insert(0, genes)
            genetic_list.append(result)

    genetic_len = len(genetic_list)
    print('looking at', genetic_len, 'genes')
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
result = simulation.round(candles_all, intervals, funds, fees, genes.signal, genes.conditions, True)

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
print('----------------------------------------')
print('candle count {:,}'.format(len(candles)))