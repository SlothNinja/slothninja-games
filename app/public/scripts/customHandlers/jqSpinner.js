define(['knockout-3.2.0', 'jquery', 'jquery-ui/spinner'], function(ko, $, spinner){
        ko.bindingHandlers.jqSpinner = {
                init: function(element, valueAccessor, allBindings, viewModel, bindingContext) {
                        var value = valueAccessor();
                        var min = 0;
                        var max = 100;
                        var $el = $(element);
                        if (allBindings.has("min")) {
                                min = allBindings.get("min");
                        };
                        if (allBindings.has("max")) {
                                max = allBindings.get("max");
                        };
                        $el.css({width: "20px"}).spinner({min: min, max: max});
                        $el.spinner("value", value());
                        $el.on( "spinstop", function() {
                                if (allBindings.has("update")) {
                                        var myUpdate = allBindings.get("update");
                                        myUpdate(element, valueAccessor, allBindings, viewModel, bindingContext);
                                } else {
                                        value($el.spinner("value"));
                                }
                        });
                },
                update: function(element, valueAccessor, allBindings, viewModel, bindingContext) {
                        if (allBindings.has("update")) {
                                var myUpdate = allBindings.get("update");
                                myUpdate(element, valueAccessor, allBindings, viewModel, bindingContext);
                        } else {
                                var value = valueAccessor();
                                value($(element).spinner("value"));
                        }
                }
        }
});
