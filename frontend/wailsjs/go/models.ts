export namespace config {
	
	export class PrivacyMeta {
	    displayVersion: number;
	    updateTime: string;
	    author: string;
	    desc: string;
	
	    static createFrom(source: any = {}) {
	        return new PrivacyMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayVersion = source["displayVersion"];
	        this.updateTime = source["updateTime"];
	        this.author = source["author"];
	        this.desc = source["desc"];
	    }
	}
	export class PrivacyFields {
	    sectionStorage: boolean;
	    sectionDevice: boolean;
	    sectionOnline: boolean;
	    sectionBehavior: boolean;
	    sectionPrivacy: boolean;
	    storageDefault: boolean;
	    standardFolder: boolean;
	    storageCustom: boolean;
	    adminPrivilege: boolean;
	    identityVerification: boolean;
	    modifySystemSettings: boolean;
	    inputMonitor: boolean;
	    microphone: boolean;
	    camera: boolean;
	    location: boolean;
	    sensor: boolean;
	    clipboard: boolean;
	    screenshot: boolean;
	    otherDevice: boolean;
	    accessOtherAppData: boolean;
	    manageOtherApps: boolean;
	    network: boolean;
	    thirdPartyService: boolean;
	    cloudAI: boolean;
	    personalization: boolean;
	    notification: boolean;
	    update: boolean;
	    autoStart: boolean;
	    contextMenu: boolean;
	    trayIcon: boolean;
	    shortcut: boolean;
	    openOtherApp: boolean;
	    personalInfo: boolean;
	    plaintextStorage: boolean;
	    analytics: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PrivacyFields(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sectionStorage = source["sectionStorage"];
	        this.sectionDevice = source["sectionDevice"];
	        this.sectionOnline = source["sectionOnline"];
	        this.sectionBehavior = source["sectionBehavior"];
	        this.sectionPrivacy = source["sectionPrivacy"];
	        this.storageDefault = source["storageDefault"];
	        this.standardFolder = source["standardFolder"];
	        this.storageCustom = source["storageCustom"];
	        this.adminPrivilege = source["adminPrivilege"];
	        this.identityVerification = source["identityVerification"];
	        this.modifySystemSettings = source["modifySystemSettings"];
	        this.inputMonitor = source["inputMonitor"];
	        this.microphone = source["microphone"];
	        this.camera = source["camera"];
	        this.location = source["location"];
	        this.sensor = source["sensor"];
	        this.clipboard = source["clipboard"];
	        this.screenshot = source["screenshot"];
	        this.otherDevice = source["otherDevice"];
	        this.accessOtherAppData = source["accessOtherAppData"];
	        this.manageOtherApps = source["manageOtherApps"];
	        this.network = source["network"];
	        this.thirdPartyService = source["thirdPartyService"];
	        this.cloudAI = source["cloudAI"];
	        this.personalization = source["personalization"];
	        this.notification = source["notification"];
	        this.update = source["update"];
	        this.autoStart = source["autoStart"];
	        this.contextMenu = source["contextMenu"];
	        this.trayIcon = source["trayIcon"];
	        this.shortcut = source["shortcut"];
	        this.openOtherApp = source["openOtherApp"];
	        this.personalInfo = source["personalInfo"];
	        this.plaintextStorage = source["plaintextStorage"];
	        this.analytics = source["analytics"];
	    }
	}
	export class PrivacyConfig {
	    options: PrivacyFields;
	    display: PrivacyFields;
	    meta: PrivacyMeta;
	
	    static createFrom(source: any = {}) {
	        return new PrivacyConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.options = this.convertValues(source["options"], PrivacyFields);
	        this.display = this.convertValues(source["display"], PrivacyFields);
	        this.meta = this.convertValues(source["meta"], PrivacyMeta);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class GeneralConfig {
	    autoStart: boolean;
	    notifications: boolean;
	    minimizeToTray: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GeneralConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.autoStart = source["autoStart"];
	        this.notifications = source["notifications"];
	        this.minimizeToTray = source["minimizeToTray"];
	    }
	}
	export class Config {
	    version: string;
	    firstLaunch: boolean;
	    language: string;
	    general: GeneralConfig;
	    privacy: PrivacyConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.firstLaunch = source["firstLaunch"];
	        this.language = source["language"];
	        this.general = this.convertValues(source["general"], GeneralConfig);
	        this.privacy = this.convertValues(source["privacy"], PrivacyConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	
	
	

}

