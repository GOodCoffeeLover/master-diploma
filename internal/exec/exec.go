package exec

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

// ExecCmd exec command on specific pod and wait the command's output.
func ExecCmdExample(podName, namespace, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	fmt.Println("Executing...")
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return fmt.Errorf("can't create kube config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("can't create client: %w", err)
	}

	client := clientset.CoreV1().RESTClient()
	cmd := []string{
		command,
	}
	req := client.Post().Namespace("default").Resource("pods").Name(podName).SubResource("exec")
	option := &v1.PodExecOptions{
		Container: podName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	fmt.Println(req.URL())
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("can't create spdy executor: %w", err)
	}

	return exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Tty:    true,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
}

type Executor struct{}

func NewExecutor() Executor {
	return Executor{}
}

func (e *Executor) Exec() {
	fmt.Println("Hello")
}
