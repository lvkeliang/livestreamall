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

        sendRequest("/user/login", `mail=${encodeURIComponent(loginMail)}&password=${encodeURIComponent(loginPassword)}`)
            .then(data => {
                if (data.status === 10000) {
                    const token = data.token; // 获取返回的 token
                    // 可以将 token 保存到本地存储或 sessionStorage 中
                    localStorage.setItem("token", token);
                    // 跳转到主页或执行其他操作
                    window.location.href = "http://47.92.137.133:9089/home";
                } else {
                    alert("登录失败：" + data.info);  // 显示后端返回的失败信息
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

        const requestBody = `mail=${encodeURIComponent(registerMail)}&password=${encodeURIComponent(registerPassword)}&nickname=${encodeURIComponent(registerNickname)}`;

        sendRequest("/user/register", requestBody)
            .then(data => {
                if (data.status === 10000) {
                    alert("注册成功！请登录");
                } else {
                    alert("注册失败：" + data.info);  // 显示后端返回的失败信息
                }
            })
            .catch(error => console.error("请求失败:", error));
    }

    // 发送请求的通用函数
    function sendRequest(url, body) {
        return fetch(`http://47.92.137.133:9089${url}`, {
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
