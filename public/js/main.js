$(document).ready(function () {
    showMarkdownText();
});

function showMarkdownText() {
    var mdElements = $('.markdown-content');
    if (mdElements.length > 0) {
        var prettify = false;
        mdElements.each(function (i, item) {
            item.innerHTML = marked(item.innerHTML.toString().replace('[break]', ''));
            if (!prettify) {
                prettify = /<pre>/.test(item.innerHTML);
            }
        });
        if (prettify) {
            mdElements.find("pre").addClass("prettyprint").addClass("linenums");
            $("<link>").attr({ rel: "stylesheet",
                type: "text/css",
                href: "/public/js/prettify/prettify.css"
            }).appendTo("head");
            $.getScript("/public/js/prettify/prettify.js", function () {
                prettyPrint();
            });
        }
    }
}