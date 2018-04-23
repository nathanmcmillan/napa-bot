import sys
import signal
import time
import json
import os.path
import patterns
import genetics
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


print('----------------------------------------')
print('|       napa genetic simulation        |')
print('----------------------------------------')

file_in = '../candles-btc-usd.txt'
historical_candles = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        historical_candles.append(candle)
historial_candle_count = len(historical_candles)

gene_epochs = 3
gene_random_limit = 10
gene_top_mix_limit = 10
gene_history_limit = 1000
cooldown = 1.0

genetic_list = []
for _ in range(gene_epochs):
    genetic_list.sort(key=itemgetter(0), reverse=True)
    genetic_list = genetic_list[:gene_history_limit]
    genetic_random = []
    for _ in range(gene_random_limit):
        genes = Genetics()
        genes.randomize()
        genetic_random.append(genes)
    genetic_todo = []
    genetic_todo.extend(genetic_random)
    top_len = min(gene_top_mix_limit, len(genetic_list))
    for index in range(0, top_len):
        top_gene = genetic_list[index][1]
        for random_gene in genetic_random:
            genetic_todo.extend(genetics.permutate(top_gene, random_gene))
        for jindex in range(index + 1, top_len):
            genetic_todo.extend(genetics.permutate(top_gene, genetic_list[jindex][1]))
    print('trying', len(genetic_todo), 'combinations')
    for genes in genetic_todo:
        orders = []
        start = 0
        end = 26
        funds = 1000.0
        print('funds ${:.2f}'.format(funds), end=' - ', flush=True)
        while end < historial_candle_count:
            candles = historical_candles[start:end]
            signal = genes.signal(candles)
            ticker_price = candles[-1].closing
            if signal == 'buy':
                if funds > 20.0:
                    continue_buy = True
                    if genes.conditions['prevent_similar']:
                        for existing_order in orders:
                            if abs(existing_order.coin_price - ticker_price) / ticker_price < 0.05:
                                continue_buy = False
                                break
                    if continue_buy:
                        buy_size = funds * genes.conditions['buy_percent']
                        funds -= buy_size
                        orders.append(SimOrder(candles[-1].closing, None, buy_size))
            elif signal == 'sell':
                for order_to_sell in orders[:]:
                    if ticker_price > order_to_sell.coin_price * genes.conditions['sell_percent']:
                        funds += ticker_price * order_to_sell.size
                        orders.remove(order_to_sell)
            start += 1
            end += 1
        worth = 0.0
        for order in orders:
            worth += order.size * historical_candles[-1].closing
        worth += funds
        print('total ${:.2f}'.format(worth))
        genetic_list.append((worth, genes))
    time.sleep(cooldown)

genetic_list.sort(key=itemgetter(0), reverse=True)
for index in range(3):
    print('----------------------------------------')
    print('top', index + 1)
    top = genetic_list[index]
    print('buy: ', end='')
    for _, criteria in top[1].buy.items():
        print(criteria.to_string(), sep='', end=' ')
    print()
    print('sell: ', end='')
    for _, criteria in top[1].sell.items():
        print(criteria.to_string(), sep='', end=' ')
    print()
    print('conditions:', top[1].conditions)
    print('worth ${:.2f}'.format(top[0]))
'''
if candle.time < 1513515600:
    continue

buy: {trend, period: 5, signal: red}
conditions: {'prevent_similar': False, 'buy_percent': 0.5447577575919434, 'sell_percent': 0.7684866533373004}
worth $1898.11

------------------
FULL

top 1
buy: {trend, period: 12, signal: green}
sell:
conditions: {'prevent_similar': False, 'buy_percent': 0.846017712025436, 'sell_percent': 0.5868770873551986}
worth $136671.96
'''