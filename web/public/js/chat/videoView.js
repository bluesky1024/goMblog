(function () {
    //前端弹幕展示相关
    class Barrage {
        constructor(id) {
            this.domList = [];
            this.dom = document.querySelector('#' + id);
            if (this.dom.style.position == '' || this.dom.style.position == 'static') {
                this.dom.style.position = 'relative';
            }
            this.dom.style.overflow = 'hidden';
            var rect = this.dom.getBoundingClientRect();
            this.domWidth = rect.right - rect.left;
            this.domHeight = rect.bottom - rect.top;
        }

        shoot(text) {
            var div = document.createElement('div');
            div.style.position = 'absolute';
            div.style.left = this.domWidth + 'px';
            div.style.top = (this.domHeight - 20) * +Math.random().toFixed(2) + 'px';
            div.style.whiteSpace = 'nowrap';
            div.style.color = '#' + Math.floor(Math.random() * 256).toString(10);
            div.innerText = text;
            this.dom.appendChild(div);

            var roll = (timer) => {
                var now = +new Date();
                roll.last = roll.last || now;
                roll.timer = roll.timer || timer;
                var left = div.offsetLeft;
                var rect = div.getBoundingClientRect();
                if (left < (rect.left - rect.right)) {
                    this.dom.removeChild(div);
                } else {
                    if (now - roll.last >= roll.timer) {
                        roll.last = now;
                        left -= 3;
                        div.style.left = left + 'px';
                    }
                    requestAnimationFrame(roll);
                }
            }
            roll(50 * +Math.random().toFixed(2));
        }

    }
    function appendList(userName,text) {
        var p = document.createElement('p');
        p.innerText = userName + ":" + text;
        document.querySelector('#content-text').appendChild(p);
    }
    var barage = new Barrage('content');

    // document.querySelector('#send').onclick = () => {
    //
    // };

    // const textList = ['弹幕', '666', '233333333', 'javascript', 'html', 'css', '前端框架', 'Vue', 'React', 'Angular',
    //     '测试弹幕效果'
    // ];
    // textList.forEach((s) => {
    //     barage.shoot(s);
    //     appendList(s);
    // })


    //后端websocket交互相关
    var socket;
    var isConnect = false;
    createWebsocket();

    function createWebsocket(){
        var scheme = document.location.protocol == "https:" ? "wss" : "ws";
        socket = new Ws(scheme + "://localhost:8080/chat/barrage/websocket/"+roomNO);

        //建立连接，通知后端，新客户进入房间，房间编号也需要发到后端
        socket.OnConnect(HandlerOnConnect);

        //接收到后端发来的新弹幕
        socket.On("clientNewMsg", function (messageData) {
            HandlerOnGetNewMessage(messageData)
        });

        //断开连接，页面样式变化
        socket.OnDisconnect(function () {
            HandlerOnDisconnect();
        });
    }
    
    function HandlerOnConnect() {
        isConnect = true;
        socket.Emit("connected",roomNO);
        console.log("on connect",roomNO);
    }

    function HandlerOnGetNewMessage(messageData) {
        console.log("get new message",JSON.parse(messageData));
        var messages = JSON.parse(messageData);
        for ( var i = 0; i <messages.length; i++){
            NewBarage(messages[i].UserName,messages[i].Message);
        }
    }

    function HandlerOnDisconnect() {
        isConnect = false;
        console.log("on disconnect");
    }

    function NewBarage(userName,message) {
        barage.shoot(message);
        appendList(userName,message);
    }

    $("#send-message").on("click", function() {
        // var audio = document.getElementById("video-view");
        // var a = audio.duration;//播放时间
        // var b = audio.currentTime;//播放进度
        // console.log("time",a,b);

        if (isConnect){
            var text = document.querySelector('#text').value;

            data = {
                roomNO:roomNO,
                message:text,
                videoTime:123
            };
            socket.Emit("newMsg",JSON.stringify(data));

            // NewBarage("我",text);
        }else{
            alert("弹幕服务连接未建立");
        }
    });
})();