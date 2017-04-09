import unittest
import requests
import sys

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
        self.assertEqual(len(json[0]), 0)

    def test030_queue_fills(self):
        for i in range(0, 100):
            res = requests.post(self.base_url + '/schedule', json={'message': 'job #' + str(i)})
            json = res.json()
            self.assertEqual(res.status_code, 200)
            self.assertEqual(json['status'], 'ok')
            self.assertEqual(json['code'], 0)

    def test031_queue_not_empty(self):
        res = requests.get(self.base_url + '/jobs')
        json = res.json()
        self.assertEqual(res.status_code, 200)
        self.assertEqual(len(json[0]), 100)


suite = unittest.TestLoader().loadTestsFromTestCase(TestDispenserd)

ret = unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
sys.exit(not ret)
