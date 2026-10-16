// Run against the primary with FINDING_ASSET_DSN. --after resumes from the last
// committed ID printed by this command. The migration must be applied first.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ca-risken/core/pkg/findingasset"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	after := flag.Uint64("after", 0, "last committed finding ID")
	size := flag.Int("chunk-size", 200, "rows per transaction")
	flag.Parse()
	if *size < 1 || *size > 1000 {
		log.Fatal("chunk-size must be between 1 and 1000")
	}
	dsn := os.Getenv("FINDING_ASSET_DSN")
	if dsn == "" {
		log.Fatal("FINDING_ASSET_DSN is required")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for {
		last, n, err := findingasset.BackfillChunk(context.Background(), db, *after, *size)
		if err != nil {
			log.Fatal(err)
		}
		*after = last
		fmt.Printf("after=%d processed=%d\n", last, n)
		if n == 0 {
			return
		}
	}
}
