import random
import math


class Network:
    def __init__(self, iterations, learning_rate):
        self.iterations = iterations
        self.learning_rate = learning_rate
        self.layers = []

    def new_hidden_layer(self, num):
        layer = [Neuron(1)]
        for i in range(num):
            layer.append(Neuron(random.random()))
        return layer

    def connect(self):
        size = len(self.layers)
        return size

    def train(self, given, answer):
        # for g in given:
        #   for
        return


class Neuron:
    def __init__(self, weight):
        self.f = []
        self.to = []

    def forward(self):
        return

    def backward(self):
        return


class Synapse:
    def __init__(self, a, b, weight):
        self.f = a
        self.to = b
        self.weight = weight


def sigmoid(num):
    return 1.0 / (1.0 + math.exp(-num))


def sigmoid_derivative(num):
    return num * (1.0 - num)


# input layer:
#   neuron for each of 26 candles
#   neuron for each variable of each candle
#   neuron for each algorithm
#
# output layer:
#   neuron for buy, sell, and wait
#
# training data:
#   did price generally continue trend
#   after a signal
#
#
#