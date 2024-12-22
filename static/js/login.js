// 获取当前页面的协议、主机名和端口，生成动态的 API 基础 URL
const baseURL = `${window.location.protocol}//${window.location.hostname}${window.location.port ? ':' + window.location.port : ''}`;

function hashPassword(password) {
    // 使用 SHA-256 哈希密码
    const hash = CryptoJS.SHA256(password);
    // 转换为十六进制字符串
    return hash.toString(CryptoJS.enc.Hex);
}

document.addEventListener("DOMContentLoaded", function() {
    // Select the card elements
    const loginCard = document.querySelector('.login');
    const registerCard = document.querySelector('.register');

    // Select the switch form link
    const switchForm = document.querySelector('#switch-form');

    // Add click event listener to the switch form link
    switchForm.addEventListener('click', () => {
        loginCard.classList.toggle('flip');
        registerCard.classList.toggle('flip');
    });

    // 登录按钮点击事件
    document.getElementById("login-submit").addEventListener("click", handleLogin);

    // 注册按钮点击事件
    document.getElementById("register-submit").addEventListener("click", handleRegister);

    // 登录处理函数
    function handleLogin() {
        const loginMail = document.getElementById("login-mail").value;
        const loginPassword = document.getElementById("login-password").value;

        if (!validateEmail(loginMail)) {
            alert("邮箱格式不正确！");
            return;
        }
        if (loginPassword.length === 0) {
            alert("密码不能为空！");
            return;
        }

        // 对密码进行加密
        const encryptedPassword = hashPassword(loginPassword); // 直接调用，不需要 .then()

        // 加密后的密码
        const requestBody = `mail=${encodeURIComponent(loginMail)}&password=${encodeURIComponent(encryptedPassword)}`;

        sendRequest("/user/login", requestBody)
            .then(data => {
                if (data.status === 10000) {
                    const token = data.token;
                    localStorage.setItem("token", token);
                    window.location.href = `${baseURL}/home`;
                } else {
                    alert("登录失败：" + data.info);
                }
            })
            .catch(error => console.error("请求失败:", error));
    }

    // 注册处理函数
    function handleRegister() {
        const registerMail = document.getElementById("register-mail").value;
        const registerPassword = document.getElementById("register-password").value;
        const registerPasswordConfirm = document.getElementById("register-password-confirm").value;
        const registerNickname = document.getElementById("register-nickname").value;

        if (!validateEmail(registerMail)) {
            alert("邮箱格式不正确！");
            return;
        }
        if (registerPassword.length < 6) {
            alert("密码不能小于6个字符！");
            return;
        }
        if (registerNickname.length >= 60) {
            alert("昵称不能超过60个字符！");
            return;
        }
        if (registerPassword !== registerPasswordConfirm) {
            alert("密码和确认密码不一致！");
            return;
        }

        // 对密码进行加密
        const encryptedPassword = hashPassword(registerPassword); // 直接调用，不需要 .then()

        const requestBody = `mail=${encodeURIComponent(registerMail)}&password=${encodeURIComponent(encryptedPassword)}&nickname=${encodeURIComponent(registerNickname)}`;

        sendRequest("/user/register", requestBody)
            .then(data => {
                if (data.status === 10000) {
                    // alert("注册成功！正在自动登录...");

                    // 自动登录，发送注册的邮箱和加密后的密码
                    handleAutoLogin(registerMail, encryptedPassword);
                } else {
                    alert("注册失败：" + data.info);
                }
            })
            .catch(error => console.error("请求失败:", error));
    }

    // 自动登录函数
    function handleAutoLogin(mail, encryptedPassword) {
        const requestBody = `mail=${encodeURIComponent(mail)}&password=${encodeURIComponent(encryptedPassword)}`;

        sendRequest("/user/login", requestBody)
            .then(data => {
                if (data.status === 10000) {
                    const token = data.token;
                    localStorage.setItem("token", token); // 保存 token
                    window.location.href = `${baseURL}/home`; // 跳转到主页
                } else {
                    alert("登录失败：" + data.info); // 登录失败时的提示
                }
            })
            .catch(error => console.error("自动登录请求失败:", error));
    }

    // 发送请求的通用函数
    function sendRequest(url, body) {
        return fetch(`${baseURL}${url}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: body
        })
            .then(response => {
                if (!response.ok) throw new Error("网络响应异常");
                return response.json();
            });
    }

    // 验证邮箱格式的函数
    function validateEmail(email) {
        const mailReg = /^\w+@[a-z0-9]+\.[a-z]+$/;
        return mailReg.test(email);
    }
});
