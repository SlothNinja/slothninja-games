define(['jquery'], function($) {
        return function put(json, update) {
                var ajaxTarget = window.location.pathname;
                $.ajax({
                        url: ajaxTarget,
                        data: json, 
                        dataType: "json",
                        success: function(data) {
                                update(data);
                        },
                        type: "PUT",
                        error: function(data) {
                                alert("Error");
                                window.location.pathname = ajaxTarget;
                        }
                }); 
        }
})
