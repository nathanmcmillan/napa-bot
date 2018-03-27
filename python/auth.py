class Auth:
    def __init__(self, auth_data):
        self.key = auth_data['key']
        self.secret = auth_data['secret']
        self.phrase = auth_data['phrase']