import requests
from threading import Thread, Event
import time

class ServiceRegistryClient:
    def __init__(self, server_address, client_name, client_address):
        self.heartbeat_interval_seconds = 20
        self.server_address = server_address
        self.client_name = client_name
        self.client_address = client_address

        self.is_registered = False
        self.client_id = ""
        self.heartbeat_thread = None
        self.stop = None


    def register(self):
        reg_data = {'name': self.client_name, 'address': self.client_address}
        r = requests.post(self.server_address + "/register", json=reg_data)
        if r.status_code == requests.codes.ok:
            resp_json = r.json()
            if 'id' in resp_json:
                self.client_id = resp_json['id']
                self.is_registered = True
                self.stop = Event()
                self.heartbeat_thread = Thread(target=self.keep_alive)
                self.heartbeat_thread.start()


    def deregister(self):
        self.stop.set()
        self.heartbeat_thread.join()
        self.is_registered = False
        
        dereg_data = {'id': self.client_id}
        requests.post(self.server_address + "/deregister", json=dereg_data)


    def keep_alive(self):
        heartbeat_data = {'id': self.client_id}
        while not self.stop.is_set():
            self.stop.wait(self.heartbeat_interval_seconds)
            requests.post(self.server_address + "/heartbeat", json=heartbeat_data)


# Example usage of the client in a service
if __name__ == '__main__':
    # Store the client object for the lifetime of the service
    c = ServiceRegistryClient('http://localhost:8080', 'srsrpy_test', 'http://localhost:3333')
    c.register()

    # Do work; live your life. The heartbeats will be sent.
    time.sleep(42.0)

    # At teardown time, deregister
    c.deregister()
