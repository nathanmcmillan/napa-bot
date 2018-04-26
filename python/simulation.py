class SimOrder:
    def __init__(self, coin_price, size, usd):
        self.coin_price = coin_price
        if size:
            self.size = size
            self.usd = coin_price * size
        else:
            self.usd = usd
            self.size = usd / coin_price


def round(candles, intervals, funds, fees, algorithm, conditions, print_trades):
    candle_count = len(candles)
    orders = []
    low = funds
    high = funds
    coins = 0.0
    buys = 0
    sells = 0
    index = intervals
    while index < candle_count:
        signal = algorithm(candles, index - intervals, index)
        ticker_price = candles[index].closing
        if signal == 'buy':
            usd = funds * conditions['fund_percent']
            if usd > 10.0:
                orders.append(SimOrder(ticker_price, None, usd))
                usd *= (1.0 + fees)
                funds -= usd
                if funds < low:
                    low = funds
                buys += 1
                if print_trades:
                    coins += orders[-1].size
                    print('time - {} - ticker ${:,.2f} - spent ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, usd, funds, coins))
        elif signal == 'sell':
            for order_to_sell in orders[:]:
                if ticker_price > order_to_sell.coin_price * conditions['min_sell']:
                    orders.remove(order_to_sell)
                    usd = (ticker_price * order_to_sell.size) * (1.0 - fees)
                    funds += usd
                    if funds > high:
                        high = funds
                    sells += 1
                    if print_trades:
                        coins -= order_to_sell.size
                        print('time - {} - ticker ${:,.2f} - made ${:,.2f} - funds ${:,.2f} - coins {:,.3f}'.format(candles[index].time, ticker_price, usd, funds, coins))
        index += 1
    total = 0.0
    coins = 0.0
    end_price = candles[-1].closing
    for order in orders:
        total += order.size * end_price
        coins += order.size
    total += funds
    print('total ${:,.2f} - coins {:,.3f}'.format(total, coins))
    return [total, coins, low, high, buys, sells]