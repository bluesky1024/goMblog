"use strict";

var follow_url = "/relation/follow";      //POST
var unfollow_url = "/relation/unfollow";  //POST

var btn_follow_lock = false;

function follow(){
    if(!btn_follow_lock) {
        btn_follow_lock = true;
        var data = {};
        data["status"] = $("#btn_follow").attr("status");
        data["uid"] = $("#btn_follow").attr("uid");
        console.log(data["status"]);
        if (data["status"] == 1) {
            if (confirm("确认取关？")){
                $.ajax({
                    type: "POST",
                    url: unfollow_url,
                    data: data,
                    dataType: "json",
                    success: function(res){
                        console.log(res);
                        if(res.Code == 1000){
                            $("#btn_follow").attr("status",0);
                            $("#btn_follow").html("关注");
                        }else{
                            alert(res.Msg)
                        }
                        btn_follow_lock = false;
                    }
                });
            }else{
                btn_follow_lock = false;
            }
        } else {
            $.ajax({
                type: "POST",
                url: follow_url,
                data: data,
                dataType: "json",
                success: function(res){
                    console.log(res);
                    if(res.Code == 1000){
                        $("#btn_follow").attr("status",1);
                        $("#btn_follow").html("取关");
                    }else{
                        alert(res.Msg)
                    }
                    btn_follow_lock = false;
                }
            });
        }
    }
}


$(function() {
    $("#btn_follow").on("click", function() {
        follow();
    });
});