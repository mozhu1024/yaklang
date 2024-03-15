package ssa

import (
	"github.com/samber/lo"
	"github.com/yaklang/yaklang/common/utils"
)

func (r *Return) calcType() Type {
	handleType := func(v Value) Type {
		t := v.GetType()
		if objTyp, ok := ToObjectType(t); ok {
			t = ParseClassBluePrint(v, objTyp)
		}
		return t
	}

	switch len(r.Results) {
	case 0:
		return BasicTypes[NullTypeKind]
	case 1:
		return handleType(r.Results[0])
	default:
		newObjTyp := NewObjectType()
		for i, v := range r.Results {
			newObjTyp.AddField(NewConst(i), handleType(v))
		}
		newObjTyp.Finish()
		newObjTyp.Kind = TupleTypeKind
		newObjTyp.Len = len(r.Results)
		return newObjTyp
	}
}

// Finish current function builder
func (b *FunctionBuilder) Finish() {
	// set defer function
	if deferLen := len(b.deferExpr); deferLen > 0 {
		endBlock := b.CurrentBlock

		deferBlock := b.GetDeferBlock()
		b.CurrentBlock = deferBlock
		for _, i := range b.deferExpr {
			if len(deferBlock.Insts) == 0 {
				deferBlock.Insts = append(deferBlock.Insts, i)
			} else {
				deferBlock.Insts = utils.InsertSliceItem(deferBlock.Insts, Instruction(i), 0)
			}
		}
		b.deferExpr = []*Call{}

		b.CurrentBlock = endBlock
	}
	// function finish
	b.Function.Finish()
}

// calculate all return instruction in function, get return type
func handlerReturnType(rs []*Return) Type {
	tmp := make(map[string]Type, len(rs))
	for _, r := range rs {
		id := ""
		typs := r.calcType()

		if _, ok := tmp[id]; !ok {
			tmp[id] = typs
		}
	}

	typs := lo.Values(tmp)
	if len(typs) == 0 {
		return BasicTypes[NullTypeKind]
	} else if len(typs) == 1 {
		return typs[0]
	} else {
		//TODO: how handler this? multiple return with different type
		// should set Warn!!
		// and ?? Type ??
		return nil
	}
}

// Finish the function, set FunctionType, set EnterBlock/ExitBlock
func (f *Function) Finish() {
	f.EnterBlock = f.Blocks[0]
	f.ExitBlock = f.Blocks[len(f.Blocks)-1]

	funType := NewFunctionType("",
		lo.Map(f.Param, func(p *Parameter, _ int) Type {
			t := p.GetType()
			return t
		}),
		handlerReturnType(f.Return),
		f.hasEllipsis,
	)
	funType.ParameterLen = f.ParamLength
	funType.ParameterValue = f.Param
	funType.SetFreeValue(f.FreeValues)
	funType.SetSideEffect(f.SideEffects)
	f.SetType(funType)
}