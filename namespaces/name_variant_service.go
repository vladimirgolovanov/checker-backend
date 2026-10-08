package namespaces

import (
	"strings"
)

func GetSaparateSymbols() []string {
	return []string{"-", "_", "."}
}

func GetBaseSeparator() string {
	return "-"
}

func GetVariants(name string, checker Checker) ([]string, []error) {
	name = checker.PrepareName(name)

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
	if containSaparates {
		separateSymbols := checker.GetSeparateSymbols()
		for _, separateSymbol := range separateSymbols {
			names = append(names, strings.ReplaceAll(name, baseSaparator, separateSymbol))
		}
		names = append(names, strings.ReplaceAll(name, baseSaparator, ""))
	} else {
		names = append(names, name)
	}

	var errs []error
	valid := []string{}
	for _, name := range names {
		if err := checker.ValidateName(name); err != nil {
			errs = append(errs, err)
		} else {
			valid = append(valid, name)
		}
	}

	if len(valid) == 0 {
		return nil, errs
	}

	return valid, nil
}
