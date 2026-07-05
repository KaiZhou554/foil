export namespace builder {
	
	export class BuildResult {
	    APKPath: string;
	    PackageName: string;
	    VersionName: string;
	    VersionCode: number;
	    Log: string;
	
	    static createFrom(source: any = {}) {
	        return new BuildResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APKPath = source["APKPath"];
	        this.PackageName = source["PackageName"];
	        this.VersionName = source["VersionName"];
	        this.VersionCode = source["VersionCode"];
	        this.Log = source["Log"];
	    }
	}

}

export namespace config {
	
	export class Config {
	    version: string;
	    firstLaunch: boolean;
	    language: string;
	    outputDir: string;
	    showFloatButton: boolean;
	    openAfterBuild: boolean;
	    useCustomCert: boolean;
	    rememberLevel: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.firstLaunch = source["firstLaunch"];
	        this.language = source["language"];
	        this.outputDir = source["outputDir"];
	        this.showFloatButton = source["showFloatButton"];
	        this.openAfterBuild = source["openAfterBuild"];
	        this.useCustomCert = source["useCustomCert"];
	        this.rememberLevel = source["rememberLevel"];
	    }
	}

}

