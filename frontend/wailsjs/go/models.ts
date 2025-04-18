export namespace main {
	
	export class FrontendSettings {
	    locale: string;
	    days: string[];
	    replacementKey: string;
	    dateFormat: string;
	
	    static createFrom(source: any = {}) {
	        return new FrontendSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.locale = source["locale"];
	        this.days = source["days"];
	        this.replacementKey = source["replacementKey"];
	        this.dateFormat = source["dateFormat"];
	    }
	}
	export class FrontendTemplate {
	    SVG: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new FrontendTemplate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SVG = source["SVG"];
	        this.name = source["name"];
	    }
	}
	export class GenerateThumbnailsJob {
	    from: string;
	    to: string;
	    time: string;
	    customDates: string[];
	
	    static createFrom(source: any = {}) {
	        return new GenerateThumbnailsJob(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	        this.time = source["time"];
	        this.customDates = source["customDates"];
	    }
	}

}

