package controller

import (
	"testing"

	flaggerv1 "github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestController_verifyCanary(t *testing.T) {
	tests := []struct {
		name    string
		canary  flaggerv1.Canary
		wantErr bool
	}{
		{
			name: "Gloo upstream in a different namespace should return an error",
			canary: flaggerv1.Canary{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cd-1",
					Namespace: "default",
				},
				Spec: flaggerv1.CanarySpec{
					UpstreamRef: &flaggerv1.CrossNamespaceObjectReference{
						Name:      "upstream",
						Namespace: "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Gloo upstream in the same namespace is allowed",
			canary: flaggerv1.Canary{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cd-1",
					Namespace: "default",
				},
				Spec: flaggerv1.CanarySpec{
					UpstreamRef: &flaggerv1.CrossNamespaceObjectReference{
						Name:      "upstream",
						Namespace: "default",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "MetricTemplate in a different namespace should return an error",
			canary: flaggerv1.Canary{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cd-1",
					Namespace: "default",
				},
				Spec: flaggerv1.CanarySpec{
					Analysis: &flaggerv1.CanaryAnalysis{
						Metrics: []flaggerv1.CanaryMetric{
							{
								TemplateRef: &flaggerv1.CrossNamespaceObjectReference{
									Name:      "mt-1",
									Namespace: "test",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "AlertProvider in a different namespace should return an error",
			canary: flaggerv1.Canary{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cd-1",
					Namespace: "default",
				},
				Spec: flaggerv1.CanarySpec{
					Analysis: &flaggerv1.CanaryAnalysis{
						Alerts: []flaggerv1.CanaryAlert{
							{
								ProviderRef: flaggerv1.CrossNamespaceObjectReference{
									Name:      "ap-1",
									Namespace: "test",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	ctrl := &Controller{
		noCrossNamespaceRefs: true,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ctrl.verifyCanary(&test.canary)
			if test.wantErr {
				require.Error(t, err)
			}
		})
	}
}
