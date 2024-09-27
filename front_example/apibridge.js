class ApiBridge {
    constructor(request, baseUrl, logger) {
        this.request = request;
        this.baseUrl = baseUrl;
        this.logger = logger;
        this.error = "Invalid response from server";

        this.logger.info(`API bridge initialized with base URL: ${baseUrl}`);
    }

    async register(login, password) {
        // this.logger.info(`Registering user: ${login}...`);
        try {
            const response = await this.request(`${this.baseUrl}/api/v1/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ login, password }),
            });
            //if (!response.ok) throw new Error('Response error: ' + response.statusCode);
            this.logger.info(`Success registration for user:${login}`);
            return
        }catch(e){
            this.logger.error(`Failed to register user: ${e.message}`);
            throw new Error(e.message);
        }
    }

    async login(login, password) {
        // this.logger.info(`Logging in user: ${login}...`);
        try {
            const response = await this.request(`${this.baseUrl}/api/v1/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ login, password }),
            });
            const token = this.getHeaderValueFromResponse(response, 'X-Allow-Session');
            // if (!response.ok) throw new Error('Response error: ' + response.statusCode);
            this.logger.info(`Success login for user:${login} with token:${token}`);
            return token;
        }catch(e){
            this.logger.error(`Failed to login user: ${e.message}`);
            throw new Error(e.message);
        }
    }
        
    // async refresh(token) {
    //     const response = await this.request(`${this.baseUrl}/api/v1/refresh`, {
    //         method: 'POST',
    //         headers: {
    //             'X-Allow-Session': `${token}`,
    //         },
    //     });
    //     if (!response.ok) {
    //         throw new Error(this.error);
    //     }
    //     newToken = response.headers.get(this.headerToken);
    //     this.logger.info(`Refreshing token from:${token} to: ${newToken}`);
    //     return newToken;
    // }
        
    async logout(token) {
        this.logger.info(`Logging out user...`);
        try {
            await this.request(`${this.baseUrl}/api/v1/logout`, {
                method: 'POST',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
        } catch (e) {
            this.logger.info(`Failed to logout user: ${e.message}`);
        }
    }

    async getInfoSelf(token) {
        try {
            const response = await this.request(`${this.baseUrl}/api/v1/user/self`, {
                method: 'GET',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
            this.logger.info(`Getting info self for user with token:${token}`);
            return JSON.parse(response.body);
        } catch (e) {
            this.logger.info(`Failed to get info self: ${e.message}`);
            throw new Error(e.message);
        }
    }
    async getInfoChats(token, offset, limit) {
        this.logger.info(`Getting chats with token: ${token} offset: ${offset} and limit: ${limit}...`);
        try {
            if (offset === undefined || limit === undefined) {
                offset = 0;
                limit = 10;
            }
            const response = await this.request(`${this.baseUrl}/api/v1/chat/self?offset=${offset}&limit=${limit}`, {
                method: 'GET',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
            this.logger.info(`Success getting chats with offset: ${offset} and limit: ${limit} body: ${response.body}`);
            return JSON.parse(response.body);
        } catch(e) {
            this.logger.error(`Failed to get chats: ${e.message}`);
            throw new Error(e.message);
        }
    }
    
    async createChat(token, chatDesc) {
		this.logger.info(`Creating chat: ${chatDesc} token: ${token}`);
        try {
            const resp = await this.request(`${this.baseUrl}/api/v1/chat/${chatDesc}`, {
                method: 'POST',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
            this.logger.info(`Success creating chat: ${chatDesc} with token: ${token}`);
            // return this.getHeaderValueFromResponse(resp, 'Location');
            return;
        } catch (e) {
            this.logger.error(`Failed to create chat: ${e.message}`);
            throw new Error(e.message);
        }
    }

    async getChatFollowers(token, chatDesc) {
        try {
            const response = await this.request(`${this.baseUrl}/api/v1/chat/${chatDesc}/followers`, {
                method: 'GET',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
            this.logger.info(`Getting followers for chat: ${chatDesc} with token: ${token} body: ${response.body}`);
            return JSON.parse(response.body);
        } catch (e) {
            this.logger.error(`Failed to get followers: ${e.message}`);
            throw new Error(e.message);
        }
    }
    async followChat(token, chatDesc) {
        try {
            await this.request(`${this.baseUrl}/api/v1/chat/${chatDesc}/follow`, {
                method: 'POST',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
            this.logger.info(`Following chat: ${chatDesc} with token: ${token}`);
            return
        } catch (e) {
            this.logger.error(`Failed to follow chat: ${e.message}`);
            throw new Error(e.message);
        }
    }
    
    async sendMessage(token, chatDesc, text) {
        try {
            this.logger.info(`Sending message to chat: ${chatDesc} with token: ${token}`);
            if (text === undefined) {
                throw new Error('Message is required');
            }
            if (text.length > 1000) {
                throw new Error('Message exceeds maximum length');
            }
            if (text.match(/<.*?>/)) {
                throw new Error('Message contains HTML tags');
            }
            if (text.match(/[A-Za-z0-9]+@[A-Za-z0-9]+\.[A-Za-z0-9]+/)) {
                throw new Error('Message contains email addresses');
            }
            if (text.match(/[0-9]{10,}/)) {
                throw new Error('Message contains phone numbers');
            }
            await this.request(`${this.baseUrl}/api/v1/chat/${chatDesc}/message`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Allow-Session': `${token}`,
                },
                body: JSON.stringify({ text }),
            });
            return;
        } catch (e) {
            this.logger.error(`Failed to send message: ${e.message}`);
            throw new Error(e.message);
        }
    }
    async getMessages(token, chatDesc, offset, limit) {
        try {
            this.logger.info(`Getting messages for chat: ${chatDesc} with token: ${token}`);
            if (offset === undefined || limit === undefined) {
                offset = 0;
                limit = 10;
            }
            const response = await this.request(`${this.baseUrl}/api/v1/chat/${chatDesc}/message?offset=${offset}&limit=${limit}`, {
                method: 'GET',
                headers: {
                    'X-Allow-Session': `${token}`,
                },
            });
            this.logger.info(`Success getting messages for chat: ${chatDesc} with token: ${token} body: ${response.body}`);
            return JSON.parse(response.body);
        } catch (e) {
            this.logger.error(`Failed to get messages: ${e.message}`);
            throw new Error(e.message);
        }
    }

    getHeaderValueFromResponse(response, strHeader) {
        const indexHeader = response.rawHeaders.indexOf(strHeader);
        if (indexHeader === -1) {
            throw new Error(`Header "${strHeader}" is missing from response`);
        }
        const headerValue = response.rawHeaders[indexHeader + 1];
        // this.logger.info(`Extracted header ${strHeader}: ${headerValue}`);
        if (headerValue === undefined) {
            throw new Error(`Value for header "${strHeader}" is missing from response`);
        }
        return headerValue;
    };
}

module.exports = ApiBridge;