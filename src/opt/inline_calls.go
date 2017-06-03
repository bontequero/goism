package opt

import (
	"sexp"
)

// InlineCalls inlines suitable functions calls inside form.
func InlineCalls(form sexp.Form) sexp.Form {
	return sexp.Rewrite(form, func(form sexp.Form) sexp.Form {
		if form, ok := form.(*sexp.Call); ok {
			return inlineCall(form.Fn, form.Args)
		}
		return nil
	})
}

func inlineCall(fn *sexp.Func, args []sexp.Form) sexp.Form {
	body := fn.Body
	if !inlineable(body) {
		return nil
	}

	expr := body.Forms[0].(*sexp.Return).Results[0]
	bindings := make([]*sexp.Bind, 0, len(fn.Params))
	for i, param := range fn.Params {
		if sexp.Cost(args[i]) == 1 {
			injectValue(param, args[i], expr)
		} else {
			bindings = append(bindings, &sexp.Bind{
				Name: param,
				Init: args[i],
			})
		}
	}

	if len(bindings) == 0 {
		// Perfect inlining, no let wrapper is needed.
		return expr
	}
	return &sexp.Let{
		Bindings: bindings,
		Expr:     expr,
	}
}

func inlineable(body *sexp.Block) bool {
	if len(body.Forms) != 1 {
		return false
	}
	form, ok := body.Forms[0].(*sexp.Return)
	if !ok {
		return false
	}
	return len(form.Results) == 1 && sexp.Cost(form) <= 20
}
