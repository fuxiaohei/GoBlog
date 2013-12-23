$(document).ready(function () {
    "use strict";
    var $msg = $('#login-msg').hide();
    var $submit = $('#login-submit');
    function validate() {
        $msg.hide();
        $submit.hide();
        var $login = $('#login');
        if ($login.val().length < 2) {
            $('#login-error').show(200);
            $login.focus();
            $submit.show();
            return false;
        } else {
            $('#login-error').hide(200);
        }
        var $pwd = $('#password');
        if ($pwd.val().length < 4) {
            $('#password-error').show(200);
            $pwd.focus();
            $submit.show();
            return false;
        } else {
            $('#password-error').hide(200);
        }
        return true;
    }

    function parseJson(json) {
        if(json.res){
            window.location.href = "/admin/";
            return;
        }else{
            $submit.show(200);
            $msg.text(json.msg).show(200);
        }
    }

    $('#login-form').ajaxForm({
        beforeSubmit: validate,
        dataType: "json",
        success: function (json) {
            parseJson(json);
        }
    })
});
