import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { TableModule } from 'primeng/table';
import { DropdownModule } from 'primeng/dropdown';
import { CommonModule } from '@angular/common';
import { InputText } from 'primeng/inputtext';
import { IconField } from 'primeng/iconfield';
import { Avatar } from 'primeng/avatar';
import { AppStore } from '../store';
import { DialogService } from 'primeng/dynamicdialog';
import { LoginDialogComponent } from '../dialogs/login-dialog.component';
import { Bosun } from '../store/models';
import { Router } from '@angular/router';

@Component({
    selector: 'bosun-catalogs',
    templateUrl: './catalogs.component.html',
    styleUrl: './catalogs.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [TableModule, DropdownModule, CommonModule, InputText, IconField, Avatar],
    providers: [DialogService],
})
export class CatalogsComponent {
    readonly store = inject(AppStore);
    readonly router = inject(Router);

    constructor(private readonly dialogService: DialogService) {}

    initiateClusterLogin(data: Bosun): void {
        const ref = this.dialogService.open(LoginDialogComponent, {
            data: data,
            width: '50vw',
            modal: true,
            draggable: false,
            focusOnShow: false,
        });
        ref.onClose.subscribe((result: boolean) => {
            if (result) {
                this.router.navigate(['/cluster', data.id]).then();
            }
        });
    }
}
