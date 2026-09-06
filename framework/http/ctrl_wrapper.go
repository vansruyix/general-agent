package http

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	"general-agent/extension/errorx"
	"general-agent/extension/logz"

	"github.com/gin-gonic/gin"
)

var ErrWrapperCall = errorx.NewError(-13001, "CtrlWrapper reflect err")

var ActualContextType = reflect.TypeOf(new(context.Context)).Elem()

type bindFunc func(ctx *gin.Context, obj interface{}) error

type CtrlWrapper struct {
	render func(ctx *gin.Context, data interface{}, err error)
}

type ctrlFn struct {
	Val    reflect.Value
	Type   reflect.Type
	NumIn  int
	NumOut int
}

func NewCtrlWrapper(render func(ctx *gin.Context, data interface{}, err error)) *CtrlWrapper {
	return &CtrlWrapper{render: render}
}

func newCtrlFn(fn interface{}) *ctrlFn {
	var (
		fnVal  = reflect.ValueOf(fn)
		fnType = reflect.TypeOf(fn)
	)

	if fnType.Kind() != reflect.Func {
		panic("CtrlWrapper support controller func only")
	}

	return &ctrlFn{
		Val:    fnVal,
		Type:   fnType,
		NumIn:  fnType.NumIn(),
		NumOut: fnType.NumOut(),
	}
}

func (cfn *ctrlFn) validate() {
	if cfn.NumIn < 1 || cfn.NumIn > 2 {
		panic("CtrlWrapper controller func support only 1 or 2 in param")
	}
	if !cfn.Type.In(0).Implements(ActualContextType) {
		panic("CtrlWrapper controller func first param must be context")
	}
	if cfn.NumOut < 1 || cfn.NumOut > 2 {
		panic(fmt.Sprintf("%s: CtrlWrapper controller func support only 1 or 2 out param", cfn.Type.Elem().Name()))
	}
}

// 绑定参数值
func (cfn *ctrlFn) bind(ctx *gin.Context, binds ...bindFunc) (argv []reflect.Value, err error) {
	argv = make([]reflect.Value, cfn.NumIn)

	argv[0] = reflect.ValueOf(ctx.Request.Context())

	if cfn.NumIn == 1 {
		return
	}
	var (
		reqType  = cfn.Type.In(1)
		reqValue reflect.Value
	)
	if reqType.Kind() == reflect.Ptr {
		if reqType.Elem() == nil {
			return
		}
		reqValue = reflect.New(reqType.Elem())
	} else {
		err = errors.New("param should be Pointer")
		logz.Error(ctx, err.Error())
		fmt.Println(string(debug.Stack()))
		return
	}
	// 第二个请求参数
	cfnReq := reqValue.Interface()

	// 默认使用 ShouldBind 做参数映射，同时支持自定义传参
	if len(binds) == 0 {
		err = defaultBind(ctx, cfnReq)
	} else {
		err = binds[0](ctx, cfnReq)
	}
	if err != nil {
		return
	}

	argv[1] = reflect.ValueOf(cfnReq)
	return
}

func defaultBind(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBind(req); err != nil {
		return err
	}
	if err := ctx.ShouldBindQuery(req); err != nil {
		return err
	}
	// if err := ctx.ShouldBindUri(req); err != nil {
	// 	return err
	// }
	return nil
}

func (cfn *ctrlFn) call(argv []reflect.Value) (result interface{}, err error) {
	var (
		resp = cfn.Val.Call(argv)
		e    interface{}
	)

	switch len(resp) {
	case 1:
		e = resp[0].Interface()
	case 2:
		result = resp[0].Interface()
		e = resp[1].Interface()
	default:
		return nil, ErrWrapperCall
	}

	err, ok := e.(error)
	if e != nil && !ok {
		return nil, ErrWrapperCall
	}

	return
}

func (c *CtrlWrapper) With(fn interface{}, binds ...bindFunc) gin.HandlerFunc {
	var cfn = newCtrlFn(fn)
	cfn.validate()

	return func(ctx *gin.Context) {
		var (
			argv   []reflect.Value
			result interface{}
			err    error
		)

		if argv, err = cfn.bind(ctx, binds...); err != nil {
			_ = ctx.Error(err)
			c.render(ctx, result, err)
			return
		}

		if result, err = cfn.call(argv); err != nil {
			_ = ctx.Error(err)
			c.render(ctx, result, err)
			return
		}

		c.render(ctx, result, err)
	}
}
