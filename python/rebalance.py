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
initial_btc = 1000.0 / 2339.01
funds = {}
usd_funds = 0.0
funds['BTC'] = initial_btc
funds['NANO'] = 0.0
funds['XLM'] = 0.0
funds['VEN'] = 0.0
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


def coins_usd(open_time):
    coins = {}
    usd = 0.0
    for coin, balance in funds.items():
        value = get_usd(open_time, coin, balance)
        coins[coin] = value
        usd += value
    coins['USD'] = usd
    return coins


def coins_percent(open_time):
    coins = coins_usd(open_time)
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

    coin_value = coins_usd(open_time)
    todo = []
    for coin in available:
        print(coin, coin_value[coin], funds[coin])
        actual_percent = coin_value[coin] / coin_value['USD']
        # if abs(target_percent - actual_percent) < 0.03:
        #    continue
        usd_amount = target_percent * coin_value['USD'] - coin_value[coin]
        todo.append((coin, usd_amount))

    todo.sort(key=lambda x: x[1])
    for fund in todo:
        coin = fund[0]
        usd_amount = min(usd_funds, fund[1])
        coin_amount = get_coin(open_time, coin, usd_amount)
        funds[coin] += coin_amount * (1.0 - fees)
        usd_funds -= usd_amount

    print('-----------------------')
    previous_time = open_time
    open_time += interval

coin_value = coins_usd(end)
total_usd = coin_value['USD']
del coin_value['USD']
for coin, usd in coin_value.items():
    percent = usd / total_usd * 100.0
    print('{0:5} {1:10,.2f} % $ {2:10,.2f} / {3:10}'.format(coin, percent, usd, funds[coin]))
print('{0:5} {1:10,.2f} % $ {2:10,.2f}'.format('USD', 100.0, total_usd))
print()
print('unused $ {:,.2f}'.format(usd_funds))
print('btc $ {:,.2f} / {}'.format(get_usd(end, 'BTC', initial_btc), initial_btc))