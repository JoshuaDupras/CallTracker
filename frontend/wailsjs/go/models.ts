export namespace db {
	
	export class Call {
	    id: number;
	    incident_number: string;
	    call_type: string;
	    mutual_aid: string;
	    address: string;
	    town: string;
	    location_notes: string;
	    // Go type: time
	    dispatched: any;
	    // Go type: time
	    enroute?: any;
	    // Go type: time
	    on_scene?: any;
	    // Go type: time
	    clear?: any;
	    narrative: string;
	    created_by: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Call(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.incident_number = source["incident_number"];
	        this.call_type = source["call_type"];
	        this.mutual_aid = source["mutual_aid"];
	        this.address = source["address"];
	        this.town = source["town"];
	        this.location_notes = source["location_notes"];
	        this.dispatched = this.convertValues(source["dispatched"], null);
	        this.enroute = this.convertValues(source["enroute"], null);
	        this.on_scene = this.convertValues(source["on_scene"], null);
	        this.clear = this.convertValues(source["clear"], null);
	        this.narrative = source["narrative"];
	        this.created_by = source["created_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Picklist {
	    id: number;
	    category: string;
	    value: string;
	    sort_order: number;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Picklist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.value = source["value"];
	        this.sort_order = source["sort_order"];
	        this.active = source["active"];
	    }
	}
	export class User {
	    id: number;
	    first_name: string;
	    last_name: string;
	    position: string;
	    ems_level: string;
	    is_admin: boolean;
	    pin?: string;
	    active: boolean;
	    // Go type: time
	    joined_date?: any;
	    // Go type: time
	    created: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.first_name = source["first_name"];
	        this.last_name = source["last_name"];
	        this.position = source["position"];
	        this.ems_level = source["ems_level"];
	        this.is_admin = source["is_admin"];
	        this.pin = source["pin"];
	        this.active = source["active"];
	        this.joined_date = this.convertValues(source["joined_date"], null);
	        this.created = this.convertValues(source["created"], null);
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

