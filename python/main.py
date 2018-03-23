
import logging
import sys
import signal
import time

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
    run = False ???????????????????????

def info(string):
    logging.info(string)
    print(string)

print('napa bot')

signal.signal(signal.SIGINT, interrupts)
signal.signal(signal.SIGTERM, interrupts)

funds = read_map('./funds.txt')
settings = read_map('./settings.txt')
orders = read_list('./orders.txt')

print(funds)
print(settings)
print(orders)

logging.basicConfig(filename='./log.txt', level=logging.DEBUG, format='%(asctime)s : %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')
info('hello python log')

while run:
    print('testing...')
    time.sleep(1)