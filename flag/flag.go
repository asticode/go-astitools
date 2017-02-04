package astiflag

import "os"

// Subcommand retrieves the subcommand from the input Args
func Subcommand() (o string) {
	if len(os.Args) >= 2 && os.Args[1][0] != '-' {
		o = os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	}
	return
}
