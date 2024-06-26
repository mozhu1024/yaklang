package sfvm

import (
	"reflect"
	"regexp"

	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/syntaxflow/sf"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
)

func (y *SyntaxFlowVisitor) VisitFilterExpr(raw sf.IFilterExprContext) error {
	if y == nil || raw == nil {
		return nil
	}
	i, ok := raw.(*sf.FilterExprContext)
	if !ok {
		return utils.Errorf("BUG: in filterExpr: %s", reflect.TypeOf(raw))
	}

	enter := y.EmitFilterExprEnter()
	defer func() {
		y.EmitFilterExprExit(enter)
	}()
	if raw := i.FilterItemFirst(); raw != nil {
		y.VisitFilterItemFirst(raw)
	}

	for _, raw := range i.AllFilterItem() {
		y.VisitFilterItem(raw)
	}
	return nil
}

func (y *SyntaxFlowVisitor) VisitFilterItem(raw sf.IFilterItemContext) error {
	switch filter := raw.(type) {
	case *sf.FirstContext:
		y.VisitFilterItemFirst(filter.FilterItemFirst())
	case *sf.FunctionCallFilterContext:
		if filter.ActualParam() != nil {
			y.VisitActualParam(filter.ActualParam())
		} else {
			y.EmitGetCall()
		}
	case *sf.FieldIndexFilterContext:
		memberRaw := filter.SliceCallItem()
		member, ok := memberRaw.(*sf.SliceCallItemContext)
		if !ok {
			panic("BUG: in fieldIndexFilter")
		}
		if member.NumberLiteral() != nil {
			y.EmitListIndex(codec.Atoi(member.NumberLiteral().GetText()))
		} else {
			y.VisitNameFilter(true, member.NameFilter())
		}
	case *sf.OptionalFilterContext:
		y.VisitConditionExpression(filter.ConditionExpression())
		y.EmitCondition()
	case *sf.NextFilterContext:
		y.EmitGetUsers()
	case *sf.DefFilterContext:
		y.EmitGetDefs()
	case *sf.DeepNextFilterContext:
		y.EmitGetBottomUsers()
	case *sf.DeepNextConfigFilterContext:
		config := []*RecursiveConfigItem{}
		if i := filter.RecursiveConfig(); i != nil {
			config = y.VisitRecursiveConfig(i.(*sf.RecursiveConfigContext))
		}
		y.EmitGetBottomUsers(config...)
	case *sf.TopDefFilterContext:
		y.EmitGetTopDefs()
	case *sf.TopDefConfigFilterContext:
		config := []*RecursiveConfigItem{}
		if i := filter.RecursiveConfig(); i != nil {
			config = y.VisitRecursiveConfig(i.(*sf.RecursiveConfigContext))
		}
		y.EmitGetTopDefs(config...)
	case *sf.UseDefCalcFilterContext:
		log.Warnf("TBD: UseDefCalcFilterContext: %v", raw.GetText())
	default:
		panic("BUG: in filterExpr")
	}
	return nil
}

func (y *SyntaxFlowVisitor) VisitFilterItemFirst(raw sf.IFilterItemFirstContext) error {

	if y == nil || raw == nil {
		return nil
	}
	switch i := raw.(type) {
	case *sf.NamedFilterContext:
		y.VisitNameFilter(false, i.NameFilter())
	case *sf.FieldCallFilterContext:
		y.VisitNameFilter(true, i.NameFilter())
	default:
		panic("BUG: in filter first")
	}

	return nil
}

func (y *SyntaxFlowVisitor) VisitNameFilter(isMember bool, i sf.INameFilterContext) error {
	if i == nil {
		return nil
	}

	ret, ok := i.(*sf.NameFilterContext)
	if !ok {
		return utils.Errorf("BUG: in nameFilter: %s", reflect.TypeOf(i))
	}

	mod := NameMatch
	if isMember {
		mod = KeyMatch
	}

	if s := ret.Star(); s != nil {
		if isMember {
			// get all member
			y.EmitSearchGlob(mod, "*")
		}
		// skip
		return nil
		// } else if id := ret.DollarOutput(); id != nil {
		// 	y.EmitSearchExact(mod, id.GetText())
		// 	return nil
	} else if id := ret.Identifier(); id != nil {
		text := ret.Identifier().GetText()
		filter, isGlob := y.FormatStringOrGlob(text) // emit field
		if isGlob {
			y.EmitSearchGlob(mod, filter)
		} else {
			y.EmitSearchExact(mod, filter)
		}
		return nil
	} else if re, ok := ret.RegexpLiteral().(*sf.RegexpLiteralContext); ok {
		text := re.RegexpLiteral().GetText()
		text = text[1 : len(text)-1]
		// log.Infof("regexp: %s", text)
		reIns, err := regexp.Compile(text)
		if err != nil {
			return err
		}
		y.EmitSearchRegexp(mod, reIns.String())
		return nil
	}
	return utils.Errorf("BUG: in nameFilter, unknown type: %s:%s", reflect.TypeOf(ret), ret.GetText())
}

func (y *SyntaxFlowVisitor) VisitActualParam(i sf.IActualParamContext) error {
	handlerStatement := func(i sf.ISingleParamContext) {
		ret, ok := (i).(*sf.SingleParamContext)
		if !ok {
			return
		}

		if ret.FilterStatement() != nil {
			y.VisitFilterStatement(ret.FilterStatement())
		}
		// TODO: handler recursive config
	}

	switch ret := i.(type) {
	case *sf.AllParamContext:
		y.EmitEnterStatement()
		y.EmitPushAllCallArgs()
		handlerStatement(ret.SingleParam())
		y.EmitExitStatement()
	case *sf.EveryParamContext:
		for i, paraI := range ret.AllActualParamFilter() {
			para, ok := paraI.(*sf.ActualParamFilterContext)
			if !ok {
				continue
			}
			single := para.SingleParam()
			if single == nil {
				continue
			}
			y.EmitEnterStatement()
			y.EmitPushCallArgs(i)
			handlerStatement(single)
			y.EmitExitStatement()
		}
		if ret.SingleParam() != nil {
			y.EmitEnterStatement()
			y.EmitPushCallArgs(len(ret.AllActualParamFilter()))
			handlerStatement(ret.SingleParam())
			y.EmitExitStatement()
		}
	default:
		return utils.Errorf("BUG: ActualParamFilter type error: %s", reflect.TypeOf(ret))
	}
	return nil
}
