const baseURL = `${window.location.origin}`; // 设置后端接口基础 URL
const token = localStorage.getItem('token'); // 从 localStorage 获取 Token
let username = '匿名用户'; // 默认用户名

// 获取 URL 参数
function getURLParameter(name) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
}

const streamId = getURLParameter('stream_id');
if (!streamId) {
    alert("缺少 stream_id 参数！");
    throw new Error("Missing stream_id");
}

// 获取用户信息
async function fetchUserInfo(token) {
    if (!token) return;
    const response = await fetch(`${baseURL}/user/info`, {
        method: 'GET',
        headers: { 'Authorization': `${token}` },
    });
    if (response.ok) {
        const data = await response.json();
        username = data.user.nickname;
        document.getElementById('user-name').textContent = `Hello, ${username}`;
    }
}

// 加载直播间信息和 HLS 地址
fetch(`/live/live-room/${streamId}`)
    .then((res) => res.json())
    .then((data) => {
        document.getElementById('live-title').textContent = data.title || "直播间";
        document.getElementById('live-description').textContent = data.description || "暂无描述";
        return fetch(`/stream/${streamId}`);
    })
    .then((res) => res.json())
    .then((data) => {
        const player = videojs('hls-player', { autoplay: true, muted: true, playsinline: true });
        player.src({ src: data.hls_url, type: 'application/x-mpegURL' });
    })
    .catch(console.error);

// 初始化 WebSocket
const ws = new WebSocket(`ws://${window.location.host}/live/ws/${streamId}`);
const chatMessages = document.getElementById('chat-messages');

let isScrolledToBottom = true; // 记录当前chat-messages是否滚动到最底部

const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

// 接收消息
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    const messageElement = document.createElement('p');
    messageElement.textContent = `${message.username}: ${message.content}`;
    messageElement.className = message.username === username ? 'me' : 'other';
    chatMessages.appendChild(messageElement);

    // 判断是否需要滚动到最底部
    if (isScrolledToBottom) {
        chatMessages.scrollTop = chatMessages.scrollHeight; // 滚动到最底部
    }
};

// 监听滚动事件
chatMessages.addEventListener('scroll', () => {
    // 检查滚动条是否在最底部
    isScrolledToBottom = chatMessages.scrollTop + chatMessages.clientHeight >= chatMessages.scrollHeight - 5;
});

// 发送消息
chatSend.addEventListener('click', () => {
    if (chatInput.value.trim() !== "") {
        // 发送的消息包括 token 和内容
        if (!token) {
            alert("你还未登录, 请先登录!");
            throw new Error("No token found in localStorage");
        }
        const message = {
            token: token,  // 需要验证的 token
            content: chatInput.value.trim()
        };

        ws.send(JSON.stringify(message));  // 发送消息
        chatInput.value = "";
    }
});

chatInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') chatSend.click();
});

function adjustChatSectionHeight() {
    const windowWidth = window.innerWidth; // 窗口内部宽度
    const videoElement = document.querySelector('.video-js');
    const chatSection = document.querySelector('.chat-section');

    if (windowWidth < 768) {
        const windowHeight = window.innerHeight;
        const videoHeight = videoElement.clientHeight;

        console.log("windowHeight: ", windowHeight)
        console.log("videoHeight: ", videoHeight)
        chatSection.style.maxHeight = `${windowHeight - videoHeight - 145}px`;
        chatSection.style.minHeight = `${windowHeight - videoHeight - 145}px`;
    }
    else {
        if (!document.fullscreenElement && videoElement && chatSection) {
            // 获取video-js的实际高度
            const videoHeight = videoElement.clientHeight;
            // 将chat-section的最大高度设置为video-js的高度
            chatSection.style.maxHeight = `${videoHeight}px`;
            chatSection.style.minHeight = `${videoHeight}px`;
        }
    }

}

// 窗口改变大小事件监听器
window.addEventListener('resize', adjustChatSectionHeight);

// 初始调整
adjustChatSectionHeight();

fetchUserInfo(token);