// API请求封装
class API {
    static async request(url, options = {}) {
        // 检查令牌是否过期
        const expires = localStorage.getItem('token_expires');
        if (expires && Date.now() > parseInt(expires)) {
            // 尝试刷新令牌
            const refreshed = await refreshToken();
            if (!refreshed) {
                // 刷新失败，重定向到登录页
                window.location.href = '/login.html';
                return;
            }
        }

        // 添加认证头
        const token = localStorage.getItem('token');
        if (token) {
            options.headers = {
                ...options.headers,
                'Authorization': `Bearer ${token}`
            };
        }

        const response = await fetch(url, options);
        
        // 处理401错误
        if (response.status === 401) {
            // 尝试刷新令牌
            const refreshed = await refreshToken();
            if (refreshed) {
                // 重试请求
                return API.request(url, options);
            } else {
                // 重定向到登录页
                window.location.href = '/login.html';
                return;
            }
        }

        return response;
    }

    // 便捷方法
    static async get(url, options = {}) {
        return API.request(url, { ...options, method: 'GET' });
    }

    static async post(url, data, options = {}) {
        return API.request(url, {
            ...options,
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            body: JSON.stringify(data)
        });
    }
} 