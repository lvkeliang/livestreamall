const baseURL = `${window.location.protocol}//${window.location.hostname}${window.location.port ? ':' + window.location.port : ''}`;

const backendApi = {
    startLive: `${baseURL}/live/start`,       // 接口：生成推流密钥
    streamInfo: (streamName) => `${baseURL}/stream/${streamName}` // 接口：获取拉流信息
};

// Utility function to make authenticated fetch requests
async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem("token");
    if (!token) {
        alert("你还未登录, 请先登录!");
        throw new Error("No token found in localStorage");
    }

    const headers = {
        ...options.headers,
        "Authorization": `${token}`,
        "Content-Type": "application/json"
    };

    const response = await fetch(url, {
        ...options,
        headers
    });

    if (response.status === 401) {
        alert("你的登录已经过期, 请再次登录!");
        throw new Error("Unauthorized");
    }

    return response;
}

const liveForm = document.getElementById("liveForm");
const streamInfo = document.getElementById("streamInfo");
const pushUrlElem = document.getElementById("pushUrl");
const streamKeyElem = document.getElementById("streamKey");
const hlsUrlElem = document.getElementById("hlsUrl");
const dashUrlElem = document.getElementById("dashUrl");
const playUrlElem = document.getElementById("playUrl");

liveForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const title = document.getElementById("title").value;
    const description = document.getElementById("description").value;

    try {
        // Step 1: 请求后端生成推流密钥和流名称
        const publishResponse = await authenticatedFetch(backendApi.startLive, {
            method: "POST",
            body: JSON.stringify({ title, description })
        });

        if (!publishResponse.ok) {
            alert("获取串流密钥失败: " + (await publishResponse.text()));
            return;
        }

        const publishData = await publishResponse.json();
        const pushStreamUrl = publishData.push_url;
        const token = publishData.stream_key;

        // 显示推流地址和密钥
        pushUrlElem.textContent = pushStreamUrl;
        streamKeyElem.textContent = token;

        // Step 2: 请求后端获取拉流信息
        const streamInfoResponse = await authenticatedFetch(backendApi.streamInfo(token));
        if (!streamInfoResponse.ok) {
            alert("获取流信息失败: " + (await streamInfoResponse.text()));
            return;
        }

        const streamData = await streamInfoResponse.json();

        // 展示拉流信息
        hlsUrlElem.textContent = streamData.hls_url;
        dashUrlElem.textContent = streamData.dash_url;
        playUrlElem.textContent = streamData.play_url;

        streamInfo.style.display = "block";
    } catch (error) {
        console.error("Error:", error);
    }
});
