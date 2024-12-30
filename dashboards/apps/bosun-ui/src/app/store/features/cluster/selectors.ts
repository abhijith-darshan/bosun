import { entityConfig, SelectEntityId } from '@ngrx/signals/entities';
import { V1Namespace } from '@kubernetes/client-node';
import { type } from '@ngrx/signals';

const namespaceSelect: SelectEntityId<V1Namespace> = (namespace) => {
    if (namespace.metadata?.name) {
        return namespace.metadata.name;
    }
    return '';
};

export const NamespaceEntityConfig = entityConfig({
    entity: type<V1Namespace>(),
    collection: 'namespaces',
    selectId: namespaceSelect,
});
