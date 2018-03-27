import unittest
from safefile import SafeFile


class TestSafeFile(unittest.TestCase):
    def test_write(self):
        content = 'this is a test line\nanother test line'
        safe = SafeFile('./test_funds.txt', './test_funds_backup.txt', './test_funds_update.txt', './test_funds_update_backup.txt')
        safe.write(content)
        with open(safe.path, "r") as f:
            self.assertEqual(content, f.read())


if __name__ == '__main__':
    unittest.main()