package checkctl

import (
	"fmt"
	"inspection/global"
)

type Resource interface {
	Get() error
	Delete() error
	Update() error
	Add() error
}

func NewResource(resourceType string) Resource {
	switch resourceType {
	case "job":
		return &Job{Name: global.ResourceName, IpString: global.NodeAddrs}
	case "desired":
		return &Desired{Name: global.ResourceName, ResourceFilePath: global.ResourceFilePath}
	case "script":
		return &Script{Name: global.ResourceName, ResourceFilePath: global.ResourceFilePath}
	}
	return nil
}
func ResourceOperate(resourceType, operate, resourceName, resourceFilePath string) error {
	resource := NewResource(resourceType)
	if resource == nil {
		return fmt.Errorf("未获取到当前的资源类型")
	}
	switch operate {
	case "add":
		if err := resource.Add(); err != nil {
			return err
		}
	case "get":
		if err := resource.Get(); err != nil {
			return err
		}
	case "update":
		if err := resource.Update(); err != nil {
			return err
		}
	case "delete":
		if err := resource.Delete(); err != nil {
			return err
		}
	}
	return nil
}
