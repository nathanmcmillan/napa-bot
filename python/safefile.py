import os
import shutil

class SafeFile:
    
    def __init__(self, path, backup_path, temp_path, temp_backup_path):
        self.path = path
        self.backup_path = backup_path
        self.temp_path = temp_path
        self.temp_backup_path = temp_backup_path
        
        
    def write(self, data):
        with open(self.temp_path, "w+") as file:
            file.write(data)
        shutil.copy(self.path, self.backup_path)
        shutil.copy(self.temp_path, self.temp_backup_path)
        os.rename(self.temp_path, self.path)