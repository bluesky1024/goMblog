"use strict";

var login_url = "/user/login";      //POST

function login(){
    var data = {};
    data["nickName"] = $("#nickName").val();
    // data["email"] = $("#email").val();
    // data["telephone"] = $("#telephone").val();
    data["password"] = $("#password").val();

    $.ajax({
        type: "POST",
        url: login_url,
        data: data,
        dataType: "json",
        success: function(res){
            console.log(res);
            if(res.Code == 1000){
                window.location.href = "/user/me";
            }
        }
    });
}

$(function() {
    $("#btn_login").on("click", function() {
        login();
    });
});