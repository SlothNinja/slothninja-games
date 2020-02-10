define(['knockout-3.2.0', 'jquery', 'jquery-ui/button'], function(ko, $, button){
        ko.bindingHandlers.jqButton = {
                init: function(element) {
                        $(element).button();
                }
        }
});
