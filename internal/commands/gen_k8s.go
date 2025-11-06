package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/riskiramdan/ShELDon/internal/textutil"
)

// NewGenK8sCommand generates Kubernetes manifests tailored for Go services.
func NewGenK8sCommand(deps Dependencies) *cobra.Command {
	var (
		app         string
		port        int
		cpuReq      string
		memReq      string
		cpuLim      string
		memLim      string
		cpuTarget   int
		minReplicas int
		maxReplicas int
		out         string
		model       string
	)

	cmd := &cobra.Command{
		Use:   "gen-k8s",
		Short: "Generate a K8s Deployment + HPA YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			if app == "" {
				return errors.New("--app is required")
			}

			prompt := fmt.Sprintf("Write a Kubernetes Deployment + HPA for a Go API. App=%s.\nConstraints: containerPort %d, requests %s/%s, limits %s/%s, HPA on CPU %d%%, min %d max %d. Include liveness/readiness on /healthz. Return only YAML.",
				app, port, cpuReq, memReq, cpuLim, memLim, cpuTarget, minReplicas, maxReplicas,
			)

			ctx, cancel := context.WithTimeout(cmd.Context(), deps.Config.Timeout)
			defer cancel()

			ans, err := deps.LLM.Generate(ctx, textutil.Choose(model, deps.Config.ModelGeneral), prompt)
			if err != nil {
				return err
			}

			if out == "" {
				out = "k8s.yaml"
			}
			return deps.Files.WriteFile(out, ans)
		},
	}

	cmd.Flags().StringVar(&app, "app", "", "App name/label")
	cmd.Flags().IntVar(&port, "port", 8080, "Container port")
	cmd.Flags().StringVar(&cpuReq, "cpu-req", "100m", "CPU request")
	cmd.Flags().StringVar(&memReq, "mem-req", "128Mi", "Memory request")
	cmd.Flags().StringVar(&cpuLim, "cpu-lim", "500m", "CPU limit")
	cmd.Flags().StringVar(&memLim, "mem-lim", "512Mi", "Memory limit")
	cmd.Flags().IntVar(&cpuTarget, "cpu-target", 60, "HPA CPU target percent")
	cmd.Flags().IntVar(&minReplicas, "min", 2, "Min replicas")
	cmd.Flags().IntVar(&maxReplicas, "max", 10, "Max replicas")
	cmd.Flags().StringVar(&out, "out", "", "Output YAML file (default k8s.yaml)")
	cmd.Flags().StringVar(&model, "model", "", "Override model")
	return cmd
}
