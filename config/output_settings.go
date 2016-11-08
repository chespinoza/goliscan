package config

type OutputSettings struct {
	OutputTemplate string `short:"t" long:"output-template" description:"Output template to use"`
	UseJSON        bool   `short:"j" long:"json" description:"Print output in JSON format"`
	DirectOnly     bool   `long:"direct-only" description:"Print only direct dependencies"`
	IndirectOnly   bool   `long:"indirect-only" description:"Print only indirect dependencies"`
}
