import sys
import signal
import time
import json
import os.path
import binance
from datetime import datetime
from datetime import timedelta


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


print('----------------------------------------')
print('|         napa binance candles         |')
print('----------------------------------------')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

file_path = '../symbols/'

filter = set(['BTC', 'ETH', 'XRP', 'BCH', 'EOS', 'LTC', 'ADA', 'XLM', 'IOTA', 'NANO', 'NEO', 'XMR', 'USDT', 'DASH', 'NEM', 'VEN', 'QTUM', 'ICX', 'OMG', 'ONT', 'STEEM', 'SC'])

run = True
interval = '1d'
granularity = 24 * 60 * 60
time_interval = float(granularity) * 500.0
time_format = '%Y-%m-%d %I:%M:%S %p'

info, status = binance.get_info()
symbols = info['symbols']
assets = set()
for symbol_data in symbols:
    if not run:
        break
    symbol = symbol_data['symbol']
    base = symbol_data['baseAsset']
    if not base in filter:
        continue
    quote = symbol_data['quoteAsset']
    if not quote in filter:
        continue
    file_out = file_path + base + '-' + quote + '.txt'
    existing_start = None
    existing_end = None
    candles = {}
    print('polling', symbol)

    if os.path.exists(file_out):
        with open(file_out, 'r') as f:
            for line in f:
                candle = binance.Candle(line.split())
                candles[candle.open_time] = candle
                if not existing_start or candle.open_time < existing_start:
                    existing_start = candle.open_time
                if not existing_end or candle.open_time > existing_end:
                    existing_end = candle.open_time

    if existing_end:
        start = datetime.utcfromtimestamp(existing_end / 1000)
        while run:
            end = start + timedelta(seconds=time_interval)
            print('{} - {}'.format(start.strftime(time_format), end.strftime(time_format)))
            start_ms = int(start.timestamp()) * 1000
            end_ms = int(end.timestamp()) * 1000
            new_candles, status = binance.get_candles(symbol, interval, start_ms, end_ms)
            if status != 200:
                print('something went wrong', status)
                run = False
                break
            if len(new_candles) == 0:
                break
            for candle in new_candles:
                candles[candle.open_time] = candle
            time.sleep(1.0)
            start = end
            if start > datetime.utcnow():
                break

    if existing_start:
        end = datetime.utcfromtimestamp(existing_start / 1000.0)
    else:
        end = datetime.utcnow()
    while run:
        start = end - timedelta(seconds=time_interval)
        print('{} - {}'.format(start.strftime(time_format), end.strftime(time_format)))
        start_ms = int(start.timestamp()) * 1000
        end_ms = int(end.timestamp()) * 1000
        new_candles, status = binance.get_candles(symbol, interval, start_ms, end_ms)
        if status != 200:
            print('something went wrong', status)
            run = False
            break
        if len(new_candles) == 0:
            break
        for candle in new_candles:
            candles[candle.open_time] = candle
        time.sleep(1.0)
        end = start

    print('writing to file')
    with open(file_out, "w+") as f:
        for _, candle in sorted(candles.items()):
            line = '{} '.format(candle.open_time)
            line += '{} '.format(candle.open)
            line += '{} '.format(candle.high)
            line += '{} '.format(candle.low)
            line += '{} '.format(candle.closing)
            line += '{} '.format(candle.volume)
            line += '{} '.format(candle.close_time)
            line += '{} '.format(candle.quote_asset_volume)
            line += '{} '.format(candle.number_of_trades)
            line += '{} '.format(candle.taker_buy_base_asset_volume)
            line += '{}\n'.format(candle.taker_buy_quote_asset_volume)
            f.write(line)
    print('finished')