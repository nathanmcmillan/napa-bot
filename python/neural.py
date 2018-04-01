import random
import math


class Network:
    def __init__(self):
        self.iterations = iterations
        self.learning_rate = learning_rate
        self.layers = []
        
        
    def new_hidden_layer(num):
        layer = [Neuron(1)]
        for i in range(num):
            layer.append(Neuron(random.random()))
        return layer
    
    
    def connect(self):
        size = len(self.layers)
        for index_a in range(0, size):
            a = self.layers[index_a]
            for index_b in range(index_a + 1, size):
                b = self.layers[index_b]
                a.to.append(Synapse(a, b, 1.0))
                b.from.append(Synapse(b, a, 1.0))
                
                
    def train(self, given, answer):
        # for g in given:
        #   for
        return
        
class Neuron:
    def __init__(self, weight):
        self.from = []
        self.to = []
       
        
    def forward():
        return
    
    
    def backward():
        return
    
    
class Synapse:
    def __init__(self, a, b, weight):
        self.from = a
        self.to = b
        self.weight = weight
        
    
def sigmoid (num):
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