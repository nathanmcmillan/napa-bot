import sys
import signal
import http.client
import time
import datetime
import re
from datetime import timedelta
from datetime import datetime


def interrupts(signal, frame):
    print()
    print('signal interrupt')
    global run
    run = False


def get_marketcap(end):
    con = http.client.HTTPSConnection(site, 443)
    con.putrequest('GET', '/historical/' + end + '/')
    con.endheaders()
    response = con.getresponse()
    raw = response.read()
    status = response.status
    con.close()
    return raw, status


run = True
time_format = '%Y%m%d'
end = datetime(2013, 4, 28)
now = datetime.utcnow()

site = 'coinmarketcap.com'
marketcap_map = {}

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

while end < now and run:
    format_end = end.strftime(time_format)
    raw, status = get_marketcap(format_end)
    if status != 200:
        break
    marketcap_map[format_end] = {}
    assets = re.split('<td class="text-left col-symbol">', str(raw))
    del assets[0]
    for asset in assets:
        index = asset.find('</td>')
        coin = asset[:index]
        data = re.split('class="no-wrap market-cap text-right" data-usd="', asset)
        if len(data) < 2:
            continue
        index = data[1].find('"')
        cap = data[1][:index]
        try:
            usd_value = float(cap)
            print('{} {} {}'.format(format_end, coin, usd_value))
            marketcap_map[format_end][coin] = usd_value
        except:
            continue
    end += timedelta(days=7)
    print('------------------', flush=True)
    time.sleep(2.0)

file_out = '../MARKET-CAP.txt'
print('writing to file')
with open(file_out, "w+") as f:
    for open_time, coin_map in sorted(marketcap_map.items()):
        for coin, usd_value in sorted(coin_map.items()):
            f.write('{} {} {}\n'.format(open_time, coin, usd_value))
print('finished')
