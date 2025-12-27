package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func Run(args []string) int {
	if len(args) == 0 {
		writeUsage(os.Stdout)
		return 2
	}

	switch args[0] {
	case "-h", "--help", "help":
		writeUsage(os.Stdout)
		return 0
	case "derive":
		return runDerive(args[1:], os.Stdout, os.Stderr)
	case "batch":
		return runBatch(args[1:], os.Stdout, os.Stderr)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		writeUsage(os.Stderr)
		return 2
	}
}

func writeUsage(w io.Writer) {
	fmt.Fprintln(w, "juno-addrgen")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Offline address derivation (UFVK + index -> j1...) for Juno Cash.")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  juno-addrgen derive --ufvk <jview1...> --index <n> [--json]")
	fmt.Fprintln(w, "  juno-addrgen batch  --ufvk <jview1...> --start <n> --count <k> [--json]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Notes:")
	fmt.Fprintln(w, "  - UFVKs are sensitive (watch-only, but reveal incoming transaction details).")
	fmt.Fprintln(w, "  - This tool is offline; it never talks to junocashd or the network.")
}

func runDerive(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("derive", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	_ = fs.String("ufvk", "", "UFVK (jview1...)")
	_ = fs.String("ufvk-file", "", "Read UFVK from file")
	_ = fs.String("ufvk-env", "", "Read UFVK from env var")
	_ = fs.Uint("index", 0, "Diversifier index")
	_ = fs.Bool("json", false, "JSON output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 2
	}

	fmt.Fprintln(stderr, "not implemented")
	return 2
}

func runBatch(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("batch", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	_ = fs.String("ufvk", "", "UFVK (jview1...)")
	_ = fs.String("ufvk-file", "", "Read UFVK from file")
	_ = fs.String("ufvk-env", "", "Read UFVK from env var")
	_ = fs.Uint("start", 0, "Start diversifier index")
	_ = fs.Uint("count", 0, "Number of addresses")
	_ = fs.Bool("json", false, "JSON output")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 2
	}

	fmt.Fprintln(stderr, "not implemented")
	return 2
}

