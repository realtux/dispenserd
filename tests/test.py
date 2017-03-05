import unittest
import requests

class TestStringMethods(unittest.TestCase):

    def test_health(self):
        res = requests.get('http://127.0.0.1:8282/')

        print res.json()

    #def test_upper(self):
    #    self.assertEqual('foo'.upper(), 'FOO')

    #def test_isupper(self):
    #    self.assertTrue('FOO'.isupper())
    #    self.assertFalse('Foo'.isupper())

    #def test_split(self):
    #    s = 'hello world'
    #    self.assertEqual(s.split(), ['hello', 'world'])
        # check that s.split fails when the separator is not a string
    #    with self.assertRaises(TypeError):
    #        s.split(2)

if __name__ == '__main__':
    unittest.main()
