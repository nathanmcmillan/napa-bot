import random
import math


def and_gate(node):
    for link in node.links:
        if link[1] and not link[0].on:
            node.on = False
            return
    node.on = True


def or_gate(node):
    node.on = False
    for link in node.links:
        if link[1] and link[0].on:
            node.on = True
            return


class Node:
    def __init__(self, layer, gate):
        self.links = []
        self.gate = gate
        self.on = False
        if layer:
            for node in layer:
                self.links.append((node, bool(random.getrandbits(1))))

    def feed(self):
        if not self.links:
            return
        self.gate(self)


def random_gate():
    if bool(random.getrandbits(1)):
        return and_gate
    else:
        return or_gate


class Net:
    def __init__(self, inputs, hidden, outputs):
        self.layers = []
        layer = []
        for _ in range(inputs):
            layer.append(Node(None, None))
        self.layers.append(layer)
        self.input_index = 0

        for count in hidden:
            layer = []
            for _ in range(count):
                layer.append(Node(self.layers[-1], random_gate()))
            self.layers.append(layer)

        layer = []
        for _ in range(outputs):
            layer.append(Node(self.layers[-1], random_gate()))
        self.layers.append(layer)

    def ready(self):
        self.input_index = 0

    def set_input(self, on):
        self.layers[0][self.input_index].on = on
        self.input_index += 1

    def feed(self):
        size = len(self.layers)
        for index in range(1, size):
            for node in self.layers[index]:
                node.feed()
        return self.layers[-1]