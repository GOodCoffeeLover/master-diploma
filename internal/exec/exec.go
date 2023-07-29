package exec

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	k8sExec "k8s.io/kubectl/pkg/cmd/exec"
)

func ExecCmdExampleV2(podName string, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	fmt.Println("Running...")
	execOpts := k8sExec.ExecOptions{}
	execOpts.PodName = podName
	execOpts.Namespace = "default"
	execOpts.Command = []string{"sh", "-c", command}
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		return fmt.Errorf("can't create kube config: %w", err)
	}
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("can't create client: %w", err)
	}

	execOpts.Config = config
	execOpts.PodClient = clientset.CoreV1()
	execOpts.Executor = &k8sExec.DefaultRemoteExecutor{}
	execOpts.TTY = true
	execOpts.Stdin = true
	execOpts.In = stdin
	execOpts.Out = stdout
	execOpts.ErrOut = stderr
	Pod, err := execOpts.PodClient.Pods(execOpts.Namespace).Get(context.TODO(), execOpts.PodName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("can't get pod: %w", err)
	}
	fmt.Println(Pod.Name, Pod.Namespace)
	return execOpts.Run()
}

// ExecCmd exec command on specific pod and wait the command's output.
func ExecCmdExample(podName string, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

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

	// config.GroupVersion = &v1.SchemeGroupVersion
	// config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})
	// config.Insecure = true

	client := clientset.CoreV1().RESTClient()
	cmd := []string{
		"sh",
		"-c",
		command,
	}
	req := client.Post().Namespace("default").Resource("pods").Name(podName).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	fmt.Println(req.URL())
	if err != nil {
		return fmt.Errorf("can't create spdy executor: %w", err)
	}
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
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
