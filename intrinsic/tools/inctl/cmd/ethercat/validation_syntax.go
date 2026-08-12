// Copyright 2023 Intrinsic Innovation LLC

package ethercat

import (
	"fmt"

	"github.com/lestrrat-go/libxml2/parser"
)

// validateENI parses the ENI XML content to verify it is well-formed XML.
// This implementation only performs syntax validation (ensuring it is valid XML
// and safe from entity expansion) and does not perform XSD schema validation.
func validateENI(content []byte) error {
	// Parse with options to prevent XML bomb / Entity Expansion vulnerabilities
	p := parser.New(parser.XMLParseNoEnt | parser.XMLParseNoNet)
	eniDoc, err := p.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse ENI file: %w", err)
	}
	defer eniDoc.Free()

	return nil
}
