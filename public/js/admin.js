$(document).ready(function(){
    showNavHover();
});

function showNavHover(){
    var $nav = $('#nav');
    var rel = $nav.data("rel");
    $nav.find('a[rel='+rel+']').addClass("hover");
}
