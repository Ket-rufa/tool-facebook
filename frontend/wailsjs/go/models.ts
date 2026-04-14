export namespace accounts {
	
	export class AccountProfile {
	    id: string;
	    accountId: string;
	    displayName: string;
	    avatar: string;
	    accountType: string;
	    provider: string;
	    sessionId: string;
	    sessionStatus: string;
	    attachedAt: string;
	    lastCheckedAt: string;
	    isDefault: boolean;
	    note: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.accountId = source["accountId"];
	        this.displayName = source["displayName"];
	        this.avatar = source["avatar"];
	        this.accountType = source["accountType"];
	        this.provider = source["provider"];
	        this.sessionId = source["sessionId"];
	        this.sessionStatus = source["sessionStatus"];
	        this.attachedAt = source["attachedAt"];
	        this.lastCheckedAt = source["lastCheckedAt"];
	        this.isDefault = source["isDefault"];
	        this.note = source["note"];
	    }
	}
	export class AttachFlowStatusResponse {
	    flowId: string;
	    state: string;
	    message: string;
	    accountPreview?: AccountProfile;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new AttachFlowStatusResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.flowId = source["flowId"];
	        this.state = source["state"];
	        this.message = source["message"];
	        this.accountPreview = this.convertValues(source["accountPreview"], AccountProfile);
	        this.updatedAt = source["updatedAt"];
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

export namespace actiontest {
	
	export class ActionLog {
	    id: string;
	    action_type: string;
	    account_id: string;
	    account_display_name: string;
	    post_id: string;
	    reaction_type: string;
	    status: string;
	    dry_run: boolean;
	    message: string;
	    // Go type: time
	    executed_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ActionLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.action_type = source["action_type"];
	        this.account_id = source["account_id"];
	        this.account_display_name = source["account_display_name"];
	        this.post_id = source["post_id"];
	        this.reaction_type = source["reaction_type"];
	        this.status = source["status"];
	        this.dry_run = source["dry_run"];
	        this.message = source["message"];
	        this.executed_at = this.convertValues(source["executed_at"], null);
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
	export class LikePostRequest {
	    post_id: string;
	    account_id: string;
	    reaction_type: string;
	    dry_run: boolean;
	    actor_source: string;
	
	    static createFrom(source: any = {}) {
	        return new LikePostRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.post_id = source["post_id"];
	        this.account_id = source["account_id"];
	        this.reaction_type = source["reaction_type"];
	        this.dry_run = source["dry_run"];
	        this.actor_source = source["actor_source"];
	    }
	}
	export class LikePostResponse {
	    success: boolean;
	    status: string;
	    action: string;
	    post_id: string;
	    account_id: string;
	    account_display_name: string;
	    reaction_type: string;
	    dry_run: boolean;
	    message: string;
	    // Go type: time
	    executed_at: any;
	
	    static createFrom(source: any = {}) {
	        return new LikePostResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.status = source["status"];
	        this.action = source["action"];
	        this.post_id = source["post_id"];
	        this.account_id = source["account_id"];
	        this.account_display_name = source["account_display_name"];
	        this.reaction_type = source["reaction_type"];
	        this.dry_run = source["dry_run"];
	        this.message = source["message"];
	        this.executed_at = this.convertValues(source["executed_at"], null);
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
	export class SessionStatus {
	    is_active: boolean;
	    provider: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_active = source["is_active"];
	        this.provider = source["provider"];
	        this.message = source["message"];
	    }
	}

}

