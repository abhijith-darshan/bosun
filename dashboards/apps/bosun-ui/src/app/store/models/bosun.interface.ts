import { pkg } from '../../../../wailsjs/go/models';
import { SelectItemGroup } from 'primeng/api';
import { V1Namespace } from '@kubernetes/client-node';
import Resource = pkg.Resource;

export interface BosunAppState {
    currentContext: Bosun;
    clusters: Bosun[];
    isClustersLoading: boolean;
    clustersLoadingErrored: boolean;
    clustersLoadingError?: string;
    isClusterLoginLoading: boolean;
    clusterLoginErrored: boolean;
    clusterLoginError: string;
}

export interface BosunNamespaceState {
    selectedNamespace: V1Namespace | null;
}

export interface Bosun extends pkg.BosunCluster {
    knownResources: SelectItemGroup<Resource>[];
}
