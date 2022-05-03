package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

var s struct {
	X int
	y float64
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) ReflectCallFunc(age int, name string) error {
	fmt.Println("json ReflectCallFunc")
	return nil
}

func (u User) GetAge(name string) int {
	fmt.Println("getAge")
	return u.Age
}

func TestReflect3(t *testing.T) {
	user := User{
		Id:   1,
		Name: "json",
		Age:  27,
	}
	getType := reflect.TypeOf(user)
	getValue := reflect.ValueOf(user)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v \n", field.Name, field.Type, value)
	}
}

func TestReflect2(t *testing.T) {
	user := User{1, "json", 25}
	getType := reflect.TypeOf(user)
	fmt.Println("get type is:", getType.Name())
	getValue := reflect.ValueOf(user)
	fmt.Println("get value is:", getValue)
	getType2 := getValue.Type()
	fmt.Println("get type is:", getType2.Name(), getType2.Field(1))
}

func TestReflect1(t *testing.T) {
	vs := reflect.ValueOf(&s).Elem()
	vx := vs.Field(1)
	vb := reflect.ValueOf(666)
	if vx.CanSet() {
		vx.Set(vb)
	}

	t.Log(s)
}

type children struct {
	Age int
}

type Nested struct {
	X     int
	Child children
}

func TestReflect4(t *testing.T) {
	n := new(Nested)
	vs := reflect.ValueOf(n).Elem()
	vz := vs.Field(1)
	if vz.CanSet() {
		vz.Set(reflect.ValueOf(children{Age: 19}))
		fmt.Println(n)
	}
}

func TestReflect5(t *testing.T) {
	getType := reflect.TypeOf(new(User))
	method := getType.Method(0)
	t.Log(method.Name, method.Name)
}

func TestReflect6(t *testing.T) {
	user := User{
		Id:   1,
		Name: "alex",
		Age:  19,
	}

	ref := reflect.ValueOf(user)
	m := ref.MethodByName("ReflectCallFunc")
	args := []reflect.Value{reflect.ValueOf(18), reflect.ValueOf("aa")}
	m.Call(args)
	tf := m.Type()
	t.Log(tf.NumIn(), tf.NumOut(), ref.NumMethod())
}

func (u *User) RefPointMethod() {
	fmt.Println("hello world")
}

func TestReflect7(t *testing.T) {
	user := User{
		Id:   1,
		Name: "Lesi",
		Age:  30,
	}
	ref := reflect.ValueOf(&user)
	m := ref.MethodByName("RefPointMethod")
	m.Call([]reflect.Value{})
}

func (u *User) PointMethodReturn(name string, age int) (string, int) {
	return name, age
}

func TestReflect8(t *testing.T) {
	user := User{
		Id:   1,
		Name: "Lesi",
		Age:  30,
	}
	ref := reflect.ValueOf(&user)
	m := ref.MethodByName("PointMethodReturn")
	args := []reflect.Value{reflect.ValueOf("json"), reflect.ValueOf(30)}
	res := m.Call(args)
	t.Log("name:", res[0].Interface())
	t.Log("age:", res[1].Interface())
}

// 反射在运行时创建结构体
func MakeStruct(vals ...interface{}) reflect.Value {
	var sfs []reflect.StructField
	for k, v := range vals {
		t := reflect.TypeOf(v)
		sf := reflect.StructField{
			Name: fmt.Sprintf("F%d", k+1),
			Type: t,
		}
		sfs = append(sfs, sf)
	}
	st := reflect.StructOf(sfs)
	os := reflect.New(st)
	return os
}

func TestReflect9(t *testing.T) {
	sr := MakeStruct(0, "", []int{})
	sr.Elem().Field(0).SetInt(20)
	sr.Elem().Field(1).SetString("reflect me")
	// 赋值数组
	v := []int{1, 2, 3}
	rv := reflect.ValueOf(v)
	sr.Elem().Field(2).Set(rv)

	t.Log(sr)
}

// 函数与反射
func Handler2(args int, reply *int) {
	fmt.Println(*reply, args)
	*reply = args
	fmt.Println(*reply, args)
}

func TestReflect10(t *testing.T) {
	v2 := reflect.ValueOf(Handler2)
	args := reflect.ValueOf(5)
	replyv := reflect.New(reflect.TypeOf(-1))
	v2.Call([]reflect.Value{args, replyv})
}

func TestReflect11(t *testing.T) {
	u := User{
		Id:   1,
		Name: "bill",
		Age:  10,
	}
	cond := User{
		Id:   1,
		Name: "bill",
		Age:  0,
	}
	updateEntity(u, cond)
}

type dbOperate interface {
	update(q interface{}) string
	where(cond interface{}) string
}

//
func updateEntity(q interface{}, cond interface{}) string {
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		update := update(q)
		if cond != nil {
			update = fmt.Sprintf("%s %s ", update, where(cond))
		}
		fmt.Printf("%v\n", update)
		return update
	}
	return ""
}

func update(q interface{}) string {
	v := reflect.ValueOf(q)
	t := reflect.TypeOf(q)

	update := fmt.Sprintf("update %s set ", t.Name())
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int:
			field := t.Field(i)
			value := v.Field(i).Int()
			if value != 0 {
				update = fmt.Sprintf("%v%v=%v,", update, field.Name, value)
			}
		case reflect.String:
			field := t.Field(i)
			value := v.Field(i).String()
			if value != "" {
				update = fmt.Sprintf("%v%v=%v,", update, field.Name, value)
			}
		}
	}
	update = update[:len(update)-1]
	return update
}

func where(cond interface{}) string {
	v := reflect.ValueOf(cond)
	t := reflect.TypeOf(cond)
	where := fmt.Sprintf(" where ")
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int:
			field := t.Field(i)
			value := v.Field(i).Int()
			if value != 0 {
				where = fmt.Sprintf("%v%v=%v", where, field.Name, value)
			}
		case reflect.String:
			field := t.Field(i)
			value := v.Field(i).String()
			if value != "" {
				where = fmt.Sprintf("%v%v=%s and ", where, field.Name, value)
			}
		}
	}
	where = where[:len(where)-4]
	return where
}
