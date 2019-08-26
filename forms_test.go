package forms

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type SimpleStruct struct {
	String string
	Int    int
	Float  float64
	Bool   bool
}

type StructBase struct {
	String   string
	Int      int
	Float    float64
	Bool     bool
	Duration time.Duration
	Strings  []string
	Bools    []bool
	Ints     []int
	Floats   []float64

	StringPtr   *string
	IntPtr      *int
	FloatPtr    *float64
	BoolPtr     *bool
	DurationPtr *time.Duration
	StringsPtr  *[]string
	BoolsPtr    *[]bool
	IntsPtr     *[]int
	FloatsPtr   *[]float64
	StructPtr   *SimpleStruct
}

func TestPtr(t *testing.T) {
	s := StructBase{}
	MakeHTML(strings.Replace(uuid.New().String(), "-", "", -1), s, os.Stdout)
}

func TestB(t *testing.T) {
	s := struct {
		Mode         *int64 `yaml:"mode"`
		Key          string `yaml:"key"`
		Period       int    `yaml:"period"`
		OnChangeOnly bool   `yaml:"onChangeOnly"`
	}{}
	data :=
		`mode:1
key: xxx
period: 10
onChangeOnly: true
`
	if err := yaml.Unmarshal([]byte(data), &s); err != nil {
		panic(err)
	}
	MakeHTML(strings.Replace(uuid.New().String(), "-", "", -1), s, os.Stderr)
	t.Error("test")

}
func TestType(t *testing.T) {
	s := StructBase{
		String:   "string",
		Int:      1,
		Float:    0.5,
		Bool:     true,
		Duration: time.Second,

		Strings: []string{"a", "b", "c"},
		Ints:    []int{10, 20, 30},
		Floats:  []float64{1.2, 1.5, 2.5},
		Bools:   []bool{true, true, false},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte("<!DOCTYPE html><html><body>"))
		MakeHTML(strings.Replace(uuid.New().String(), "-", "", -1), s, w)
		w.Write([]byte("</body></html>"))
	})

	http.ListenAndServe(":8011", nil)
}

func TestMain(t *testing.T) {

	type Struct struct {
		StructBase

		StructOne StructBase
		Structs   []StructBase

		Textarea    string `htmlForm:"type: textarea; rows: 10"`
		CustomLabel string `htmlForm:"label: Custom; description: Custom Description"`

		ReadonlyString string  `htmlForm:"readonly: true"`
		ReadonlyInt    int     `htmlForm:"readonly: true"`
		ReadonlyFloat  float64 `htmlForm:"readonly: true"`
		ReadonlyBool   bool    `htmlForm:"readonly: true"`

		ReadonlyStrings []string  `htmlForm:"readonly: true"`
		ReadonlyInts    []int     `htmlForm:"readonly: true"`
		ReadonlyFloats  []float64 `htmlForm:"readonly: true"`
		ReadonlyBools   []bool    `htmlForm:"readonly: true"`

		ReadonlyStruct  StructBase   `htmlForm:"readonly: true"`
		ReadonlyStructs []StructBase `htmlForm:"readonly: true; itemLabel: Struct"`

		ReadonlyTextarea    string `htmlForm:"type: textarea; rows: 20; readonly: true"`
		ReadonlyCustomLabel string `htmlForm:"label: Custom; description: Custom Description; readonly: true"`
	}

	data := Struct{

		StructBase: StructBase{
			String:   "string",
			Int:      1,
			Float:    0.5,
			Bool:     true,
			Duration: time.Second,

			Strings: []string{"a", "b", "c"},
			Ints:    []int{10, 20, 30},
			Floats:  []float64{1.2, 1.5, 2.5},
			Bools:   []bool{true, true, false},
		},
		StructOne: StructBase{
			String: "string",
			Int:    1,
			Float:  0.5,
			Bool:   true,

			Strings: []string{"a", "b", "c"},
			Ints:    []int{10, 20, 30},
			Floats:  []float64{1.2, 1.5, 2.5},
			Bools:   []bool{true, true, false},
		},
		Textarea:    "hello\n1\n2\n3\n",
		CustomLabel: "custom label",

		ReadonlyString: "readonly string",
		ReadonlyInt:    100,
		ReadonlyFloat:  22.5,
		ReadonlyBool:   false,

		ReadonlyStrings: []string{"A", "B", "C"},
		ReadonlyInts:    []int{1, 2, 3, 4, 5, 6},
		ReadonlyFloats:  []float64{100.1, 100.2, 100.3},
		ReadonlyBools:   []bool{true, true, false},

		ReadonlyStruct: StructBase{
			String: "string",
			Int:    1,
			Float:  0.5,
			Bool:   true,

			Strings: []string{"a", "b", "c"},
			Ints:    []int{10, 20, 30},
			Floats:  []float64{1.2, 1.5, 2.5},
			Bools:   []bool{true, true, false},
		},
		ReadonlyStructs: []StructBase{
			StructBase{
				String: "string",
				Int:    1,
				Float:  0.5,
				Bool:   true,

				Strings: []string{"a", "b", "c"},
				Ints:    []int{10, 20, 30},
				Floats:  []float64{1.2, 1.5, 2.5},
				Bools:   []bool{true, true, false},
			},
		},

		ReadonlyTextarea:    "readonly textarea",
		ReadonlyCustomLabel: "readonly custom label",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte("<!DOCTYPE html><html><body>"))
		MakeHTML(strings.Replace(uuid.New().String(), "-", "", -1), data, w)
		w.Write([]byte("</body></html>"))
	})

	http.ListenAndServe(":8011", nil)
}
