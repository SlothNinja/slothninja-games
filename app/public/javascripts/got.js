function handler(data) {
        var $ajaxTarget = window.location.href;
        $.ajax({
                url: $ajaxTarget,
                dataType: "html",
                data: data,
                success: function(data) {
                        AjaxUpdate(data);
                        DisplayDialog(data);
                },
                type: "PUT",
                error: function(data) {
                        alert("Error");
                        window.location.href = target;
                },
        });
};
