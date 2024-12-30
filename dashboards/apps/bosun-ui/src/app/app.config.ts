import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter, withHashLocation } from '@angular/router';

import { routes } from './app.routes';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { providePrimeNG } from 'primeng/config';
import { AppPreset } from './app.preset';

export const appConfig: ApplicationConfig = {
    providers: [
        provideZoneChangeDetection({ eventCoalescing: true }),
        provideRouter(routes, withHashLocation()),
        provideAnimationsAsync(),
        providePrimeNG({
            ripple: true,
            theme: {
                preset: AppPreset,
                options: {
                    darkModeSelector: '.bosun-dark',
                },
            },
        }),
    ],
};
