import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { CatalogsComponent } from './catalogs/catalogs.component';
import { ClusterComponent } from './cluster/cluster.component';
import { NamespacesComponent } from './cluster/grouped/namespaces.component';
import { InitiateDaemonSetsTracking, InitiateDeploymentsTracking, InitiatePodsTracking } from './store/route.resolver';
import { PodsComponent } from './cluster/workloads/pods/pods.component';
import { DeploymentsComponent } from './cluster/workloads/deployments/deployments.component';
import { DaemonSetsComponent } from './cluster/workloads/daemonSets/daemon-sets.component';

export const routes: Routes = [
    {
        path: '',
        redirectTo: '/home',
        pathMatch: 'full',
    },
    {
        path: 'home',
        component: HomeComponent,
    },
    {
        path: 'catalogs',
        component: CatalogsComponent,
    },
    {
        path: 'cluster/:id',
        component: ClusterComponent,
        children: [
            {
                path: 'pods',
                resolve: { resolveData: InitiatePodsTracking },
                component: PodsComponent,
            },
            {
                path: 'deployments',
                resolve: { resolveData: InitiateDeploymentsTracking },
                component: DeploymentsComponent,
            },
            {
                path: 'daemonsets',
                resolve: { resolveData: InitiateDaemonSetsTracking },
                component: DaemonSetsComponent,
            },
            {
                path: 'namespaces',
                component: NamespacesComponent,
            },
        ],
    },
    {
        path: '**',
        redirectTo: '/home',
    },
];
