import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { AppStore } from '../../store';
import { IconField } from 'primeng/iconfield';
import { InputText } from 'primeng/inputtext';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { ResourceAgePipe } from '../../utils';

@Component({
    templateUrl: './namespaces.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [IconField, InputText, TableModule, Tag, ResourceAgePipe],
})
export class NamespacesComponent {
    readonly store = inject(AppStore);
}
