export namespace lunii {
	
	export class DiskUsage {
	    free: number;
	    used: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskUsage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.free = source["free"];
	        this.used = source["used"];
	        this.total = source["total"];
	    }
	}
	export class Device {
	    mountPoint: string;
	    uuid: number[];
	    uuidHex: string;
	    specificKey: number[];
	    serialNumber: string;
	    firmwareVersionMajor: number;
	    firmwareVersionMinor: number;
	    sdCardSize: number;
	    sdCardUsed: number;
	    diskUsage?: DiskUsage;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mountPoint = source["mountPoint"];
	        this.uuid = source["uuid"];
	        this.uuidHex = source["uuidHex"];
	        this.specificKey = source["specificKey"];
	        this.serialNumber = source["serialNumber"];
	        this.firmwareVersionMajor = source["firmwareVersionMajor"];
	        this.firmwareVersionMinor = source["firmwareVersionMinor"];
	        this.sdCardSize = source["sdCardSize"];
	        this.sdCardUsed = source["sdCardUsed"];
	        this.diskUsage = this.convertValues(source["diskUsage"], DiskUsage);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Job {
	    isComplete: boolean;
	    totalImages: number;
	    totalAudios: number;
	    imagesDone: number;
	    audiosDone: number;
	    hasError: string;
	    initDone: boolean;
	    binGenerationDone: boolean;
	    unpackDone: boolean;
	    imagesConversionDone: boolean;
	    audiosConversionDone: boolean;
	    metadataDone: boolean;
	    copyingDone: boolean;
	    indexDone: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Job(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isComplete = source["isComplete"];
	        this.totalImages = source["totalImages"];
	        this.totalAudios = source["totalAudios"];
	        this.imagesDone = source["imagesDone"];
	        this.audiosDone = source["audiosDone"];
	        this.hasError = source["hasError"];
	        this.initDone = source["initDone"];
	        this.binGenerationDone = source["binGenerationDone"];
	        this.unpackDone = source["unpackDone"];
	        this.imagesConversionDone = source["imagesConversionDone"];
	        this.audiosConversionDone = source["audiosConversionDone"];
	        this.metadataDone = source["metadataDone"];
	        this.copyingDone = source["copyingDone"];
	        this.indexDone = source["indexDone"];
	    }
	}
	export class Metadata {
	    uuid: number[];
	    ref: string;
	    title: string;
	    description: string;
	    packType: string;
	
	    static createFrom(source: any = {}) {
	        return new Metadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.ref = source["ref"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.packType = source["packType"];
	    }
	}

}

export namespace main {
	
	export class CheckUpdateResponse {
	    canUpdate: boolean;
	    latestVersion: string;
	    releaseNotes: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckUpdateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.canUpdate = source["canUpdate"];
	        this.latestVersion = source["latestVersion"];
	        this.releaseNotes = source["releaseNotes"];
	    }
	}
	export class Infos {
	    version: string;
	    machineId: string;
	    os: string;
	    arch: string;
	
	    static createFrom(source: any = {}) {
	        return new Infos(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.machineId = source["machineId"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	    }
	}

}

