import websocket
from threading import Thread
import json
import time

FEED_SITE = 'wss://ws-feed.gdax.com'

class Subscription:
    def __init__(self):
        self.thred = Thread(target=self.run_feed)
        self.thred.start()
        
    def feed_message(feed, message):
        print(message)
        feed.close()


    def feed_close(feed):
        print('socket closed')


    def feed_open(feed):
        products = ['BTC-USD', 'ETH-USD']
        channels = ['ticker']
        params = {'type': 'subscribe', 'product_ids': products, 'channels': channels}
        feed.send(json.dumps(params))


    def run_feed():
        feed = websocket.WebSocketApp(FEED_SITE, on_message=self.feed_message, on_close=self.feed_close)
        feed.on_open = self.feed_open
        feed.run_forever()
        print('thread shutdown')


if __name__ == '__main__':
    sub = Subscription()
    print('hello feed')
    sub.thred.join()