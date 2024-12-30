import { Component } from '@angular/core';
import { NgOptimizedImage } from '@angular/common';
import { RouterLink } from '@angular/router';
import { ButtonDirective } from 'primeng/button';

@Component({
    selector: 'bosun-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.scss'],
    imports: [NgOptimizedImage, RouterLink, ButtonDirective],
})
export class HomeComponent {}
