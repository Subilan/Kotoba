import config from '../config';

interface Payload {
	username: string;
	avatar: string;
	website: string;
}

type ConfigName = keyof typeof config;

function el() {
	return document.getElementById('kotoba') as HTMLDivElement;
}

class Local {
	public static getToken() {
		return localStorage.getItem('kotoba-login-token');
	}

	public static setToken(value: string) {
		localStorage.setItem('kotoba-login-token', value);
	}
}

function getPayloadFromToken(token: string): Payload {
	return JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString());
}

class Limit {
	public target: any;
	public defaultValue: any;

	constructor(input: any) {
		this.target = input;
	}

	public default(value: any) {
		this.defaultValue = value;
		return this;
	}

	public setDefault() {
		this.target = this.defaultValue;
	}

	public mustType(type: 'string' | 'number' | 'bigint' | 'boolean' | 'symbol' | 'undefined' | 'object' | 'function') {
		if (typeof this.target != type) this.setDefault();
		return this;
	}

	public min(num: number) {
		if (this.target < num) this.setDefault();
		return this;
	}

	public max(num: number) {
		if (this.target > num) this.setDefault();
		return this;
	}

	public collect() {
		return this.target;
	}
}

function c(configName: ConfigName) {
    const lim = new Limit(config[configName]);
    switch (configName) {
        case 'commentAllowGuest':
            return lim.default(false).mustType('boolean').collect();
        case 'commentTextareaPlaceholder':
            return lim.default('在此输入评论').mustType('string').collect();
        case 'commentTextareaMaxlength':
            return lim.default(500).mustType('number').collect();
        case 'commentTextareaMinlength':
            return lim.default(0).mustType('number').collect();
        case 'commentSameIPLimit':
            return lim.default(1).mustType('number').collect();
        case 'commentShowGuestIP':
            return lim.default(false).mustType('boolean').collect();
        case 'commentShowIPOrigin':
            return lim.default(false).mustType('boolean').collect();
    }
}

async function main() {
	const token = Local.getToken();
	let template = `[[KOTOBA_TEMPLATE_HTML]]`;

	if (!token) {
		console.warn('Token not found.');
	}

	const payload = getPayloadFromToken(token);

	template = template
    .replace('macro:user-avatar', payload.avatar)
    .replace('macro:input-placeholder', c("commentTextareaPlaceholder"))
    .replace('macro:input-maxlength', c("commentTextareaMaxlength"))
    .replace('macro:input-minlength', c("commentTextareaMinlength"));
}

main();
