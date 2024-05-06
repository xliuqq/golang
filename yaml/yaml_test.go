package yaml

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

type student struct {
	Name     string    `yaml:"name"`
	Birth    int       `yaml:"birth"`
	ArgsName string    `yaml:"argsName"`
	Args     extension `yaml:"args,omitempty"`
}

type extension struct {
	Raw    []byte `yaml:"-"`
	Object any    `yaml:"-"`
}

type StudentArgs struct {
	Skill string `yaml:"skill"`
	Born  string `yaml:"born"`
}

func init() {
	Register(&StudentArgs{})
}

var mappings = map[string]reflect.Type{}

func Register(v interface{}) {
	t := reflect.TypeOf(v)
	// 使用指针
	if t.Kind() != reflect.Pointer {
		panic("All types must be pointers to structs.")
	}
	// 指针需要用 Elem 获取所指向的对象的 Type
	elem := t.Elem()
	mappings[elem.Name()] = t
}

func (e *extension) UnmarshalYAML(value *yaml.Node) error {
	if e == nil {
		return errors.New("extension: UnmarshalYAML on nil pointer")
	}

	var str interface{}
	err := value.Decode(&str)
	e.Raw, err = yaml.Marshal(str)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalYAML(data string) (*student, error) {
	var stu student
	err := yaml.Unmarshal([]byte(data), &stu)
	if err != nil {
		return nil, err
	}

	if len(stu.Args.Raw) != 0 {
		name := stu.ArgsName
		r, ok := mappings[name]
		if !ok {
			return &stu, errors.New(fmt.Sprintf("type %s is not registered", name))
		}
		// 指针需要用 Elem 获取所指向的对象的Type，构造对象，然后再用 Interface 变成指针
		v := reflect.New(r.Elem()).Interface()
		err = yaml.Unmarshal(stu.Args.Raw, v)
		if err != nil {
			return &stu, err
		}
		stu.Args.Object = v
	}

	return &stu, nil
}

func TestYaml_Unmarshal(t *testing.T) {

	stuStrings := []string{
		`
name: xliu1992
birth: 1992
argsName: StudentArgs
args:
  skill: "IT"
  born: "China"
`,
		`
name: xliu1992
birth: 1992
`,
	}
	for _, stuString := range stuStrings {
		stu, err := UnmarshalYAML(stuString)
		if err != nil {
			t.Errorf("yaml unmarshal error, got %v", err)
		}
		fmt.Printf("args %v\n", stu.Args.Object)
	}
}
