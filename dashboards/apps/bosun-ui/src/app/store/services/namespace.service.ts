import { Injectable } from '@angular/core';
import { ListNamespaces, StartNamespaceSync, StopNamespaceSync } from '../../../../wailsjs/go/pkg/App';
import { V1Namespace } from '@kubernetes/client-node';

@Injectable({ providedIn: 'root' })
export class NamespaceService {
    async startNamespaceSync(contextId: string): Promise<void> {
        try {
            await StartNamespaceSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async stopNamespaceSync(contextId: string): Promise<void> {
        try {
            await StopNamespaceSync(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async listNamespaces(contextId: string): Promise<V1Namespace[]> {
        try {
            return (await ListNamespaces(contextId)) as V1Namespace[];
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }
}
