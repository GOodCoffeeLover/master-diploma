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
func ExecCmdExample(podName string, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	fmt.Println("Executing...")

	config := &restclient.Config{
		Host:        "https://192.168.49.2:8443",
		APIPath:     "/api",
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6ImdmeDRGUU1LM0ExdFBpUTloWERrRjdLcjJac0Z0UVhxMmxGbTF3dlcyQTQifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNjg5OTQxMDE4LCJpYXQiOjE2ODk5Mzc0MTgsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImFkbWluLXVzZXIiLCJ1aWQiOiJjNjViYTI1NC1hMzQ5LTRmZWItODJkZi0wNzJlNjdiYjYyMzIifX0sIm5iZiI6MTY4OTkzNzQxOCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6YWRtaW4tdXNlciJ9.l2O9vuTARa5cdFDILInfLrOTgrQGaPn6drGKu8YtxotgtOpWi_RFLPj-C1Pwap9cezt2bTQpQbVbuOu5_LshMBaF4VG6E6j5H843_8AEUrBrZ_9XegL48MX3wOKhpVcWaMpY812QnDUIz-zg6FgMyUuZ31s3Vq_rYqo7vrl_3N906zGdmx2sdJYaaLYiqKwLe1lZoUoqfMb7DPM1IhzjTTt3Lyqh0Ak8BnCzznWiolKuniS7H_X2pU_QI9Ini-vMvF34AjK4OmgviJR708wsYEDsSITwhn-yq2BziDWy1UgmCS4foix0wdgnMrcnEeJ2uoMYUVFsE7lzKpxh3_xFBQ",
	}
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})
	config.Insecure = true

	client, err := restclient.RESTClientFor(config)
	if err != nil {
		return fmt.Errorf("can't create client: %w", err)
	}
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
