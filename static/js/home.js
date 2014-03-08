$(document).ready(function () {
    //fixHeader();
    topButton();
    renderMarkdown();
    initComment();
});

function fixHeader() {
    var $nav = $('#nav');
    var top = $nav.offset().top;
    console.log(top);
    $(window).scroll(function () {
        if (top < $(this).scrollTop()) {
            $nav.addClass("fixed").removeClass("text-center");
        } else {
            $nav.removeClass("fixed").addClass("text-center");
        }
    });
}

function topButton() {
    var top = $('#nav').offset().top;
    var $top = $('#go-top');
    $(window).scroll(function () {
        if (top < $(this).scrollTop()) {
            $top.removeClass("hide");
        } else {
            $top.addClass('hide');
        }
    });
    $top.on("click", function () {
        $('body,html').animate({scrollTop: 0}, 500);
        return false;
    })
}

function renderMarkdown() {
    var $md = $('.markdown');
    $md.each(function (i, item) {
        $(item).html(marked($(item).html().replace(/&gt;/g, '>')));
    });
    var code = $md.find('pre code');
    if (code.length) {
        $("<link>").attr({ rel: "stylesheet", type: "text/css", href: "/static/css/highlight.css"}).appendTo("head");
        $.getScript("/static/lib/highlight.min.js", function () {
            code.each(function (i, item) {
                hljs.highlightBlock(item)
            });
        });
    }
}

function initComment() {
    var $list = $('#comment-list');
    if (!$list.length) {
        return;
    }
    if (localStorage.getItem("comment-author")) {
        $('#comment-author').val(localStorage.getItem("comment-author"));
        $('#comment-email').val(localStorage.getItem("comment-email"));
        $('#comment-url').val(localStorage.getItem("comment-url"));
        $('#comment-avatar').attr("src", localStorage.getItem("comment-avatar"));
        $('.c-avatar').removeClass("null");
    }
    $('#comment-content').on("focus", function () {
        if ($('.c-avatar').hasClass("null")) {
            $('.c-avatar-field').remove();
            $('.c-info-fields').removeClass("hide");
        }
    });
    $('.not-me').on("click", function () {
        $('.c-avatar-field').remove();
        $('.c-info-fields').removeClass("hide");
        return false;
    });
    $('#comment-show').on("click", ".enable", function () {
        $("#comment-show").remove();
        $('#comment-form').removeClass("hide");
    });
    $('#comment-form').ajaxForm(function (json) {
        if (json.res) {
            localStorage.setItem("comment-author", $('#comment-author').val());
            localStorage.setItem("comment-email", $('#comment-email').val());
            localStorage.setItem("comment-url", $('#comment-url').val());
            localStorage.setItem("comment-avatar", json.comment.avatar);
            var tpl = $($('#comment-tpl').html());
            tpl.find(".c-avatar").attr("src", json.comment.avatar).attr("alt", json.comment.author);
            tpl.find(".c-author").attr("href", json.comment.url).text(json.comment.author);
            tpl.find(".c-reply").attr("rel", json.comment.id);
            tpl.find(".c-content").html(json.comment.content);
            if (json.comment.parent_md) {
                tpl.find(".c-p-md").html(marked(json.comment.parent_md));
            } else {
                tpl.find(".c-p-md").remove();
            }
            tpl.attr("id", "comment-" + json.comment.id);
            if (json.comment.status == "approved") {
                tpl.find(".c-check").remove();
            }
            $list.append(tpl);
            $('.cancel-reply').trigger("click");
            $('#comment-content').val("");
        } else {
            alert("提交失败!");
        }
    });
    $list.on("click", ".c-reply", function () {
        $('.reply-md').remove();
        var id = $(this).attr("rel");
        var pc = $('#comment-' + id);
        var md = "> @" + pc.find(".c-author").text() + "\n\n";
        md += "> " + pc.find(".c-content").html() + "\n";
        $('#comment-content').before('<div class="reply-md markdown">' + marked(md) + '</div>');
        $('#comment-parent').val(id);
        if ($('#comment-show').length) {
            $('#comment-show .enable').trigger("click");
        }
        $('.cancel-reply').show();
        var top = $('#comment-form').offset().top;
        $('body,html').animate({scrollTop: top}, 500);
        return false;
    });
    $('.cancel-reply').on("click", function () {
        $('.reply-md').remove();
        $('#comment-parent').val(0);
        $(this).hide();
        return false;
    });
    if (!$list.hasClass("comment-pager-true")) {
        return;
    }
    var isPager = false, index = 0;
    $list.find(".comment").each(function (i, item) {
        if (i >= 6 && !isPager) {
            isPager = true;
        }
        index = parseInt(i / 6 + 1);
        $(item).addClass("comment-index-" + index).hide();
    });
    if (!isPager) {
        $('.comment').show();
        return;
    }
    var html = ['<div id="comment-list-pager" class="pager">'];
    for (var i = 1; i <= index; i++) {
        if (i == index) {
            html.push('<li class="item current"><a href="#">' + i + '</a></li>');
            continue;
        }
        html.push('<li class="item"><a href="#">' + i + '</a></li>');
    }
    html.push('</div>');
    $('#comment-title').after(html.join(""));
    $('.comment-index-' + index).show();
    var $pager = $('#comment-list-pager');
    $pager.on("click", "a", function (e) {
        $('.comment').hide();
        console.log('.comment-index-' + $(e.target).text());
        $('.comment-index-' + $(e.target).text()).show();
        $pager.find(".current").removeClass("current");
        $(e.target).parent().addClass("current");
        return false;
    });
}
