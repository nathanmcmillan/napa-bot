import os
import binance
import gdax
from operator import itemgetter
from collections import OrderedDict
from fractions import Fraction
from datetime import datetime

print('----------------------------------------')
print('|      napa index fund back test       |')
print('----------------------------------------')


def fmt_frac(frac):
    tup = str(frac).split('/')
    if len(tup) > 1:
        return float(tup[0]) / float(tup[1])
    return float(tup[0])


btc_file = '../BTC-USD.txt'
eth_file = '../ETH-USD.txt'
cap_file = '../MARKET-CAP.txt'
coin_folder = '../symbols/'
epoch = datetime(1970, 1, 1)

btc_candles = {}
with open(btc_file, 'r') as f:
    for line in f:
        candle = gdax.Candle(line.split())
        btc_candles[candle.time] = candle

eth_candles = {}
with open(eth_file, 'r') as f:
    for line in f:
        candle = gdax.Candle(line.split())
        eth_candles[candle.time] = candle

market_caps = {}
with open(cap_file, 'r') as f:
    for line in f:
        data = line.split()
        str_date = data[0]
        year = int(str_date[0:4])
        month = int(str_date[4:6])
        day = int(str_date[6:9])
        str_datetime = datetime(year, month, day)
        coin = data[1]
        cap = Fraction(data[2])
        open_time = int((str_datetime - epoch).total_seconds())
        if not open_time in market_caps:
            market_caps[open_time] = {}
        market_caps[open_time][coin] = cap

coins = set()
symbols = set()
exchanges = {}
candles = {}
one_day = 24 * 60 * 60
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

while not start in market_caps:
    start += one_day

fees = Fraction(0.001)
fee_purchase = Fraction(1) - fees
initial_btc = Fraction(1000.0 / 2339.01)
min_rebalance = Fraction(0.005)
hundred_percent = Fraction(100)
zero = Fraction(0)
one = Fraction(1)


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


def coins_usd(funds, open_time):
    coins = {}
    usd = zero
    for coin, balance in funds.items():
        value = get_usd(open_time, coin, balance)
        coins[coin] = value
        usd += value
    coins['USD'] = usd
    return coins


def balancing(funds, to_buy, to_sell):
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


def target_simple(open_time, previous_time, available):
    percent = Fraction(1, len(available))
    targets = {}
    for coin in available:
        targets[coin] = percent
    return targets


def target_trend(open_time, previous_time, available):
    trends = []
    for coin in available:
        to = exchanges[open_time][coin]
        for to_coin in to:
            symbol = coin + to_coin
            if not symbol in candles[open_time] or not symbol in candles[previous_time]:
                symbol = to_coin + coin
            if not symbol in candles[open_time] or not symbol in candles[previous_time]:
                continue
            was = candles[open_time][symbol].closing
            now = candles[previous_time][symbol].closing
            trends.append((coin, to_coin, (now - was) / was))
    trends.sort(key=lambda x: x[2])

    top = set()
    above = 0
    for trend in trends[:]:
        coin = trend[0]
        if coin in top:
            trends.remove(trend)
        else:
            top.add(trend)
            if trend[2] > 0.0:
                above += 1

    percent = Fraction(1, max(1, above))
    targets = {}
    for coin in available:
        targets[coin] = 0.0
    for trend in trends:
        coin = trend[0]
        rate = trend[2]
        if rate > 0.0:
            targets[coin] = percent
    return targets


def target_cap(open_time, previous_time, available):
    if open_time in market_caps:
        targets = {}
        total_market_cap = Fraction(0)
        for coin in available:
            total_market_cap += market_caps[open_time][coin]
        for coin in available:
            targets[coin] = market_caps[open_time][coin] / total_market_cap
        return targets
    else:
        return target_simple(open_time, previous_time, available)


todo = []
todo.append(('simple', target_simple))
todo.append(('trend', target_trend))
todo.append(('market cap', target_cap))

intervals = [('one day', one_day), ('two days', one_day * 2), ('three days', one_day * 3), ('five days', one_day * 5), ('one week', one_day * 7), ('two weeks', one_day * 14), ('four weeks', one_day * 28)]
portfolio = [2, 5, 10, 20, 30]

ls = []
for holding in portfolio:
    for interval_pair in intervals:
        interval = interval_pair[1]
        for test in todo:
            print('testing...', end=' ', flush=True)

            name = test[0]
            algo = test[1]

            previous_time = start
            open_time = start + interval

            funds = {}
            total_market_cap = Fraction(0)
            for coin in exchanges[start]:
                if coin in market_caps[start]:
                    total_market_cap += market_caps[start][coin]
            cap = []
            for coin in exchanges[start]:
                if coin in market_caps[start]:
                    cap.append((coin, market_caps[start][coin] / total_market_cap))
            cap.sort(key=itemgetter(1), reverse=True)

            for index in range(min(holding, len(cap))):
                funds[cap[index][0]] = Fraction(0)  # need to run every interval for newly added
            funds['BTC'] = initial_btc

            while open_time < end:

                available = set()
                for coin in funds:
                    if coin in exchanges[open_time]:
                        available.add(coin)

                targets = algo(open_time, previous_time, available)
                coin_value = coins_usd(funds, open_time)

                to_buy = []
                to_sell = []
                for coin in available:
                    actual_percent = coin_value[coin] / coin_value['USD']
                    if abs(targets[coin] - actual_percent) < min_rebalance:
                        continue
                    target_usd = targets[coin] * coin_value['USD']
                    target_coin_amount = get_coin_amount(open_time, coin, target_usd)
                    amount = target_coin_amount - funds[coin]
                    if amount > zero:
                        to_buy.append([coin, amount])
                    else:
                        to_sell.append([coin, -amount])
                balancing(funds, to_buy, to_sell)

                previous_time = open_time
                open_time += interval

            coin_value = coins_usd(funds, end)
            ls.append((interval_pair, name, coin_value, coin_value['USD'], holding))
            print('$', coin_value['USD'])

print('----------------------------------------')
print('btc $ {:,.2f} / {}'.format(get_usd(end, 'BTC', initial_btc), fmt_frac(initial_btc)))

ls.sort(key=itemgetter(3), reverse=True)

for ls_tuple in ls:
    print('----------------------------------------')
    interval = ls_tuple[0][0]
    holding = ls_tuple[4]
    name = ls_tuple[1]
    print(name, '|', interval, '| holding', holding)
    coin_value = ls_tuple[2]
    total_usd = coin_value['USD']
    del coin_value['USD']
    for coin, usd in coin_value.items():
        percent = usd / total_usd * hundred_percent
        print('{0:5} {1:10,.2f} % $ {2:10,.2f} / {3:10}'.format(coin, fmt_frac(percent), fmt_frac(usd), fmt_frac(funds[coin])))
    print('{0:5} {1:10,.2f} % $ {2:10,.2f}'.format('USD', fmt_frac(hundred_percent), fmt_frac(total_usd)))

print('calculate fee + tax for every order placed')
print('download full 5 minute ETH price history')
print('time sensitivity analysis / how does start date affect end value')
print('search for best symbol deal when exchanging coins / have to buy extra of one / map of shortest path')