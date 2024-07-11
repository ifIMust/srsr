# posterboy is a test utility for poking json POST routes.
# It should be replaced with actual tests.

import json
import requests

req = {'name': 'flardman', 'address': '127.0.0.1:4951'}
res = requests.post('http://localhost:8080/register', json=req)
if res.ok:
    print(json.dumps(res.json()))

req2 = {'name': 'flardman'}
res2 = requests.post('http://localhost:8080/lookup', json=req2)
if res2.ok:
    print(json.dumps(res2.json()))
