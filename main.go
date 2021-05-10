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

	artifact1 := regexp.MustCompile(`{{< artifact org="(.*?)" repo="(.*?)" file="(.*?)" >}}`)
	artifact2 := regexp.MustCompile(`{{< artifact repo="(.*?)" file="(.*?)" >}}`)

	tabs := regexp.MustCompile(`(?s){{< tabs .*?>}}(.*?)({{<? /tabs >?}}|$)`)
	tabstag := regexp.MustCompile(`({{< tabs .*?>}}|{{< /tabs >}}|{{< /tab >}})|{{ /tab }}|{{ /tabs }}`)
	tab := regexp.MustCompile(`{{% tab name="(.*?)" .*?%}}`)

	if err := filepath.WalkDir(*base, func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		in, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Fix artifact macros
		out := artifact1.ReplaceAllString(string(in), `{{ artifact(org="$1", repo="$2", file="$3") }}`)
		out = artifact2.ReplaceAllString(out, `{{ artifact( repo="$1", file="$2") }}`)

		// Fix branch tag
		out = strings.ReplaceAll(out, "{{< branch >}}", "{{ branch }}")

		// Fix tabs
		out = tabs.ReplaceAllStringFunc(out, func(t string) string {
			lines := strings.Split(t, "\n")
			r := strings.Builder{}
			for _, l := range lines {
				l = tabstag.ReplaceAllString(l, "")

				l = tab.ReplaceAllString(l, `=== "$1"`)
				padding := ""
				if !strings.HasPrefix(strings.TrimSpace(l), "===") && strings.TrimSpace(l) != "" {
					padding = "      "
				}

				l = strings.TrimSpace(l)
				r.WriteString(padding + l + "\n")
			}

			return r.String()
		})

		if err := ioutil.WriteFile(path, []byte(out), 0); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal("FAILED:", err)
	}
}
