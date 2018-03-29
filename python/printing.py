import logging


def init():
    logging.basicConfig(filename='./log.txt', level=logging.DEBUG, format='%(asctime)s : %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')


def log(string):
    logging.info(string)
    print(string)
