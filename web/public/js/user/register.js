"use strict";

var register_url = "/user/register";      //POST

function register_submit(){
	var data = {};
	data["nickName"] = $("#nickName").val();
	data["email"] = $("#email").val();
	data["telephone"] = $("#telephone").val();
	data["password"] = $("#password").val();
	var password_confirm = $("#password_confirm").val();
	if(data["password"] != password_confirm){
		alert("密码不一致");
		return false;
	}

	$.ajax({
		type: "POST",
		url: register_url,
		data: data,
		dataType: "json",
		success: function(res){
			console.log(res);
			if(res.Code == 1000){
				alert("注册成功,尝试登录");
				window.location.href = "/user/login";
			}
			alert(res.Msg);
		}
	});
}

$(function() {	
	$("#register_submit").on("click", function() {
		register_submit();
	});
});