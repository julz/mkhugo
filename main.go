package main

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	base := flag.String("path", "docs", "")
	flag.Parse()

	r1 := regexp.MustCompile(`{{< artifact org="(.*?)" repo="(.*?)" file="(.*?)" >}}`)
	r2 := regexp.MustCompile(`{{< artifact repo="(.*?)" file="(.*?)" >}}`)

	if err := filepath.WalkDir(*base, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		out := r1.ReplaceAllString(string(in), `{{ artifact(org="$1", repo="$2", file="$3") }}`)
		out = r2.ReplaceAllString(out, `{{ artifact( repo="$1", file="$2") }}`)
		if err := ioutil.WriteFile(path, []byte(out), 0); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal("FAILED:", err)
	}
}
