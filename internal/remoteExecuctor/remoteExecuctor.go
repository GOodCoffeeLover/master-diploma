package remoteExecuctor

import (
	"context"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func PrintlnRaw(output io.Writer, msg interface{}) {
	// m := fmt.Sprint(msg)
	// fullMessage := m + "\n" + strings.Repeat("\b", len(m))
	// fmt.Fprint(output, fullMessage)
}

type RemoteExecutor struct {
	config    *rest.Config
	k8sClient rest.Interface
	podName   string
	namespace string
}

func NewRemoteExecutor(config *rest.Config, namespace, podName string) (*RemoteExecutor, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := clientset.CoreV1().RESTClient()
	return &RemoteExecutor{
		config:    config,
		k8sClient: client,
		namespace: namespace,
		podName:   podName,
	}, nil
}

func (re *RemoteExecutor) Exec(command string, stdin io.Reader, stdout io.WriteCloser) error {
	if stdin == nil && stdout == nil {
		return fmt.Errorf("can't execute command(%v) with nil stdin and stdout", command)
	}

	req := re.k8sClient.Post().Namespace(re.namespace).Resource("pods").Name(re.podName).SubResource("exec")

	option := &v1.PodExecOptions{
		Container: re.podName,
		Command:   []string{command},
	}
	if stdin != nil {
		option.Stdin = true
		option.TTY = true
	}
	if stdout != nil {
		defer stdout.Close()
		option.Stdout = true
		option.Stderr = true
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	exec, err := remotecommand.NewSPDYExecutor(re.config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("can't create spdy executor: %w", err)
	}

	return exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Tty:    option.TTY,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stdout,
	})
}
