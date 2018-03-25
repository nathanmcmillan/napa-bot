import logging
import sys
import signal
import time
import json
import gdax
from macd import ConvergeDiverge
from ema import MovingAverage
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
    print(' signal interrupt')
    global run
    run = False


def info(string):
    logging.info(string)
    print(string)
   

def analyze_and_trade():
    print('stuff')


print('napa bot')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

auth = read_map('../../private.txt')
funds = read_map('./funds.txt')
settings = read_map('./settings.txt')
order_id_list = read_list('./orders.txt')

print('funds', funds)
print('settings', settings)
print('orders', order_id_list)

# logging.basicConfig(filename='./log.txt', level=logging.DEBUG, format='%(asctime)s : %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')
# info('hello python log')

ema_short = int(settings['ema-short'])
ema_long = int(settings['ema-long'])
time_interval = int(settings['granularity'])
time_offset = ema_long * time_interval

accounts, status = gdax.get_accounts(auth)
for key, current_account in accounts.items():
    print('account', current_account.currency, current_account.available)

orders = []
for order_id in order_id_list:
    current_order, status = gdax.get_order(auth, order_id)
    print('order', current_order.id, current_order.side, current_order.executed_value, current_order.fill_fees)
    orders.append(current_order)

last_candle_time = 0
    
while run:
    end = datetime.utcnow()
    start = end - timedelta(seconds=time_offset)
    candles, status = gdax.get_candles('BTC-USD', start.isoformat(), end.isoformat(), settings['granularity'])
    candle_num = len(candles)
    if candle_num > 0 and candles[-1].time > last_candle_time:
        last_candle_time = candles[-1].time
        macd = ConvergeDiverge(ema_short, ema_long, candles[0].closing)
        for index in range(1, candle_num):
            current_candle = candles[index]
            macd.update(current_candle.closing)
        print(macd.current, macd.signal)
        analyze_and_trade()
        wait_til = time.time() + time_interval
    else:
        wait_til = time.time() + timedelta(seconds=10)
    print('sleeping til', wait_til)
    while run and time.time() < wait_til:
        time.sleep(2)
print('close')
