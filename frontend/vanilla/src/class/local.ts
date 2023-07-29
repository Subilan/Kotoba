export class Local {
	public static getToken() {
		return localStorage.getItem('kotoba-login-token');
	}

	public static setToken(value: string) {
		localStorage.setItem('kotoba-login-token', value);
	}
}