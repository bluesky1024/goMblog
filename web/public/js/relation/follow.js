"use strict";

var followUrl = "/relation/follow";      //POST
var unfollowUrl = "/relation/unfollow";  //POST
var setGroupUrl = "/relation/set/group"; //POST

var lastGroupSel = 0;

var btnLock = false;

function follow(uid){
    if(!btnLock){
        btnLock = true;

        var data = {};
        data["status"] = 1;
        data["uid"] = uid;
        $.ajax({
            type: "POST",
            url: followUrl,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    $("#btn_follow_" + uid).attr("status",1);
                    $("#btn_follow_" + uid).html("取关");
                }else{
                    alert(res.Msg)
                }
                btnLock = false;
            }
        });
    }
}

function unfollow(uid){
    if(!btnLock){
        btnLock = true;

        var data = {};
        data["status"] = 0;
        data["uid"] = uid;
        $.ajax({
            type: "POST",
            url: unfollowUrl,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    $("#btn_follow_" + uid).attr("status",0);
                    $("#btn_follow_" + uid).html("关注");
                }else{
                    alert(res.Msg)
                }
                btnLock = false;
            }
        });
    }
}

function setGroup(uid,groupId,oriGroupId){
    if(!btnLock){
        btnLock = true;

        var data = {};
        data["uidFollow"] = uid;
        data["groupId"] = groupId;
        $.ajax({
            type: "POST",
            url: setGroupUrl,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    $("#group_sel_" + uid).attr("ori_group_id",groupId)
                }else{
                    alert(res.Msg);
                    var oriStr = "option[value=" + oriGroupId + "]";
                    var newStr = "option[value=" + groupId + "]";
                    $("#group_sel_" + uid).find(newStr).attr("selected",false);
                    $("#group_sel_" + uid).find(oriStr).attr("selected",true);
                }
                btnLock = false;
            },
            error: function(){
                alert("error happend");
                var oriStr = "option[value=" + oriGroupId + "]";
                var newStr = "option[value=" + groupId + "]";
                $("#group_sel_" + uid).find(newStr).attr("selected",false);
                $("#group_sel_" + uid).find(oriStr).attr("selected",true);
                btnLock = false;
            }
        });
    }
}

$(function () {
    $(".group-sel").change(function(){
        // var uid = $(this).attr("uid");
        // var group_id = $(this).val();
        // var oriGroupId = $(this).attr("ori_group_id");
        // setGroup(uid,group_id,oriGroupId);
    });

    $(".btn-follow").on("click",function(){
        if($("this").attr("status") == 1){
            unfollow($("this").attr("uid"));
        }else{
            follow($("this").attr("uid"));
        }
    })
});