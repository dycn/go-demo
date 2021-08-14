package process

import (
	"context"
	"fmt"
	"zdy_demo/common/rpc/thrift/gen-go/students"
)

type ClassMembers struct {
	Name string
	All  []*students.Student
}

func (c *ClassMembers) Ping(_ context.Context, num int64) error {
	fmt.Println("收到ping", num)
	return nil
}

func (c *ClassMembers) GetClassName(_ context.Context) (string, error) {
	return c.Name, nil
}

func (c *ClassMembers) List(_ context.Context) ([]*students.Student, error) {
	return c.All, nil
}

func (c *ClassMembers) Add(_ context.Context, s *students.Student) error {
	c.All = append(c.All, s)
	return nil
}

func (c *ClassMembers) IsNameExist(_ context.Context, name string) (bool, error) {
	var exists bool
	for _, v := range c.All {
		if v.Name == name {
			exists = true
			break
		}
	}
	return exists, nil
}

func (c *ClassMembers) Count(_ context.Context) (int16, error) {
	return int16(len(c.All)), nil
}
