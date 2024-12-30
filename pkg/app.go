package pkg

import (
	"bosun/pkg/internal"
	"bosun/pkg/internal/cluster"
	"bosun/pkg/internal/workloads"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"slices"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
	"strings"
	"unicode"
)

// App struct
type App struct {
	ctx           context.Context
	bosunClusters map[string]BosunCluster
}

type BosunCluster struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Version     string `json:"version"`
	client      client.Client
	clientSet   *kubernetes.Clientset
	namespace   *cluster.BosunNamespace
	pod         *workloads.BosunPod
	deployment  *workloads.BosunDeployment
	daemonSet   *workloads.BosunDaemonSet
	statefulSet *workloads.BosunStatefulSet
	replicaSet  *workloads.BosunReplicaSet
	job         *workloads.BosunJob
	cronJob     *workloads.BosunCronJob
}

type Resource struct {
	Key          string   `json:"key"`
	Kind         string   `json:"kind"`
	Name         string   `json:"name"`
	Namespaced   bool     `json:"namespaced"`
	ShortNames   []string `json:"shortNames"`
	SingularName string   `json:"singularName"`
	PluralName   string   `json:"pluralName"`
	DisplayName  string   `json:"displayName"`
	Verbs        []string `json:"verbs"`
	Version      string   `json:"version"`
	Group        string   `json:"group"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Start is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Start(ctx context.Context) {
	a.ctx = ctx
	err := a.prepareKubeConfigs()
	if err != nil {
		writeLogs(err.Error())
		os.Exit(1)
	}
}

func (a *App) prepareKubeConfigs() error {
	home := homedir.HomeDir()
	defaultConfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.LoadFromFile(defaultConfig)
	if err != nil {
		writeLogs(err.Error())
		return err
	}
	a.bosunClusters = make(map[string]BosunCluster)
	for contextKey, configContext := range config.Contexts {
		if strings.TrimSpace(configContext.Cluster) == "" || strings.TrimSpace(configContext.AuthInfo) == "" {
			continue
		}
		cl, cs, err := a.getKubeClients(config, contextKey)
		if err != nil {
			writeLogs(err.Error())
			return err
		}
		clusterID := uuid.New().String()
		a.bosunClusters[clusterID] = BosunCluster{
			ID:          clusterID,
			ShortName:   generateShortName(contextKey),
			Name:        contextKey,
			client:      cl,
			clientSet:   cs,
			namespace:   cluster.NewBosunNamespace(cl, cs),
			pod:         workloads.NewBosunPod(cl, cs),
			deployment:  workloads.NewBosunDeployment(cl, cs),
			daemonSet:   workloads.NewBosunDaemonSet(cl, cs),
			statefulSet: workloads.NewBosunStatefulSet(cl, cs),
			replicaSet:  workloads.NewBosunReplicaSet(cl, cs),
			job:         workloads.NewBosunJob(cl, cs),
			cronJob:     workloads.NewBosunCronJob(cl, cs),
		}
	}
	return nil
}

func (a *App) GetKnownResources(contextID string) ([]Resource, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	cl := kCtx.clientSet
	// get all known resources
	apiResources, err := cl.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	groupedResources := make([]Resource, 0)

	// Iterate through the fetched API resources and populate the map
	for _, group := range apiResources {
		for _, resource := range group.APIResources {
			// Extract the groupKey using a helper method
			groupKey := determineGroupKey(group.GroupVersion, internal.BosunKind(resource.Kind))

			// If a valid group key is found, append the resource to the corresponding group
			if groupKey != "" {
				// create a Resource struct for each valid API resource
				resourceData := Resource{
					Key:          groupKey,
					Kind:         resource.Kind,
					Name:         resource.Name,
					Namespaced:   resource.Namespaced,
					ShortNames:   resource.ShortNames,
					SingularName: resource.SingularName,
					PluralName:   resource.Name,
					DisplayName:  internal.GetDisplayName(internal.BosunKind(resource.Kind)),
					Verbs:        resource.Verbs,
					Version:      group.APIVersion,
					Group:        group.GroupVersion,
				}
				groupedResources = append(groupedResources, resourceData)
			}
		}
	}
	slices.SortFunc(groupedResources, func(i, j Resource) int {
		return strings.Compare(i.Kind, j.Kind)
	})
	return groupedResources, nil
}

// determineGroupKey extracts the groupKey based on the group version and resource kind
func determineGroupKey(groupVersion string, kind internal.BosunKind) string {
	switch groupVersion {
	case internal.GroupApps, internal.GroupBatch:
		if kind == internal.KindDeployment || kind == internal.KindDaemonSet ||
			kind == internal.KindReplicaSet || kind == internal.KindStatefulSet ||
			kind == internal.KindJob || kind == internal.KindCronJob {
			return internal.CategoryWorkload
		}
	case internal.GroupV1, internal.GroupAdmissionRegistration:
		if kind == internal.KindSecret || kind == internal.KindConfigMap || kind == internal.KindMutatingWebhookConfiguration || kind == internal.KindValidatingWebhookConfiguration {
			return internal.CategoryConfig
		}
		if kind == internal.KindPod {
			return internal.CategoryWorkload
		}
		if kind == internal.KindNamespaces {
			return internal.CategoryNamespaces
		}
		if kind == internal.KindNodes {
			return internal.CategoryNodes
		}
	}
	// Return an empty string for kinds that don't match, so they won't be added to the map
	return ""
}

func (a *App) GetVersion(contextID string) (string, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return "", errors.New("context not found")
	}
	info, err := kCtx.clientSet.DiscoveryClient.ServerVersion()
	if err != nil {
		writeLogs(err.Error())
		return "", err
	}
	kCtx.Version = info.String()
	return info.String(), nil
}

func (a *App) getKubeClients(config *api.Config, contextName string) (client.Client, *kubernetes.Clientset, error) {
	restConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{CurrentContext: contextName}).ClientConfig()
	if err != nil {
		writeLogs(err.Error())
		return nil, nil, err
	}
	restConfig.Timeout = 30 * time.Second
	k8sClient, err := client.New(restConfig, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		writeLogs(err.Error())
		return nil, nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		writeLogs(err.Error())
		return nil, nil, err
	}
	return k8sClient, clientSet, nil
}

func (a *App) ReadKubeConfigs() ([]BosunCluster, error) {
	return a.kubeConfigToKubeContexts(), nil
}

func (a *App) kubeConfigToKubeContexts() []BosunCluster {
	var clusters []BosunCluster
	for k, ctx := range a.bosunClusters {
		runtime.LogPrintf(a.ctx, "bosun cluster ID %s", k)
		clusters = append(clusters, BosunCluster{ID: k, Name: ctx.Name, ShortName: ctx.ShortName})
	}
	// Sort the contexts by Name in ascending order
	sort.Slice(clusters, func(i, j int) bool {
		// Prioritize contexts starting with "kind"
		if strings.HasPrefix(clusters[i].Name, "kind") && !strings.HasPrefix(clusters[j].Name, "kind") {
			return true
		}
		if !strings.HasPrefix(clusters[i].Name, "kind") && strings.HasPrefix(clusters[j].Name, "kind") {
			return false
		}
		return clusters[i].Name < clusters[j].Name
	})
	return clusters
}

func generateShortName(word string) string {
	splitWords := strings.Split(word, "-")
	var shortName string
	for _, w := range splitWords {
		if len(w) > 0 {
			if len(shortName) == 3 {
				break
			}
			shortName += string(unicode.ToUpper(rune(w[0])))
		}
	}
	return shortName
}

func writeLogs(content string) {
	home := homedir.HomeDir()
	t := time.Now().Unix()
	logsPath := filepath.Join(home, ".bosun", fmt.Sprintf("app_logs_%d.log", t))
	file, _ := os.Create(logsPath)
	defer file.Close()
	_, _ = file.WriteString(content)
}
