define(['knockout'], function(ko){
        ko.bindingHandlers.booleanValue = {
                init: function(element, valueAccessor, allBindingsAccessor) {
                        var observable = valueAccessor(),
                            interceptor = ko.computed({
                                read: function() {
                                    return observable().toString();
                                },
                                write: function(newValue) {
                                    observable(newValue === "true");
                                }                   
                            });

                        ko.applyBindingsToNode(element, { value: interceptor });
                }
        };
});
