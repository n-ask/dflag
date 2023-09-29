package dflag

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
)

func init() {
	flag.Bool("help", false, "more information")
}

var errNotAStruct = errors.New("interface is not a struct")

type cmdFlags struct {
	dest    any
	element reflect.Value
	parsed  bool
}

func (c cmdFlags) parse() error {
	t := c.element.Type()
	for i := 0; i < t.NumField(); i++ {
		f := c.element.Field(i)
		if f.Kind() == reflect.Ptr {
			return fmt.Errorf("invalid pointer field %s of type %v... don't use pointers in your cli struct", c.element.Type().Field(i).Name, f.Type())
		}
		field := t.Field(i)
		key := field.Tag.Get(cli.String())
		if key == "" {
			continue
		}
		if !f.CanSet() {
			return fmt.Errorf("unable to set value for field %s of type %v. Is field exported?", c.element.Type().Field(i).Name, f.Type())
		}
		kind := f.Kind()
		switch kind {
		case reflect.String:
			ptr := f.Addr().Interface().(*string)
			flag.StringVar(ptr, key, getDefaultValueForString(field), getDefaultUsage(field, kind))
		case reflect.Bool:
			ptr := f.Addr().Interface().(*bool)
			flag.BoolVar(ptr, key, getDefaultValueForBool(field), getDefaultUsage(field, kind))
		case reflect.Int:
			ptr := f.Addr().Interface().(*int)
			flag.IntVar(ptr, key, int(getDefaultValueForInt(field)), getDefaultUsage(field, kind))
		case reflect.Int64:
			ptr := f.Addr().Interface().(*int64)
			flag.Int64Var(ptr, key, getDefaultValueForInt(field), getDefaultUsage(field, kind))
		case reflect.Float64:
			ptr := f.Addr().Interface().(*float64)
			flag.Float64Var(ptr, key, getDefaultValueForFloat(field), getDefaultUsage(field, kind))
		default:
			return fmt.Errorf("unhandled field %s of type %v", c.element.Type().Field(i).Name, f.Type())
		}
	}
	flag.Parse()
	c.parsed = true

	return nil
}

func loadStruct(structValue any) (*cmdFlags, error) {
	value := reflect.ValueOf(structValue).Elem()
	if value.Kind() != reflect.Struct {
		return nil, errNotAStruct
	}
	return &cmdFlags{
		dest:    structValue,
		element: value,
	}, nil
}

func Load(dest any) error {
	cf, err := loadStruct(dest)
	if err != nil {
		return err
	}
	err = cf.parse()
	if err != nil {
		return err
	}

	if help := flag.Lookup("help"); help != nil {
		if help.Value.String() == "true" {
			flag.PrintDefaults()

			flag.VisitAll(func(f *flag.Flag) {
				fmt.Printf("%+v\n", f)
			})
			os.Exit(0)
		}
	}
	return nil
}
