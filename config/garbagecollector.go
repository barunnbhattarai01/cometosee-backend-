package config

import (
	"log"
	"os"
	"runtime/debug"
	"strconv"
)

// go knowws meory limit in bytes
const MB = 1 << 20

func SetupGarbageCollector() {

	if gc, err := strconv.Atoi(os.Getenv("GOGC")); err == nil {
		debug.SetGCPercent(gc)
		log.Printf("GC=%d", gc)
	}

	//memoryy limitt
	if mem, err := strconv.Atoi(os.Getenv("MEM_LIMIT_MB")); err == nil {
		debug.SetMemoryLimit(int64(mem) * MB)
		log.Printf(" MemoryLimit=%dMB", mem)
	}

}
