import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { ButtonDirective } from 'primeng/button';

@Component({
    selector: 'bosun-root',
    imports: [RouterLink, RouterLinkActive, RouterOutlet, ButtonDirective],
    templateUrl: './app.component.html',
    styleUrl: './app.component.scss',
})
export class AppComponent {}
