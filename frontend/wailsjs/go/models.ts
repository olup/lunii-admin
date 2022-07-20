export namespace lunii {
	
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
	    }
	}
	export class Metadata {
	    uuid: number[];
	    ref: string;
	    title: string;
	    description: string;
	    isOfficialPack: boolean;
	    isUnknown: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Metadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.ref = source["ref"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.isOfficialPack = source["isOfficialPack"];
	        this.isUnknown = source["isUnknown"];
	    }
	}

}

