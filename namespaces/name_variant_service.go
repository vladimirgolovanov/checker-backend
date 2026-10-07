package namespaces

import "strings"

func GetSaparateSymbols() []string {
	return []string{"-", "_", "."}
}

func GetBaseSeparator() string {
	return "-"
}

func GetVariants(name string, checker Checker) []string {
	names := []string{}
	baseSaparator := GetBaseSeparator()
	separates := GetSaparateSymbols()
	containSaparates := false
	for _, symbol := range separates {
		if strings.Contains(name, symbol) {
			containSaparates = true
			name = strings.ReplaceAll(name, symbol, baseSaparator)
		}
	}
	name = checker.PrepareName(name)
	if containSaparates {
		separateSymbols := checker.GetSeparateSymbols()
		for _, separateSymbol := range separateSymbols {
			names = append(names, strings.ReplaceAll(name, baseSaparator, separateSymbol))
		}
		names = append(names, strings.ReplaceAll(name, baseSaparator, ""))
	} else {
		names = append(names, name)
	}

	return names
}
