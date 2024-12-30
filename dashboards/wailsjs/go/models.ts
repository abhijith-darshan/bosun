export namespace pkg {
	
	export class BosunCluster {
	    id: string;
	    name: string;
	    shortName: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new BosunCluster(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.shortName = source["shortName"];
	        this.version = source["version"];
	    }
	}
	export class Resource {
	    key: string;
	    kind: string;
	    name: string;
	    namespaced: boolean;
	    shortNames: string[];
	    singularName: string;
	    pluralName: string;
	    displayName: string;
	    verbs: string[];
	    version: string;
	    group: string;
	
	    static createFrom(source: any = {}) {
	        return new Resource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.kind = source["kind"];
	        this.name = source["name"];
	        this.namespaced = source["namespaced"];
	        this.shortNames = source["shortNames"];
	        this.singularName = source["singularName"];
	        this.pluralName = source["pluralName"];
	        this.displayName = source["displayName"];
	        this.verbs = source["verbs"];
	        this.version = source["version"];
	        this.group = source["group"];
	    }
	}

}

