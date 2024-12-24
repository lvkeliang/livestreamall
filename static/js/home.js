// 获取当前页面的协议、主机名和端口，生成动态的 API 基础 URL
const baseURL = `${window.location.protocol}//${window.location.hostname}${window.location.port ? ':' + window.location.port : ''}`;

async function fetchLiveRooms() {
    const token = localStorage.getItem('token');

    const response = await fetch(`${baseURL}/live/live_rooms`, {
        method: 'GET',
        headers: { 'Authorization': `${token}` }
    });

    if (!response.ok) {
        console.error('Failed to fetch live rooms');
        return;
    }

    const data = await response.json();
    const liveRooms = data.live_rooms;

    const liveRoomsContainer = document.getElementById('live-rooms');
    liveRoomsContainer.innerHTML = '';

    liveRooms.forEach(room => {
        const roomCard = document.createElement('div');
        roomCard.className = 'room-card';
        roomCard.innerHTML = `
        <h3>${room.Title}</h3>
        <p>${room.Description}</p>
        <button onclick="window.location.href='${baseURL}/live/play?stream_id=${room.StreamName}'">观看</button>
      `;
        liveRoomsContainer.appendChild(roomCard);
    });
}

function handleAuthState() {
    const token = localStorage.getItem('token');
    const loginBtn = document.getElementById('login-btn');
    const logoutBtn = document.getElementById('logout-btn');
    const userInfo = document.getElementById('user-info');

    if (token) {
        loginBtn.style.display = 'none';
        logoutBtn.style.display = 'inline-block';
        userInfo.style.display = 'inline-block';
        fetchUserInfo(token);
    } else {
        loginBtn.style.display = 'inline-block';
        logoutBtn.style.display = 'none';
        userInfo.style.display = 'none';
    }
}

async function fetchUserInfo(token) {
    const response = await fetch(`${baseURL}/user/info`, {
        method: 'GET',
        headers: { 'Authorization': `${token}` }
    });

    if (!response.ok) {
        console.error('Failed to fetch user info');
        return;
    }

    const data = await response.json();
    document.getElementById('user-name').textContent = `Hello, ${data.user.nickname}`;
}

async function startLive() {
    const token = localStorage.getItem('token');
    if (!token) {
        alert('你还未登录, 请先登录!');
        window.location.href = `${baseURL}/user/login`;
        return;
    }
    window.open(`${baseURL}/live/start`, '_blank');
}

function logout() {
    localStorage.removeItem('token');
    handleAuthState();
    window.location.reload();
}

window.onload = () => {
    handleAuthState();
    fetchLiveRooms();
};