define(['knockout-3.2.0'], function(ko){
        ko.bindingHandlers.intValue = {
                init: function(element, valueAccessor, allBindingsAccessor) {
                        var observable = valueAccessor();
                        var interceptor = ko.computed({
                                read: function() {
                                        return observable().toString();
                                },
                            write: function(newValue) {
                                    observable(parseInt(newValue, 10));
                            },
                        });
                        ko.applyBindingsToNode(element, { value: interceptor });
                }
        }
});
