import { AfterViewInit, ChangeDetectionStrategy, Component, inject, OnDestroy, OnInit, viewChild } from '@angular/core';
import { AppStore } from '../../../store';
import { V1DaemonSet, V1Namespace } from '@kubernetes/client-node';
import { filterTableByNamespace, ResourceAgePipe, ResourceEvents } from '../../../utils';
import { EventsOn } from '../../../../../wailsjs/runtime';
import { IconField } from 'primeng/iconfield';
import { InputText } from 'primeng/inputtext';
import { Select, SelectChangeEvent } from 'primeng/select';
import { Table, TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { FormsModule } from '@angular/forms';

@Component({
    templateUrl: 'daemon-sets.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [IconField, InputText, ResourceAgePipe, Select, TableModule, Tag, FormsModule],
})
export class DaemonSetsComponent implements OnInit, AfterViewInit, OnDestroy {
    readonly store = inject(AppStore);
    dataTable = viewChild<Table>('dt');

    async ngOnInit() {
        await this.store.listDaemonSets(this.store.currentContext().id);
        EventsOn(ResourceEvents.DaemonSetUpdate, (data: V1DaemonSet) => {
            this.processDaemonSetEvent(ResourceEvents.DaemonSetUpdate, data);
        });
        EventsOn(ResourceEvents.DaemonSetDelete, (data: V1DaemonSet) => {
            this.processDaemonSetEvent(ResourceEvents.DaemonSetDelete, data);
        });
    }

    ngAfterViewInit() {
        const ns = this.store.selectedNamespace();
        if (ns) {
            filterTableByNamespace(ns, this.dataTable() as Table);
        }
    }

    setSelectedNamespace(daemonSet: V1DaemonSet, dt: Table) {
        const ns = this.store.namespacesEntities().find((n) => n.metadata?.name === daemonSet.metadata?.namespace);
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

    async ngOnDestroy() {
        await this.store.stopDaemonSetTracking(this.store.currentContext().id);
    }

    private processDaemonSetEvent(event: ResourceEvents, data: V1DaemonSet) {
        switch (event) {
            case ResourceEvents.DaemonSetUpdate:
                this.store.addUpdateDaemonSet(data);
                break;
            case ResourceEvents.DaemonSetDelete:
                this.store.deleteDaemonSet(data);
                break;
            default:
                console.error('Unknown event: ', event);
        }
    }
}
