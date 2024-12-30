import { Injectable } from '@angular/core';
import {
    ListDaemonSets,
    ListDeployments,
    ListPods,
    StartDaemonSetsSync,
    StartDeploymentsSync,
    StartPodSync,
    StopDaemonSetsSync,
    StopDeploymentsSync,
    StopPodSync,
} from '../../../../wailsjs/go/pkg/App';
import { V1DaemonSet, V1Deployment, V1Pod } from '@kubernetes/client-node';

@Injectable({ providedIn: 'root' })
export class WorkloadService {
    async startPodSync(contextId: string) {
        try {
            await StartPodSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async stopPodSync(contextId: string) {
        try {
            await StopPodSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async listPods(contextId: string) {
        try {
            return (await ListPods(contextId)) as V1Pod[];
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async startDeploymentSync(contextId: string) {
        try {
            await StartDeploymentsSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async stopDeploymentSync(contextId: string) {
        try {
            await StopDeploymentsSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async listDeployments(contextId: string) {
        try {
            return (await ListDeployments(contextId)) as V1Deployment[];
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async startDaemonSync(contextId: string) {
        try {
            await StartDaemonSetsSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async stopDaemonSync(contextId: string) {
        try {
            await StopDaemonSetsSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async listDaemonSets(contextId: string) {
        try {
            return (await ListDaemonSets(contextId)) as V1DaemonSet[];
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }
}
