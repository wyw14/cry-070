package config

import "os"

type Config struct{ HTTPAddr, RuntimeDir string }

func Load() Config {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8092"
	}
	dir := os.Getenv("RUNTIME_DIR")
	if dir == "" {
		dir = "./runtime"
	}
	return Config{HTTPAddr: addr, RuntimeDir: dir}
}
