package command

import (
	"flag"
	"fmt"
	"strings"
	"text/tabwriter"
)

func usageFunc(c *Command) string {
	var b strings.Builder

	if c.ShortHelp != "" {
		fmt.Fprintf(&b, "%s\n", c.ShortHelp)
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "Usage:")
	if c.ShortUsage != "" {
		fmt.Fprintf(&b, " %s\n", c.ShortUsage)
	} else {
		fmt.Fprintf(&b, " %s\n", c.Name)
	}
	fmt.Fprintf(&b, "\n")

	if c.LongHelp != "" {
		fmt.Fprintf(&b, "%s\n\n", c.LongHelp)
	}

	if countFlags(c.Flags) > 0 {
		fmt.Fprintf(&b, "Flags\n")
		tw := tabwriter.NewWriter(&b, 0, 2, 2, ' ', 0)
		c.Flags.VisitAll(func(f *flag.Flag) {
			_, usage := flag.UnquoteUsage(f)

			fmt.Fprintf(tw, "  --%s\t%s\n", f.Name, usage)
		})
		tw.Flush()
		fmt.Fprintf(&b, "\n")
	}

	return strings.TrimSpace(b.String()) + "\n"
}

func countFlags(fs *flag.FlagSet) (n int) {
	fs.VisitAll(func(*flag.Flag) { n++ })
	return n
}
