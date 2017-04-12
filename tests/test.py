import unittest
import requests
import sys
import random

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
            text = res.text
            self.assertEqual(res.status_code, 200)
            self.assertEqual(text.startswith('job #'), True)

    def test034_queue_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json['main']), 0)



suite = unittest.TestLoader().loadTestsFromTestCase(TestDispenserd)

ret = unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
sys.exit(not ret)
