package utils

import (
	"Gopo/utils/proto"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"math/rand"
	"strings"
)

type options struct {
	envOptions     []cel.EnvOption
	programOptions []cel.ProgramOption
}

func newEnvOption() *options {
	opt := &options{}
	opt.envOptions = []cel.EnvOption{
		cel.Types(
			&proto.Response{},&proto.Reverse{}),
		cel.Declarations(
			decls.NewVar("response", decls.NewObjectType("proto.Response")),
			//decls.NewVar("reverse", decls.NewObjectType("proto.Reverse")),
		),
		//自定义函数声明
		cel.Declarations(bcontansDec, md5Dec, randomIntDec, reverseWaitDec),
	}
	opt.programOptions = []cel.ProgramOption{
		//函数实现
		cel.Functions(bcontainsFunc, md5Func, randomIntFunc, reverseWaitFunc),
	}
	return opt
}

// UpdateCompileOptions 用于把变量注册到env中，在expression的时候可以索引到
func (opt *options) UpdateCompileOptions(args map[string]string) {
	for k, v := range args {
		var d *exprpb.Decl
		if strings.HasPrefix(v, "randomInt") {
			d = decls.NewVar(k, decls.Int)
			} else if strings.HasPrefix(v, "newReverse") {
				d = decls.NewVar(k, decls.NewObjectType("proto.Reverse"))
		} else {
			d = decls.NewVar(k, decls.String)
		}
		opt.envOptions = append(opt.envOptions, cel.Declarations(d))
	}
}

//UpdateFunctionOptions 用来预先处理rule的键名，加载到env中
//后续处理类似 r0()&&r1()这类的expression，可以索引到env中执行
func (opt *options) UpdateFunctionOptions(name string, isTrue ref.Val) {
	//expression:=v.Expression
	//declarations
	dec := decls.NewFunction(name, decls.NewOverload(name, []*exprpb.Type{}, decls.Bool))
	opt.envOptions = append(opt.envOptions, cel.Declarations(dec))
	function := &functions.Overload{
		Operator: name,
		Function: func(values ...ref.Val) ref.Val {
			return isTrue
		},
	}
	opt.programOptions = append(opt.programOptions, cel.Functions(function))
}

func Evaluate(env *cel.Env, expression string, params map[string]interface{}) (ref.Val, error) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		Error("compile: ", iss.Err())
		return nil, iss.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		ErrorF("Program creation error: %v", err)
		return nil, err
	}

	out, _, err := prg.Eval(params)
	if err != nil {
		ErrorF("Evaluation error: %v", err)
		return nil, err
	}
	return out, nil
}
func (opt *options) CompileOptions() []cel.EnvOption {
	return opt.envOptions
}
func (opt *options) ProgramOptions() []cel.ProgramOption {
	return opt.programOptions
}

/*
	Declarations
	Functions
*/
var bcontansDec = decls.NewFunction("bcontains", decls.NewInstanceOverload("bytes_contains_bytes", []*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bool))
var bcontainsFunc = &functions.Overload{
	Operator: "bcontains",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.Bool(bytes.Contains(lhs.(types.Bytes), rhs.(types.Bytes)))
	},
}

var md5Dec = decls.NewFunction("md5", decls.NewOverload("md5_string", []*exprpb.Type{decls.String}, decls.String))
var md5Func = &functions.Overload{
	Operator: "md5_string",
	Unary: func(value ref.Val) ref.Val {
		v, ok := value.(types.String)
		if !ok {
			return types.ValOrErr(value, "unexpected type '%v' passed to md5_string", value.Type())
		}
		return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
	},
}

var randomIntDec = decls.NewFunction("randomInt", decls.NewOverload("randomInt_int_int", []*exprpb.Type{decls.Int, decls.Int}, decls.Int))
var randomIntFunc = &functions.Overload{
	Operator: "randomInt_int_int",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		from, ok := lhs.(types.Int)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to randomInt", lhs.Type())
		}
		to, ok := rhs.(types.Int)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to randomInt", rhs.Type())
		}
		min, max := int(from), int(to)
		return types.Int(rand.Intn(max-min) + min)
	},
}
var reverseWaitDec = decls.NewFunction("wait", decls.NewInstanceOverload("reverse_wait_int", []*exprpb.Type{decls.Any, decls.Int}, decls.Bool))
var reverseWaitFunc = &functions.Overload{
	Operator: "reverse_wait_int",
	Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
		reverse, ok := lhs.Value().(*proto.Reverse)
		if !ok {
			return types.ValOrErr(lhs, "unexpected type '%v' passed to 'wait'", lhs.Type())
		}
		timeout, ok := rhs.Value().(int64)
		if !ok {
			return types.ValOrErr(rhs, "unexpected type '%v' passed to 'wait'", rhs.Type())
		}
		return types.Bool(reverseCheck(reverse, timeout))
	},
}
