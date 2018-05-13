import os
import binance
import gdax
from operator import itemgetter
from collections import OrderedDict
from fractions import Fraction

print('----------------------------------------')
print('|      napa index fund back test       |')
print('----------------------------------------')


def fmt_frac(frac):
    tup = str(frac).split('/')
    if len(tup) > 1:
        return float(tup[0]) / float(tup[1])
    return float(tup[0])


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

fees = Fraction(0.001)
fee_purchase = Fraction(1) - fees
initial_btc = Fraction(1000.0 / 2339.01)
min_rebalance = Fraction(0.005)
hundred_percent = Fraction(100)
zero = Fraction(0)
one = Fraction(1)
funds = {}
funds['BTC'] = initial_btc
funds['NANO'] = Fraction(0)
funds['XLM'] = Fraction(0)
funds['VEN'] = Fraction(0)
interval = 24 * 60 * 60


def get_usd(open_time, coin, balance):
    btc = 'BTC'
    usd = zero
    if coin == 'USDT':
        usd = one
    elif coin == btc:
        usd = btc_candles[open_time].closing
    elif coin in exchanges[open_time] and btc in exchanges[open_time][coin]:
        symbol = coin + btc
        if symbol in candles[open_time]:
            btc = candles[open_time][symbol].closing
            usd = btc * btc_candles[open_time].closing
    return usd * balance


def get_coin_amount(open_time, coin, usd):
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
    return zero


def coins_usd(open_time):
    coins = {}
    usd = zero
    for coin, balance in funds.items():
        value = get_usd(open_time, coin, balance)
        coins[coin] = value
        usd += value
    coins['USD'] = usd
    return coins


def balancing(to_buy, to_sell):
    for fund in to_sell[:]:
        coin = fund[0]
        if coin != 'BTC' and coin != 'ETH':
            symbol = coin + 'BTC'
            if not symbol in candles[open_time]:
                print(symbol, 'oh shit')
                continue
            amount = fund[1]
            to_sell.remove(fund)
            usd = get_usd(open_time, coin, amount)
            btc_amount = get_coin_amount(open_time, 'BTC', usd) * fee_purchase
            funds[coin] -= amount
            funds['BTC'] += btc_amount
            found = False
            for sell_fund in to_sell:
                nest_coin = sell_fund[0]
                if nest_coin == 'BTC':
                    sell_fund[1] += btc_amount
                    found = True
                    break
            if not found:
                for buy_fund in to_buy[:]:
                    nest_coin = buy_fund[0]
                    if nest_coin == 'BTC':
                        if buy_fund[1] > btc_amount:
                            buy_fund[1] -= btc_amount
                            found = True
                        else:
                            to_buy.remove(buy_fund)
                            btc_amount -= buy_fund[1]
                        break
            if not found:
                to_sell.append(['BTC', btc_amount])

    to_buy.sort(key=lambda x: x[1], reverse=True)
    to_sell.sort(key=lambda x: x[1], reverse=True)

    len_sell_coins = len(to_sell)
    for fund in to_buy:
        coin = fund[0]
        amount = fund[1]
        usd = get_usd(open_time, coin, amount)
        sell_index = 0
        while sell_index < len_sell_coins:
            sell_tuple = to_sell[sell_index]
            sell_coin = sell_tuple[0]
            symbol = coin + sell_coin
            if not symbol in candles[open_time]:
                symbol = sell_coin + coin
                if not symbol in candles[open_time]:
                    sell_index += 1
                    print(symbol, 'does not exist')
                    continue
            sell_coin_amount = sell_tuple[1]
            coins_to_sell_for_buying = get_coin_amount(open_time, sell_coin, usd)
            if coins_to_sell_for_buying > sell_coin_amount:
                usd_of_sell_coin_amount = get_usd(open_time, sell_coin, sell_coin_amount)
                funds[coin] += get_coin_amount(open_time, coin, usd_of_sell_coin_amount) * fee_purchase
                funds[sell_coin] -= sell_coin_amount
                del to_sell[sell_index]
                len_sell_coins -= 1
                continue
            funds[coin] += amount * fee_purchase
            funds[sell_coin] -= coins_to_sell_for_buying
            to_sell[sell_index][1] -= coins_to_sell_for_buying
            break


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
    target_percent = Fraction(1, len(available))

    coin_value = coins_usd(open_time)
    target_usd = target_percent * coin_value['USD']
    to_buy = []
    to_sell = []
    for coin in available:
        print(coin, fmt_frac(coin_value[coin]), fmt_frac(funds[coin]))
        actual_percent = coin_value[coin] / coin_value['USD']
        if abs(target_percent - actual_percent) < min_rebalance:
            continue
        target_coin_amount = get_coin_amount(open_time, coin, target_usd)
        amount = target_coin_amount - funds[coin]
        if amount > zero:
            to_buy.append([coin, amount])
        else:
            to_sell.append([coin, -amount])
    balancing(to_buy, to_sell)

    print('-----------------------')
    previous_time = open_time
    open_time += interval

coin_value = coins_usd(end)
total_usd = coin_value['USD']
del coin_value['USD']
for coin, usd in coin_value.items():
    percent = usd / total_usd * hundred_percent
    print('{0:5} {1:10,.2f} % $ {2:10,.2f} / {3:10}'.format(coin, fmt_frac(percent), fmt_frac(usd), fmt_frac(funds[coin])))
print('{0:5} {1:10,.2f} % $ {2:10,.2f}'.format('USD', fmt_frac(hundred_percent), fmt_frac(total_usd)))
print()
print('btc $ {:,.2f} / {}'.format(get_usd(end, 'BTC', initial_btc), fmt_frac(initial_btc)))

print('calculate fee + tax for every order placed')
print('download full 5 minute ETH price history')
print('time sensitivity analysis / how does start date affect end value')
print('different rebalance periods')
print('percent change / weighted market cap / simple percent')
print('search for best symbol deal when exchanging coins / have to buy extra of one / map of shortest path')