import os
import binance
from operator import itemgetter
from collections import OrderedDict

print('----------------------------------------')
print('|      napa index fund back test       |')
print('----------------------------------------')

location = './symbols/'
coins = set()
symbols = set()
exchanges = {}
candles = {}
start = None
end = None
for file_in in os.listdir(location):
    with open(os.path.join(location, file_in), 'r') as open_file:
        symbol = file_in.split('.')[0]
        coin_pair = symbol.split('-')
        base = coin_pair[0]
        quote = coin_pair[1]
        coins.add(base)
        coins.add(quote)
        symbol = base + quote
        symbols.add(symbol)
        if not base in exchanges:
            exchanges[base] = set()
        if not quote in exchanges:
            exchanges[quote] = set()
        exchanges[base].add(quote)
        exchanges[quote].add(base)
        for line in open_file:
            candle = binance.Candle(line.split())
            open_time = int(candle.open_time / 1000)
            if not open_time in candles:
                candles[open_time] = {}
            candles[open_time][symbol] = candle
            if not start or open_time < start:
                start = open_time
            if not end or open_time > end:
                end = open_time
candles = OrderedDict(sorted(candles.items(), key=lambda t: t[0]))

fees = 0.001
funds = {}
funds['BTC'] = 5.0
funds['NANO'] = 5.0
funds['XLM'] = 5.0
funds['VEN'] = 5.0
interval = 24 * 60 * 60


def convert_to_usdt(open_time, coin, balance):
    btc = 'BTC'
    usdt = 'USDT'
    btcusdt = 'BTCUSDT'
    closing = -1.0

    if coin == usdt:
        closing = 1.0
    elif usdt in exchanges[coin]:
        symbol = coin + usdt
        if symbol in symbols:
            closing = candles[open_time][symbol].closing
    elif btc in exchanges[coin]:
        symbol = coin + btc
        if symbol in symbols:
            btc = candles[open_time][symbol].closing
            closing = btc * candles[open_time][btcusdt].closing

    return closing * balance


open_time = start
while open_time < end:
    prices = {}
    for symbol in symbols:
        if symbol in candles[open_time]:
            prices['pair'] = candles[open_time][symbol].closing
    open_time += interval

usd = 0.0
for coin, balance in funds.items():
    value = convert_to_usdt(end, coin, balance)
    print('{0:10} $ {1:,.2f}'.format(coin, value))
    usd += value
print('{0:10} $ {1:,.2f}'.format('USD', usd))