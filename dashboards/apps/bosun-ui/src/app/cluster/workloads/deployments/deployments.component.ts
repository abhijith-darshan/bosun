import { AfterViewInit, ChangeDetectionStrategy, Component, inject, OnDestroy, OnInit, viewChild } from '@angular/core';
import { IconField } from 'primeng/iconfield';
import { InputText } from 'primeng/inputtext';
import { filterTableByNamespace, ResourceAgePipe, ResourceEvents } from '../../../utils';
import { Select, SelectChangeEvent } from 'primeng/select';
import { Table, TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { AppStore } from '../../../store';
import { EventsOn } from '../../../../../wailsjs/runtime';
import { V1Deployment, V1Namespace } from '@kubernetes/client-node';
import { FormsModule } from '@angular/forms';

@Component({
    templateUrl: 'deployments.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [IconField, InputText, ResourceAgePipe, Select, TableModule, Tag, FormsModule],
})
export class DeploymentsComponent implements OnInit, AfterViewInit, OnDestroy {
    readonly store = inject(AppStore);
    dataTable = viewChild<Table>('dt');

    async ngOnInit() {
        await this.store.listDeployments(this.store.currentContext().id);
        EventsOn(ResourceEvents.DeploymentUpdate, (data: V1Deployment) => {
            this.processDeploymentEvent(ResourceEvents.DeploymentUpdate, data);
        });
        EventsOn(ResourceEvents.DeploymentDelete, (data: V1Deployment) => {
            this.processDeploymentEvent(ResourceEvents.DeploymentDelete, data);
        });
    }

    ngAfterViewInit() {
        const ns = this.store.selectedNamespace();
        if (ns) {
            filterTableByNamespace(ns, this.dataTable() as Table);
        }
    }

    setSelectedNamespace(deployment: V1Deployment, dt: Table) {
        const ns = this.store.namespacesEntities().find((n) => n.metadata?.name === deployment.metadata?.namespace);
        if (ns) {
            filterTableByNamespace(ns, dt);
            this.store.setSelectedNamespace(ns);
        }
    }

    filterByNamespace(event: SelectChangeEvent, dt: Table) {
        const ns = event.value as V1Namespace;
        filterTableByNamespace(ns, dt);
        this.store.setSelectedNamespace(ns);
    }

    getStatusSeverity(statusType: string) {
        switch (statusType) {
            case 'Available':
                return 'success';
            case 'Progressing':
                return 'info';
            case 'ReplicaFailure':
                return 'danger';
            default:
                return null;
        }
    }

    async ngOnDestroy() {
        await this.store.stopDeploymentTracking(this.store.currentContext().id);
    }

    private processDeploymentEvent(event: ResourceEvents, data: V1Deployment) {
        switch (event) {
            case ResourceEvents.DeploymentUpdate:
                this.store.addUpdateDeployment(data);
                break;
            case ResourceEvents.DeploymentDelete:
                this.store.deleteDeployment(data);
                break;
            default:
                console.log('Unknown event: ', event);
        }
    }
}
