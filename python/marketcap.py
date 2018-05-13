import http.client
import time
import datetime
import re
from datetime import timedelta
from datetime import datetime

site = 'coinmarketcap.com'


def get_marketcap(end):
    con = http.client.HTTPSConnection(site, 443)
    con.putrequest('GET', '/historical/' + end + '/')
    con.endheaders()
    response = con.getresponse()
    raw = response.read()
    status = response.status
    con.close()
    return raw, status


time_format = '%Y%m%d'
end = datetime(2013, 4, 28)
now = datetime.utcnow()

marketcap_map = {}

while end < now:
    format_end = end.strftime(time_format)
    raw, status = get_marketcap(format_end)
    if status != 200:
        break
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
            print('{} {} {:,.2f}'.format(format_end, coin, usd_value))
        except:
            continue
    end += timedelta(days=7)
    print('------------------', flush=True)
    time.sleep(5.0)

file_out = 'market-cap.txt'
print('writing to file')
with open(file_out, "w+") as f:
    for key, candle in sorted(marketcap_map.items()):
        f.write('{} {} {:.2f}\n'.format(coin.time, coin.name, cap.cap))
print('finished')
