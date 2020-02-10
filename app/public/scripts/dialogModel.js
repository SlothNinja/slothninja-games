define(['knockout-3.2.0', 'jquery'], function(ko, $) {
        return function dialogModel(data) {
                var self = this;

                self.initialized = ko.observable(data.initialized);
                self.title = ko.observable(data.title);
                self.open = ko.observable(data.open);
                self.minWidth = ko.observable(data.minWidth);
                self.content = ko.observable(data.content);
        }
})
