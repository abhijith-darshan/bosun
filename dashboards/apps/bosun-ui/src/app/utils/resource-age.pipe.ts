import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'resourceAge',
})
export class ResourceAgePipe implements PipeTransform {
    transform(value: string): string {
        const creationDate = new Date(value);

        // get the current date
        const now = new Date();

        // calc the difference in ms
        const diffInMs = now.getTime() - creationDate.getTime();

        // handle case where the resource was created less than a second ago
        if (diffInMs < 1000) {
            return `${Math.floor(diffInMs / 100)}s`; // less than 1 second
        }

        // convert ms to s, m, h, and d
        const seconds = Math.floor(diffInMs / 1000);
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);
        const years = Math.floor(days / 365);

        // select time format
        if (years > 2) {
            const remainingDays = days % 365;
            return `${years}y${remainingDays}d`;
        } else if (days > 0) {
            return `${days}d`;
        } else if (hours > 0) {
            const remainingMinutes = minutes % 60;
            return `${hours}h${remainingMinutes}m`;
        } else if (minutes > 0) {
            return `${minutes}m`;
        } else {
            return `${seconds}s`;
        }
    }
}
