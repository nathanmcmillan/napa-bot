import time
import http.client
import hmac
import hashlib
import time
import base64
import json


def prepare_request(method, site, path, body):
    con = http.client.HTTPSConnection(site, 443)
    if body:
        con.putrequest(method, path, body)
    else:
        con.putrequest(method, path)
    con.putheader('Accept', 'application/json')
    con.putheader('Content-Type', 'application/json')
    con.putheader('User-Agent', 'napa')
    return con


def request(method, site, path, body):
    con = prepare_request(method, site, path, body)
    con.endheaders()
    response = con.getresponse()
    raw_js = response.read()
    status = response.status
    con.close()
    time.sleep(0.5)
    try:
        return json.loads(raw_js.decode()), status
    except Exception:
        return raw_js, status

def private_request(auth, method, site, path, body):
    con = prepare_request(method, site, path, body)
    timestamp = str(time.time())
    message = (timestamp + method + path + body).encode()
    hmac_key = base64.b64decode(auth.secret)
    signature = hmac.new(hmac_key, message, hashlib.sha256)
    signature_b64 = base64.b64encode(signature.digest()).decode()
    con.putheader('CB-ACCESS-KEY', auth.key)
    con.putheader('CB-ACCESS-SIGN', signature_b64)
    con.putheader('CB-ACCESS-TIMESTAMP', timestamp)
    con.putheader('CB-ACCESS-PASSPHRASE', auth.phrase)
    con.endheaders()
    response = con.getresponse()
    raw_js = response.read()
    status = response.status
    con.close()
    time.sleep(0.5)
    try:
        return json.loads(raw_js.decode()), status
    except Exception:
        return raw_js, status
    
    
    