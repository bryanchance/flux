package compiler

import (
	"fmt"
	"regexp"

	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/values"
	"github.com/pkg/errors"
)

type Func interface {
	Type() semantic.Type
	Eval(input values.Object) (values.Value, error)
	EvalString(input values.Object) (string, error)
	EvalInt(input values.Object) (int64, error)
	EvalUInt(input values.Object) (uint64, error)
	EvalFloat(input values.Object) (float64, error)
	EvalBool(input values.Object) (bool, error)
	EvalTime(input values.Object) (values.Time, error)
	EvalDuration(input values.Object) (values.Duration, error)
	EvalRegexp(input values.Object) (*regexp.Regexp, error)
	EvalArray(input values.Object) (values.Array, error)
	EvalObject(input values.Object) (values.Object, error)
	EvalFunction(input values.Object) (values.Function, error)
}

type Evaluator interface {
	Type() semantic.Type
	EvalString(scope Scope) (values.Value, error)
	EvalInt(scope Scope) (values.Value, error)
	EvalUInt(scope Scope) (values.Value, error)
	EvalFloat(scope Scope) (values.Value, error)
	EvalBool(scope Scope) (values.Value, error)
	EvalTime(scope Scope) (values.Value, error)
	EvalDuration(scope Scope) (values.Duration, error)
	EvalRegexp(scope Scope) (*regexp.Regexp, error)
	EvalArray(scope Scope) (values.Array, error)
	EvalObject(scope Scope) (values.Object, error)
	EvalFunction(scope Scope) (values.Function, error)
}

type ValueEvaluator interface {
	EvalValue(scope Scope) (values.Value, error)
}

type compiledFn struct {
	root       Evaluator
	fnType     semantic.Type
	inputScope Scope
}

func (c compiledFn) validate(input values.Object) error {
	sig := c.fnType.FunctionSignature()
	properties := input.Type().Properties()
	if len(properties) != len(sig.Parameters) {
		return errors.New("mismatched parameters and properties")
	}
	for k, v := range sig.Parameters {
		if properties[k] != v {
			return fmt.Errorf("parameter %q has the wrong type, expected %v got %v", k, v, properties[k])
		}
	}
	return nil
}

func (c compiledFn) buildScope(input values.Object) error {
	if err := c.validate(input); err != nil {
		return err
	}
	input.Range(func(k string, v values.Value) {
		c.inputScope[k] = v
	})
	return nil
}

func (c compiledFn) Type() semantic.Type {
	return c.fnType.FunctionSignature().Return
}

func (c compiledFn) Eval(input values.Object) (values.Value, error) {
	if err := c.buildScope(input); err != nil {
		return nil, err
	}

	return eval(c.root, c.inputScope)
}

func (c compiledFn) EvalString(input values.Object) (string, error) {
	if err := c.buildScope(input); err != nil {
		return "", err
	}
	v, err := c.root.EvalString(c.inputScope)
	return v.Str(), err
}
func (c compiledFn) EvalBool(input values.Object) (bool, error) {
	if err := c.buildScope(input); err != nil {
		return false, err
	}
	v, err := c.root.EvalBool(c.inputScope)
	return v.Bool(), err
}
func (c compiledFn) EvalInt(input values.Object) (int64, error) {
	if err := c.buildScope(input); err != nil {
		return 0, err
	}
	v, err := c.root.EvalInt(c.inputScope)
	return v.Int(), err
}
func (c compiledFn) EvalUInt(input values.Object) (uint64, error) {
	if err := c.buildScope(input); err != nil {
		return 0, err
	}
	v, err := c.root.EvalUInt(c.inputScope)
	return v.UInt(), err
}
func (c compiledFn) EvalFloat(input values.Object) (float64, error) {
	if err := c.buildScope(input); err != nil {
		return 0, err
	}
	v, err := c.root.EvalFloat(c.inputScope)
	return v.Float(), err
}
func (c compiledFn) EvalTime(input values.Object) (values.Time, error) {
	if err := c.buildScope(input); err != nil {
		return 0, err
	}
	v, err := c.root.EvalTime(c.inputScope)
	return v.Time(), err
}
func (c compiledFn) EvalDuration(input values.Object) (values.Duration, error) {
	if err := c.buildScope(input); err != nil {
		return 0, err
	}
	return c.root.EvalDuration(c.inputScope)
}
func (c compiledFn) EvalRegexp(input values.Object) (*regexp.Regexp, error) {
	if err := c.buildScope(input); err != nil {
		return nil, err
	}
	return c.root.EvalRegexp(c.inputScope)
}
func (c compiledFn) EvalArray(input values.Object) (values.Array, error) {
	if err := c.buildScope(input); err != nil {
		return nil, err
	}
	return c.root.EvalArray(c.inputScope)
}
func (c compiledFn) EvalObject(input values.Object) (values.Object, error) {
	if err := c.buildScope(input); err != nil {
		return nil, err
	}
	return c.root.EvalObject(c.inputScope)
}
func (c compiledFn) EvalFunction(input values.Object) (values.Function, error) {
	if err := c.buildScope(input); err != nil {
		return nil, err
	}
	return c.root.EvalFunction(c.inputScope)
}

type Scope map[string]values.Value

func (s Scope) Type(name string) semantic.Type {
	return s[name].Type()
}
func (s Scope) Set(name string, v values.Value) {
	s[name] = v
}

func (s Scope) GetString(name string) string {
	return s[name].Str()
}
func (s Scope) GetInt(name string) int64 {
	return s[name].Int()
}
func (s Scope) GetUInt(name string) uint64 {
	return s[name].UInt()
}
func (s Scope) GetFloat(name string) float64 {
	return s[name].Float()
}
func (s Scope) GetBool(name string) bool {
	return s[name].Bool()
}
func (s Scope) GetTime(name string) values.Time {
	return s[name].Time()
}
func (s Scope) GetDuration(name string) values.Duration {
	return s[name].Duration()
}
func (s Scope) GetRegexp(name string) *regexp.Regexp {
	return s[name].Regexp()
}
func (s Scope) GetArray(name string) values.Array {
	return s[name].Array()
}
func (s Scope) GetObject(name string) values.Object {
	return s[name].Object()
}
func (s Scope) GetFunction(name string) values.Function {
	return s[name].Function()
}

func (s Scope) Copy() Scope {
	n := make(Scope, len(s))
	for k, v := range s {
		n[k] = v
	}
	return n
}

func eval(e Evaluator, scope Scope) (values.Value, error) {
	var v values.Value
	var err error
	switch e.Type().Nature() {
	case semantic.String:
		var v0 values.Value
		v0, err = e.EvalString(scope)
		if err == nil {
			v = v0
		}
	case semantic.Int:
		var v0 values.Value
		v0, err = e.EvalInt(scope)
		if err == nil {
			v = v0
		}
	case semantic.UInt:
		var v0 values.Value
		v0, err = e.EvalUInt(scope)
		if err == nil {
			v = v0
		}
	case semantic.Float:
		var v0 values.Value
		v0, err = e.EvalFloat(scope)
		if err == nil {
			v = v0
		}
	case semantic.Bool:
		var v0 values.Value
		v0, err = e.EvalBool(scope)
		if err == nil {
			v = v0
		}
	case semantic.Time:
		var v0 values.Value
		v0, err = e.EvalTime(scope)
		if err == nil {
			v = v0
		}
	case semantic.Duration:
		var v0 values.Duration
		v0, err = e.EvalDuration(scope)
		if err == nil {
			v = values.NewDuration(v0)
		}
	case semantic.Regexp:
		var v0 *regexp.Regexp
		v0, err = e.EvalRegexp(scope)
		if err == nil {
			v = values.NewRegexp(v0)
		}
	case semantic.Array:
		v, err = e.EvalArray(scope)
	case semantic.Object:
		v, err = e.EvalObject(scope)
	case semantic.Function:
		v, err = e.EvalFunction(scope)
	case semantic.Nil:
		return nil, nil
	default:
		err = fmt.Errorf("eval: unknown type: %v", e.Type())
	}

	return v, err
}

type blockEvaluator struct {
	t     semantic.Type
	body  []Evaluator
	value values.Value
}

func (e *blockEvaluator) Type() semantic.Type {
	return e.t
}

func (e *blockEvaluator) eval(scope Scope) error {
	var err error
	for _, b := range e.body {
		e.value, err = eval(b, scope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *blockEvaluator) EvalString(scope Scope) (values.Value, error) {
	values.CheckKind(e.t.Nature(), semantic.String)
	err := e.eval(scope)
	if err != nil {
		return values.NewString(""), err
	}
	return e.value, nil
}
func (e *blockEvaluator) EvalInt(scope Scope) (values.Value, error) {
	values.CheckKind(e.t.Nature(), semantic.Int)
	err := e.eval(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	return e.value, nil
}
func (e *blockEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	values.CheckKind(e.t.Nature(), semantic.UInt)
	err := e.eval(scope)
	if err != nil {
		return values.NewUInt(0), err
	}
	return e.value, nil
}
func (e *blockEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	values.CheckKind(e.t.Nature(), semantic.Float)
	err := e.eval(scope)
	if err != nil {
		return values.NewFloat(0), err
	}
	return e.value, nil
}
func (e *blockEvaluator) EvalBool(scope Scope) (values.Value, error) {
	values.CheckKind(e.t.Nature(), semantic.Bool)
	err := e.eval(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return e.value, nil
}
func (e *blockEvaluator) EvalTime(scope Scope) (values.Value, error) {
	values.CheckKind(e.t.Nature(), semantic.Time)
	err := e.eval(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	return e.value, nil
}
func (e *blockEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	values.CheckKind(e.t.Nature(), semantic.Duration)
	err := e.eval(scope)
	if err != nil {
		return 0, err
	}
	return e.value.Duration(), nil
}
func (e *blockEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	values.CheckKind(e.t.Nature(), semantic.Regexp)
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return e.value.Regexp(), nil
}
func (e *blockEvaluator) EvalArray(scope Scope) (values.Array, error) {
	values.CheckKind(e.t.Nature(), semantic.Object)
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return e.value.Array(), nil
}
func (e *blockEvaluator) EvalObject(scope Scope) (values.Object, error) {
	values.CheckKind(e.t.Nature(), semantic.Object)
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return e.value.Object(), nil
}
func (e *blockEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	values.CheckKind(e.t.Nature(), semantic.Object)
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return e.value.Function(), nil
}

type returnEvaluator struct {
	Evaluator
}

type declarationEvaluator struct {
	t    semantic.Type
	id   string
	init Evaluator
}

func (e *declarationEvaluator) Type() semantic.Type {
	return e.t
}

func (e *declarationEvaluator) eval(scope Scope) error {
	v, err := eval(e.init, scope)
	if err != nil {
		return err
	}

	scope.Set(e.id, v)
	return nil
}

func (e *declarationEvaluator) EvalString(scope Scope) (values.Value, error) {
	err := e.eval(scope)
	if err != nil {
		return values.NewString(""), err
	}
	return values.NewString(scope.GetString(e.id)), nil
}
func (e *declarationEvaluator) EvalInt(scope Scope) (values.Value, error) {
	err := e.eval(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	return values.NewInt(scope.GetInt(e.id)), nil
}
func (e *declarationEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	err := e.eval(scope)
	if err != nil {
		return values.NewUInt(0), err
	}

	return values.NewUInt(scope.GetUInt(e.id)), nil
}
func (e *declarationEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	err := e.eval(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}

	return values.NewFloat(scope.GetFloat(e.id)), nil
}
func (e *declarationEvaluator) EvalBool(scope Scope) (values.Value, error) {
	err := e.eval(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return values.NewBool(scope.GetBool(e.id)), nil
}
func (e *declarationEvaluator) EvalTime(scope Scope) (values.Value, error) {
	err := e.eval(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	return values.NewTime(scope.GetTime(e.id)), nil
}
func (e *declarationEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	err := e.eval(scope)
	if err != nil {
		return 0, err
	}

	return scope.GetDuration(e.id), nil
}
func (e *declarationEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return scope.GetRegexp(e.id), nil
}
func (e *declarationEvaluator) EvalArray(scope Scope) (values.Array, error) {
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return scope.GetArray(e.id), nil
}
func (e *declarationEvaluator) EvalObject(scope Scope) (values.Object, error) {
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return scope.GetObject(e.id), nil
}
func (e *declarationEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return scope.GetFunction(e.id), nil
}

type objEvaluator struct {
	t          semantic.Type
	properties map[string]Evaluator
}

func (e *objEvaluator) Type() semantic.Type {
	return e.t
}

func (e *objEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *objEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *objEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *objEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *objEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *objEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *objEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *objEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *objEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *objEvaluator) EvalObject(scope Scope) (values.Object, error) {
	obj := values.NewObject()
	for k, node := range e.properties {
		v, err := eval(node, scope)
		if err != nil {
			return nil, err
		}
		obj.Set(k, v)
	}
	return obj, nil
}
func (e *objEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type arrayEvaluator struct {
	t     semantic.Type
	array []Evaluator
}

func (e *arrayEvaluator) Type() semantic.Type {
	return e.t
}

func (e *arrayEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *arrayEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *arrayEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *arrayEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *arrayEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *arrayEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *arrayEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *arrayEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *arrayEvaluator) EvalArray(scope Scope) (values.Array, error) {
	arr := values.NewArray(e.t)
	for _, ev := range e.array {
		v, err := eval(ev, scope)
		if err != nil {
			return nil, err
		}
		arr.Append(v)
	}
	return arr, nil
}
func (e *arrayEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *arrayEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type logicalEvaluator struct {
	t           semantic.Type
	operator    ast.LogicalOperatorKind
	left, right Evaluator
}

func (e *logicalEvaluator) Type() semantic.Type {
	return e.t
}

func (e *logicalEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *logicalEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *logicalEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *logicalEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *logicalEvaluator) EvalBool(scope Scope) (values.Value, error) {
	l, err := e.left.EvalBool(scope)
	if err != nil {
		return values.NewBool(false), err
	}

	switch e.operator {
	case ast.AndOperator:
		if !l.Bool() {
			return values.NewBool(false), nil
		}
	case ast.OrOperator:
		if l.Bool() {
			return values.NewBool(true), nil
		}
	default:
		panic(fmt.Errorf("unknown logical operator %v", e.operator))
	}

	r, err := e.right.EvalBool(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return r, nil
}
func (e *logicalEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *logicalEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *logicalEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *logicalEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *logicalEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *logicalEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type conditionalEvaluator struct {
	t          semantic.Type
	test       Evaluator
	consequent Evaluator
	alternate  Evaluator
}

func (e *conditionalEvaluator) Type() semantic.Type {
	return e.t
}

func (e *conditionalEvaluator) eval(scope Scope) (values.Value, error) {
	t, err := eval(e.test, scope)
	if err != nil {
		return nil, err
	}

	if t.Bool() {
		return eval(e.consequent, scope)
	} else {
		return eval(e.alternate, scope)
	}
}

func (e *conditionalEvaluator) EvalString(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewString(""), err
	}
	return v, nil
}
func (e *conditionalEvaluator) EvalInt(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	return v, nil
}
func (e *conditionalEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewUInt(0), err
	}
	return v, nil
}
func (e *conditionalEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}
	return v, nil
}
func (e *conditionalEvaluator) EvalBool(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return v, nil
}
func (e *conditionalEvaluator) EvalTime(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	return v, nil
}
func (e *conditionalEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	v, err := e.eval(scope)
	if err != nil {
		return 0, err
	}
	return v.Duration(), nil
}
func (e *conditionalEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Regexp(), nil
}
func (e *conditionalEvaluator) EvalArray(scope Scope) (values.Array, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Array(), nil
}
func (e *conditionalEvaluator) EvalObject(scope Scope) (values.Object, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Object(), nil
}
func (e *conditionalEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Function(), nil
}

type binaryEvaluator struct {
	t           semantic.Type
	left, right Evaluator
	f           values.BinaryFunction
}

func (e *binaryEvaluator) Type() semantic.Type {
	return e.t
}

func (e *binaryEvaluator) eval(scope Scope) (values.Value, values.Value, error) {
	l, err := eval(e.left, scope)
	if err != nil {
		return nil, nil, err
	}
	r, err := eval(e.right, scope)
	if err != nil {
		return nil, nil, err
	}
	return l, r, nil
}

func (e *binaryEvaluator) EvalString(scope Scope) (values.Value, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return values.NewString(""), err
	}
	return e.f(l, r), nil
}
func (e *binaryEvaluator) EvalInt(scope Scope) (values.Value, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	return e.f(l, r), nil
}
func (e *binaryEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return values.NewUInt(0), err
	}
	return e.f(l, r), nil
}
func (e *binaryEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}
	return e.f(l, r), nil
}
func (e *binaryEvaluator) EvalBool(scope Scope) (values.Value, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return e.f(l, r), nil
}
func (e *binaryEvaluator) EvalTime(scope Scope) (values.Value, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	return e.f(l, r), nil
}
func (e *binaryEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	l, r, err := e.eval(scope)
	if err != nil {
		return 0, err
	}
	return e.f(l, r).Duration(), nil
}
func (e *binaryEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *binaryEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *binaryEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *binaryEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type unaryEvaluator struct {
	t    semantic.Type
	node Evaluator
}

func (e *unaryEvaluator) Type() semantic.Type {
	return e.t
}

func (e *unaryEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *unaryEvaluator) EvalInt(scope Scope) (values.Value, error) {
	v, err := e.node.EvalInt(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	// There is only one integer unary operator
	return values.NewInt(-v.Int()), nil
}
func (e *unaryEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *unaryEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	v, err := e.node.EvalFloat(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}
	// There is only one float unary operator
	return values.NewFloat(-v.Float()), nil
}
func (e *unaryEvaluator) EvalBool(scope Scope) (values.Value, error) {
	v, err := e.node.EvalBool(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	// There is only one bool unary operator
	return values.NewBool(!v.Bool()), nil
}
func (e *unaryEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *unaryEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	v, err := e.node.EvalDuration(scope)
	if err != nil {
		return 0, err
	}
	// There is only one duration unary operator
	return -v, nil
}
func (e *unaryEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *unaryEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *unaryEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *unaryEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type integerEvaluator struct {
	t semantic.Type
	i int64
}

func (e *integerEvaluator) Type() semantic.Type {
	return e.t
}

func (e *integerEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *integerEvaluator) EvalInt(scope Scope) (values.Value, error) {
	return values.NewInt(e.i), nil
}
func (e *integerEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	return values.NewUInt(uint64(e.i)), nil
}
func (e *integerEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *integerEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *integerEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *integerEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *integerEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *integerEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *integerEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *integerEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type stringEvaluator struct {
	t semantic.Type
	s string
}

func (e *stringEvaluator) Type() semantic.Type {
	return e.t
}

func (e *stringEvaluator) EvalString(scope Scope) (values.Value, error) {
	return values.NewString(e.s), nil
}
func (e *stringEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *stringEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *stringEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *stringEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *stringEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *stringEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *stringEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *stringEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *stringEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *stringEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type regexpEvaluator struct {
	t semantic.Type
	r *regexp.Regexp
}

func (e *regexpEvaluator) Type() semantic.Type {
	return e.t
}

func (e *regexpEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *regexpEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *regexpEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *regexpEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *regexpEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *regexpEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *regexpEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *regexpEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	return e.r, nil
}
func (e *regexpEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *regexpEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *regexpEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type booleanEvaluator struct {
	t semantic.Type
	b bool
}

func (e *booleanEvaluator) Type() semantic.Type {
	return e.t
}

func (e *booleanEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *booleanEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *booleanEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *booleanEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *booleanEvaluator) EvalBool(scope Scope) (values.Value, error) {
	return values.NewBool(e.b), nil
}
func (e *booleanEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *booleanEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *booleanEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *booleanEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *booleanEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *booleanEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type floatEvaluator struct {
	t semantic.Type
	f float64
}

func (e *floatEvaluator) Type() semantic.Type {
	return e.t
}

func (e *floatEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *floatEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *floatEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *floatEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	return values.NewFloat(e.f), nil
}
func (e *floatEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *floatEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *floatEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *floatEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *floatEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *floatEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *floatEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type timeEvaluator struct {
	t    semantic.Type
	time values.Time
}

func (e *timeEvaluator) Type() semantic.Type {
	return e.t
}

func (e *timeEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *timeEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *timeEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *timeEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *timeEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *timeEvaluator) EvalTime(scope Scope) (values.Value, error) {
	return values.NewTime(e.time), nil
}
func (e *timeEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *timeEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *timeEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *timeEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *timeEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type durationEvaluator struct {
	t        semantic.Type
	duration values.Duration
}

func (e *durationEvaluator) Type() semantic.Type {
	return e.t
}

func (e *durationEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *durationEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *durationEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *durationEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *durationEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *durationEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *durationEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	return e.duration, nil
}
func (e *durationEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *durationEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *durationEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *durationEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Function))
}

type identifierEvaluator struct {
	t    semantic.Type
	name string
}

func (e *identifierEvaluator) Type() semantic.Type {
	return e.t
}

func (e *identifierEvaluator) EvalString(scope Scope) (values.Value, error) {
	return values.NewString(scope.GetString(e.name)), nil
}
func (e *identifierEvaluator) EvalInt(scope Scope) (values.Value, error) {
	return values.NewInt(scope.GetInt(e.name)), nil
}
func (e *identifierEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	return values.NewUInt(scope.GetUInt(e.name)), nil
}
func (e *identifierEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	return values.NewFloat(scope.GetFloat(e.name)), nil
}
func (e *identifierEvaluator) EvalBool(scope Scope) (values.Value, error) {
	return values.NewBool(scope.GetBool(e.name)), nil
}
func (e *identifierEvaluator) EvalTime(scope Scope) (values.Value, error) {
	return values.NewTime(scope.GetTime(e.name)), nil
}
func (e *identifierEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	return scope.GetDuration(e.name), nil
}
func (e *identifierEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	return scope.GetRegexp(e.name), nil
}
func (e *identifierEvaluator) EvalArray(scope Scope) (values.Array, error) {
	return scope.GetArray(e.name), nil
}
func (e *identifierEvaluator) EvalObject(scope Scope) (values.Object, error) {
	return scope.GetObject(e.name), nil
}
func (e *identifierEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	return scope.GetFunction(e.name), nil
}

type valueEvaluator struct {
	value values.Value
}

func (e *valueEvaluator) Type() semantic.Type {
	return e.value.Type()
}

func (e *valueEvaluator) EvalString(scope Scope) (values.Value, error) {
	return e.value, nil
}
func (e *valueEvaluator) EvalInt(scope Scope) (values.Value, error) {
	return e.value, nil
}
func (e *valueEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	return e.value, nil
}
func (e *valueEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	return e.value, nil
}
func (e *valueEvaluator) EvalBool(scope Scope) (values.Value, error) {
	return e.value, nil
}
func (e *valueEvaluator) EvalTime(scope Scope) (values.Value, error) {
	return e.value, nil
}
func (e *valueEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	return e.value.Duration(), nil
}
func (e *valueEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	return e.value.Regexp(), nil
}
func (e *valueEvaluator) EvalArray(scope Scope) (values.Array, error) {
	return e.value.Array(), nil
}
func (e *valueEvaluator) EvalObject(scope Scope) (values.Object, error) {
	return e.value.Object(), nil
}
func (e *valueEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	return e.value.Function(), nil
}

type memberEvaluator struct {
	t        semantic.Type
	object   Evaluator
	property string
}

func (e *memberEvaluator) Type() semantic.Type {
	return e.t
}

func (e *memberEvaluator) EvalString(scope Scope) (values.Value, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return values.NewString(""), err
	}
	v, _ := o.Get(e.property)
	return v, nil
}
func (e *memberEvaluator) EvalInt(scope Scope) (values.Value, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	v, _ := o.Get(e.property)
	return v, nil
}
func (e *memberEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return values.NewUInt(0), err
	}
	v, _ := o.Get(e.property)
	return v, nil
}
func (e *memberEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}
	v, _ := o.Get(e.property)
	return v, nil
}
func (e *memberEvaluator) EvalBool(scope Scope) (values.Value, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	v, _ := o.Get(e.property)
	return v, nil
}
func (e *memberEvaluator) EvalTime(scope Scope) (values.Value, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	v, _ := o.Get(e.property)
	return v, nil
}
func (e *memberEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return 0, err
	}
	v, _ := o.Get(e.property)
	return v.Duration(), nil
}
func (e *memberEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return nil, err
	}
	v, _ := o.Get(e.property)
	return v.Regexp(), nil
}
func (e *memberEvaluator) EvalArray(scope Scope) (values.Array, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return nil, nil
	}
	v, _ := o.Get(e.property)
	return v.Array(), nil
}
func (e *memberEvaluator) EvalObject(scope Scope) (values.Object, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return nil, nil
	}
	v, _ := o.Get(e.property)
	return v.Object(), nil
}
func (e *memberEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	o, err := e.object.EvalObject(scope)
	if err != nil {
		return nil, err
	}
	v, _ := o.Get(e.property)
	return v.Function(), nil
}

type arrayIndexEvaluator struct {
	t     semantic.Type
	array Evaluator
	index Evaluator
}

func (e *arrayIndexEvaluator) Type() semantic.Type {
	return e.t
}

func (e *arrayIndexEvaluator) eval(scope Scope) (values.Value, error) {
	a, err := e.array.EvalArray(scope)
	if err != nil {
		return nil, err
	}
	i, err := e.index.EvalInt(scope)
	if err != nil {
		return nil, err
	}
	return a.Get(int(i.Int())), nil
}

func (e *arrayIndexEvaluator) EvalString(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewString(""), err
	}
	return v, nil
}
func (e *arrayIndexEvaluator) EvalInt(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	return v, nil
}
func (e *arrayIndexEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewUInt(0), err
	}
	return v, nil
}
func (e *arrayIndexEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}
	return v, nil
}
func (e *arrayIndexEvaluator) EvalBool(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return v, nil
}
func (e *arrayIndexEvaluator) EvalTime(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	return v, nil
}
func (e *arrayIndexEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	v, err := e.eval(scope)
	if err != nil {
		return 0, err
	}
	return v.Duration(), nil
}
func (e *arrayIndexEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Regexp(), nil
}
func (e *arrayIndexEvaluator) EvalArray(scope Scope) (values.Array, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Array(), nil
}
func (e *arrayIndexEvaluator) EvalObject(scope Scope) (values.Object, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Object(), nil
}
func (e *arrayIndexEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Function(), nil
}

type callEvaluator struct {
	t      semantic.Type
	callee Evaluator
	args   Evaluator
}

func (e *callEvaluator) Type() semantic.Type {
	return e.t
}

func (e *callEvaluator) eval(scope Scope) (values.Value, error) {
	args, err := e.args.EvalObject(scope)
	if err != nil {
		return nil, err
	}
	f, err := e.callee.EvalFunction(scope)
	if err != nil {
		return nil, err
	}
	return f.Call(args)
}

func (e *callEvaluator) EvalString(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewString(""), err
	}
	return v, nil
}
func (e *callEvaluator) EvalInt(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewInt(0), err
	}
	return v, nil
}
func (e *callEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewUInt(0), err
	}
	return v, nil
}
func (e *callEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewFloat(0.0), err
	}
	return v, nil
}
func (e *callEvaluator) EvalBool(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewBool(false), err
	}
	return v, nil
}
func (e *callEvaluator) EvalTime(scope Scope) (values.Value, error) {
	v, err := e.eval(scope)
	if err != nil {
		return values.NewTime(0), err
	}
	return v, nil
}
func (e *callEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	v, err := e.eval(scope)
	if err != nil {
		return 0, err
	}
	return v.Duration(), nil
}
func (e *callEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Regexp(), nil
}
func (e *callEvaluator) EvalArray(scope Scope) (values.Array, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Array(), nil
}
func (e *callEvaluator) EvalObject(scope Scope) (values.Object, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Object(), nil
}
func (e *callEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	v, err := e.eval(scope)
	if err != nil {
		return nil, err
	}
	return v.Function(), nil
}

type functionEvaluator struct {
	t      semantic.Type
	body   Evaluator
	params []functionParam
}

func (e *functionEvaluator) Type() semantic.Type {
	return e.t
}

func (e *functionEvaluator) EvalString(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.String))
}
func (e *functionEvaluator) EvalInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Int))
}
func (e *functionEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.UInt))
}
func (e *functionEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Float))
}
func (e *functionEvaluator) EvalBool(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Bool))
}
func (e *functionEvaluator) EvalTime(scope Scope) (values.Value, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Time))
}
func (e *functionEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Duration))
}
func (e *functionEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Regexp))
}
func (e *functionEvaluator) EvalArray(scope Scope) (values.Array, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Array))
}
func (e *functionEvaluator) EvalObject(scope Scope) (values.Object, error) {
	panic(values.UnexpectedKind(e.t.Nature(), semantic.Object))
}
func (e *functionEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	return &functionValue{
		t:      e.t,
		body:   e.body,
		params: e.params,
		scope:  scope,
	}, nil
}

type functionValue struct {
	t      semantic.Type
	body   Evaluator
	params []functionParam
	scope  Scope
}

type functionParam struct {
	Key     string
	Default Evaluator
	Type    semantic.Type
}

func (f *functionValue) Type() semantic.Type {
	return f.t
}
func (f *functionValue) PolyType() semantic.PolyType {
	return f.t.PolyType()
}

func (f *functionValue) IsNull() bool {
	return false
}
func (f *functionValue) Str() string {
	panic(values.UnexpectedKind(semantic.Function, semantic.String))
}
func (f *functionValue) Int() int64 {
	panic(values.UnexpectedKind(semantic.Function, semantic.Int))
}
func (f *functionValue) UInt() uint64 {
	panic(values.UnexpectedKind(semantic.Function, semantic.UInt))
}
func (f *functionValue) Float() float64 {
	panic(values.UnexpectedKind(semantic.Function, semantic.Float))
}
func (f *functionValue) Bool() bool {
	panic(values.UnexpectedKind(semantic.Function, semantic.Bool))
}
func (f *functionValue) Time() values.Time {
	panic(values.UnexpectedKind(semantic.Function, semantic.Time))
}
func (f *functionValue) Duration() values.Duration {
	panic(values.UnexpectedKind(semantic.Function, semantic.Duration))
}
func (f *functionValue) Regexp() *regexp.Regexp {
	panic(values.UnexpectedKind(semantic.Function, semantic.Regexp))
}
func (f *functionValue) Array() values.Array {
	panic(values.UnexpectedKind(semantic.Function, semantic.Array))
}
func (f *functionValue) Object() values.Object {
	panic(values.UnexpectedKind(semantic.Function, semantic.Object))
}
func (f *functionValue) Function() values.Function {
	return f
}
func (f *functionValue) Equal(rhs values.Value) bool {
	if f.Type() != rhs.Type() {
		return false
	}
	v, ok := rhs.(*functionValue)
	return ok && (f == v)
}
func (f *functionValue) HasSideEffect() bool {
	return false
}

func (f *functionValue) Call(args values.Object) (values.Value, error) {
	scope := f.scope.Copy()
	for _, p := range f.params {
		a, ok := args.Get(p.Key)
		if !ok && p.Default != nil {
			v, err := eval(p.Default, f.scope)
			if err != nil {
				return nil, err
			}
			a = v
		}
		scope.Set(p.Key, a)
	}
	return eval(f.body, scope)
}

type noopEvaluator struct {
}

func (noopEvaluator) Type() semantic.Type {
	return semantic.Nil
}

func (noopEvaluator) EvalString(scope Scope) (values.Value, error) {
	return values.NewString(""), nil
}

func (noopEvaluator) EvalInt(scope Scope) (values.Value, error) {
	return values.NewInt(0), nil
}

func (noopEvaluator) EvalUInt(scope Scope) (values.Value, error) {
	return values.NewUInt(0), nil
}

func (noopEvaluator) EvalFloat(scope Scope) (values.Value, error) {
	return values.NewFloat(0.0), nil
}

func (noopEvaluator) EvalBool(scope Scope) (values.Value, error) {
	return values.NewBool(false), nil
}

func (noopEvaluator) EvalTime(scope Scope) (values.Value, error) {
	return values.NewTime(0), nil
}

func (noopEvaluator) EvalDuration(scope Scope) (values.Duration, error) {
	return 0, nil
}

func (noopEvaluator) EvalRegexp(scope Scope) (*regexp.Regexp, error) {
	return nil, nil
}

func (noopEvaluator) EvalArray(scope Scope) (values.Array, error) {
	return nil, nil
}

func (noopEvaluator) EvalObject(scope Scope) (values.Object, error) {
	return nil, nil
}

func (noopEvaluator) EvalFunction(scope Scope) (values.Function, error) {
	return nil, nil
}
