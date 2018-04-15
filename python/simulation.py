import neural
from gdax import Candle
from operator import itemgetter

debug = False


class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price


def read_map(path):
    map = {}
    with open(path, 'r') as open_file:
        for line in open_file:
            (key, value) = line.split()
            map[key] = value
    return map


def get_parameters(candles, start, end):
    low = candles[start].closing
    high = candles[start].closing
    for index in range(start, end):
        candle = candles[index]
        if candle.closing < low:
            low = candle.closing
        elif candle.closing > high:
            high = candle.closing
    price_range = high - low
    parameters = []
    for index in range(start, end):
        candle = candles[index]
        percent = (candle.closing - low) / price_range
        parameters.append(percent)
    return parameters


print('----------------------------------------')
print('|           napa simulation            |')
print('----------------------------------------')

file_in = '../candles-btc-usd.txt'
candles = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        candles.append(candle)
candle_count = len(candles)

parameters = 6 # 672 (28 days of hours)
networks = []
end_price = candles[-1].closing

epochs = 5
random_samples = 1
top_samples = 5

# todo: dropout

for _ in range(epochs):
    
    todo = []
    
    random_networks = []
    for _ in range(random_samples):
        network = neural.Network(parameters, [parameters], 2)
        todo.append(network)
        random_networks.append(network)
        
    top = min(top_samples, len(networks))
    for index in range(0, top):
        top_network = networks[index][1]
        for random_network in random_networks:
            mix = neural.combine_networks(top_network, random_network)
            todo.append(mix)
        for jindex in range(index + 1, top):
            mix = neural.combine_networks(top_network, networks[jindex][1])
            todo.append(mix)

    print('testing', len(todo), 'networks')    
    for network in todo:
        start = 0
        end = parameters
        funds = 1000.0
        orders = []
        print('funds ${:.2f}'.format(funds), end=' - ', flush=True)
        while end < candle_count:
            signal = network.predict(get_parameters(candles, start, end))
            if signal[0] > 0.5 and funds > 20.0:
                ticker = candles[end - 1]
                buy_size = funds * 0.6
                funds -= buy_size
                orders.append(SimOrder(ticker.closing, None, buy_size))
            elif signal[1] > 0.5:
                ticker = candles[end - 1]
                for order_to_sell in orders[:]:
                    if ticker.closing > order_to_sell.coin_price:
                        funds += ticker.closing * order_to_sell.size
                        orders.remove(order_to_sell)
            start += 1
            end += 1
        worth = 0.0
        for order in orders:
            worth += order.size * end_price
        worth += funds
        print('total ${:.2f}'.format(worth))
        networks.append((worth, network))

networks.sort(key=itemgetter(0), reverse=True)
for index in range(3):
    print('top', index + 1, 'funds ${:.2f}'.format(networks[index][0]))
for layer in networks[0][1].layers:
    for neuron in layer:
        for synapse in neuron.synapses:
            print(str(synapse.weight), end=' ')
        print()
    print()
