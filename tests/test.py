import unittest
import requests
import sys
import random
import json

class TestDispenserd(unittest.TestCase):

    base_url = 'http://127.0.0.1:8282'

    def test010_is_running(self):
        res = requests.get(self.base_url + '/')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(json['status'], 'ok')

    def test020_queue_is_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['main']), 0)

    def test030_queue_fills(self):
        for i in range(0, 100):
            res = requests.post(self.base_url + '/schedule', \
                json={'priority': random.randint(0, 125), 'message': 'job #' + str(i)})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['status'], 'ok')
            self.assertEqual(json['code'], 0)

    def test031_queue_not_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['main']), 100)

    def test032_queue_properly_ordered(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        previous_priority = -1
        previous_date = ''
        for job in json['main']:
            self.assertLessEqual(previous_priority, job['priority'])
            if previous_priority == job['priority']:
                self.assertLessEqual(previous_date, job['timestamp'])
            previous_priority = job['priority']
            previous_date = job['timestamp']

    def test033_queue_drains(self):
        for i in range(0, 100):
            res = requests.post(self.base_url + '/receive_noblock')
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['message'].startswith('job #'), True)

    def test034_queue_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['main']), 0)

    def test040_queues_fill(self):
        for i in range(0, 30):
            res = requests.post(self.base_url + '/schedule', \
                json={'lane': 'lane1', 'priority': random.randint(0, 125), 'message': 'job #' + str(i)})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['status'], 'ok')
            self.assertEqual(json['code'], 0)
        for i in range(0, 50):
            res = requests.post(self.base_url + '/schedule', \
                json={'lane': 'lane2', 'priority': random.randint(0, 125), 'message': 'job #' + str(i)})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['status'], 'ok')
            self.assertEqual(json['code'], 0)
        for i in range(0, 70):
            res = requests.post(self.base_url + '/schedule', \
                json={'lane': 'lane3', 'priority': random.randint(0, 125), 'message': 'job #' + str(i)})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['status'], 'ok')
            self.assertEqual(json['code'], 0)

    def test041_queues_not_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['lane1']), 30)
        self.assertEqual(len(json['lane2']), 50)
        self.assertEqual(len(json['lane3']), 70)

    def test042_queues_properly_ordered(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        previous_priority = -1
        previous_date = ''
        for job in json['lane1']:
            self.assertLessEqual(previous_priority, job['priority'])
            if previous_priority == job['priority']:
                self.assertLessEqual(previous_date, job['timestamp'])
            previous_priority = job['priority']
            previous_date = job['timestamp']
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        previous_priority = -1
        previous_date = ''
        for job in json['lane2']:
            self.assertLessEqual(previous_priority, job['priority'])
            if previous_priority == job['priority']:
                self.assertLessEqual(previous_date, job['timestamp'])
            previous_priority = job['priority']
            previous_date = job['timestamp']
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        previous_priority = -1
        previous_date = ''
        for job in json['lane3']:
            self.assertLessEqual(previous_priority, job['priority'])
            if previous_priority == job['priority']:
                self.assertLessEqual(previous_date, job['timestamp'])
            previous_priority = job['priority']
            previous_date = job['timestamp']

    def test043_queue1_drains(self):
        for i in range(0, 30):
            res = requests.post(self.base_url + '/receive_noblock', \
                json={'lane': 'lane1'})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['message'].startswith('job #'), True)

    def test044_queue1_empty_queue23_full(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['lane1']), 0)
        self.assertEqual(len(json['lane2']), 50)
        self.assertEqual(len(json['lane3']), 70)

    def test045_queue2_drains(self):
        for i in range(0, 50):
            res = requests.post(self.base_url + '/receive_noblock', \
                json={'lane': 'lane2'})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['message'].startswith('job #'), True)

    def test046_queue12_empty_queue3_full(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['lane1']), 0)
        self.assertEqual(len(json['lane2']), 0)
        self.assertEqual(len(json['lane3']), 70)

    def test047_queue3_drains(self):
        for i in range(0, 70):
            res = requests.post(self.base_url + '/receive_noblock', \
                json={'lane': 'lane3'})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['message'].startswith('job #'), True)

    def test048_queue123_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['lane1']), 0)
        self.assertEqual(len(json['lane2']), 0)
        self.assertEqual(len(json['lane3']), 0)

suite = unittest.TestLoader().loadTestsFromTestCase(TestDispenserd)

ret = unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
sys.exit(not ret)
