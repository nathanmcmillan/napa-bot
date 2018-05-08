import binance

print('----------------------------------------')
print('|       napa binance index fund        |')
print('----------------------------------------')

info, status = binance.get_info()
symbols = info['symbols']
assets = set()
for symbol in symbols:
    assets.add(symbol['baseAsset'])
for asset in assets:
    print(asset)