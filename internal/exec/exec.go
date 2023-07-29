package exec

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	k8sExec "k8s.io/kubectl/pkg/cmd/exec"
)

func ExecCmdExampleV2(podName, namespace, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	fmt.Println("Running...")
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
	execOpts := k8sExec.ExecOptions{
		Command:          []string{command},
		EnforceNamespace: true,
		StreamOptions: k8sExec.StreamOptions{
			PodName:   podName,
			Namespace: namespace,
			// ContainerName: podName,
			TTY:   true,
			Stdin: true,
			IOStreams: genericclioptions.IOStreams{
				In:     stdin,
				Out:    stdout,
				ErrOut: stdout,
			},
		},
		Config:    config,
		PodClient: clientset.CoreV1(),
		Executor:  &k8sExec.DefaultRemoteExecutor{},
	}
	Pod, err := execOpts.PodClient.Pods(execOpts.Namespace).Get(context.TODO(), execOpts.PodName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("can't get pod: %w", err)
	}
	fmt.Println(Pod.Name, Pod.Namespace)
	return execOpts.Run()
}

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

	// config.Insecure = true

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
	// exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	// if err != nil {
	// 	return fmt.Errorf("can't create spdy executor: %w", err)
	// }
	// err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
	// 	Tty:    true,
	// 	Stdin:  stdin,
	// 	Stdout: stdout,
	// 	Stderr: stderr,
	// })
	e := k8sExec.DefaultRemoteExecutor{}
	e.Execute("POST", req.URL(), config, stdin, stdout, stderr, false, nil)
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
