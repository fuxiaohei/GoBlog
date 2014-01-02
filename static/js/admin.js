function initArticleTable() {
    $('.delete').on("click", function () {
        return confirm("删除文章将清除相关标签评论数据？确认继续？");
    });
}

function initArticleForm() {
    $("#content").markdown({
        autofocus: false,
        savable: false,
        onPreview: function (e) {
            var content = marked(e.getContent(), {
                breaks: true
            });
            if (content.indexOf("</pre>") >= 1) {
                setTimeout(function () {
                    $(".md-preview pre").addClass("prettyprint").addClass("linenums");
                    prettyPrint();
                }, 200);
            }
            return content;
        }
    }).autosize({"append": "\n"});
    $('#article-form').ajaxForm({
        dataType: "json",
        success: function (json) {
            if (json.res) {
                alert("保存完成");
                window.location.href = "/admin/article/";
            } else {
                alert(json.msg);
            }
        }
    })
}

function initCategoryTable() {
    $('#category-table-form').on("submit", function () {
        var category = $('input[name=category]:checked').val();
        if (category == $('#move').val()) {
            alert("不能移动到自己");
            return false;
        }
        return confirm("删除将不可恢复!");
    });
}

function initCommentList() {
    $('#comment-list-container').on("click", ".approve", function (e) {
        e.preventDefault();
        var id = $(this).attr("rel");
        $.post("/admin/comment/status", {"id": id, "status": "approved"}, function (json) {
            if (json.res) {
                var $c = $('#comment-' + id);
                $c.find(".author").removeClass("gray").addClass("blue");
                $c.find(".approve").remove();
            } else {
                if (json.msg) {
                    alert(json.msg);
                }
            }
        });
    })
        .on("click", ".spam",function (e) {
            e.preventDefault();
            var id = $(this).attr("rel");
            $.post("/admin/comment/status", {"id": id, "status": "spam"}, function (json) {
                if (json.res) {
                    var $c = $('#comment-' + id);
                    $c.find(".author").removeClass("blue").addClass("gray");
                    $c.find(".spam").remove();
                } else {
                    if (json.msg) {
                        alert(json.msg);
                    }
                }
            });
        }).on("click", ".del", function (e) {
            if (!confirm("删除评论将使回复无效？确认继续？")) {
                return false;
            }
            e.preventDefault();
            var id = $(this).attr("rel");
            $.post("/admin/comment/delete", {"id": id}, function (json) {
                if (json.res) {
                    $('#comment-' + id).remove();
                } else {
                    if (json.msg) {
                        alert(json.msg);
                    }
                }
            });
        })
        .on("click", ".reply", function (e) {
            e.preventDefault();
            var id = $(this).attr("rel");
            $('#comment-' + id).append($('#comment-reply-form').detach().show());
            $('#comment-parent').val(id);
        });
    $('#comment-reply-form').ajaxForm({
        success: function (json) {
            if (json.res) {
                var comment = json.comment;
                var html = '<div class="comment-reply">';
                html += '<span class="reply-author">' + comment.Author + ': </span>';
                html += comment.Content + '</div>';
                $('#comment-' + comment.Pid).find("section").after(html);
            } else {
                alert(json.msg);
            }
        }
    });
    $('#comment-reply-hide').on("click", function () {
        $('#comment-reply-form').hide();
    })
}