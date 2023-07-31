package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
)

func PrintlnRaw(output io.Writer, msg interface{}) {
	m := fmt.Sprint(msg)
	fullMessage := m + "\n" + strings.Repeat("\b", len(m))
	fmt.Fprint(output, fullMessage)
}

// ExecCmd exec command on specific pod and wait the command's output.
func ExecCmdExample(podName, namespace, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	PrintlnRaw(os.Stderr, "Executing...")
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
	PrintlnRaw(os.Stderr, req.URL())
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
