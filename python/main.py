import logging
import sys
import signal
import time
import json
import gdax
import trading
import patterns
import printing
from trends import MovingAverage, ConvergeDiverge
from momentum import MoneyFlow, RelativeStrength, OnBalanceVolume
from safefile import SafeFile
from auth import Auth
from datetime import datetime
from datetime import timedelta

run = True


def read_map(path):
    map = {}
    with open(path, 'r') as open_file:
        for line in open_file:
            (key, value) = line.split()
            map[key] = value
    return map


def read_float_map(path):
    map = {}
    with open(path, 'r') as open_file:
        for line in open_file:
            (key, value) = line.split()
            map[key] = float(value)
    return map


def read_list(path):
    ls = []
    with open(path, 'r') as open_file:
        for line in open_file:
            ls.append(line.strip())
    return ls


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


print('----------------------------------------')
print('|               napa bot               |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

printing.init()

funds_file = SafeFile('../funds.txt', '../funds_backup.txt', '../funds_update.txt', '../funds_update_backup.txt')
orders_file = SafeFile('../orders.txt', '../orders_backup.txt', '../orders_update.txt', '../orders_update_backup.txt')

auth = Auth(read_map('../../private.txt'))
settings = read_map('../settings.txt')
funds = read_float_map(funds_file.path)
order_id_list = read_list(orders_file.path)

print('funds', funds)
print('settings', settings)
print('orders', order_id_list)

ema_short = int(settings['ema-short'])
ema_long = int(settings['ema-long'])
time_interval = float(settings['granularity'])
time_offset = ema_long * time_interval
product = settings['product']
granularity = settings['granularity']

orders = []
for order_id in order_id_list:
    current_order, status = gdax.get_order(auth, order_id)
    print('order', current_order.id, current_order.side, current_order.executed_value, current_order.fill_fees)
    orders.append(current_order)

quick_time = '%I:%M:%S %p'
expanded_time = '%m-%d %I:%M:%S %p'
last_candle_time = 0.0
first_iteration = True
analysis_text = ''

money_flow_index = MoneyFlow(ema_short)
relative_strength_index = RelativeStrength(ema_short)
balance_volume = OnBalanceVolume()

while run:
    end = datetime.utcnow()
    start = end - timedelta(seconds=time_offset)
    print_out = '{} - {}'.format(start.strftime(expanded_time), end.strftime(expanded_time))
    candles, status = gdax.get_candles(product, start.isoformat(), end.isoformat(), granularity)
    candle_num = len(candles)
    if candle_num >= ema_short and candles[-1].time > last_candle_time:
        last_candle_time = candles[-1].time
        macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
        money_flow_index.update(candles)
        relative_strength_index.update(candles)
        balance_volume.update(candles)
        for index in range(1, candle_num):
            current_candle = candles[index]
            macd.update(current_candle.closing)
        analysis_text = '{:.2f} | {}'.format(candles[-1].closing, patterns.trend(candles))
        analysis_text += ' | macd {:.2f}'.format(macd.current)
        analysis_text += ' | flow {:.2f}'.format(money_flow_index.current)
        analysis_text += ' | obv {:.2f}'.format(balance_volume.current)
        analysis_text += ' | rsi {:.2f}'.format(relative_strength_index.current)
        analysis_text += ' | hammer {}'.format(patterns.hammer(candles[-1]))
        analysis_text += ' | star {}'.format(patterns.shooting_star(candles[-1]))
        analysis_text += ' | marubozu {}'.format(patterns.marubozu(candles[-1]))
        trading.process(auth, product, orders, orders_file, funds, funds_file, macd.signal)
        if first_iteration:
            wait = time_interval - (time.time() - candles[-1].time)
            if wait < 0.0:
                wait_til = time.time() + time_interval
            else:
                wait_til = time.time() + wait
                first_iteration = False
        else:
            wait_til = time.time() + time_interval
    else:
        wait_til = time.time() + 10.0
    print_out += ' | sleeping - {}'.format(datetime.fromtimestamp(wait_til).strftime(quick_time))
    print(analysis_text)
    # print(print_out)
    while run and time.time() < wait_til:
        time.sleep(2)
print('close')
