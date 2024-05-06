package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

/**
  K8s Scheduler Plugin 中的 PluginConfig 中的 Args 字段是允许任意类型的。
  K8s 的实现机制，v1/KubeSchedulerConfiguration 是对外的对象（用于反序列化），给 KubeScheduler 的是内部对象 type/KubeSchedulerConfiguration
*/
type student struct {
	Name     string `json:"name" yaml:"name"`
	Birth    int    `json:"birth" yaml:"birth"`
	ArgsName string `json:"argsName" yaml:"argsName"`
	// K8s 中会将 Args runtime.RawExtension 转为另一个同名 的 Args runtime.Object 对象
	Args extension `json:"args,omitempty" yaml:"args,omitempty"`
}

// 分装一层的好处是可以通过自定义 UnmarshalJSON ，如果 Args 的类型是 any，则需要对 struct student 自定义反序列化实现（会比较难实现）
type extension struct {
	Raw    []byte `json:"-" yaml:"-"`
	Object any    `json:"-" yaml:"-"`
}

type StudentArgs struct {
	Skill string `json:"skill" yaml:"skill"`
	Born  string `json:"born" yaml:"born"`
}

var mappings = map[string]reflect.Type{}

func init() {
	Register(&StudentArgs{})
}

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

func (e *extension) UnmarshalJSON(in []byte) error {
	if e == nil {
		return errors.New("runtime.RawExtension: UnmarshalJSON on nil pointer")
	}
	if !bytes.Equal(in, []byte("null")) {
		e.Raw = append(e.Raw[0:0], in...)
	}
	return nil
}

func UnmarshalJSON(data string) (*student, error) {
	var stu student
	err := json.Unmarshal([]byte(data), &stu)
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
		err = json.Unmarshal(stu.Args.Raw, v)
		if err != nil {
			return &stu, err
		}
		stu.Args.Object = v
	}

	return &stu, nil
}

func TestJson_Unmarshal(t *testing.T) {

	stuStrings := []string{
		`
{
	"name": "xliu1992",
	"birth" : 1992,
	"argsName": "StudentArgs",
	"args" : {
		"skill" : "IT",
		"born" : "China"
	}
}
`,
		`
{
	"name": "xliu1992",
	"birth" : 1992,
	"argsName": "StudentArgs",
	"args": null
}
`,
	}

	for _, stuString := range stuStrings {
		stu, err := UnmarshalJSON(stuString)
		if err != nil {
			t.Errorf("json unmarshal error, got %v", err)
		}
		fmt.Printf("args %v\n", stu.Args.Object)
	}
}
