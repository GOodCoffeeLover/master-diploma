package exec

import (
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecCmd exec command on specific pod and wait the command's output.
func ExecCmdExample(podName string, command string,
	stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	fmt.Println("Executing...")

	config := &restclient.Config{
		Host:    "http://localhost:8080",
		APIPath: "/api",
	}
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})

	client, err := restclient.RESTClientFor(config)
	if err != nil {
		return fmt.Errorf("can't create client: %w", err)
	}
	cmd := []string{
		"sh",
		"-c",
		command,
	}
	req := client.Post().Resource("pods").Name(podName).Namespace("default").SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if stdin == nil {
		option.Stdin = false
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	url := req.URL()
	fmt.Println(url)
	if err != nil {
		return fmt.Errorf("can't create spdy executor: %w", err)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return fmt.Errorf("error while executing: %w", err)
	}

	return nil
}

type Executor struct{}

func NewExecutor() Executor {
	return Executor{}
}

func (e *Executor) Exec() {
	fmt.Println("Hello")
}
