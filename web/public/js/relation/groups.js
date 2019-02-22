"use strict";

var addGroupUrl = "/relation/add/group";         //POST
var updateGroupUrl = "/relation/update/group";   //POST
var delGroupUrl = "/relation/del/group";         //POST


var lastActivedTag;


$(function() {
    $(".label_group").click(function () {
        if (lastActivedTag != undefined && lastActivedTag.length > 0){
            lastActivedTag.removeClass("label-primary");
            lastActivedTag.removeClass("actived_tag");
            lastActivedTag.addClass("label-info");
        }
        $(this).addClass("actived_tag");
        $(this).addClass("label-primary");
        $(this).removeClass("label-info");
        lastActivedTag = $(this)
    });

    $("#btn_update_group").click(function () {
        if ($(".actived_tag").length <= 0){
            alert("please choose span first");
        }

        var inputGroup = $("#input_update_group");
        if(inputGroup.val() == "" || inputGroup.val().length > 20){
            alert("please write suitable group name");
        }

        var data = {};
        data["groupId"] = $(".actived_tag").attr("group_id");
        data["groupName"] = inputGroup.val();

        $.ajax({
            type: "POST",
            url: updateGroupUrl,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    window.location.reload();
                }
                alert(res.Msg);
            }
        });

    });

    $("#btn_del_group").click(function () {
        if ($(".actived_tag").length <= 0){
            alert("please choose span first");
        }

        var data = {};
        data["groupId"] = $(".actived_tag").attr("group_id");

        $.ajax({
            type: "POST",
            url: delGroupUrl,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    window.location.reload();
                }
                alert(res.Msg);
            }
        });
    });

    $("#btn_add_group").click(function () {
        var inputGroup = $("#input_new_group");
        if(inputGroup.val() == "" || inputGroup.val().length > 20){
            alert("please write suitable group name");
        }

        var data = {};
        data["groupName"] = inputGroup.val();

        $.ajax({
            type: "POST",
            url: addGroupUrl,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    window.location.reload();
                }
                alert(res.Msg);
            }
        });
    });
});