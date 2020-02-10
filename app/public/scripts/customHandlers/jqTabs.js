define(['knockout-3.2.0', 'jquery', 'jquery-ui/tabs'], function(ko, $, tabs){
        ko.bindingHandlers.jqTabs = {
                init: function(element) {
                        $(element).tabs({ selected: 0});
                }
        };
});
