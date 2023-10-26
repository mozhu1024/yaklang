package js2ssa

import (
	"github.com/google/uuid"
	"github.com/yaklang/yaklang/common/log"
	"golang.org/x/exp/slices"

	JS "github.com/yaklang/yaklang/common/yak/antlr4JS/parser"
	"github.com/yaklang/yaklang/common/yak/ssa"
)

// entry point
func (b *astbuilder) build(ast *JS.ProgramContext) {
	if s, ok := ast.SourceElements().(*JS.SourceElementsContext); ok {
		b.buildSourceElements(s)
	}
}

// statement list
func (b *astbuilder) buildStatementList(stmtlist *JS.StatementListContext) {
	recoverRange := b.SetRange(&stmtlist.BaseParserRuleContext)
	defer recoverRange()
	allstmt := stmtlist.AllStatement()
	if len(allstmt) == 0 {
		b.NewError(ssa.Warn, TAG, "empty statement list")
	} else {
		for _, stmt := range allstmt {
			if stmt, ok := stmt.(*JS.StatementContext); ok {
				b.buildStatement(stmt)
			}
		}
	}
}

func (b *astbuilder) buildStatement(stmt *JS.StatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	// var
	if s, ok := stmt.VariableStatement().(*JS.VariableStatementContext); ok {
		b.buildVariableStatement(s)
		return
	}

	// expr
	if s, ok := stmt.ExpressionStatement().(*JS.ExpressionStatementContext); ok {
		b.buildExpressionStatement(s)
	}

	// if
	if s, ok := stmt.IfStatement().(*JS.IfStatementContext); ok {
		b.buildIfStatementContext(s)
	}

	// block
	if s, ok := stmt.Block().(*JS.BlockContext); ok {
		b.buildBlock(s)
	}

	// do while
	if s, ok := stmt.IterationStatement().(*JS.DoStatementContext); ok {
		b.buildDoStatement(s)
	}

	// for
	if s, ok := stmt.IterationStatement().(*JS.ForStatementContext); ok {
		b.buildForStatement(s)
	}

	// while
	if s, ok := stmt.IterationStatement().(*JS.WhileStatementContext); ok {
		b.buildwhileStatement(s)
	}

	// forIn
	if s, ok := stmt.IterationStatement().(*JS.ForInStatementContext); ok {
		b.buildForInStatement(s)
	}

	// forOf
	if s, ok := stmt.IterationStatement().(*JS.ForOfStatementContext); ok {
		b.buildForOfStatement(s)
	}

	// function
	if s, ok := stmt.FunctionDeclaration().(*JS.FunctionDeclarationContext); ok {
		b.buildFunctionDeclaration(s)
	}

	// ret
	if s, ok := stmt.ReturnStatement().(*JS.ReturnStatementContext); ok {
		b.buildReturnStatement(s)
	}

	// break
	if s, ok := stmt.BreakStatement().(*JS.BreakStatementContext); ok {
		b.buildBreakStatement(s)
	}

	if s, ok := stmt.LabelledStatement().(*JS.LabelledStatementContext); ok {
		b.buildLabelledStatement(s)
	}

	// try
	if s, ok := stmt.TryStatement().(*JS.TryStatementContext); ok {
		b.buildTryStatement(s)
	}

	// switch
	if s, ok := stmt.SwitchStatement().(*JS.SwitchStatementContext); ok {
		b.buildSwitchStatement(s)
	}

}

func (b *astbuilder) buildVariableStatement(stmt *JS.VariableStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s, ok := stmt.VariableDeclarationList().(*JS.VariableDeclarationListContext); ok {
		b.buildAllVariableDeclaration(s)
		return
	}

}

func (b *astbuilder) buildAllVariableDeclaration(stmt *JS.VariableDeclarationListContext) []ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()
	var ret []ssa.Value

	// TODO: 如何去实现一个不可以被重复赋值的变量 complete

	// checking varModifier - decorator (var / let / const)
	// think `var a = 1`, `let a = 1`, `const a = 1`;

	declare := ""
	if m := stmt.VarModifier(); m != nil {
		if m.Const() != nil {
			declare = "c"
		} else if m.Var() != nil {
			// 定义域特殊，允许重赋值，宽松的很
			declare = "v"
		} else if m.Let_() != nil {
			// 脑子正常的定义域处理，不允许重复定义
			declare = "l"
		} else {
			// strict mode ?
			b.NewError(ssa.Error, TAG, "wrong declare varmodifier")
			return nil
		}
		for _, jsstmt := range stmt.AllVariableDeclaration() {
			v, _ := b.buildVariableDeclaration(jsstmt, declare)
			ret = append(ret, v)
		}
		// fmt.Println(ret)
		return ret
	}
	return nil
}

func (b *astbuilder) buildVariableDeclaration(stmt JS.IVariableDeclarationContext, Type string) (ssa.Value, ssa.LeftValue) {
	a := stmt.Assign()
	varText := stmt.Assignable().GetText()

	if a == nil {
		if Type == "c" {
			v := b.GetFromCmap(varText)
			if v {
				b.NewError(ssa.Error, TAG, "the const have been declared in the block")
			} else {
				b.NewError(ssa.Error, TAG, "const must have value")
			}
			return nil, nil
		} else if Type == "l" {
			v := b.GetFromLmap(varText)
			if v {
				b.NewError(ssa.Error, TAG, "the let have been declared in the block")
				return nil, nil
			} else {
				b.AddToLmap(varText)
			}
		}

		// 返回一个any
		return ssa.NewAny(), nil
	} else {
		assignValue := func() (ssa.Value, ssa.LeftValue) {
			var lValue ssa.LeftValue

			// 得到一个左值
			if as, ok := stmt.Assignable().(*JS.AssignableContext); ok {
				lValue = b.buildAssignableContext(as)
			}

			x := stmt.SingleExpression()
			result, _ := b.buildSingleExpression(x, false)
			// fmt.Println("result :", result)

			lValue.Assign(result, b.FunctionBuilder)
			return lValue.GetValue(b.FunctionBuilder), lValue
		}

		if Type == "c" {
			v := b.GetFromCmap(varText)
			if v {
				b.NewError(ssa.Error, TAG, "the const have been declared in the block")
				return nil, nil
			} else {
				rv, lv := assignValue()
				b.AddToCmap(varText)
				return rv, lv
			}
		} else if Type == "l" {
			v := b.GetFromLmap(varText)
			if v {
				b.NewError(ssa.Error, TAG, "the let have been declared in the block")
				return nil, nil
			} else {
				rv, lv := assignValue()
				b.AddToLmap(varText)
				return rv, lv
			}
		} else {
			return assignValue()
		}

	}
}

func (b *astbuilder) buildAssignableContext(stmt *JS.AssignableContext) ssa.LeftValue {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if i := stmt.Identifier(); i != nil {
		text := i.GetText()
		_, lv := b.buildIdentifierExpression(text, true)
		return lv
	}

	return nil
}

type getSingleExpr interface {
	SingleExpression(i int) JS.ISingleExpressionContext
}

func (b *astbuilder) buildSingleExpression(stmt JS.ISingleExpressionContext, IslValue bool) (ssa.Value, ssa.LeftValue) {
	// TODO: singleExpressions unfinish

	if v := b.buildOnlyRightSingleExpression(stmt); v != nil {
		return v, nil
	} else {
		if IslValue {
			_, lValue := b.buildSingleExpressionEx(stmt, IslValue)
			return nil, lValue
		} else {
			rValue, _ := b.buildSingleExpressionEx(stmt, IslValue)
			return rValue, nil
		}
	}
}

func (b *astbuilder) buildOnlyRightSingleExpression(stmt JS.ISingleExpressionContext) ssa.Value {

	// fmt.Println("build single expression: ", stmt.GetText())

	getValue := func(single getSingleExpr, i int) ssa.Value {
		if s := single.SingleExpression(i); s != nil {
			v, _ := b.buildSingleExpression(s, false)
			return v
		} else {
			return nil
		}
	}
	getBinaryOp := func() (single getSingleExpr, Op ssa.BinaryOpcode, IsBinOp bool) {
		single, Op, IsBinOp = nil, 0, false
		for {
			// a := stmt
			// fmt.Println(a.GetText())

			// +/-
			if s, ok := stmt.(*JS.AdditiveExpressionContext); ok {
				if op := s.Plus(); op != nil {
					single, Op, IsBinOp = s, ssa.OpAdd, true
				} else if op := s.Minus(); op != nil {
					single, Op, IsBinOp = s, ssa.OpSub, true
				} else {
					break
				}
			}

			// TODO: need more expressions
			// ('==' | '!=' | '===' | '!==')
			if s, ok := stmt.(*JS.EqualityExpressionContext); ok {
				if op := s.Equals_(); op != nil {
					single, Op, IsBinOp = s, ssa.OpEq, true
				} else if op := s.NotEquals(); op != nil {
					single, Op, IsBinOp = s, ssa.OpNotEq, true
				} else {
					break
				}
			}

			// ('<' | '>' | '<=' | '>=')
			if s, ok := stmt.(*JS.RelationalExpressionContext); ok {
				if op := s.LessThan(); op != nil {
					single, Op, IsBinOp = s, ssa.OpLt, true
				} else if op := s.MoreThan(); op != nil {
					single, Op, IsBinOp = s, ssa.OpGt, true
				} else if op := s.LessThanEquals(); op != nil {
					single, Op, IsBinOp = s, ssa.OpLtEq, true
				} else if op := s.GreaterThanEquals(); op != nil {
					single, Op, IsBinOp = s, ssa.OpGtEq, true
				} else {
					break
				}
			}

			// ('<<' | '>>' | '>>>') 缺>>>
			if s, ok := stmt.(*JS.BitShiftExpressionContext); ok {
				if op := s.LeftShiftArithmetic(); op != nil {
					single, Op, IsBinOp = s, ssa.OpShl, true
				} else if op := s.RightShiftArithmetic(); op != nil {
					single, Op, IsBinOp = s, ssa.OpShr, true
				} else {
					break
				}
			}

			// ('*' | '/' | '%')
			if s, ok := stmt.(*JS.MultiplicativeExpressionContext); ok {
				if op := s.Multiply(); op != nil {
					single, Op, IsBinOp = s, ssa.OpMul, true
				} else if op := s.Divide(); op != nil {
					single, Op, IsBinOp = s, ssa.OpDiv, true
				} else if op := s.Modulus(); op != nil {
					single, Op, IsBinOp = s, ssa.OpMod, true
				} else {
					break
				}
			}

			// '^'
			if s, ok := stmt.(*JS.BitXOrExpressionContext); ok {
				if op := s.BitXOr(); op != nil {
					single, Op, IsBinOp = s, ssa.OpXor, true
				} else {
					break
				}
			}

			// '&'
			if s, ok := stmt.(*JS.BitAndExpressionContext); ok {
				if op := s.BitAnd(); op != nil {
					single, Op, IsBinOp = s, ssa.OpAnd, true
				} else {
					break
				}
			}

			// '|'
			if s, ok := stmt.(*JS.BitOrExpressionContext); ok {
				if op := s.BitOr(); op != nil {
					single, Op, IsBinOp = s, ssa.OpOr, true
				} else {
					break
				}
			}

			return
		}

		b.NewError(ssa.Error, TAG, "binary operator not support: %s", stmt.GetText())
		return
	}

	// advanced expression
	// 处理的时候需要知道哪些是高级逻辑：
	// 高级逻辑需要处理成类似 “分支” 的行为，一般都会牵扯类似“短路”特性；
	// 也不是说长得像二元运算就一定是二元运算
	// 例如：false && dump() 这个表达式，dump()是不会执行的，因为false && dump()的结果一定是false
	handlePrimaryBinaryOperation := func() ssa.Value {
		// 数学运算
		single, opcode, IsBinOp := getBinaryOp()
		if IsBinOp {
			op1 := getValue(single, 0)
			op2 := getValue(single, 1)
			if op1 == nil || op2 == nil {
				b.NewError(ssa.Error, TAG, "in operator need two expression")
				return nil
			}
			return b.EmitBinOp(opcode, op1, op2)
		}

		// fallback is right?
		b.NewError(ssa.Error, TAG, "error binary operator")
		return nil
	}

	//advanced expression
	handlerAdvancedExpression := func(cond func(string) ssa.Value, trueExpr, falseExpr func() ssa.Value) ssa.Value {
		// 逻辑运算聚合产生phi指令
		id := uuid.NewString()

		ifb := b.BuildIf()
		ifb.BuildCondition(
			func() ssa.Value {
				return cond(id)
			},
		)

		ifb.BuildTrue(
			func() {
				v := trueExpr()
				b.WriteVariable(id, v)
			},
		)

		if falseExpr != nil {
			ifb.BuildFalse(func() {
				v := falseExpr()
				b.WriteVariable(id, v)
			})
		}

		ifb.Finish()
		return b.ReadVariable(id, true)

	}

	// fmt.Println("the expression: ", stmt.GetText())

	switch s := stmt.(type) {
	case *JS.FunctionExpressionContext:
		return b.buildFunctionExpression(s)
	case *JS.ClassExpressionContext:
	case *JS.OptionalChainExpressionContext:
		// advanced
		// let c = a?.b
		// roughly means: c = a ? a.b : undefined
		// roughly means: let c = undefined; if (a) {c = a.b }
	case *JS.MemberIndexExpressionContext:
	case *JS.MemberDotExpressionContext:
	case *JS.NewExpressionContext:
	case *JS.ArgumentsExpressionContext:
		// function call
		return b.EmitCall(b.buildArgumentsExpression(s))
	case *JS.MetaExpressionContext:
	case *JS.PostIncrementExpressionContext:
		if expr := s.SingleExpression(); expr != nil {
			_, lValue := b.buildSingleExpression(expr, true)
			if v := lValue.GetValue(b.FunctionBuilder); v == nil {
				b.NewError(ssa.Error, TAG, AssignLeftSideEmpty())
				return nil
			} else {
				rValue := b.EmitBinOp(ssa.OpAdd, lValue.GetValue(b.FunctionBuilder), ssa.NewConst(1))
				lValue.Assign(rValue, b.FunctionBuilder)
				// fmt.Println("++ result: ", lValue.GetValue(b.FunctionBuilder))
				return lValue.GetValue(b.FunctionBuilder)
			}
		}
	case *JS.PostDecreaseExpressionContext:
		if expr := s.SingleExpression(); expr != nil {
			_, lValue := b.buildSingleExpression(expr, true)
			if v := lValue.GetValue(b.FunctionBuilder); v == nil {
				b.NewError(ssa.Error, TAG, AssignLeftSideEmpty())
				return nil
			} else {
				rValue := b.EmitBinOp(ssa.OpSub, lValue.GetValue(b.FunctionBuilder), ssa.NewConst(1))
				lValue.Assign(rValue, b.FunctionBuilder)
				return lValue.GetValue(b.FunctionBuilder)
			}
		}
	case *JS.DeleteExpressionContext:
	case *JS.VoidExpressionContext:
	case *JS.TypeofExpressionContext:
	case *JS.PreIncrementExpressionContext:
		if expr := s.SingleExpression(); expr != nil {
			_, lValue := b.buildSingleExpression(expr, true)
			if v := lValue.GetValue(b.FunctionBuilder); v == nil {
				b.NewError(ssa.Error, TAG, AssignLeftSideEmpty())
				return nil
			} else {
				rValue := b.EmitBinOp(ssa.OpAdd, lValue.GetValue(b.FunctionBuilder), ssa.NewConst(1))
				lValue.Assign(rValue, b.FunctionBuilder)
				// fmt.Println("++ result: ", lValue.GetValue(b.FunctionBuilder))
				return lValue.GetValue(b.FunctionBuilder)
			}
		}
	case *JS.PreDecreaseExpressionContext:
		if expr := s.SingleExpression(); expr != nil {
			_, lValue := b.buildSingleExpression(expr, true)
			if v := lValue.GetValue(b.FunctionBuilder); v == nil {
				b.NewError(ssa.Error, TAG, AssignLeftSideEmpty())
				return nil
			} else {
				rValue := b.EmitBinOp(ssa.OpSub, lValue.GetValue(b.FunctionBuilder), ssa.NewConst(1))
				lValue.Assign(rValue, b.FunctionBuilder)
				return lValue.GetValue(b.FunctionBuilder)
			}
		}
	case *JS.UnaryPlusExpressionContext:
	case *JS.UnaryMinusExpressionContext:
	case *JS.BitNotExpressionContext:
	case *JS.NotExpressionContext:
	case *JS.AwaitExpressionContext:
	case *JS.PowerExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.MultiplicativeExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.AdditiveExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.CoalesceExpressionContext:
		// advanced
		if expr := s.SingleExpression(0); expr != nil {
			rv, _ := b.buildSingleExpression(expr, false)
			if rv != nil {
				return rv
			} else {
				v, _ := b.buildSingleExpression(expr, false)
				return v
			}
		}
	case *JS.BitShiftExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.RelationalExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.InstanceofExpressionContext:
	case *JS.InExpressionContext:
	case *JS.EqualityExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.BitAndExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.BitXOrExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.BitOrExpressionContext:
		return handlePrimaryBinaryOperation()
	case *JS.LogicalAndExpressionContext:
		// advanced
		return handlerAdvancedExpression(
			func(id string) ssa.Value {
				v := getValue(s, 0)
				b.WriteVariable(id, v)
				return v
			},
			func() ssa.Value {
				return getValue(s, 1)
			},
			nil,
		)
	case *JS.LogicalOrExpressionContext:
		// advanced
		return handlerAdvancedExpression(
			func(id string) ssa.Value {
				v := getValue(s, 0)
				b.WriteVariable(id, v)
				return b.EmitUnOp(ssa.OpNot, v)
			},
			func() ssa.Value {
				return getValue(s, 1)
			},
			nil,
		)
	case *JS.TernaryExpressionContext:
		// advanced
		return handlerAdvancedExpression(
			func(_ string) ssa.Value {
				return getValue(s, 0)
			},
			func() ssa.Value {
				return getValue(s, 1)
			},
			func() ssa.Value {
				return getValue(s, 2)
			},
		)
	case *JS.AssignmentExpressionContext:
		return b.buildAssignmentExpression(s)
	case *JS.AssignmentOperatorExpressionContext:
		_, lValue := b.buildSingleExpression(s.SingleExpression(0), true)
		rValue, _ := b.buildSingleExpression(s.SingleExpression(1), false)

		if lValue == nil || rValue == nil {
			b.NewError(ssa.Error, TAG, "in operator need two expression")
			return nil
		}

		if f, ok := s.AssignmentOperator().(*JS.AssignmentOperatorContext); ok {
			return b.buildAssignmentOperatorContext(f, lValue, rValue)
		}
	case *JS.ImportExpressionContext:
	case *JS.TemplateStringExpressionContext:
	case *JS.YieldExpressionContext:
	case *JS.ThisExpressionContext:

	case *JS.IdentifierExpressionContext:
	// identify是左值那边的
	// 	rv, _ :=  b.buildIdentifierExpression(s.GetText(), false)
	// 	return rv
	case *JS.SuperExpressionContext:
	case *JS.LiteralExpressionContext:
		return b.buildLiteralExpression(s)
	case *JS.ArrayLiteralExpressionContext:
		if expr, ok := s.ArrayLiteral().(*JS.ArrayLiteralContext); ok {
			return b.buildArrayLiteral(expr)
		}
	case *JS.ObjectLiteralExpressionContext:
		if expr, ok := s.ObjectLiteral().(*JS.ObjectLiteralContext); ok {
			return b.buildObjectLiteral(expr)
		}
	case *JS.ParenthesizedExpressionContext:
		if expr, ok := s.ExpressionSequence().(*JS.ExpressionSequenceContext); ok {
			exprs := b.buildExpressionSequence(expr)
			return exprs[len(exprs)-1]
		}
	default:
		log.Warnf("not support expression: %s", stmt.GetText())
		return nil
	}
	// log.Warnf("unfinished expression")
	return nil
}

func (b *astbuilder) buildSingleExpressionEx(stmt JS.ISingleExpressionContext, IslValue bool) (ssa.Value, ssa.LeftValue) {
	//可能是左值也可能是右值的

	//标识符
	if s, ok := stmt.(*JS.IdentifierExpressionContext); ok {
		i := s.GetText()
		value, lValue := b.buildIdentifierExpression(i, IslValue)
		return value, lValue
	}

	if s, ok := stmt.(*JS.MemberIndexExpressionContext); ok {
		value, lValue := b.buildMemberIndexExpression(s, IslValue)
		return value, lValue
	}

	if s, ok := stmt.(*JS.MemberDotExpressionContext); ok {
		value, lValue := b.buildMemberDotExpression(s, IslValue)
		return value, lValue
	}

	return nil, nil
}

func (b *astbuilder) buildFunctionExpression(stmt *JS.FunctionExpressionContext) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s, ok := stmt.AnonymousFunction().(*JS.ArrowFunctionContext); ok {
		funcName := ""
		if a, ok := s.ArrowFunctionParameters().(*JS.ArrowFunctionParametersContext); ok {
			if Name := a.Identifier(); Name != nil {
				funcName = Name.GetText()
			}
		}

		newFunc, symbol := b.NewFunc(funcName)
		current := b.CurrentBlock
		buildFunc := func() {
			b.FunctionBuilder = b.PushFunction(newFunc, symbol, current)

			if p, ok := s.ArrowFunctionParameters().(*JS.ArrowFunctionParametersContext); ok {
				if f, ok := p.FormalParameterList().(*JS.FormalParameterListContext); ok {
					b.buildFormalParameterList(f)
				}
			}

			if f, ok := s.ArrowFunctionBody().(*JS.ArrowFunctionBodyContext); ok {
				if fb, ok := f.FunctionBody().(*JS.FunctionBodyContext); ok {
					b.buildFunctionBody(fb)
				} else if s := f.SingleExpression(); s != nil {
					rv, _ := b.buildSingleExpression(s, false)
					var values []ssa.Value
					values = append(values, rv)
					b.EmitReturn(values)
				}
			}

			b.Finish()
			b.FunctionBuilder = b.PopFunction()
		}

		b.AddSubFunction(buildFunc)

		if funcName != "" {
			b.WriteVariable(funcName, newFunc)
		}

		return newFunc
	} else {
		if s, ok := stmt.AnonymousFunction().(*JS.AnonymousFunctionDeclContext); ok {
			funcName := ""
			newFunc, symbol := b.NewFunc(funcName)
			current := b.CurrentBlock

			buildFunc := func() {
				b.PushFunction(newFunc, symbol, current)

				if f, ok := s.FormalParameterList().(*JS.FormalParameterListContext); ok {
					b.buildFormalParameterList(f)
				}

				if f, ok := s.FunctionBody().(*JS.FunctionBodyContext); ok {
					b.buildFunctionBody(f)
				}

				b.Finish()
				b.PopFunction()
			}

			b.AddSubFunction(buildFunc)

			return newFunc
		}
	}

	return nil
}

func (b *astbuilder) buildArgumentsExpression(stmt *JS.ArgumentsExpressionContext) *ssa.Call {
	Iscall := false
	var args []ssa.Value
	isEllipsis := false

	if s := stmt.SingleExpression(); s != nil {
		rv, _ := b.buildSingleExpression(s, false)
		if rv != nil {
			if s, ok := stmt.Arguments().(*JS.ArgumentsContext); ok {
				args, isEllipsis = b.buildArguments(s)
			}
			Iscall = true
		}
		if Iscall {
			c := b.NewCall(rv, args)
			if isEllipsis {
				c.IsEllipsis = true
			}

			return c
		}
	}
	b.NewError(ssa.Error, TAG, "call target is nil")
	return nil
}

func (b *astbuilder) buildArguments(stmt *JS.ArgumentsContext) ([]ssa.Value, bool) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()
	hasEll := false
	var v []ssa.Value

	for _, i := range stmt.AllArgument() {
		if a, ok := i.(*JS.ArgumentContext); ok {
			if a.Ellipsis() != nil {
				hasEll = true
			}

			if s := a.SingleExpression(); s != nil {
				rv, _ := b.buildSingleExpression(s, false)
				v = append(v, rv)
			} else if s := a.Identifier(); s != nil {
				text := a.Identifier().GetText()
				rv, _ := b.buildIdentifierExpression(text, false)
				v = append(v, rv)
			}
		}
	}
	return v, hasEll
}

func (b *astbuilder) buildAssignmentOperatorContext(stmt *JS.AssignmentOperatorContext, lValue ssa.LeftValue, rValue ssa.Value) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	var Op ssa.BinaryOpcode
	if op := stmt.PlusAssign(); op != nil {
		Op = ssa.OpAdd
	} else if op := stmt.MinusAssign(); op != nil {
		Op = ssa.OpSub
	} else if op := stmt.DivideAssign(); op != nil {
		Op = ssa.OpDiv
	} else if op := stmt.ModulusAssign(); op != nil {
		Op = ssa.OpMod
	} else if op := stmt.DivideAssign(); op != nil {
		Op = ssa.OpDiv
	} else if op := stmt.MultiplyAssign(); op != nil {
		Op = ssa.OpMul
	} else if op := stmt.LeftShiftArithmeticAssign(); op != nil {
		Op = ssa.OpShl
	} else if op := stmt.RightShiftArithmeticAssign(); op != nil {
		Op = ssa.OpShr
	} else if op := stmt.BitOrAssign(); op != nil {
		Op = ssa.OpOr
	} else if op := stmt.BitXorAssign(); op != nil {
		Op = ssa.OpXor
	} else if op := stmt.BitAndAssign(); op != nil {
		Op = ssa.OpAnd
	}

	// TODO:powerAssign **=, RightShiftLogicalAssign >>>=

	value := b.EmitBinOp(Op, lValue.GetValue(b.FunctionBuilder), rValue)
	lValue.Assign(value, b.FunctionBuilder)

	// fmt.Println("test assignOpreator: ", lValue.GetValue(b.FunctionBuilder))
	return lValue.GetValue(b.FunctionBuilder)
}

func (b *astbuilder) buildIdentifierExpression(text string, IslValue bool) (ssa.Value, ssa.LeftValue) {
	// recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	// defer recoverRange()

	if IslValue {
		if b.GetFromCmap(text) {
			b.NewError(ssa.Error, TAG, "const cannot be assigned")
			return nil, nil
		}

		// leftValue
		if v := b.ReadVariable(text, false); v != nil {
			switch value := v.(type) {
			case *ssa.Parameter:
				if value.IsFreeValue {
					field := b.NewCaptureField(text)
					var tmp ssa.Value = field
					ssa.ReplaceValue(v, tmp)
					if index := slices.Index(b.FreeValues, v); index != -1 {
						b.FreeValues[index] = tmp
					}
					b.SetReg(field)
					b.ReplaceVariable(text, value, field)
					return nil, field
				}
			default:
			}
		}

		lValue := ssa.NewIdentifierLV(text, b.CurrentPos)
		return nil, lValue
	} else {
		rValue := b.ReadVariable(text, true)
		// fmt.Println(rValue)
		return rValue, nil
	}
}

func (b *astbuilder) buildMemberIndexExpression(stmt *JS.MemberIndexExpressionContext, IsValue bool) (ssa.Value, ssa.LeftValue) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	// fmt.Println("memberIndex: ", stmt.GetText())

	var expr ssa.Value

	if IsValue {
		if s := stmt.SingleExpression(0); s != nil {
			expr, _ = b.buildSingleExpression(s, false)
		} else {
			b.NewError(ssa.Error, TAG, AssignLeftSideEmpty())
			return nil, nil
		}

		// left
		var index ssa.Value
		if s := stmt.SingleExpression(1); s != nil {
			index, _ = b.buildSingleExpression(s, false)
		}

		lv := b.EmitFieldMust(expr, index)
		lv.GetValue(b.FunctionBuilder)

		return nil, lv
	} else {

		if s := stmt.SingleExpression(0); s != nil {
			expr, _ = b.buildSingleExpression(s, false)
		}

		var value ssa.Value
		if s := stmt.SingleExpression(1); s != nil {
			value, _ = b.buildSingleExpression(s, false)
		}
		return b.EmitField(expr, value), nil
	}
}

func (b *astbuilder) buildMemberDotExpression(stmt *JS.MemberDotExpressionContext, IsValue bool) (ssa.Value, ssa.LeftValue) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	var expr ssa.Value

	if s := stmt.SingleExpression(); s != nil {
		expr, _ = b.buildSingleExpression(s, false)
	} else {
		b.NewError(ssa.Error, TAG, AssignLeftSideEmpty())
		return nil, nil
	}

	var index ssa.Value
	if s, ok := stmt.IdentifierName().(*JS.IdentifierNameContext); ok {
		index = ssa.NewConst(s.GetText())
	}

	if IsValue {
		lv := b.EmitFieldMust(expr, index)
		lv.GetValue(b.FunctionBuilder)
		return nil, lv
	} else {
		return b.EmitField(expr, index), nil
	}

}

func (b *astbuilder) buildAssignmentExpression(stmt *JS.AssignmentExpressionContext) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	_, op1 := b.buildSingleExpression(stmt.SingleExpression(0), true)
	op2, _ := b.buildSingleExpression(stmt.SingleExpression(1), false)

	if op1 != nil && op2 != nil {
		// text := stmt.SingleExpression(0).GetText()
		// lValue := ssa.NewIdentifierLV(text, b.CurrentPos)
		op1.Assign(op2, b.FunctionBuilder)
		// fmt.Print(text)
		// fmt.Print("=")
		// fmt.Println(op1.GetValue(b.FunctionBuilder))
	} else {
		b.NewError(ssa.Error, TAG, "AssignmentExpression cannot get right assignable: %s", stmt.GetText())
	}

	return op2
}

func (b *astbuilder) buildExpressionStatement(stmt *JS.ExpressionStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s, ok := stmt.ExpressionSequence().(*JS.ExpressionSequenceContext); ok {
		b.buildExpressionSequence(s)
	}
}

func (b *astbuilder) buildArrayLiteral(stmt *JS.ArrayLiteralContext) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	var value []ssa.Value

	if s, ok := stmt.ElementList().(*JS.ElementListContext); ok {
		for _, i := range s.AllArrayElement() {
			if e := i.Ellipsis(); e != nil {
				b.HandlerEllipsis()
			}
			if s := i.SingleExpression(); s != nil {
				rv, _ := b.buildSingleExpression(s, false)
				value = append(value, rv)
			}
		}
	}

	return b.CreateInterfaceWithVs(nil, value)
}

func (b *astbuilder) buildObjectLiteral(stmt *JS.ObjectLiteralContext) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	// TODO: ObjectLiteral propertyAssignment remain 2

	var value []ssa.Value
	var keys []ssa.Value

	for _, p := range stmt.AllPropertyAssignment() {
		var rv ssa.Value
		var key ssa.Value

		if pro, ok := p.(*JS.PropertyExpressionAssignmentContext); ok {
			if s, ok := pro.PropertyName().(*JS.PropertyNameContext); ok {
				key = b.buildPropertyName(s)
			}

			if s := pro.SingleExpression(); s != nil {
				rv, _ = b.buildSingleExpression(s, false)
			}
		} else if pro, ok := p.(*JS.ComputedPropertyExpressionAssignmentContext); ok {
			if s := pro.SingleExpression(0); s != nil {
				key, _ = b.buildSingleExpression(s, false)
			}
			if s := pro.SingleExpression(1); s != nil {
				rv, _ = b.buildSingleExpression(s, false)
			}
		} else if pro, ok := p.(*JS.FunctionPropertyContext); ok {
			var funcName string
			if s, ok := pro.PropertyName().(*JS.PropertyNameContext); ok {
				funcName = s.GetText()
			}

			newFunc, symbol := b.NewFunc(funcName)
			current := b.CurrentBlock

			buildFunc := func() {
				b.FunctionBuilder = b.PushFunction(newFunc, symbol, current)

				if s, ok := pro.FormalParameterList().(*JS.FormalParameterListContext); ok {
					b.buildFormalParameterList(s)
				}

				if f, ok := pro.FunctionBody().(*JS.FunctionBodyContext); ok {
					b.buildFunctionBody(f)
				}

				b.Finish()
				b.FunctionBuilder = b.PopFunction()

			}

			b.AddSubFunction(buildFunc)

			if funcName != "" {
				b.WriteVariable(funcName, newFunc)
			}
			return newFunc

		} else if pro, ok := p.(*JS.PropertyGetterContext); ok {
			_ = pro
			// fmt.Println(pro)
		} else if pro, ok := p.(*JS.PropertySetterContext); ok {
			_ = pro
			// fmt.Println(pro)
		} else if pro, ok := p.(*JS.PropertyShorthandContext); ok {
			if s := pro.SingleExpression(); s != nil {
				rv, _ = b.buildSingleExpression(s, false)
			}

			if pro.Ellipsis() != nil {
				b.HandlerEllipsis()
			}
		} else {
			b.NewError(ssa.Error, TAG, "Not propertyAssignment")
		}

		value = append(value, rv)
		keys = append(keys, key)
	}

	if keys[0] == nil {
		return b.CreateInterfaceWithVs(nil, value)
	}

	return b.CreateInterfaceWithVs(keys, value)
}

func (b *astbuilder) buildPropertyName(stmt *JS.PropertyNameContext) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s, ok := stmt.IdentifierName().(*JS.IdentifierNameContext); ok {
		return b.buildIdentifierName(s)
	} else if s := stmt.StringLiteral(); s != nil {
		return b.buildStringLiteral(s)
	} else if s, ok := stmt.NumericLiteral().(*JS.NumericLiteralContext); ok {
		return b.buildNumericLiteral(s)
	} else if s := stmt.SingleExpression(); s != nil {
		rv, _ := b.buildSingleExpression(s, false)
		return rv
	} else {
		b.NewError(ssa.Error, TAG, "Not support the propertyName")
	}

	return nil
}

func (b *astbuilder) buildIdentifierName(stmt *JS.IdentifierNameContext) ssa.Value {
	if s, ok := stmt.Identifier().(*JS.IdentifierContext); ok {
		text := s.GetText()
		_, lv := b.buildIdentifierExpression(text, true)
		return lv.GetValue(b.FunctionBuilder)
	} else if s := stmt.ReservedWord(); s != nil {
		if v := s.Keyword(); v != nil {
			text := v.GetText()
			return ssa.NewConst(text)
		} else if v := s.NullLiteral(); v != nil {
			return b.buildNullLiteral()
		} else if v := s.BooleanLiteral(); v != nil {
			return b.buildBooleanLiteral(stmt.GetText())
		} else {
			b.NewError(ssa.Error, TAG, "not support the format")
		}
	}
	return nil
}

func (b *astbuilder) buildExpressionSequence(stmt *JS.ExpressionSequenceContext) []ssa.Value {
	// 需要修改改函数及引用，不存在if中存在多个singleExpression的情况
	// compelte

	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	var values []ssa.Value

	for _, expr := range stmt.AllSingleExpression() {
		if s := expr; s != nil {
			rv, _ := b.buildSingleExpression(s, false)
			values = append(values, rv)
		}
	}
	return values
}

func (b *astbuilder) buildIfStatementContext(stmt *JS.IfStatementContext) {
	// var buildIf func(stmt *JS.IfStatementContext) *ssa.IfBuilder
	buildIf := func(stmt *JS.IfStatementContext) *ssa.IfBuilder {
		recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
		defer recoverRange()

		i := b.BuildIf()

		// if instruction condition
		i.BuildCondition(
			func() ssa.Value {
				if s := stmt.SingleExpression(); s != nil {
					value, _ := b.buildSingleExpression(s, false)
					return value
				}
				return nil
			})

		i.BuildTrue(
			func() {
				if s, ok := stmt.Statement(0).(*JS.StatementContext); ok {
					b.buildStatement(s)
				}
			},
		)

		if s, ok := stmt.Statement(1).(*JS.StatementContext); ok {
			if !ok {
				return i
			} else {
				i.BuildFalse(
					func() {
						b.buildStatement(s)
					},
				)
			}
		}

		return i
	}

	i := buildIf(stmt)
	i.Finish()
}

func (b *astbuilder) buildBlock(stmt *JS.BlockContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()
	b.CurrentBlock.SetPosition(b.CurrentPos)

	if s, ok := stmt.StatementList().(*JS.StatementListContext); ok {
		b.PushBlockSymbolTable()
		b.buildStatementList(s)
		b.PopBlockSymbolTable()
	} else {
		b.NewError(ssa.Warn, TAG, "empty block")
	}
}

// do while
func (b *astbuilder) buildDoStatement(stmt *JS.DoStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	// do while需要分次

	// 先进行一次do
	if s, ok := stmt.Statement().(*JS.StatementContext); ok {
		b.buildStatement(s)
	}

	// 构建循环进行条件判断
	loop := b.BuildLoop()

	var cond JS.ISingleExpressionContext

	if s := stmt.SingleExpression(); s != nil {
		cond = s
	}

	loop.BuildCondition(func() ssa.Value {
		var condition ssa.Value
		if cond == nil {
			condition = ssa.NewConst(true)
		} else {
			condition, _ = b.buildSingleExpression(cond, false)
			if condition == nil {
				condition = ssa.NewConst(true)
			}
		}
		return condition
	})

	loop.BuildBody(func() {
		if s, ok := stmt.Statement().(*JS.StatementContext); ok {
			b.buildStatement(s)
		}
	})

	loop.Finish()

}

// while
func (b *astbuilder) buildwhileStatement(stmt *JS.WhileStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	// 构建循环进行条件判断
	loop := b.BuildLoop()

	var cond JS.ISingleExpressionContext

	if s := stmt.SingleExpression(); s != nil {
		cond = s
	}

	loop.BuildCondition(func() ssa.Value {
		var condition ssa.Value
		if cond == nil {
			condition = ssa.NewConst(true)
		} else {
			condition, _ = b.buildSingleExpression(cond, false)
			if condition == nil {
				condition = ssa.NewConst(true)
			}
		}
		return condition
	})

	loop.BuildBody(func() {
		if s, ok := stmt.Statement().(*JS.StatementContext); ok {
			b.buildStatement(s)
		}
	})

	loop.Finish()

}

// for
func (b *astbuilder) buildForStatement(stmt *JS.ForStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	loop := b.BuildLoop()

	var cond JS.ISingleExpressionContext

	// fmt.Println("---------------------")
	if first, ok := stmt.ForFirst().(*JS.ForFirstContext); ok {
		if f, ok := first.VariableDeclarationList().(*JS.VariableDeclarationListContext); ok {
			loop.BuildFirstExpr(func() []ssa.Value {
				recoverRange := b.SetRange(&f.BaseParserRuleContext)
				defer recoverRange()
				return b.buildAllVariableDeclaration(f)
			})
		} else if f := first.SingleExpression(); f != nil {
			loop.BuildFirstExpr(func() []ssa.Value {
				// recoverRange := b.SetRange(&f.BaseParserRuleContext)
				// defer recoverRange()
				var ret []ssa.Value
				value, _ := b.buildSingleExpression(f, false)
				ret = append(ret, value)
				return ret
			})
		}
	}

	if expr, ok := stmt.ForSecond().(*JS.ForSecondContext); ok {
		if e, ok := expr.SingleExpression().(JS.ISingleExpressionContext); ok {
			cond = e
		}
	}

	if third, ok := stmt.ForThird().(*JS.ForThirdContext); ok {
		if t, ok := third.SingleExpression().(JS.ISingleExpressionContext); ok {
			loop.BuildThird(func() []ssa.Value {
				// build third expression in loop.latch
				// recoverRange := b.SetRange(&t.BaseParserRuleContext)
				// defer recoverRange()
				var ret []ssa.Value
				value, _ := b.buildSingleExpression(t, false)
				ret = append(ret, value)
				return ret
			})
		}
	}

	// 构建条件
	loop.BuildCondition(func() ssa.Value {
		var condition ssa.Value
		// 没有条件就是永真
		if cond == nil {
			condition = ssa.NewConst(true)
		} else {
			condition, _ = b.buildSingleExpression(cond, false)
			if condition == nil {
				condition = ssa.NewConst(true)
			}
		}
		return condition
	})

	// build body
	loop.BuildBody(func() {
		if s, ok := stmt.Statement().(*JS.StatementContext); ok {
			b.buildStatement(s)
		}
	})

	loop.Finish()
}

// for in 取key
func (b *astbuilder) buildForInStatement(stmt *JS.ForInStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	loop := b.BuildLoop()

	loop.BuildCondition(func() ssa.Value {
		var left ssa.LeftValue
		var value ssa.Value

		if s, ok := stmt.VariableDeclaration().(*JS.VariableDeclarationContext); ok {
			_, left = b.buildVariableDeclaration(s, "v")
			value, _ = b.buildSingleExpression(stmt.SingleExpression(0), false)
		} else {
			_, left = b.buildSingleExpression(stmt.SingleExpression(0), true)
			value, _ = b.buildSingleExpression(stmt.SingleExpression(1), false)
		}

		key, _, ok := b.EmitNext(value, false)
		left.Assign(key, b.FunctionBuilder)

		return ok
	})

	loop.BuildBody(func() {
		if s, ok := stmt.Statement().(*JS.StatementContext); ok {
			b.buildStatement(s)
		}
	})

	loop.Finish()
}

// for of 取值
func (b *astbuilder) buildForOfStatement(stmt *JS.ForOfStatementContext) {
	// todo: handle await

	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	loop := b.BuildLoop()

	loop.BuildCondition(func() ssa.Value {
		var left ssa.LeftValue
		var value ssa.Value

		if s, ok := stmt.VariableDeclaration().(*JS.VariableDeclarationContext); ok {
			_, left = b.buildVariableDeclaration(s, "v")
			value, _ = b.buildSingleExpression(stmt.SingleExpression(0), false)
		} else {
			_, left = b.buildSingleExpression(stmt.SingleExpression(0), true)
			value, _ = b.buildSingleExpression(stmt.SingleExpression(1), false)
		}

		_, field, ok := b.EmitNext(value, true)
		left.Assign(field, b.FunctionBuilder)

		return ok
	})

	loop.BuildBody(func() {
		if s, ok := stmt.Statement().(*JS.StatementContext); ok {
			b.buildStatement(s)
		}
	})

	loop.Finish()
}

func (b *astbuilder) buildFunctionDeclaration(stmt *JS.FunctionDeclarationContext) ssa.Value {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	funcName := ""
	if Name := stmt.Identifier(); Name != nil {
		funcName = Name.GetText()
	}

	// fmt.Println("funcName: ", funcName)

	newFunc, symbol := b.NewFunc(funcName)
	current := b.CurrentBlock
	buildFunc := func() {
		b.FunctionBuilder = b.PushFunction(newFunc, symbol, current)

		if s, ok := stmt.FormalParameterList().(*JS.FormalParameterListContext); ok {
			b.buildFormalParameterList(s)
		}

		if f, ok := stmt.FunctionBody().(*JS.FunctionBodyContext); ok {
			b.buildFunctionBody(f)
		}

		b.Finish()
		b.FunctionBuilder = b.PopFunction()

	}

	b.AddSubFunction(buildFunc)

	if funcName != "" {
		b.WriteVariable(funcName, newFunc)
	}
	return newFunc
}

func (b *astbuilder) buildFunctionBody(stmt *JS.FunctionBodyContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s, ok := stmt.SourceElements().(*JS.SourceElementsContext); ok {
		b.buildSourceElements(s)
		return
	}
}

func (b *astbuilder) buildSourceElements(stmt *JS.SourceElementsContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s := stmt.AllSourceElement(); s != nil {
		for _, i := range s {
			b.buildSourceElement(i)
		}
	}
}

func (b *astbuilder) buildSourceElement(stmt JS.ISourceElementContext) {
	if s, ok := stmt.Statement().(*JS.StatementContext); ok {
		b.buildStatement(s)
		return
	}
}

func (b *astbuilder) buildFormalParameterList(stmt *JS.FormalParameterListContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if f := stmt.AllFormalParameterArg(); f != nil {
		for _, i := range f {
			if a, ok := i.(*JS.FormalParameterArgContext); ok {
				b.buildFormalParameterArg(a)
			}
		}

		if l, ok := stmt.LastFormalParameterArg().(*JS.LastFormalParameterArgContext); ok {
			b.buildLastFormalParameterArg(l)
		}
		return
	}

	if l, ok := stmt.LastFormalParameterArg().(*JS.LastFormalParameterArgContext); ok {
		b.buildLastFormalParameterArg(l)
		return
	}

	b.NewError(ssa.Error, TAG, ArrowFunctionNeedExpressionOrBlock())
}

func (b *astbuilder) buildFormalParameterArg(stmt *JS.FormalParameterArgContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	a := stmt.Assign()

	if a == nil {
		b.NewParam(stmt.GetText())
	} else {
		p := b.NewParam(stmt.Assignable().GetText())

		x := stmt.SingleExpression()
		result, _ := b.buildSingleExpression(x, false)

		p.SetDefault(result)
		return
	}
}

func (b *astbuilder) buildLastFormalParameterArg(stmt *JS.LastFormalParameterArgContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if e := stmt.Ellipsis(); e != nil {
		b.HandlerEllipsis()
	}

	if s := stmt.SingleExpression(); s != nil {
		b.buildSingleExpression(s, false)
	}
}

func (b *astbuilder) buildReturnStatement(stmt *JS.ReturnStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()
	if s, ok := stmt.ExpressionSequence().(*JS.ExpressionSequenceContext); ok {
		values := b.buildExpressionSequence(s)
		b.EmitReturn(values)
	} else {
		b.EmitReturn(nil)
	}
}

func (b *astbuilder) buildImportStatement(stmt *JS.ImportStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	if s := stmt.ImportFromBlock(); s != nil {
		if str := s.StringLiteral(); s != nil {
			b.buildStringLiteral(str)
		}
	}
}

func (b *astbuilder) buildBreakStatement(stmt *JS.BreakStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()

	var _break *ssa.BasicBlock

	if s, ok := stmt.Identifier().(*JS.IdentifierContext); ok {
		text := s.GetText()
		if _break = b.GetLabel(text); _break != nil {
			b.EmitJump(_break)
		} else {
			b.NewError(ssa.Error, TAG, UndefineLabelstmt())
		}
		return

		// fmt.Println("want break to :", text)
	} else {
		if _break = b.GetBreak(); _break != nil {
			b.EmitJump(_break)
		} else {
			b.NewError(ssa.Error, TAG, UnexpectedBreakStmt())
		}
		return
	}

}

func (b *astbuilder) buildLabelledStatement(stmt *JS.LabelledStatementContext) {
	recoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer recoverRange()
	text := ""
	if s, ok := stmt.Identifier().(*JS.IdentifierContext); ok {
		text = s.GetText()
	}

	// unsealed block
	block := b.NewBasicBlockUnSealed(text)
	b.AddLabel(text, block)
	// to block
	b.EmitJump(block)
	b.CurrentBlock = block
	if s, ok := stmt.Statement().(*JS.StatementContext); ok {
		b.buildStatement(s)
	}
	b.DeleteLabel(text)
	block.Sealed()
	// // to done
	// done := b.NewBasicBlock("done")
	// b.EmitJump(done)
	// b.CurrentBlock = done

}

func (b *astbuilder) buildTryStatement(stmt *JS.TryStatementContext) {
	revcoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer revcoverRange()

	try := b.BuildTry()

	try.BuildTryBlock(func() {
		if s, ok := stmt.Block().(*JS.BlockContext); ok {
			b.buildBlock(s)
		}
	})

	try.BuildCatch(func() string {
		var id string
		// TODO: Assignable could be wrong, need to fix
		if s, ok := stmt.CatchProduction().(*JS.CatchProductionContext); ok {
			if a, ok := s.Assignable().(*JS.AssignableContext); ok {
				b.buildAssignableContext(a)
				id = a.GetText()
			}

			if bl, ok := stmt.Block().(*JS.BlockContext); ok {
				b.buildBlock(bl)
			}
		}
		return id
	})

	if s, ok := stmt.FinallyProduction().(*JS.FinallyProductionContext); ok {

		try.BuildFinally(func() {
			if bl, ok := s.Block().(*JS.BlockContext); ok {
				b.buildBlock(bl)
			}
		})
	}

	try.Finish()

}

func (b *astbuilder) buildSwitchStatement(stmt *JS.SwitchStatementContext) {
	revcoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer revcoverRange()

	Switchb := b.BuildSwitch()
	Switchb.DefaultBreak = false

	if s := stmt.SingleExpression(); s != nil {
		Switchb.BuildCondition(func() ssa.Value {
			rv, _ := b.buildSingleExpression(s, false)
			return rv
		})
	} else {
		b.NewError(ssa.Warn, TAG, "switch expression is nil")
	}

	if s, ok := stmt.CaseBlock().(*JS.CaseBlockContext); ok {
		b.buildCaseBlock(s, Switchb)
	}
}

func (b *astbuilder) buildCaseBlock(stmt *JS.CaseBlockContext, Switchb *ssa.SwitchBuilder) {
	revcoverRange := b.SetRange(&stmt.BaseParserRuleContext)
	defer revcoverRange()

	var stList []*JS.StatementListContext

	var caseNum int
	var exprs []ssa.Value
	caseNum = 0
	for _, s := range stmt.AllCaseClauses() {
		if cs, ok := s.(*JS.CaseClausesContext); ok {
			for _, i := range cs.AllCaseClause() {
				if c, ok := i.(*JS.CaseClauseContext); ok {
					rv, _ := b.buildSingleExpression(c.SingleExpression(), false)
					exprs = append(exprs, rv)

					if st, ok := c.StatementList().(*JS.StatementListContext); ok {
						stList = append(stList, st)
					} else {
						stList = append(stList, nil)
					}

					caseNum += 1
				}
			}
		}
	}

	Switchb.BuildHanlder(func() (int, []ssa.Value) {
		return caseNum, exprs
	})

	Switchb.BuildBody(func(i int) {
		if stList[i] != nil {
			b.buildStatementList(stList[i])
		}
	})

	if s, ok := stmt.DefaultClause().(*JS.DefaultClauseContext); ok {
		if st, ok := s.StatementList().(*JS.StatementListContext); ok {
			Switchb.BuildDefault(func() {
				b.buildStatementList(st)
			})
		}
	}

	Switchb.Finsh()

}