import { ChangeDetectionStrategy, Component, inject, OnDestroy, OnInit } from '@angular/core';
import { AppStore } from '../store';
import { Accordion, AccordionContent, AccordionHeader, AccordionPanel } from 'primeng/accordion';
import { FormsModule } from '@angular/forms';
import { Listbox, ListboxChangeEvent } from 'primeng/listbox';
import { ActivatedRoute, Router, RouterOutlet } from '@angular/router';
import { pkg } from '../../../wailsjs/go/models';
import { SelectItem, SelectItemGroup } from 'primeng/api';
import { EventsOn } from '../../../wailsjs/runtime';
import { V1Namespace } from '@kubernetes/client-node';
import { ResourceEvents } from '../utils';
import Resource = pkg.Resource;

@Component({
    templateUrl: './cluster.component.html',
    styleUrl: './cluster.component.scss',
    changeDetection: ChangeDetectionStrategy.Default,
    imports: [Accordion, AccordionPanel, AccordionHeader, AccordionContent, FormsModule, Listbox, RouterOutlet],
})
export class ClusterComponent implements OnInit, OnDestroy {
    readonly store = inject(AppStore);
    readonly router = inject(Router);
    readonly route = inject(ActivatedRoute);

    async ngOnInit() {
        EventsOn(ResourceEvents.NamespaceUpdate, (data: V1Namespace) => {
            this.processNamespaceEvent(ResourceEvents.NamespaceUpdate, data);
        });
        EventsOn(ResourceEvents.NamespaceDelete, (data: V1Namespace) => {
            console.log('delete event', data);
            this.processNamespaceEvent(ResourceEvents.NamespaceDelete, data);
        });
    }

    async switchView(event: ListboxChangeEvent) {
        const resource = event.value as SelectItem<Resource>;
        if (resource && resource.label) {
            await this.router.navigate([resource.label.toLowerCase()], { relativeTo: this.route });
        }
    }

    async switchGroupView(group: SelectItemGroup<Resource>): Promise<void> {
        if (group.items.length === 0) {
            await this.router.navigate([group.label.toLowerCase()], { relativeTo: this.route });
        }
    }

    ngOnDestroy() {
        this.store.stopNamespaceTracking(this.store.currentContext().id);
    }

    private processNamespaceEvent(event: ResourceEvents, data: V1Namespace): void {
        switch (event) {
            case ResourceEvents.NamespaceUpdate:
                this.store.addUpdateNamespace(data);
                break;
            case ResourceEvents.NamespaceDelete:
                this.store.deleteNamespace(data);
                break;
            default:
                console.log('Unknown event: ', event);
        }
    }
}
