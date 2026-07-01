package main

func extractBonusFlag(args []string) ([]string, bool) {
	if len(args) > 0 && (args[0] == "--bonus" || args[0] == "-b") {
		return args[1:], true
	}
	return args, false
}

//each train has it own color, and each station has its own color??
