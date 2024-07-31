package generator

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	luc "github.com/PlayerR9/lib_units/common"
	lus "github.com/PlayerR9/lib_units/slices"
)

// InitLogger initializes the logger with the given prefix.
//
// Parameters:
//   - prefix: The prefix to use for the logger.
//
// Returns:
//   - *log.Logger: The initialized logger. Never nil.
//
// If the prefix is empty, it defaults to "go_generator".
func InitLogger(prefix string) *log.Logger {
	if prefix == "" {
		prefix = "go_generator"
	}

	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(prefix)
	builder.WriteString("]: ")

	logger_prefix := builder.String()

	logger := log.New(os.Stdout, logger_prefix, log.Lshortfile)
	return logger
}

// align_generics is a helper function that aligns the generics in the given StructFieldsVal and GenericsSignVal.
//
// Parameters:
//   - fv: The StructFieldsVal.
//   - gv: The GenericsSignVal.
//
// Returns:
//   - error: An error if the alignment fails (i.e., either the StructFieldsVal or the GenericsSignVal is nil when the other is not).
func align_generics(fv *StructFieldsVal, tv *TypeListVal, gv *GenericsSignVal) error {
	if gv == nil {
		if fv != nil && tv != nil {
			return luc.NewErrInvalidUsage(
				errors.New("not specified the *StructFieldsVal and *TypeListVal but specified the *GenericsSignVal"),
				"Make sure to call go_generator.SetStructFieldsFlag() and go_generator.SetTypeListFlag() as well",
			)
		}
	} else {
		if fv == nil && tv == nil {
			return luc.NewErrInvalidUsage(
				errors.New("not specified the *StructFieldsVal and *TypeListVal but specified the *GenericsSignVal"),
				"Make sure to call go_generator.SetStructFieldsFlag() and go_generator.SetTypeListFlag() as well",
			)
		}
	}

	var all_generics []rune

	if fv != nil {
		iter := fv.generics.KeyIterator()
		luc.Assert(iter != nil, "iter must not be nil")

		for {
			id, err := iter.Consume()
			if err != nil {
				break
			}

			all_generics = lus.TryInsert(all_generics, id)
		}
	}

	if tv != nil {
		iter := tv.generics.KeyIterator()
		luc.Assert(iter != nil, "iter must not be nil")

		for {
			id, err := iter.Consume()
			if err != nil {
				break
			}

			all_generics = lus.TryInsert(all_generics, id)
		}
	}

	for _, generic_id := range all_generics {
		pos, ok := slices.BinarySearch(gv.letters, generic_id)
		if ok {
			continue
		}

		gv.letters = slices.Insert(gv.letters, pos, generic_id)
		gv.types = slices.Insert(gv.types, pos, "any")
	}

	return nil
}

// ParseFlags parses the command line flags.
//
// Returns:
//   - error: An error if any.
func ParseFlags() error {
	flag.Parse()

	err := align_generics(StructFieldsFlag, TypeListFlag, GenericsSigFlag)
	if err != nil {
		return err
	}

	return nil
}

// Generater is the interface that all generators must implement.
type Generater interface {
	// SetPackageName sets the package name for the generated code.
	//
	// Parameters:
	//   - pkg_name: The package name to use for the generated code.
	SetPackageName(pkg_name string)
}

// DoFunc is the type of the function to perform on the data before generating the code.
//
// Parameters:
//   - T: The data to perform the function on.
//
// Returns:
//   - error: An error if occurred.
type DoFunc[T Generater] func(T) error

// CodeGenerator is the code generator.
type CodeGenerator[T Generater] struct {
	// t is the template to use for the generated code.
	templ *template.Template

	// do_funcs is the list of functions to perform on the data before generating the code.
	do_funcs []DoFunc[T]
}

// NewCodeGenerator creates a new code generator.
//
// Parameters:
//   - templ: The template to use for the generated code.
//
// Returns:
//   - *CodeGenerator: The code generator.
//   - error: An error of type *common.ErrInvalidParameter if the templ is nil.
func NewCodeGenerator[T Generater](templ *template.Template) (*CodeGenerator[T], error) {
	if templ == nil {
		return nil, luc.NewErrNilParameter("templ")
	}

	return &CodeGenerator[T]{
		templ:    templ,
		do_funcs: make([]DoFunc[T], 0),
	}, nil
}

// NewCodeGeneratorFromTemplate creates a new code generator from a template. Panics
// if the template is not valid.
//
// Parameters:
//   - name: The name of the template.
//   - templ: The template to use for the generated code.
//
// Returns:
//   - *CodeGenerator: The code generator.
//   - error: An error of type *common.ErrInvalidParameter if the templ is invalid.
func NewCodeGeneratorFromTemplate[T Generater](name, templ string) (*CodeGenerator[T], error) {
	t, err := template.New(name).Parse(templ)
	if err != nil {
		return nil, luc.NewErrInvalidParameter("templ", err)
	}

	luc.AssertNil(t, "t")

	return &CodeGenerator[T]{
		templ:    t,
		do_funcs: make([]DoFunc[T], 0),
	}, nil
}

// AddDoFunc adds a function to perform on the data before generating the code.
//
// Parameters:
//   - do_func: The function to perform on the data before generating the code.
//
// Does nothing if the do_func is nil.
func (cg *CodeGenerator[T]) AddDoFunc(do_func DoFunc[T]) {
	if do_func == nil {
		return
	}

	cg.do_funcs = append(cg.do_funcs, do_func)
}

// Generate generates code using the given generator and writes it to the given destination file.
//
// WARNING:
//   - Remember to call this function iff the function go_generator.SetOutputFlag() was called
//     and only after the function flag.Parse() was called.
//
// Parameters:
//   - file_name: The file name to use for the generated code.
//   - suffix: The suffix to use for the generated code. This should end with the ".go" extension.
//   - data: The data to use for the generated code.
//
// Returns:
//   - string: The output location of the generated code.
//   - error: An error if occurred.
//
// Errors:
//   - *common.ErrInvalidParameter: If the file_name or suffix is an empty string.
//   - error: Any other type of error that may have occurred.
func (cg *CodeGenerator[T]) Generate(file_name, suffix string, data T) (string, error) {
	luc.AssertNil(cg.templ, "cg.templ")

	output_loc, err := FixOutputLoc(file_name, suffix)
	if err != nil {
		return output_loc, fmt.Errorf("failed to fix output location: %w", err)
	}

	pkg_name, err := FixImportDir(output_loc)
	if err != nil {
		return output_loc, fmt.Errorf("failed to fix import path: %w", err)
	}

	data.SetPackageName(pkg_name)

	for _, f := range cg.do_funcs {
		if f == nil {
			continue
		}

		err := f(data)
		if err != nil {
			return output_loc, err
		}
	}

	var buff bytes.Buffer

	err = cg.templ.Execute(&buff, data)
	if err != nil {
		return output_loc, err
	}

	res := buff.Bytes()

	dir := filepath.Dir(output_loc)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return output_loc, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	err = os.WriteFile(output_loc, res, 0644)
	if err != nil {
		return output_loc, err
	}

	return output_loc, nil
}
