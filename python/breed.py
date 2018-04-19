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

debug = False


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
print('|             napa breed               |')
print('----------------------------------------')

file_in = '../candles-btc-usd.txt'
historical_candles = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        historical_candles.append(candle)
historial_candle_count = len(historical_candles)

gene_epochs = 25
gene_random_limit = 10
gene_top_mix_limit = 10
gene_history_limit = 1000

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
            genetic_todo.append(genetics.mix(top_gene, random_gene))
        for jindex in range(index + 1, top_len):
            genetic_todo.append(genetics.mix(top_gene, genetic_list[jindex][1]))
    print('trying', len(genetic_todo), 'combinations')
    for genes in genetic_todo:
        orders = []
        start = 0
        end = 26
        funds = 1000.0
        print('funds ${:.2f}'.format(funds), end=' - ', flush=True)
        while end < historial_candle_count:
            if historical_candles[end].time < 1513515600:
                start += 1
                end += 1
                continue
            candles = historical_candles[start:end]
            signal = genes.signal(candles)
            ticker_price = candles[-1].closing
            if signal == 'buy':
                if funds > 20.0:
                    continue_buy = True
                    if genes.conditions['prevent_similar']:
                        for existing_order in orders:
                            if abs(existing_order.coin_price - ticker_price) / ticker_price < genes.conditions['similarity']:
                                if debug:
                                    print('not buying due to existing order bought at ${}'.format(ticker_price))
                                continue_buy = False
                                break
                    if continue_buy:
                        buy_size = funds * genes.conditions['buy_percent']
                        funds -= buy_size
                        orders.append(SimOrder(candles[-1].closing, None, buy_size))
                        if debug:
                            print('buy | {} | coin price ${:.2f} using ${:.2f}'.format(candles[-1].time, candles[-1].closing, buy_size))
                elif debug:
                    print('not enough funds ${:.2f}'.format(funds))
            elif signal == 'sell':
                for order_to_sell in orders[:]:
                    if ticker_price > order_to_sell.coin_price * genes.conditions['sell_percent']:
                        funds += ticker_price * order_to_sell.size
                        orders.remove(order_to_sell)
                        if debug:
                            print('sell | {} | ${:.2f} -> ${:.2f} | funds ${:.2f}'.format(candles[-1].time, order_to_sell.coin_price, ticker_price, funds))
            start += 1
            end += 1
        worth = 0.0
        for order in orders:
            worth += order.size * historical_candles[-1].closing
        worth += funds
        print('total ${:.2f}'.format(worth))
        genetic_list.append((worth, genes))

genetic_list.sort(key=itemgetter(0), reverse=True)
for index in range(3):
    print('----------------------------------------')
    print('top', index + 1)
    top = genetic_list[index]
    print(top[1].buy.to_string())
    print(top[1].sell.to_string())
    print(top[1].conditions)
    print('funds ${:.2f}'.format(top[0]))
