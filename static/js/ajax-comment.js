var Comment = {};

Comment.InitReplyEvent = function () {
    var $form = $('#comment-form');
    var $noReply = $('#comment-no-reply');
    $('#comment-list').on("click", ".reply", function (e) {
        e.preventDefault();
        var $this = $(this);
        $this.parent().after($form.detach());
        $noReply.show();
    });
    $noReply.on("click",function(e){
        e.preventDefault();
        $('#comment-container').append($form.detach());
        $noReply.hide();
    })
}
