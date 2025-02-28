package datasciencepipelines

import (
	"encoding/json"
	"fmt"
	"path"

	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"

	"github.com/opendatahub-io/opendatahub-operator/v2/apis/common"
	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/apis/components/v1alpha1"
	"github.com/opendatahub-io/opendatahub-operator/v2/controllers/status"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
	odhdeploy "github.com/opendatahub-io/opendatahub-operator/v2/pkg/deploy"
)

const (
	ArgoWorkflowCRD = "workflows.argoproj.io"
	ComponentName   = componentApi.DataSciencePipelinesComponentName

	ReadyConditionType = conditionsv1.ConditionType(componentApi.DataSciencePipelinesKind + status.ReadySuffix)

	// LegacyComponentName is the name of the component that is assigned to deployments
	// via Kustomize. Since a deployment selector is immutable, we can't upgrade existing
	// deployment to the new component name, so keep it around till we figure out a solution.
	LegacyComponentName = "data-science-pipelines-operator"

	managedPipelineParamsKey = "MANAGEDPIPELINES"
	platformVersionParamsKey = "PLATFORMVERSION"
)

var (
	imageParamMap = map[string]string{
		"IMAGES_DSPO":                    "RELATED_IMAGE_ODH_DATA_SCIENCE_PIPELINES_OPERATOR_CONTROLLER_IMAGE",
		"IMAGES_APISERVER":               "RELATED_IMAGE_ODH_ML_PIPELINES_API_SERVER_V2_IMAGE",
		"IMAGES_PERSISTENCEAGENT":        "RELATED_IMAGE_ODH_ML_PIPELINES_PERSISTENCEAGENT_V2_IMAGE",
		"IMAGES_SCHEDULEDWORKFLOW":       "RELATED_IMAGE_ODH_ML_PIPELINES_SCHEDULEDWORKFLOW_V2_IMAGE",
		"IMAGES_ARGO_EXEC":               "RELATED_IMAGE_ODH_DATA_SCIENCE_PIPELINES_ARGO_ARGOEXEC_IMAGE",
		"IMAGES_ARGO_WORKFLOWCONTROLLER": "RELATED_IMAGE_ODH_DATA_SCIENCE_PIPELINES_ARGO_WORKFLOWCONTROLLER_IMAGE",
		"IMAGES_DRIVER":                  "RELATED_IMAGE_ODH_ML_PIPELINES_DRIVER_IMAGE",
		"IMAGES_LAUNCHER":                "RELATED_IMAGE_ODH_ML_PIPELINES_LAUNCHER_IMAGE",
		"IMAGES_MLMDGRPC":                "RELATED_IMAGE_ODH_MLMD_GRPC_SERVER_IMAGE",
		"IMAGES_PIPELINESRUNTIMEGENERIC": "RELATED_IMAGE_ODH_ML_PIPELINES_RUNTIME_GENERIC_IMAGE",
	}

	overlaysSourcePaths = map[common.Platform]string{
		cluster.SelfManagedRhoai: "overlays/rhoai",
		cluster.ManagedRhoai:     "overlays/rhoai",
		cluster.OpenDataHub:      "overlays/odh",
		cluster.Unknown:          "overlays/odh",
	}

	paramsPath = path.Join(odhdeploy.DefaultManifestPath, ComponentName, "base")
)

func manifestPath(p common.Platform) types.ManifestInfo {
	return types.ManifestInfo{
		Path:       odhdeploy.DefaultManifestPath,
		ContextDir: ComponentName,
		SourcePath: overlaysSourcePaths[p],
	}
}

func computeParamsMap(rr *types.ReconciliationRequest) (map[string]string, error) {
	dsp, ok := rr.Instance.(*componentApi.DataSciencePipelines)
	if !ok {
		return nil, fmt.Errorf("resource instance %v is not a componentApi.DataSciencePipelines", rr.Instance)
	}

	data, err := json.Marshal(dsp.Spec.PreloadedPipelines)
	if err != nil {
		return nil, fmt.Errorf("marshalling preloaded pipelines failed: %w", err)
	}

	data, err = json.Marshal(string(data))
	if err != nil {
		return nil, fmt.Errorf("marshalling preloaded pipelines failed: %w", err)
	}

	extraParamsMap := map[string]string{
		managedPipelineParamsKey: string(data),
	}

	return extraParamsMap, nil
}
