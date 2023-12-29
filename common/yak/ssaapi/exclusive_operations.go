package ssaapi

import (
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/yak/ssa"
)

// GetContextValue can handle context
func (v *Value) GetContextValue(i string) (*Value, bool) {
	return v.runtimeCtx.Get(i)
}

func (v *Value) SetContextValue(i string, values *Value) *Value {
	v.runtimeCtx.Set(i, values)
	return v
}

func (v *Value) GetParent() (*Value, bool) {
	return v.GetContextValue("parent")
}

func (v *Value) SetParent(value *Value) *Value {
	v.runtimeCtx.Set("parent", value)
	return v
}

// GetTopDefs desc all of 'Defs' is not used by any other value
func (i *Value) GetTopDefs() Values {
	return i.getTopDefs(nil)
}

func (v Values) GetTopDefs() Values {
	ret := make(Values, 0, len(v))
	var m = make(map[*Value]struct{})
	v.WalkDefs(func(i *Value) {
		if !i.HasOperands() {
			if _, ok := m[i]; ok {
				return
			}
			m[i] = struct{}{}
			ret = append(ret, i)
		}
	})
	return ret
}

func (i *Value) visitedDefsDefault(actx *AnalyzeContext) Values {
	var vals Values
	for _, def := range i.node.GetValues() {
		if ret := NewValue(def).SetParent(i).getTopDefs(actx); len(ret) > 0 {
			vals = append(vals, ret...)
		}
	}

	if len(vals) <= 0 {
		vals = append(vals, i)
	}
	if maskable, ok := i.node.(ssa.Maskable); ok {
		for _, def := range maskable.GetMask() {
			if ret := NewValue(def).SetParent(i).getTopDefs(actx); len(ret) > 0 {
				vals = append(vals, ret...)
			}
		}
	}
	return vals
}

func (i *Value) getTopDefs(actx *AnalyzeContext) Values {
	if i == nil {
		return nil
	}

	if actx == nil {
		actx = NewAnalyzeContext()
	}

	switch ret := i.node.(type) {
	case *ssa.Phi:
		if !actx.ThePhiShouldBeVisited(i) {
			// phi is visited...
			return Values{}
		}
		actx.VisitPhi(i)
		return i.visitedDefsDefault(actx)
	case *ssa.Call:
		caller := ret.Method
		if caller == nil {
			return Values{i} // return self
		}

		err := actx.PushCall(i)
		if err != nil {
			log.Errorf("push call error: %v", err)
			return Values{i}
		}

		// TODO: trace the specific return-values
		return NewValue(caller).SetParent(i).getTopDefs(actx)
	case *ssa.Function:
		log.Info("ssa.Function checking...")
		var vals Values
		// handle return
		for _, r := range ret.Return {
			for _, subVal := range r.GetValues() {
				if ret := NewValue(subVal).SetParent(i).getTopDefs(actx); len(ret) > 0 {
					vals = append(vals, ret...)
				}
			}
		}
		if len(vals) == 0 {
			return Values{i} // no return, use undefined
		}
		return vals
	case *ssa.Parameter:
		log.Infof("checking ssa.Parameters...: %v", ret.String())
		parent, ok := i.GetParent()
		if !ok {
			log.Warn("topdefs parameter context error, skip")
			return Values{i}
		}
		if parent.IsFunction() {
			called, ok := parent.GetParent()
			if !ok {
				log.Infof("parent function is not called by any other function, skip")
				return Values{i}
			}
			if !called.IsCall() {
				log.Infof("parent function is not called by any other function, skip (%T)", called)
				return Values{i}
			}
			var vals Values
			calledInstance := called.node.(*ssa.Call)
			for _, i := range calledInstance.Args {
				traced := NewValue(i).SetParent(called)
				if ret := traced.getTopDefs(actx); len(ret) > 0 {
					vals = append(vals, ret...)
				} else {
					vals = append(vals, traced)
				}
			}
			if len(vals) == 0 {
				return Values{NewValue(ssa.NewUndefined("_")).SetParent(i)} // no return, use undefined
			}
			return vals
		} else if parent != i {
			var vals Values
			if ret.IsFreeValue {
				// free value
				// fetch parent
				fun := ret.GetFunc().GetParent() // func.parent
				for _, va := range fun.GetValuesByName(ret.GetName()) {
					_, isSideEffect := va.(*ssa.SideEffect)
					if isSideEffect {
						continue
					}

					if ret := NewValue(va).SetParent(i).getTopDefs(actx); len(ret) > 0 {
						vals = append(vals, ret...)
					}
				}
			}
			if len(vals) <= 0 {
				return Values{i}
			}
			return vals
		}
		return Values{i}
	}
	return i.visitedDefsDefault(actx)
}