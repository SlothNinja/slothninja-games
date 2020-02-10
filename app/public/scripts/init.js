require(['knockout-3.2.0', 'jquery', 'lodash', 'viewModel', 'customHandlers/jqButton', 'customHandlers/jqTabs', 'customHandlers/jqDialog', 'customHandlers/jqSpinner', 'customHandlers/booleanValue', 'customHandlers/intValue', 'customHandlers/imagemapster'], function(ko, $, _, viewModel) {
        var target = window.location.pathname + "/json";
        $.ajax({
                dataType: "json",
                url: target,
                data: { action: "json" },
                success: function(data) {
                        new viewModel(data).init();
                }, 
        });
});
