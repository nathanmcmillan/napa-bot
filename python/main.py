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
from datetime import datetime

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
    print('signal interrupt')
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


def gdax_get_time():
    request('GET', SITE, '/time', '')


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
orders = read_list('./orders.txt')

print(funds)
print(settings)
print(orders)

# logging.basicConfig(filename='./log.txt', level=logging.DEBUG, format='%(asctime)s : %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')
# info('hello python log')

gdax_get_time()
gdax_get_order(auth, '96a93c65-f207-41a3-95f2-a23b083a1be1')
gdax_get_candles('BTC-USD', datetime.utcnow(), datetime.utcnow(), settings['granularity'])

# while run:
# print('testing...')
# time.sleep(1)
