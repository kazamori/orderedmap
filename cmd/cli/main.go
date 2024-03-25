package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/kazamori/orderedmap"
)

var (
	data = flag.String("data", "{}", "target json data")
)

func main() {
	flag.Parse()

	b := []byte(*data)
	var m orderedmap.OrderedMap[string, any]
	if err := json.Unmarshal(b, &m); err != nil {
		slog.Error("failed to unmarshal as json", "err", err)
		return
	}

	// pretty-print
	b2, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		slog.Error("pretty-print", "value", m, "err", err)
		return
	}
	fmt.Fprintln(os.Stderr, string(b2))
}
