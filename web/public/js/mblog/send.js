"use strict";

var mblogSendUrl = "/mblog/send";      //POST

function mblogSend(){
    var data = {};
    data["content"] = $("#content").val();
    data["readAble"] = $("#read_able").val();

    $.ajax({
        type: "POST",
        url: mblogSendUrl,
        data: data,
        dataType: "json",
        success: function(res){
            console.log(res);
            if(res.Code == 1000){
                alert("发布成功");
                window.location.href = "/user/me";
                return true;
            }
            alert(res.Msg);
        }
    });
}

$(function() {
    $("#btn_send").on("click", function() {
        mblogSend();
    });
});