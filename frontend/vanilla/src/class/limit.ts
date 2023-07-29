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