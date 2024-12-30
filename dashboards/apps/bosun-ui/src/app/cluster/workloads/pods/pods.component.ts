import { AfterViewInit, ChangeDetectionStrategy, Component, inject, OnDestroy, OnInit, viewChild } from '@angular/core';
import { AppStore } from '../../../store';
import { IconField } from 'primeng/iconfield';
import { InputText } from 'primeng/inputtext';
import { Table, TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { V1Namespace, V1Pod } from '@kubernetes/client-node';
import { EventsOn } from '../../../../../wailsjs/runtime';
import { filterTableByNamespace, ResourceAgePipe, ResourceEvents } from '../../../utils';
import { Select, SelectChangeEvent } from 'primeng/select';
import { FormsModule } from '@angular/forms';

@Component({
    templateUrl: 'pods.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [ResourceAgePipe, IconField, InputText, TableModule, Tag, Select, FormsModule],
})
export class PodsComponent implements OnInit, AfterViewInit, OnDestroy {
    readonly store = inject(AppStore);
    dataTable = viewChild<Table>('dt');

    async ngOnInit() {
        await this.store.listPods(this.store.currentContext().id);
        EventsOn(ResourceEvents.PodUpdate, (data: V1Pod) => {
            this.processPodEvent(ResourceEvents.PodUpdate, data);
        });
        EventsOn(ResourceEvents.PodDelete, (data: V1Pod) => {
            this.processPodEvent(ResourceEvents.PodDelete, data);
        });
    }

    ngAfterViewInit() {
        const ns = this.store.selectedNamespace();
        if (ns) {
            filterTableByNamespace(ns, this.dataTable() as Table);
        }
    }

    setSelectedNamespace(pod: V1Pod, dt: Table) {
        const ns = this.store.namespacesEntities().find((n) => n.metadata?.name === pod.metadata?.namespace);
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

    getStatusSeverity(pod: V1Pod) {
        if (pod?.status) {
            switch (pod.status.phase) {
                case 'Running':
                    return 'success';
                case 'Succeeded':
                    return 'success';
                case 'Pending':
                    return 'warn';
                case 'Failed':
                    return 'danger';
                case 'Terminating':
                    return 'secondary';
                default:
                    return 'contrast';
            }
        }
        return null;
    }

    async ngOnDestroy() {
        await this.store.stopPodTracking(this.store.currentContext().id);
    }

    private processPodEvent(event: ResourceEvents, data: V1Pod) {
        switch (event) {
            case ResourceEvents.PodUpdate:
                this.store.addUpdatePod(data);
                break;
            case ResourceEvents.PodDelete:
                this.store.deletePod(data);
                break;
            default:
                console.log('Unknown event: ', event);
        }
    }
}
