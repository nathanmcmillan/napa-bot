
from safe_file import SafeFile

safe = SafeFile('./test_funds.txt', './test_funds_backup.txt', './test_funds_update.txt', './test_funds_update_backup.txt')
safe.write('this is a test line\nanother test line')