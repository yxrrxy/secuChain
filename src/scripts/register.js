document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const messageEl = document.getElementById('message');
    const togglePasswordBtn = document.querySelector('.toggle-password');
    const passwordInput = document.getElementById('password');
    const body = document.body;

    // 表单验证规则
    const validationRules = {
        username: {
            pattern: /^[a-zA-Z0-9_\u4e00-\u9fa5]{3,50}$/,
            message: '用户名必须是3-50个字符，只能包含字母、数字、下划线和中文'
        },
        password: {
            pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{6,32}$/,
            message: '密码必须包含大小写字母和数字，长度在6-32个字符之间'
        },
        email: {
            pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
            message: '请输入有效的邮箱地址'
        }
    };

    if (togglePasswordBtn) {
        togglePasswordBtn.addEventListener('click', () => {
            const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
            passwordInput.setAttribute('type', type);
            togglePasswordBtn.querySelector('i').classList.toggle('fa-eye');
            togglePasswordBtn.querySelector('i').classList.toggle('fa-eye-slash');

            if (type === 'text') {
                enableDarkMode();
            } else {
                disableDarkMode();
            }
        });
    }

    function enableDarkMode() {
        body.classList.add('dark-mode');
    }

    function disableDarkMode() {
        body.classList.remove('dark-mode');
    }

    // 表单提交处理
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const email = document.getElementById('email').value;
        
        // 表单验证
        if (!validateForm()) {
            return;
        }
        
        // 显示加载状态
        const submitBtn = e.target.querySelector('button[type="submit"]');
        setLoadingState(submitBtn, true);
        
        try {
            const response = await fetch('/api/v1/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password, email })
            });
            
            const data = await response.json();
            
            handleResponse(response.status, data);
        } catch (error) {
            console.error('注册错误:', error);
            showMessage('网络错误，请检查连接', 'error');
        } finally {
            setLoadingState(submitBtn, false);
        }
    });

    // 字段验证函数
    function validateField(field) {
        const input = document.getElementById(field);
        const value = input.value;
        const rule = validationRules[field];
        
        if (!rule) return true;
        
        const isValid = rule.pattern.test(value);
        if (!isValid) {
            showMessage(rule.message, 'error');
            input.focus();
            return false;
        }
        
        return true;
    }

    // 完整表单验证
    function validateForm() {
        return ['username', 'password', 'email'].every(validateField);
    }

    // 响应处理函数
    function handleResponse(status, data) {
        switch (status) {
            case 200:
            case 201:
                showMessage('注册成功！正在跳转到登录页面...', 'success');
                setTimeout(() => {
                    window.location.href = 'login.html';
                }, 1500);
                break;
                
            case 400:
                showMessage(data.message || '请检查输入信息是否正确', 'error');
                break;
                
            case 409:
                showMessage('用户名或邮箱已被注册', 'error');
                break;
                
            case 500:
                showMessage('服务器错误，请稍后重试', 'error');
                break;
                
            default:
                showMessage(data.message || '注册失败，请重试', 'error');
        }
    }

    // 显示消息提示
    function showMessage(text, type = 'error') {
        messageEl.textContent = text;
        messageEl.className = `message ${type}`;
        
        if (type === 'success') {
            setTimeout(() => {
                messageEl.textContent = '';
                messageEl.className = 'message';
            }, 3000);
        }
    }

    // 设置按钮加载状态
    function setLoadingState(button, isLoading) {
        if (isLoading) {
            button.disabled = true;
            button.innerHTML = '<i class="fas fa-spinner fa-spin"></i> 注册中...';
        } else {
            button.disabled = false;
            button.innerHTML = '<span>注册</span><i class="fas fa-user-plus"></i>';
        }
    }

    // 实时表单验证
    document.querySelectorAll('input').forEach(input => {
        input.addEventListener('blur', () => {
            validateField(input.id);
        });
        
        input.addEventListener('input', () => {
            messageEl.textContent = '';
            messageEl.className = 'message';
        });
    });
});