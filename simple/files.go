package main

func readList(path string) []string {
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return list(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file")
		time.Sleep(time.Second)
	}
}

func readMap(path string) map[string]string {
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return hashmap(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file")
		time.Sleep(time.Second)
	}
}

func writeList(path string, list []string) {
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return list(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file")
		time.Sleep(time.Second)
	}
}