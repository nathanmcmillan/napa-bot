import neural
import time
import patterns
from trends import ConvergeDiverge
from gdax import Candle
from operator import itemgetter


class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price


def get_parameters(candles, start, end):
    '''volume_low = candles[start].volume
    volume_high = candles[start].volume
    low = candles[start].low
    high = candles[start].high
    for index in range(start, end):
        candle = candles[index]
        if candle.low < low:
            low = candle.low
        if candle.high > high:
            high = candle.high
        if candle.volume < volume_low:
            volume_low = candle.volume
        elif candle.volume > volume_high:
            volume_high = candle.volume
    price_range = high - low
    volume_range = volume_high - volume_low
    parameters = []
    for index in range(start, end):
        candle = candles[index]
        parameters.append((candle.low - low) / price_range)
        parameters.append((candle.open - low) / price_range)
        parameters.append((candle.closing - low) / price_range)
        parameters.append((candle.high - low) / price_range)
        parameters.append((candle.volume - volume_low) / volume_range)
    return parameters'''
    macd = ConvergeDiverge(12, 26, candles[start].closing)
    for index in range(start + 1, end):
        candle = candles[index]
        macd.update(candle.closing)
    parameters = []
    parameters.append(float(macd.signal == 'buy'))
    parameters.append(float(macd.signal == 'sell'))
    parameters.append(float(patterns.trend(candles, end - 12, end) == 'green'))
    parameters.append(float(patterns.trend(candles, end - 12, end) == 'red'))
    for index in range(end - 12, end):
        candle = candles[index]
        parameters.append(float(patterns.marubozu(candle) == 'green'))
        parameters.append(float(patterns.marubozu(candle) == 'red'))
        parameters.append(float(patterns.hammer(candle) == 'green'))
        parameters.append(float(patterns.hammer(candle) == 'red'))
        parameters.append(float(patterns.shooting_star(candle) == 'green'))
        parameters.append(float(patterns.shooting_star(candle) == 'red'))
        parameters.append(float(patterns.color(candle) == 'green'))
        parameters.append(float(patterns.color(candle) == 'red'))
    return parameters


print('----------------------------------------')
print('|           napa simulation            |')
print('----------------------------------------')

file_in = '../candles-btc-usd.txt'
file_out = '../network-btc-usd.txt'

candles = []
with open(file_in, 'r') as open_file:
    for line in open_file:
        candle = Candle(line.split())
        if candle.time < 1513515600:
            continue
        candles.append(candle)
candle_count = len(candles)

batch = 26
parameters = 2 + 2 + 8 * 12
networks = []
end_price = candles[-1].closing

epochs = 20
random_samples = 5
top_samples = 5
cooldown = 1.0
intra_cooldown = 0.01

for epoch in range(epochs):

    todo = []

    random_networks = []
    for _ in range(random_samples):
        network = neural.Network(parameters, [], 2)
        todo.append(network)
        random_networks.append(network)

    top = min(top_samples, len(networks))
    for index in range(0, top):
        top_network = networks[index][1]
        for random_network in random_networks:
            mix = neural.combine_networks(top_network, random_network)
            todo.append(mix)

    print('testing', len(todo), 'networks (', epoch, '/', epochs, ')')
    for network in todo:
        start = 0
        end = batch
        funds = 1000.0
        funds_high = funds
        funds_low = funds
        orders = []
        print('funds ${:.2f}'.format(funds), end=' - ', flush=True)
        while end < candle_count:
            signal = network.predict(get_parameters(candles, start, end))
            if signal[0] > 0.5 and funds > 20.0:
                ticker = candles[end - 1]
                buy_size = funds * 0.6
                funds -= buy_size
                if funds < funds_low:
                    funds_low = funds
                orders.append(SimOrder(ticker.closing, None, buy_size))
            elif signal[1] > 0.5:
                ticker = candles[end - 1]
                for order_to_sell in orders[:]:
                    if ticker.closing > order_to_sell.coin_price:
                        funds += ticker.closing * order_to_sell.size
                        if funds > funds_high:
                            funds_high = funds
                        orders.remove(order_to_sell)
            start += 1
            end += 1
        worth = 0.0
        for order in orders:
            worth += order.size * end_price
        worth += funds
        print('total ${:.2f} - low ${:.2f} - high ${:.2f}'.format(worth, funds_low, funds_high))
        networks.append((worth, network))
        time.sleep(intra_cooldown)
    time.sleep(cooldown)

networks.sort(key=itemgetter(0), reverse=True)
for index in range(3):
    print('top', index + 1, 'funds ${:.2f}'.format(networks[index][0]))

print('writing to file... ', end='', flush=True)
with open(file_out, "w+") as f:
    for layer in networks[0][1].layers:
        for neuron in layer:
            for synapse in neuron.synapses:
                f.write(str(synapse.weight) + ' ')
            f.write('\n')
        f.write('\n')
print('done')
