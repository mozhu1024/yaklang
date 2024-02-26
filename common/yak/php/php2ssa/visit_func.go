package php2ssa

import phpparser "github.com/yaklang/yaklang/common/yak/php/parser"

func (y *builder) VisitFunctionDeclaration(raw phpparser.IFunctionDeclarationContext) interface{} {
	if y == nil || raw == nil {
		return nil
	}

	i, _ := raw.(*phpparser.FunctionDeclarationContext)
	if i == nil {
		return nil
	}

	var attr string
	if ret := i.Attributes(); ret != nil {
		y.VisitAttributes(ret)
		_ = attr
	}

	isRef := i.Ampersand() != nil
	_ = isRef

	funcName := i.Identifier().GetText()
	ir := y.ir
	if funcName != "" {
		ir.SetMarkedFunction(funcName)
	}

	newFunction := ir.NewFunc(funcName)
	val := ir.ReadOrCreateVariable(funcName)
	variable := val.GetVariable(funcName)
	ir.AssignVariable(variable, newFunction)

	y.ir = ir.PushFunction(newFunction)

	{
		//y.ir = y.ir.PushFunction(funcDec, symbolTable, current)
		for _, formal := range i.FormalParameterList().(*phpparser.FormalParameterListContext).AllFormalParameter() {
			_ = formal
			// build param
		}
		y.VisitBlockStatement(i.BlockStatement())
		y.ir.Finish()
		y.ir = ir.PopFunction()
		if y.ir == nil {
			y.ir = ir
		}
	}

	//y.ir.WriteVariable(funcName, funcDec)
	return nil
}
