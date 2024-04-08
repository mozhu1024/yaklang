package ssa

import (
	"github.com/yaklang/yaklang/common/utils"
)

func NewCall(target Value, args []Value, binding map[string]Value, block *BasicBlock) *Call {
	// handler "this" in parameter
	{
		AddThis := func(this Value) {
			if len(args) == 0 {
				args = append(args, this)
			} else {
				args = utils.InsertSliceItem(args, this, 0)
			}
		}
		switch t := target.(type) {
		case *ClassMethod: // for class instance
			AddThis(t.This)
			target = t.Func
		case *Function: // for object inner function
			if len(t.Param) > 0 {
				if para := t.Param[0]; para != nil && para.IsObject() && para.GetDefault() != nil {
					AddThis(para.GetDefault())
				}
			}
		default:
		}
	}

	if binding == nil {
		binding = make(map[string]Value)
	}

	c := &Call{
		anValue:     NewValue(),
		Method:      target,
		Args:        args,
		Binding:     binding,
		Async:       false,
		Unpack:      false,
		IsDropError: false,
		IsEllipsis:  false,
	}
	return c
}

func (f *FunctionBuilder) NewCall(target Value, args []Value) *Call {
	call := NewCall(target, args, nil, f.CurrentBlock)
	return call
}

func (f *FunctionBuilder) EmitCall(c *Call) *Call {
	if f.CurrentBlock.finish {
		return nil
	}

	f.emit(c)
	c.handlerReturnType()
	c.handleCalleeFunction()

	return c
}

// handler Return type, and handle drop error
func (c *Call) handlerReturnType() {
	// get function type
	funcTyp, ok := ToFunctionType(c.Method.GetType())
	if !ok {
		return
	}
	// inference call instruction type
	if c.IsDropError {
		if retType, ok := funcTyp.ReturnType.(*ObjectType); ok {
			if retType.Combination && retType.FieldTypes[len(retType.FieldTypes)-1].GetTypeKind() == ErrorTypeKind {
				if len(retType.FieldTypes) == 1 {
					c.SetType(BasicTypes[NullTypeKind])
				} else if len(retType.FieldTypes) == 2 {
					// if len(t.FieldTypes) == 2 {
					c.SetType(retType.FieldTypes[0])
				} else {
					ret := NewStructType()
					ret.FieldTypes = retType.FieldTypes[:len(retType.FieldTypes)-1]
					ret.Keys = retType.Keys[:len(retType.Keys)-1]
					ret.KeyTyp = retType.KeyTyp
					ret.Combination = true
					ret.Len = len(ret.FieldTypes)
					ret.Kind = TupleTypeKind
					c.SetType(ret)
				}
				return
			}
		} else if t, ok := funcTyp.ReturnType.(*BasicType); ok && t.Kind == ErrorTypeKind {
			// pass
			c.SetType(BasicTypes[NullTypeKind])
			for _, variable := range c.GetAllVariables() {
				variable.NewError(Error, SSATAG, ValueIsNull())
			}
			return
		}
		c.NewError(Warn, SSATAG, FunctionContReturnError())
	} else {
		c.SetType(funcTyp.ReturnType)
	}
}

// handler if method, set object for first argument
func (c *Call) handleCalleeFunction() {

	// get function type
	funcTyp, ok := ToFunctionType(c.Method.GetType())
	if !ok {
		return
	}

	{
		builder := c.GetFunc().builder
		recoverBuilder := builder.SetCurrent(c)
		currentScope := c.GetBlock().ScopeTable

		for true {
			if len(c.Args) == len(funcTyp.ParameterValue) {
				break
			}
			for _, p := range funcTyp.ParameterValue {
				//TODO:  i don't know why this condition, it should be clearer
				if len(c.Args) == p.FormalParameterIndex {
					// free-value member-call will be set to parameter,
					// this parameter will contain default value, use this.
					if p.GetDefault() != nil {
						c.Args = append(c.Args, p.GetDefault())
					}
				}
				if !p.IsMemberCall {
					continue
				}
				if p.MemberCallObjectIndex >= len(c.Args) {
					// log.Errorf("handleCalleeFunction: memberCallObjectIndex out of range %d vs len: %d", p.MemberCallObjectIndex, len(c.Args))
					continue
				}

				if _, typ := checkCanMemberCall(c.Args[p.MemberCallObjectIndex], p.MemberCallKey); typ == nil {
					builder.NewErrorWithPos(Error, SSATAG,
						p.GetRange(),
						FreeValueNotMember(
							c.Args[p.MemberCallObjectIndex].GetName(),
							p.MemberCallKey.String(),
							c.GetRange(),
						),
					)
					c.NewError(Error, SSATAG,
						FreeValueNotMemberInCall(
							c.Args[p.MemberCallObjectIndex].GetName(),
							p.MemberCallKey.String(),
						),
					)
					continue
				}
				c.Args = append(c.Args,
					builder.ReadMemberCallVariable(c.Args[p.MemberCallObjectIndex], p.MemberCallKey),
				)
			}
			break
		}

		// handle side effect
		for _, se := range funcTyp.SideEffects {
			var variable *Variable
			if se.IsMemberCall {
				if se.ParameterIndex >= len(c.Args) {
					// log.Errorf("handleCalleeFunction: ParameterIndex out of range %d", se.ParameterIndex)
					continue
				}
				// if side-effect is member call, create member call variable
				variable = builder.CreateMemberCallVariable(c.Args[se.ParameterIndex], se.Key)
			} else {

				// side-effect only create in scope that lower or same than modify's scope
				if !se.forceCreate && !currentScope.IsSameOrSubScope(se.Variable.GetScope()) {
					continue
				}
				variable = builder.CreateVariable(se.Name)
			}

			// TODO: handle side effect in loop scope,
			// will replace value in scope and create new phi
			sideEffect := builder.EmitSideEffect(se.Name, c, se.Modify)
			if sideEffect != nil {
				builder.AssignVariable(variable, sideEffect)
				sideEffect.SetVerboseName(se.VerboseName)
			}
		}
		recoverBuilder()
	}

	// only handler in method call
	if !funcTyp.IsMethod {
		return
	}

	is := c.Method.IsMember()
	if !is {
		// this function is method Function, but no member call get this.
		// error
		return
	}

}

func (c *Call) HandleFreeValue(fvs []*Parameter) {
	builder := c.GetFunc().builder
	recoverBuilder := builder.SetCurrent(c)
	defer recoverBuilder()

	for _, fv := range fvs {
		// if freeValue has default value, skip
		if fv.GetDefault() != nil {
			continue
		}

		v := builder.PeekValue(fv.GetName())

		if v != nil {
			c.Binding[fv.GetName()] = v
		} else {
			// mark error in freeValue.Variable
			// get freeValue
			if variable := fv.GetVariable(fv.GetName()); variable != nil {
				variable.NewError(Error, SSATAG, BindingNotFound(fv.GetName(), c.GetRange()))
			}
			// skip instance function, or `go` with instance function,
			// this function no variable, and code-range of call-site same as function.
			// we don't mark error in call-site.
			if fun, ok := ToFunction(c.Method); ok {
				if len(fun.GetAllVariables()) == 0 {
					continue
				}
			}
			// other code will mark error in function call-site
			c.NewError(Error, SSATAG, BindingNotFoundInCall(fv.GetName()))
		}
	}
}
