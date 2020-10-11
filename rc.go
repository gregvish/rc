package rc

import(
    "fmt"
    "runtime"
)


type printfFunc = func(string, ...interface{})


func getCallerLine() string {
    var pc, file, line, ok = runtime.Caller(2)
    if !ok {
        panic("Can't get runtime.Caller")
    }

    var func_name = runtime.FuncForPC(pc).Name()
    return fmt.Sprintf("%s:%d [%s]", file, line, func_name)
}

func OnErr(err *error, args ...interface{}) bool {
    var print_func printfFunc = nil
    var extra_info string

    if *err == nil {
        return false
    }

    var trace = getCallerLine()
    var new_err = fmt.Errorf("%s\n%w", trace, *err)

    if len(args) == 0 {
        goto Exit
    }

    switch second_arg := args[0].(type) {
        case printfFunc:
            print_func = second_arg
            args = args[1:]
    }

    if len(args) == 0 {
        goto Exit
    }

    switch next_arg := args[0].(type) {
        case string:
            extra_info = fmt.Sprintf(next_arg, args[1:]...)
            new_err = fmt.Errorf("%s %s\n%w", trace, extra_info, *err)

        case error:
            if len(args) > 1 {
                if extra_info_fmt, ok := args[1].(string); ok {
                    extra_info = fmt.Sprintf(extra_info_fmt, args[2:]...)
                    new_err = fmt.Errorf("%s (%w) %s\n%v", trace, next_arg, extra_info, *err)
                } else {
                    panic("Argument after wrap error not an fmt sring");
                }

            } else {
                new_err = fmt.Errorf("%s (%w)\n%v", trace, next_arg, *err)
            }

        default:
            panic("Wrap with something not error or string")
    }

Exit:
    if print_func != nil {
        print_func("\n------ Traceback ------\n%v\n-----------------------\n", new_err)
    }

    *err = new_err
    return true
}

func MakeErr(err error, args ...interface{}) error {
    var new_err error
    var extra_info string

    var trace = getCallerLine()

    if len(args) > 0 {
        if extra_info_fmt, ok := args[0].(string); ok {
            extra_info = fmt.Sprintf(extra_info_fmt, args[1:]...)
            new_err = fmt.Errorf("%s (%w) %s", trace, err, extra_info)

        } else {
            panic("Second argument not an fmt string")
        }

    } else {
        new_err = fmt.Errorf("%s (%w)", trace, err)
    }

    return new_err
}

func Assert(condition bool, args ...interface{}) {
    if !condition {
        if len(args) == 0 {
            panic("Assert")
        }

        extra_fmt, ok := args[0].(string)
        if !ok {
            panic("Assert argument 2 not fmt string")
        }
        panic(fmt.Sprintf(extra_fmt, args[1:]...))
    }
}
