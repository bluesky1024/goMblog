"use strict";

var feed_more_url = "/feed/more";       //POST
var feed_newer_url = "feed/newer";      //POST
var first_mid = 0;
var last_mid = 0;
var cur_group = 0;

var is_scroll_valid = true;  //防止请求过程中反复获取数据

function render_new_card(mblog){
    var mblog_element = $("#mblog-demo").clone();
    mblog_element.show();
    mblog_element.attr("id","mblog-" + mblog.Mid);
    mblog_element.attr("mid",mblog.Mid);
    mblog_element.attr("uid",mblog.Uid);
    mblog_element.find(".mblog-base-info").html(mblog.NickName + "   " + mblog.CreateTime);
    mblog_element.find(".mblog-content").html(mblog.Content);
    mblog_element.find(".mblog-trans").html(mblog.TransCnt);
    mblog_element.find(".mblog-comment").html(mblog.CommentCnt);
    mblog_element.find(".mblog-likes").html(mblog.LikesCnt);
    return mblog_element;
}

/*在feed尾部扩充微博*/
function append_feed_cards(mblog_datas,user_datas){
    $.each(mblog_datas,function (ind,mblog) {
        mblog["NickName"] = user_datas[mblog.Uid]["NickName"];
        var tempHtml = render_new_card(mblog)
        $("#feed-block").append(tempHtml);
        if(first_mid == 0){
            first_mid = mblog.Mid;
        }
        last_mid = mblog.Mid;
    });
}

/*在feed顶部插入微博*/
function prepend_feed_cards(mblog_datas,user_datas){
    $.each(mblog_datas,function (ind,mblog) {
        mblog["NickName"] = user_datas[mblog.Uid]["NickName"];
        var tempHtml = render_new_card(mblog)
        $("#feed-block").prepend(tempHtml);
        first_mid = mblog.Mid;
    });
}

function get_more_feed(){
    var data = {};
    data["lastMid"] = last_mid;
    data["groupId"] = cur_group;

    if(is_scroll_valid){
        is_scroll_valid = false;

        $.ajax({
            type: "POST",
            url: feed_more_url,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    append_feed_cards(res.Data.MblogsInfo,res.Data.UsersInfo)
                }else{
                    alert(res.Msg);
                }
                is_scroll_valid = true;
            },
            error: function(){
                alert("error happend");
                is_scroll_valid = true;
            }
        });
    }
}

function get_newer_feed() {
    var data = {};
    data["firstMid"] = first_mid;
    data["groupId"] = cur_group;

    if(is_scroll_valid){
        is_scroll_valid = false;

        $.ajax({
            type: "POST",
            url: feed_newer_url,
            data: data,
            dataType: "json",
            success: function(res){
                console.log(res);
                if(res.Code == 1000){
                    prepend_feed_cards(res.Data.MblogsInfo,res.Data.UsersInfo)
                }else{
                    alert(res.Msg);
                }
                is_scroll_valid = true;
            },
            error: function(){
                alert("error happend");
                is_scroll_valid = true;
            }
        });
    }
}

/*切换分组*/
function switch_group(group_id){
    if(group_id == cur_group){
        return;
    }

    //清空变量
    cur_group = group_id;
    first_mid = 0;
    last_mid = 0;

    //清空feed列表
    $("#feed-block").html("");

    //重新拉取
    get_more_feed();
}
/*切换分组*/


$(function() {
    $(window).bind("scroll",function(){
        var win_height=$(window).height();
        var doc_height=$(document).height();
        var scroll_top=$(document).scrollTop();

        //获取更新数据
        if(scroll_top == 0){
            get_newer_feed();
        }

        //获取更旧数据
        if(scroll_top + win_height >= doc_height){
            get_more_feed();
        }
    });

    $("#btn-get-more-feed").on("click", function() {
        // $("#feed-block").html("");
        // get_more_feed(cur_page,cur_group);
        var mblog = {};
        mblog['NickName'] = "fang";
        mblog['Content'] = "content";
        mblog['TransCnt'] = "1";
        mblog['CommentCnt'] = "2";
        mblog['LikesCnt'] = "3";
        $("#feed-block").append(render_new_card(mblog));

        mblog['NickName'] = "fang2";
        mblog['Content'] = "content";
        mblog['TransCnt'] = "11";
        mblog['CommentCnt'] = "22";
        mblog['LikesCnt'] = "33";
        $("#feed-block").append(render_new_card(mblog));
    });

    $("#group-sel").change(function(){
        var group_id = $(this).val();
        switch_group(group_id);
    });

    $(document).ready(function(){
        get_more_feed();
    });
});