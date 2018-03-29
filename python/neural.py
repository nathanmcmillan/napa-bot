class Synapse:
    def __init__(self, periods):
        self.link_in = []
        self.link_out = []
        self.weights = 0
        self.bias = 0
        self.value = 0
        self.threshold_high = 0
        self.threshold_low = 0

    def result(self):
        if self.value < self.threshold_low:
            return -1.0
        if self.value > self.threshold_high:
            return 1.0
        return 0.0
