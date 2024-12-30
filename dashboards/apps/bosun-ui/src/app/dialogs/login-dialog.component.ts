import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { AppStore } from '../store';
import { ProgressSpinner } from 'primeng/progressspinner';
import { DynamicDialogConfig, DynamicDialogRef } from 'primeng/dynamicdialog';
import { Button } from 'primeng/button';
import { Bosun } from '../store/models';

@Component({
    templateUrl: './login-dialog.component.html',
    styleUrl: './login-dialog.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
    imports: [ProgressSpinner, Button],
})
export class LoginDialogComponent implements OnInit {
    readonly store = inject(AppStore);

    constructor(
        private readonly dialogConfig: DynamicDialogConfig,
        private readonly ref: DynamicDialogRef,
    ) {}

    async ngOnInit() {
        await this.login();
    }

    async login() {
        const data = this.dialogConfig.data as Bosun;
        await this.store.loginToCluster(data, this.ref);
    }

    cancel() {
        this.ref.close(false);
    }
}
