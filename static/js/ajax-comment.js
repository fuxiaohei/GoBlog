var Comment = {};

Comment.InitReplyEvent = function () {
    var $form = $('#comment-form');
    var $noReply = $('#comment-no-reply');
    $('#comment-list').on("click", ".reply", function (e) {
        e.preventDefault();
        var $this = $(this);
        if ($this.parent().hasClass("comment-child")) {
            $this.parent().find(">section").after($form.detach());
        } else {
            $this.parent().after($form.detach());
        }
        $noReply.show();
        $('#comment-parent').val($this.attr("rel"));
    });
    $noReply.on("click", function (e) {
        e.preventDefault();
        $('#comment-container').append($form.detach());
        $noReply.hide();
        $('#comment-parent').val(0);
    })
};

Comment.InitFormEvent = function () {
    $('#author').val(localStorage.getItem("comment-author"));
    $('#email').val(localStorage.getItem("comment-email"));
    $('#site').val(localStorage.getItem("comment-site"));
    $('#comment-form').ajaxForm({
        beforeSubmit: function () {
            localStorage.setItem("comment-author", $('#author').val());
            localStorage.setItem("comment-email", $('#email').val());
            localStorage.setItem('comment-site', $('#site').val());
        },
        dataType: "json",
        success: function (json) {
            if (json.res) {
                var c = json.comment, html;
                c.CreateTime = (new Date(c.CreateTime * 1000)).format("MM.dd hh:mm");
                if (c.Pid > 0) {
                    html = juicer($('#comment-child-template').html(), c);
                    var $html = $(html);
                    $html.find(".parent").remove();
                    var $p = $('#comment-' + c.Pid);
                    if($p.hasClass("comment-item")){
                        $p.find(">p.meta").after($html);
                    }else{
                        $p.find(">section").after($html);
                    }
                } else {
                    html = juicer($('#comment-item-template').html(), c);
                    $('#comment-list').append(html);
                }
                $('#content').val("");
            }
        }
    });
};

Comment.InitCommentList = function () {
    var url = window.location.href;
    var $itemTemplate = juicer($('#comment-item-template').html());
    var $childTemplate = juicer($('#comment-child-template').html());
    $.getJSON(url, function (json) {
        if (json.res) {
            var childHtml = {};
            var topHtml = [];
            $(json.comments.reverse()).each(function (i, item) {
                item.CreateTime = (new Date(item.CreateTime * 1000)).format("MM.dd hh:mm");
                var $html;
                if (item.Pid < 1) {
                    $html = $($itemTemplate.render(item));
                    if (childHtml["c-" + item.Id]) {
                        $html.append(childHtml["c-" + item.Id].reverse())
                            .find("> .comment-child > .content > .parent").text("@" + item.Author);
                    }
                    topHtml.push($html);
                } else {
                    $html = $($childTemplate.render(item));
                    if (childHtml["c-" + item.Id]) {
                        $html.append(childHtml["c-" + item.Id].reverse())
                            .find("> .comment-child > .content > .parent").text("@" + item.Author);
                    }
                    if (!childHtml["c-" + item.Pid]) {
                        childHtml["c-" + item.Pid] = [];
                    }
                    childHtml["c-" + item.Pid].push($html);
                }
            });
            $('#comment-list').append(topHtml.reverse());
        }
    })
};

Date.prototype.format = function (format) {
    var o = {
        "M+": this.getMonth() + 1, //month
        "d+": this.getDate(), //day
        "h+": this.getHours(), //hour
        "m+": this.getMinutes(), //minute
        "s+": this.getSeconds(), //second
        "q+": Math.floor((this.getMonth() + 3) / 3), //quarter
        "S": this.getMilliseconds() //millisecond
    };
    if (/(y+)/.test(format)) format = format.replace(RegExp.$1,
        (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)if (new RegExp("(" + k + ")").test(format)) {
        format = format.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length));
    }
    return format;
};
