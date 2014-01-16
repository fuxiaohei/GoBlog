function initUpload(p) {
    $('#attach-show').on("click", function () {
        $('#attach-upload').trigger("click");
    });
    $('#attach-upload').on("change", function () {
        if (confirm("立即上传?")) {
            var bar = $('<p class="file-progress inline-block">0%</p>');
            $('#attach-form').ajaxSubmit({
                "beforeSubmit": function () {
                    $(p).before(bar);
                },
                "uploadProgress": function (event, position, total, percentComplete) {
                    var percentVal = percentComplete + '%';
                    bar.css("width", percentVal).html(percentVal);
                },
                "success": function (json) {
                    if (!json.res) {
                        bar.html(json.msg).addClass("err");
                        setTimeout(function () {
                            bar.remove();
                        }, 5000);
                    } else {
                        bar.html("/" + json.file.Url + "&nbsp;&nbsp;&nbsp;(@" + json.file.Name + ")");
                    }
                    $('#attach-upload').val("");
                }
            });
        } else {
            $(this).val("");
        }
    });
}