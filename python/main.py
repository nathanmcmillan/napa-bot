import logging
import sys
import signal
import time
import http.client
import hmac
import hashlib
import time
import base64
import json
from ema import MovingAverage
from macd import ConvergeDiverge
from datetime import datetime
from datetime import timedelta

SITE = 'api.gdax.com'
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


def request(method, site, path, body):
    con = http.client.HTTPSConnection(site, 443)
    if body:
        con.putrequest(method, path, body)
    else:
        con.putrequest(method, path)
    con.putheader('Accept', 'application/json')
    con.putheader('Content-Type', 'application/json')
    con.putheader('User-Agent', 'napa')
    con.endheaders()
    response = con.getresponse()
    print(response.read(), response.status, response.reason)
    con.close()
    time.sleep(0.5)


def private_request(auth, method, site, path, body):
    con = http.client.HTTPSConnection(site, 443)
    if body:
        con.putrequest(method, path, body)
    else:
        con.putrequest(method, path)

    con.putheader('Accept', 'application/json')
    con.putheader('Content-Type', 'application/json')
    con.putheader('User-Agent', 'napa')

    timestamp = str(time.time())
    message = (timestamp + method + path + body).encode('utf-8')
    hmac_key = base64.b64decode(auth['secret'])
    signature = hmac.new(hmac_key, message, hashlib.sha256)
    signature_b64 = base64.b64encode(signature.digest()).decode('utf-8')

    con.putheader('CB-ACCESS-KEY', auth['key'])
    con.putheader('CB-ACCESS-SIGN', signature_b64)
    con.putheader('CB-ACCESS-TIMESTAMP', timestamp)
    con.putheader('CB-ACCESS-PASSPHRASE', auth['phrase'])

    con.endheaders()
    response = con.getresponse()
    print(response.read(), response.status, response.reason)
    con.close()
    time.sleep(0.5)


def gdax_get_order(auth, id):
    private_request(auth, 'GET', SITE, '/orders/' + id, '')


def gdax_get_candles(product, start, end, granularity):
    request('GET', SITE, '/products/' + product + '/candles?start=' + start + '&end=' + end + '&granularity=' + granularity, '')


print('napa bot')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

auth = read_map('../../private.txt')
funds = read_map('./funds.txt')
settings = read_map('./settings.txt')
order_ids = read_list('./orders.txt')

print('funds', funds)
print('settings', settings)
print('orders', order_ids)

# logging.basicConfig(filename='./log.txt', level=logging.DEBUG, format='%(asctime)s : %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')
# info('hello python log')

ema_short = int(settings['ema-short'])
ema_long = int(settings['ema-long'])
time_interval = int(settings['granularity'])
time_offset = ema_long * time_interval

for id in order_ids:
    gdax_get_order(auth, id)

ema = MovingAverage(ema_short, 5)
ema.update(7)
print(ema.current)

macd = ConvergeDiverge(12, 26, 10)
macd.update(45)
print(macd.current, macd.signal)

#while run:
#    end = datetime.utcnow()
#    start = end - timedelta(seconds=time_offset)# time.sleep(1)
#    gdax_get_candles('BTC-USD', start.isoformat(), end.isoformat(), settings['granularity'])
#    time.sleep(time_interval)
#print('close')
