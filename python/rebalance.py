import os
import binance
import gdax
from operator import itemgetter
from collections import OrderedDict

print('----------------------------------------')
print('|      napa index fund back test       |')
print('----------------------------------------')

btc_file = '../BTC-USD.txt'
coin_folder = '../symbols/'

btc_candles = {}
with open(btc_file, 'r') as f:
    for line in f:
        candle = gdax.Candle(line.split())
        btc_candles[candle.time] = candle

coins = set()
symbols = set()
exchanges = {}
candles = {}
start = None
end = None
for file_in in os.listdir(coin_folder):
    with open(os.path.join(coin_folder, file_in), 'r') as open_file:
        symbol = file_in.split('.')[0]
        coin_pair = symbol.split('-')
        base = coin_pair[0]
        quote = coin_pair[1]
        coins.add(base)
        coins.add(quote)
        symbol = base + quote
        symbols.add(symbol)
        for line in open_file:
            candle = binance.Candle(line.split())
            open_time = int(candle.open_time / 1000)
            if not open_time in exchanges:
                exchanges[open_time] = {}
            if not base in exchanges[open_time]:
                exchanges[open_time][base] = set()
            if not quote in exchanges[open_time]:
                exchanges[open_time][quote] = set()
            exchanges[open_time][base].add(quote)
            exchanges[open_time][quote].add(base)
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
usd_funds = 0.0
funds['BTC'] = 5.0
funds['NANO'] = 5.0
funds['XLM'] = 5.0
funds['VEN'] = 5.0
interval = 24 * 60 * 60


def get_usd(open_time, coin, balance):
    btc = 'BTC'
    usd = 0.0

    if coin == 'USDT':
        usd = 1.0
    elif coin == btc:
        usd = btc_candles[open_time].closing
    elif coin in exchanges[open_time] and btc in exchanges[open_time][coin]:
        symbol = coin + btc
        if symbol in candles[open_time]:
            btc = candles[open_time][symbol].closing
            usd = btc * btc_candles[open_time].closing

    return usd * balance


def get_coin(open_time, coin, usd):
    btc = 'BTC'

    if coin == 'USDT':
        return usd
    elif coin == btc:
        return usd / btc_candles[open_time].closing
    elif coin in exchanges[open_time] and btc in exchanges[open_time][coin]:
        symbol = coin + btc
        if symbol in candles[open_time]:
            btc = candles[open_time][symbol].closing
            return usd / (btc * btc_candles[open_time].closing)
    return 0.0


def coins_usd():
    coins = {}
    usd = 0.0
    for coin, balance in funds.items():
        value = get_usd(end, coin, balance)
        coins[coin] = value
        usd += value
    coins['USD'] = usd
    return coins


def coins_percent():
    coins = coins_usd()
    usd = coins['USD']
    del coins['USD']
    percents = {}
    for coin, balance in coins.items():
        percents[coin] = (balance / usd) * 100.0
    return percents


previous_time = start
open_time = start + interval
while open_time < end:
    ''' trends = []
    for symbol in symbols:
        if symbol in candles[open_time] and symbol in candles[previous_time]:
            was = candles[open_time][symbol].closing
            now = candles[previous_time][symbol].closing
            trends.append((symbol, (now - was) / was))
    trends.sort(key=lambda x: x[1]) '''

    available = set()
    for coin in funds:
        if coin in exchanges[open_time]:
            available.add(coin)
    target_percent = 1.0 / len(available)

    todo = []
    for coin, percent in coins_percent().items():
        if not coin in available:
            continue
        difference = (percent / 100.0) - target_percent
        amount = difference * funds[coin]
        print(coin, target_percent, percent, difference, amount)
        todo.append((coin, amount))

    todo.sort(key=lambda x: x[1], reverse=True)
    for fund in todo:
        coin = fund[0]
        amount = fund[1]
        funds[coin] -= amount
        usd_funds += get_usd(open_time, coin, amount)
        print('UPDATE', coin, funds[coin], usd_funds)

    if usd_funds > 0.0:
        equal = usd_funds * target_percent
        for coin in available:
            funds[coin] += get_coin(open_time, coin, equal)
            usd_funds -= equal

    previous_time = open_time
    open_time += interval

for coin, usd in coins_usd().items():
    print('{0:10} $ {1:,.2f}'.format(coin, usd))
print()
for coin, percent in coins_percent().items():
    print('{0:10} {1:,.2f} %'.format(coin, percent))

print('unused usd', usd_funds)