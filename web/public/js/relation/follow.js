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
                    $("#group_sel_" + uid).val(oriGroupId);
                }
                btnLock = false;
            },
            error: function(){
                alert("error happend");
                $("#group_sel_" + uid).val(oriGroupId);
                btnLock = false;
            }
        });
    }
}

function refresh_sel(uid,groupId){
    $("#group_sel_" + uid).val(groupId);
    // $("#group_sel_" + uid).children("option").each(function(){
    //     alert($(this).text());
    //     if($(this).val() != groupId){
    //         $(this).removeAttr("selected");
    //     }else{
    //         $(this).attr("selected","selected");
    //     }
    // });
}

$(function () {
    $(".group-sel").change(function(){
        // alert(123);
        var uid = $(this).attr("uid");
        var group_id = $(this).val();
        var oriGroupId = $(this).attr("ori_group_id");
        setGroup(uid,group_id,oriGroupId);
    });

    $(".btn-follow").on("click",function(){
        if($("this").attr("status") == 1){
            unfollow($("this").attr("uid"));
        }else{
            follow($("this").attr("uid"));
        }
    })
});