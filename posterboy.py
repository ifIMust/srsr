# posterboy is a test utility for poking json POST routes.
# It should be replaced with actual tests.

import json
import requests
import time

id = ''

req = {'name': 'flardman', 'address': '127.0.0.1:4951'}
res = requests.post('http://localhost:8080/register', json=req)
if res.ok:
    j = res.json()
    print(json.dumps(j))
    id = j['id']

# print('Sleeping')
# time.sleep(5.0)


req2 = {'name': 'flardman'}
res2 = requests.post('http://localhost:8080/lookup', json=req2)

if res2.ok:
    j = res2.json()
    print(json.dumps(j))



# req3 = {'id': id}
# res3 = requests.post('http://localhost:8080/deregister', json=req3)
# if res3.ok:
#     print(json.dumps(res3.json()))
