package yak2ssa

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/yaklang/yaklang/common/log"
	yak "github.com/yaklang/yaklang/common/yak/antlr4yak/parser"
	"os"
	"path/filepath"
	"strconv"
)

func (s *astbuilder) buildInclude(i *yak.IncludeStmtContext) {
	targetFile := i.StringLiteral().GetText()
	targetFile, _ = strconv.Unquote(targetFile)

	spew.Dump(targetFile)
	var newCode string
	if filepath.IsAbs(targetFile) {
		codeRaw, _ := os.ReadFile(targetFile)
		newCode = string(codeRaw)
	} else {
		filename, err := filepath.Abs(targetFile)
		if err != nil {
			log.Warnf("yaklang builder include %v failed: %v", targetFile, err)
		}
		codeRaw, _ := os.ReadFile(filename)
		newCode = string(codeRaw)
	}

	if newCode == "" {
		log.Warnf("yaklang builder include %v failed: %v", targetFile, "empty file")
		return
	}

	err := frontEnd(newCode, false, func(ast *yak.ProgramContext) {
		s.build(ast)
	})
	if err != nil {
		log.Errorf("yaklang builder include %v failed: %v", targetFile, err)
	}
}
