import random
import math

learning_rate = 0.09
learning_momentum = 0.015


def combine_networks(a, b):
    num_layers = len(a.layers)
    inputs = len(a.layers[0])
    outputs = len(a.layers[-1])
    
    hidden = []
    for index in range(1, num_layers - 1):
        layer_len = len(a.layers[index])
        hidden.append(layer_len)
    
    n = Network(inputs, hidden, outputs)
    
    for index in range(1, num_layers - 1):
        layer_a = a.layers[index]
        layer_b = b.layers[index]
        layer_n = n.layers[index]
        layer_len = len(layer_a)
        for jindex in range(layer_len):
            neuron_a = layer_a[jindex]
            neuron_b = layer_b[jindex]
            neuron_n = layer_n[jindex]
            synapse_len = len(neuron_a.synapses)
            for kindex in range(synapse_len):
                synapse_a = neuron_a.synapses[kindex]
                synapse_b = neuron_b.synapses[kindex]
                synapse_n = neuron_n.synapses[kindex]
                synapse_n.weight = (synapse_a.weight + synapse_b.weight) * 0.5
                
    return n
    
def tanh(x):
    return math.tanh(x)


def tanh_derivative(x):
    return 1.0 - pow(math.tanh(x), 2)


def sigmoid(x):
    return 1.0 / (1.0 + math.exp(-x))


def sigmoid_derivative(x):
    return x * (1.0 - x)


def rectify(x):
    return max(0.0, x)


def rectify_derivative(x):
    if x > 0.0:
        return 1.0
    return 0.0


def neuron_mean(layer):
    sum = 0.0
    for neuron in layer:
        sum += neuron.output
    return sum / float(len(layer))


def neuron_variance(layer, mean):
    sum = 0.0
    for neuron in layer:
        sum += pow(neuron.output - mean, 2)
    return sum / float(len(layer))


class Synapse:
    def __init__(self, neuron):
        self.neuron = neuron
        self.weight = 2.0 * random.random() - 1.0
        self.derivative_weight = 0.0


class Neuron:
    def __init__(self, previous_layer):
        self.synapses = []
        self.error = 0.0
        self.gradient = 0.0
        self.output = 0.0
        if previous_layer:
            for neuron in previous_layer:
                self.synapses.append(Synapse(neuron))

    def feed_forward(self, activate):
        if not self.synapses:
            return
        sum = 0.0
        for synapse in self.synapses:
            sum += synapse.neuron.output * synapse.weight
        self.output = activate(sum)

    def back_propagate(self, d_activate):
        self.gradient = self.error * d_activate(self.output)
        for synapse in self.synapses:
            synapse.derivative_weight = learning_rate * synapse.neuron.output * self.gradient + learning_momentum * synapse.derivative_weight
            synapse.weight += synapse.derivative_weight
            synapse.neuron.error += synapse.weight * self.gradient
        self.error = 0.0


class Network:
    def __init__(self, inputs, hidden, outputs, activation='sigmoid'):
        self.layers = []
        if activation == 'rectify':
            self.activation = rectify
            self.activation_derivative = rectify_derivative
        elif activation == 'sigmoid':
            self.activation = sigmoid
            self.activation_derivative = sigmoid_derivative
        elif activation == 'tanh':
            self.activation = tanh
            self.activation_derivative = tanh_derivative

        current_layer = []
        for _ in range(inputs):
            current_layer.append(Neuron(None))
        self.layers.append(current_layer)

        for neuron_count in hidden:
            bias = Neuron(None)
            bias.output = 1.0
            current_layer = []
            for _ in range(neuron_count):
                current_layer.append(Neuron(self.layers[-1]))
            self.layers.append(current_layer)

        current_layer = []
        for _ in range(outputs):
            current_layer.append(Neuron(self.layers[-1]))
        self.layers.append(current_layer)

    def set_input(self, inputs):
        for index in range(len(inputs)):
            self.layers[0][index].output = inputs[index]

    def feed_forward(self):
        size = len(self.layers)
        for index in range(1, size):
            current_layer = self.layers[index]
            for neuron in current_layer:
                neuron.feed_forward(self.activation)

    def back_propagate(self, actual):
        size = len(actual)
        for index in range(size):
            self.layers[-1][index].error = actual[index] - self.layers[-1][index].output
        for layer in reversed(self.layers):
            for neuron in layer:
                neuron.back_propagate(self.activation_derivative)
            
    def get_error(self, actual):
        error = 0.0
        size = len(actual)
        for i in range(size):
            error += pow(actual[i] - self.layers[-1][i].output, 2)
        error /= size
        return math.sqrt(error)

    def get_results(self):
        results = []
        for neuron in self.layers[-1]:
            results.append(neuron.output)
        return results

    def predict(self, data):
        self.set_input(data)
        self.feed_forward()
        return self.get_results()

