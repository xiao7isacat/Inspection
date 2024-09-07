package checkctl

import "k8s.io/klog/v2"

type Job struct {
}

func (this *Job) Get() error {
	klog.Info("get job")
	return nil
}

func (this *Job) Add(resourceFilePath string) error {
	klog.Info("add job")
	return nil
}

func (this *Job) Update(resourceFilePath string) error {
	klog.Info("update job")
	return nil
}

func (this *Job) Delete() error {
	klog.Info("delete job")
	return nil
}
