import logging
import sys
import signal
import time
import json
import gdax
import trading
from macd import ConvergeDiverge
from ema import MovingAverage
from safefile import SafeFile
from auth import Auth
from datetime import datetime
from datetime import timedelta

run = True


def read_map(path):
    map = {}
    with open(path, "r") as file:
        for line in file:
            (key, value) = line.split()
            map[key] = value
    return map


def read_list(path):
    ls = []
    with open(path, "r") as file:
        for line in file:
            ls.append(line.strip())
    return ls


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


def info(string):
    logging.info(string)
    print(string)


print('----------------------------------------')
print('|               napa bot               |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

funds_file = SafeFile('./funds.txt', './funds_backup.txt', './funds_update.txt', './funds_update_backup.txt')
orders_file = SafeFile('./orders.txt', './orders_backup.txt', './orders_update.txt', './orders_update_backup.txt')

auth = Auth(read_map('../../private.txt'))
settings = read_map('./settings.txt')
funds = read_map(funds_file.path)
order_id_list = read_list(orders_file.path)

print('funds', funds)
print('settings', settings)
print('orders', order_id_list)

# logging.basicConfig(filename='./log.txt', level=logging.DEBUG, format='%(asctime)s : %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')
# info('hello python log')

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

last_candle_time = 0.0
first_iteration = True

while run:
    end = datetime.utcnow()
    start = end - timedelta(seconds=time_offset)
    print('polling from', start.isoformat(), 'to', end.isoformat())
    candles, status = gdax.get_candles(product, start.isoformat(), end.isoformat(), granularity)
    candle_num = len(candles)
    if candle_num > 0 and candles[-1].time > last_candle_time:
        last_candle_time = candles[-1].time
        macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
        for index in range(1, candle_num):
            current_candle = candles[index]
            macd.update(current_candle.closing)
        print('macd', macd.current, 'signal', macd.signal)
        trading.process(auth, product, orders, orders_file, funds, funds_file, macd)
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
    print('sleeping til', datetime.fromtimestamp(wait_til).strftime('%I:%M:%S %p'))
    while run and time.time() < wait_til:
        time.sleep(2)
print('close')
