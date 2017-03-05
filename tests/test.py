import unittest
import requests

class TestDispenserd(unittest.TestCase):

    def test_is_running(self):
        res = requests.get('http://127.0.0.1:8282/').json()

        self.assertEqual(res['status'], 'ok')

    def test_queue_is_empty(self):
        res = requests.get('http://127.0.0.1:8282/jobs').json()

        self.assertEqual(len(res), 0)

suite = unittest.TestLoader().loadTestsFromTestCase(TestDispenserd)
unittest.TextTestRunner(verbosity=2).run(suite)
